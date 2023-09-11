package main

import (
	"flag"
	"fmt"

	"d1y.io/neovideo/config"
	baseControllers "d1y.io/neovideo/controllers/base"
	"github.com/kataras/iris/v12"
)

type NeovideoApp struct {
	App *iris.Application
}

func (na *NeovideoApp) Init() {
	var confPath = flag.String("conf", "./config/conf.example.yml", "config file path")
	flag.Parse()
	if len(*confPath) <= 0 {
		panic("config file path is required")
	}
	_, err := config.InitWithFile(*confPath)
	if err != nil {
		panic(err)
	}
	na.App = iris.New()
	na.Register()
}

func (na *NeovideoApp) Register() {
	na.App.PartyFunc("/api/v1", func(u iris.Party) {
		u.PartyFunc("/base", func(u iris.Party) {
			baseControllers.Register(u)
		})
	})
}

func (na *NeovideoApp) getPort() string {
	port := config.Get().Port
	return fmt.Sprintf(":%d", port)
}

func (na *NeovideoApp) Run() {
	na.App.Listen(na.getPort())
}

func newNeovideoApp() *NeovideoApp {
	return &NeovideoApp{}
}
