package main

func main() {
	app := newNeovideoApp()
	app.Init()
	defer app.GetDB().Close()
	app.Run()
}
