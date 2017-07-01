package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/igungor/tlbot"
)

// flags
var (
	flagConfig = flag.String("c", "./mercimek.conf", "configuration file path")
)

var config *Config

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
	config, err = readConfig(*flagConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}

	b := tlbot.New(config.Token)
	err = b.SetWebhook(config.Webhook)
	if err != nil {
		log.Fatalf("error while setting webhook: %v", err)
	}

	log.Printf("Webhook set to %v\n", config.Webhook)

	messages := b.Listen(net.JoinHostPort(config.Host, config.Port))
	for msg := range messages {
		if msg.IsService() {
			continue
		}

		go handleMercimek(&b, &msg)
	}
}

type Config struct {
	Token   string `json:"token"`
	Webhook string `json:"webhook"`
	Host    string `json:"host"`
	Port    string `json:"port"`

	BinaryPath          string `json:"binary-path"`
	ParticleSize        string `json:"particle-size"`
	ParticleCircularity string `json:"particle-circularity"`
}

func readConfig(configpath string) (config *Config, err error) {
	f, err := os.Open(configpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	if config.Token == "" {
		return nil, fmt.Errorf("token field can not be empty")
	}
	if config.Webhook == "" {
		return nil, fmt.Errorf("webhook field can not be empty")
	}
	return config, nil
}
