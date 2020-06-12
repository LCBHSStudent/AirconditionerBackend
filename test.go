package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	var a = []int{1, 2, 3, 4}
	b, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
	}
	var result = string(b)
	fmt.Println(result)

	var c []int
	json.Unmarshal([]byte(result), &c)
	fmt.Printf("%T", c)
}
