package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/wrfly/cachet-monitor/cachet"

	"gopkg.in/yaml.v2"
)

func main() {
	var config = flag.String("config", "config.yml", "config file path")
	flag.Parse()

	bs, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Fatalf("read config file %s err: %s", *config, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(bs, &cfg); err != nil {
		log.Fatalf("parse config file %s err: %s", *config, err)
	}

	// connect to cache
	if err := cachet.Init(cfg.Cachet); err != nil {
		log.Fatalf("init cachet client err: %s", err)
	}

	// do the real work
	for _, target := range cfg.Targets {
		go check(target)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	log.Printf("interrupt signal received, exit")
}

func check(t Target) {
	_check := func(t Target) {
		log.Printf("checking %s", t.Component)
		resp, err := http.Get(t.URL)
		if err != nil {
			log.Printf("check %s err: %s", t.URL, err)
			cachet.UpdateStatus(t.Component, cachet.StateMajorOutage)
			return
		}

		if resp.StatusCode != t.ExpectedCode {
			log.Printf("check %s err: %d != %d", t.Component, resp.StatusCode, t.ExpectedCode)
			cachet.UpdateStatus(t.Component, cachet.StateMajorOutage)
		} else {
			cachet.UpdateStatus(t.Component, cachet.StateOperational)
			log.Printf("%s is operational as excepted", t.Component)
		}
	}

	_check(t)
	for range time.Tick(t.Interval) {
		_check(t)
	}
}
