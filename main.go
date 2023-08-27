package main

func main() {
	MustInit()
	go InitSystemTray()
	InitServer()
}
