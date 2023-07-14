package main

import (
	"fmt"
	"sync"
)

func main() {
	// Добавить переменную, которая позволит нашему коду 
    // синхронизировать доступ к переменным в памяти
	var memoryAccess sync.Mutex
	var value int
	go func() {
		// Здесь мы объявляем что пока мы не объявили противоположное, 
		// т.е. Unlock(), наша горутина должна иметь эксклюзивный 
		// доступ к этой памяти
		memoryAccess.Lock()
		value++
		// Здесь мы объявляем, что горутина завершила работу с 
		// этой областью памяти
		memoryAccess.Unlock()
	}()
	
	// Здесь мы единожды снова объявляем что следующие условные операторы 
	// должны иметь эксклюзивный доступ к переменной в памяти
	memoryAccess.Lock()
	if value == 0 {
		fmt.Printf("the value is %v.\n", value)
	} else {
		fmt.Printf("the value is %v.\n", value)
	}
	// Здесь мы объявляем, что снова завершили работу с этой памятью
	memoryAccess.Unlock()
}
