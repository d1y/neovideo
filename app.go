package main

import (
	"database/sql"
	"flag"
	"fmt"

	"d1y.io/neovideo/config"
	"d1y.io/neovideo/controllers/handler"
	jiexiControllers "d1y.io/neovideo/controllers/jiexi"
	maccmsControllers "d1y.io/neovideo/controllers/maccms"
	"d1y.io/neovideo/sqls"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
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
	na.initDB()
	na.App = iris.New()
	na.Register()
}

func (na *NeovideoApp) initDB() {
	if err := sqls.Open(config.Get().Db, &gorm.Config{}); err != nil {
		panic(err)
	}
	sqls.AutoMigrate()
}

func (na *NeovideoApp) GetDB() *sql.DB {
	return sqls.RealDb()
}

func (na *NeovideoApp) Register() {
	na.App.Get("/siteinfo", func(i iris.Context) {
		routeMeta := na.App.GetRoutesReadOnly()
		handler.Siteinfo(i, routeMeta)
	})
	na.App.Get("/ping", handler.Ping)
	na.App.PartyFunc("/api/v1", func(u iris.Party) {
		u.PartyFunc("/maccms", maccmsControllers.Register)
		u.PartyFunc("/jiexi", jiexiControllers.Register)
	})
}

func (na *NeovideoApp) getPort() string {
	port := config.Get().Port
	if port == 0 {
		panic("port is required")
	}
	return fmt.Sprintf(":%d", port)
}

func (na *NeovideoApp) Run() {
	na.App.Listen(na.getPort())
}

func newNeovideoApp() *NeovideoApp {
	return &NeovideoApp{}
}
