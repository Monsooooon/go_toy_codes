package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go_toy_codes/protobuf_demo/todo"
	"google.golang.org/protobuf/proto"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "miss subcommand")
		os.Exit(1)
	}

	var err error
	switch cmd := flag.Arg(0); cmd {
	case "add":
		err = add(strings.Join(flag.Args()[1:], " "))
	case "list":
		err = list()
	default:
		err = fmt.Errorf("unknown subcommand: %s", cmd)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const dbPath = "mydb.pb"

// add a string into the pb db
func add(text string) error {
	task := &todo.Task{
		Text: text,
		Done: false,
	}

	// convert struct to bytes iwth protocal buffer
	b, err := proto.Marshal(task)
	if err != nil {
		return fmt.Errorf("could not encode task: %v", err)
	}

	// open pb db
	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}

	// write byte length first
	proto.
		err = gob.NewEncoder(f).Encode(int32(len(b)))
	if err != nil {
		return fmt.Errorf("could not encode length of message: %v", err)
	}

	// then write bytes
	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could not write task to file: %v", err)
	}

	// remember to close file
	err = f.Close()
	if err != nil {
		return fmt.Errorf("could not close file: %v", err)
	}

	return nil
}

func list() error {
	// read ALL contents from file mydb.pb using ioutil.ReadFile
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("could not read file %s: %v", dbPath, err)
	}
	for {

		if len(b) == 0 {
			break
		} else if len(b) < 4 {
			return fmt.Errorf("remaining odd %d bytes", len(b))
		}

		var length int32
		err = gob.NewDecoder(bytes.NewReader(b[:4])).Decode(&length)
		if err != nil {
			return fmt.Errorf("could not decode message length: %v", err)
		}

		// move b forward
		b = b[4:]

		var task todo.Task
		if err := proto.Unmarshal(b[:length], &task); err != nil {
			return fmt.Errorf("could not unmarshal task: %v", err)
		}

		// move b  forward again
		b = b[length:]

		// print task
		if task.Done {
			fmt.Print("[√] ")
		} else {
			fmt.Print("[x] ")
		}
		fmt.Println(task.Text)
	}

	return nil
}
