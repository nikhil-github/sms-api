package main

import (
	"github.com/nikhil-github/sms-api/pkg/wiring"
)

func main() {
	var cfg *wiring.Config
	a := wiring.App{
		Config: cfg,
	}
	a.Run()
}
