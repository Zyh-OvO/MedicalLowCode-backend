package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Age       []int  `json:"age"`
	Name      string `json:"name"`
	Niubility bool   `json:"niubility"`
	An        any    `json:"an"`
}

func test(map[string]any) {
	return
}

func main() {
	b := []byte(`{"age":18,"name":"5lmh.com","marry":false,"an":"test"}`)
	p := Person{
		Name: "12312312313",
		Age:  []int{1, 2, 3},
		An:   1,
	}
	json.Unmarshal(b, &p)
	fmt.Printf("%+v", p)
}
