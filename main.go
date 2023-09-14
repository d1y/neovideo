package main

import "d1y.io/neovideo/sqls"

func main() {
	app := newNeovideoApp()
	app.Init()
	defer sqls.Close()
	app.Run()
}
