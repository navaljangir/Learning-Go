package main

import (
	"fmt"
	"time"
)

func main(){
	var time time.Time = time.Now()
	fmt.Printf("%v\n", time)
}