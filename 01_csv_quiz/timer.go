package main

import (
	"fmt"
	"time"
)

func timeMe() {
	timer := time.NewTimer(time.Duration(1) * time.Second)

	// the select needs to be inside the for loop
	// otherwise timer will be checked only once at the begging of execution
	// and obviously fail
	for i := 1; i <= 1000000; i++ {
		select {
			case <- timer.C:
				fmt.Println("Time's up!")
				return
			default:
				fmt.Println("test line")
		}
	}
}
