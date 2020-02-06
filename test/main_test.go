package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

var sample = "./sample.txt"
var size = 100

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

// Create the test database
func setup() {
	f, err := os.Create(sample)
	if err != nil {
		fmt.Println(err)
		return
	}

	b := bufio.NewWriter(f)
	_, err = fmt.Fprintf(b, "Email addressess below, this first line is only for title!\n")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	for i := 0; i < size; i++ {
		_, err := fmt.Fprintf(b, "%010dHello@World\n", i)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Clean up
func shutdown() {
	os.Remove(sample)
}

func Test