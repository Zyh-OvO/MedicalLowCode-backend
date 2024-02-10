package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hs, _ := bcrypt.GenerateFromPassword([]byte("aewqwd2323542sdv463sdgs45"), bcrypt.DefaultCost)
	fmt.Println(string(hs))
}
