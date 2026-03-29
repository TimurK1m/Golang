package main

import (
	"fmt"
	"sync"
)

func main() {
	safeMap := make(map[string]int)
	var mu sync.RWMutex

	for i := 0; i < 100; i++ {
		
		go func(key int) {

			mu.Lock()
			safeMap["key"] = key
			mu.Unlock()
		}(i)
	}

	

	mu.RLock()
	value := safeMap["key"]
	mu.RUnlock()

	fmt.Printf("Value: %d\n", value)
}