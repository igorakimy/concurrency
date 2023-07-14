package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	type value struct {
		mu sync.Mutex
		value int
	}

	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		// Здесь мы пытаемся войти в критическую секцию входящего значения
		v1.mu.Lock()
		// Здесь мы используем конструкцию defer чтобы выйти из критической 
		// секции до возвращения printSum
		defer v1.mu.Unlock()

		// Здесь мы спим определенное кол-во времени, чтобы симулировать 
		// работу и вызвать deadlock
		time.Sleep(2 * time.Second)
		v2.mu.Lock()
		defer v2.mu.Unlock()
		
		fmt.Printf("sum=%v\n", v1.value + v2.value)
	}

	var a, b value
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}