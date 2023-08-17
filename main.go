package main

import (
	"orangeadd.com/clipboard-client/client"
	"sync"
)

func main() {
	err := client.MustInit()
	if err != nil {
		return
	}
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	client.ListenClipboardText()
	client.ListenClipboardImage()

	waitGroup.Wait()
}
