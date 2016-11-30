package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/igungor/tlbot"
)

var httpclient = http.Client{Timeout: 30 * time.Second}

func handleMercimek(b *tlbot.Bot, msg *tlbot.Message) {
	opts := &tlbot.SendOptions{}

	var fileID string
	var err error
	if len(msg.Photos) != 0 {
		// last photo has the max resolution
		fileID = msg.Photos[len(msg.Photos)-1].FileID
	} else if msg.Document.Exists() {
		if !strings.HasPrefix(msg.Document.MimeType, "image") {
			err = fmt.Errorf("gonderdigin dosya fotografa benzemiyor")
		} else {
			fileID = msg.Document.FileID
		}
	} else {
		err = fmt.Errorf("mercimeklerin fotografini cekip gonder")
	}

	if err != nil {
		opts.ParseMode = tlbot.ModeMarkdown
		_, _ = b.SendMessage(msg.Chat.ID, err.Error(), opts)
		return
	}

	u, err := b.GetFileDownloadURL(fileID)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := httpclient.Get(u)
	if err != nil {
		errmsg := fmt.Sprintf("gonderdigin dosyayi su sebepten indiremedim: %v", err)
		_, _ = b.SendMessage(msg.Chat.ID, errmsg, nil)
		return
	}
	defer resp.Body.Close()

	tempdir, err := ioutil.TempDir("", "mercimek-")
	if err != nil {
		errmsg := fmt.Sprintf("bir takim hatalar: %v", err)
		_, _ = b.SendMessage(msg.Chat.ID, errmsg, nil)
		return
	}

	origImage, err := os.OpenFile(filepath.Join(tempdir, "orig.jpg"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		errmsg := fmt.Sprintf("bir takim hatalar: %v", err)
		_, _ = b.SendMessage(msg.Chat.ID, errmsg, nil)
		return

	}
	defer origImage.Close()

	const macroName = "macro.ijm"
	resultImagePath := filepath.Join(tempdir, "result.jpg")

	macro, err := os.OpenFile(filepath.Join(tempdir, macroName), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		errmsg := fmt.Sprintf("bir takim hatalar: %v", err)
		_, _ = b.SendMessage(msg.Chat.ID, errmsg, nil)
		return
	}
	defer macro.Close()

	_, err = io.Copy(origImage, resp.Body)
	if err != nil {
		errmsg := fmt.Sprintf("bir takim hatalar: %v", err)
		_, _ = b.SendMessage(msg.Chat.ID, errmsg, nil)
		return
	}

	tmpl, err := template.New(macroName).Parse(macroTemplate)
	if err != nil {
		errmsg := fmt.Sprintf("template ile ilgili bi hata yaptim: %v", err)
		_, _ = b.SendMessage(msg.Chat.ID, errmsg, nil)
		return
	}

	r := struct {
		ImagePath           string
		ResultImagePath     string
		ParticleSize        string
		ParticleCircularity string
	}{
		ImagePath:           origImage.Name(),
		ResultImagePath:     resultImagePath,
		ParticleSize:        config.ParticleSize,
		ParticleCircularity: config.ParticleCircularity,
	}

	err = tmpl.Execute(macro, r)
	if err != nil {
		errmsg := fmt.Sprintf("template ile ilgili bi hata yaptim: %v", err)
		_, _ = b.SendMessage(msg.Chat.ID, errmsg, nil)
		return
	}

	count, err := executeMacro(macro.Name())
	if err != nil {
		errmsg := fmt.Sprintf("mercimekleri sayamadim cunku: %v", err)
		_, _ = b.SendMessage(msg.Chat.ID, errmsg, nil)
		return
	}

	resultImage, err := os.OpenFile(resultImagePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		errmsg := fmt.Sprintf("bir takim hatalar: %v", err)
		_, _ = b.SendMessage(msg.Chat.ID, errmsg, nil)
		return
	}
	defer resultImage.Close()

	photo := tlbot.Photo{
		File: tlbot.File{
			Body: resultImage,
			Name: resultImage.Name(),
		},
		Caption: count,
	}
	_, err = b.SendPhoto(msg.Chat.ID, photo, opts)
	if err != nil {
		errmsg := fmt.Sprintf("hersey hazirdi ama son anda bir hata olustu ya: %v", err)
		_, _ = b.SendMessage(msg.Chat.ID, errmsg, nil)
		return

	}
}

func executeMacro(macroPath string) (string, error) {
	cmd := exec.Command(config.BinaryPath, "--headless", "--console", "-macro", macroPath)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	s := strings.TrimSpace(buf.String())
	lines := strings.Split(s, "\n")
	lastLine := lines[len(lines)-1]
	words := strings.Split(lastLine, "\t")

	return words[0], nil
}

var macroTemplate = `
open("{{.ImagePath}}");

run("Subtract Background...", "rolling=50");
run("Enhance Local Contrast (CLAHE)", "blocksize=127 histogram=256 maximum=3 mask=*None*");
run("8-bit");
setAutoThreshold("Default dark");
call("ij.plugin.frame.ThresholdAdjuster.setMode", "Red");
//run("Threshold...");
setOption("BlackBackground", false);
run("Convert to Mask");
run("Remove Outliers...", "radius=2 threshold=50 which=Dark");
run("Watershed");
run("Analyze Particles...", "size={{.ParticleSize}} circularity={{.ParticleCircularity}} show=Overlay display exclude clear");

saveAs("jpeg", "{{.ResultImagePath}}");
`
