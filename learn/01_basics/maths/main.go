package main

import (
	"fmt"
)

func main(){
	r := struct{ Error error }{Error: fmt.Errorf("Error")}
	 fmt.Println("Error:", r.Error)
}