package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/igungor/telegram"
)

// flags
var (
	flagConfig = flag.String("c", "./mercimek.conf", "configuration file path")
)

var cfg *config

func usage() {
	fmt.Fprintf(os.Stderr, "mercimek is a niche Telegram bot. It counts mercimek. Yep.\n\n")
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  mercimek -c <path of mercimek.conf>\n\n")
	fmt.Fprintf(os.Stderr, "flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetPrefix("mercimek: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Usage = usage
	flag.Parse()

	var err error
	cfg, err = readConfig(*flagConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}

	bot := telegram.New(cfg.Token)
	err = bot.SetWebhook(cfg.Webhook)
	if err != nil {
		log.Fatalf("error while setting webhook: %v", err)
	}

	log.Printf("Webhook set to %v\n", cfg.Webhook)

	http.HandleFunc("/", bot.Handler())

	go func() {
		addr := net.JoinHostPort(cfg.Host, cfg.Port)
		log.Fatal(http.ListenAndServe(addr, nil))
	}()

	for msg := range bot.Messages() {
		if msg.IsService() {
			continue
		}

		go handleMercimek(bot, msg)
	}
}

type config struct {
	Token   string `json:"token"`
	Webhook string `json:"webhook"`
	Host    string `json:"host"`
	Port    string `json:"port"`

	BinaryPath          string `json:"binary-path"`
	ParticleSize        string `json:"particle-size"`
	ParticleCircularity string `json:"particle-circularity"`
}

func readConfig(configpath string) (*config, error) {
	f, err := os.Open(configpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("token field can not be empty")
	}
	if cfg.Webhook == "" {
		return nil, fmt.Errorf("webhook field can not be empty")
	}
	if cfg.BinaryPath == "" {
		return nil, fmt.Errorf("binary-path can not be empty")
	}

	_, err = os.Stat(cfg.BinaryPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("binary could not be found: %v", err)
	}

	return &cfg, nil
}
