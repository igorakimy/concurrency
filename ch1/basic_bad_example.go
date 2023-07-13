package main

import (
	"fmt"
	"time"
)

func main() {
	var data int
	go func() {
		data++
	}()
	// Это плохо!
	time.Sleep(1 * time.Second)
	if data == 0 {
		fmt.Printf("the value is %v\n", data)
	}
}
