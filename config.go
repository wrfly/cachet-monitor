package main

import (
	"time"

	"github.com/wrfly/cachet-monitor/cachet"
)

type Target struct {
	Component    string
	URL          string
	ExpectedCode int `yaml:"code"`
	Interval     time.Duration
}

type Config struct {
	Cachet  cachet.Config
	Targets []Target
}
