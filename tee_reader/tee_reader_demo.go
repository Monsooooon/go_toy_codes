package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	// echo the stdio to the file
	f, err := os.Create("test.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	scanner := bufio.NewScanner(io.TeeReader(os.Stdin, w)) // write content is not buffered!
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println("[Scan]", text)
	}
	w.Flush()
}
