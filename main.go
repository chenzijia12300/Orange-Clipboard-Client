package main

import (
	"sync"
)

func main() {
	err := MustInit()
	if err != nil {
		return
	}
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	waitGroup.Wait()
}
