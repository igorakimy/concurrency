package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	cadence := sync.NewCond(&sync.Mutex{})
	go func() {
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}

	// tryDir - попытка движения человека в каком-либо направлении и
	// возврат на исходную позицию. Каждое направление представлено
	// в виде кол-ва попыток движения в этом направлении
	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool {
		fmt.Fprintf(out, " %v", dirName)
		// Объявить о намерении двигаться в каком-то направлении
		atomic.AddInt32(dir, 1)
		// Чтобы продемонстрировать livelock, каждый человек должен
		// двигаться синхронно, так что takeStep() имитирует константу
		// для синхронного движения
		takeStep()
		if atomic.LoadInt32(dir) == 1 {
			fmt.Fprint(out, ". Success!")
			return true
		}
		takeStep()
		// Здесь человек понимает, что не может идти в этом направлении
		// и сдается. Мы указываем это, уменьшая направление на единицу
		atomic.AddInt32(dir, -1)
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool { return tryDir("left", &left, out) }
	tryRight := func(out *bytes.Buffer) bool { return tryDir("right", &right, out) }

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer func() { fmt.Println(out.String()) }()
		defer walking.Done()
		fmt.Fprintf(&out, "%v is trying to scoot:", name)
		// Установлено ограниченное кол-во попыток, чтобы программа завершилась.
		// В программах с livelock такого ограничения может и не быть, пот почему
		// это проблема!
		for i := 0; i < 5; i++ {
			// Сначала человек пытается двигаться по левой стороне, если терпит
			// неудачу, тогда пытается двигаться по правой стороне
			if tryLeft(&out) || tryRight(&out) {
				return
			}
		}
		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	}

	// Эта переменная позволяет программе ожидать, пока оба человека не завершат
	// движение или не сдадутся
	var peopleInHallway sync.WaitGroup
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")
	peopleInHallway.Wait()
}
