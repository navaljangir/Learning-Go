package main

import "fmt"

func main(){
	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Code   int	`json:"code"`
	}
	resp := &Response{
		Status: "success",
		Code: 200,
	}
	fmt.Printf("Response: %+v\n", resp)
}