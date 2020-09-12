package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		fmt.Println("Hello World!")
		time.Sleep(time.Second * 1)
	}
}
