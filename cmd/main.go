package main

import (
	"flag"
	"github.com/hyperjumptech/httptarget/model"
	"github.com/hyperjumptech/httptarget/server"
	"github.com/sirupsen/logrus"
)

var (
	portFlag = flag.Int("p", 51423, "Listen port")
	hostFlag = flag.String("h", "0.0.0.0", "Bind host")

	bodyFlag     = flag.String("body", "OK", "Bind host")
	pathFlag     = flag.String("path", "/", "Base path")
	codeFlag     = flag.Int("code", 200, "Response code")
	minDelayFlag = flag.Int("minDelay", 0, "Minimum Delay Millisecond")
	maxDelayFlag = flag.Int("maxDelay", 200, "Maximum Delay Millisecond")
)

func main() {
	flag.Parse()

	initEp := &model.EndPoint{
		ID:            0,
		BasePath:      *pathFlag,
		DelayMinMs:    *minDelayFlag,
		DelayMaxMs:    *maxDelayFlag,
		ReturnCode:    *codeFlag,
		ReturnHeaders: map[string][]string{"Content-Type": {"text/plain"}},
		ReturnBody:    *bodyFlag,
	}

	err := server.Start(*hostFlag, *portFlag, initEp)
	if err != nil {
		logrus.Error(err)
	}
}
