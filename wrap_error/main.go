package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	err := readConfig()
	/*
		if err == os.ErrNotExist {
			fmt.Printf("err == os.ErrNotExist")
		}
	*/
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("errors.Is(err, os.ErrNotExist) == true")
	}
	fmt.Printf("err msg: %s\n", err.Error())
}

func readConfig() error {
	f, err := os.Open("config.json")
	if err != nil {
		// add additional msg to the err's msg
		// fmt.Printf("%T", err) -> *err.wrapError
		// %w
		return fmt.Errorf("Cannot open config.json: %w", err)
	}
	f.Close()
	return nil
}
