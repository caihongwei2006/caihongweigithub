package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println("failed opening file", err)
		return
	} else {
		fmt.Println("success", file)
	}
	defer file.Close()
	fout, err := os.Openfile()

}
