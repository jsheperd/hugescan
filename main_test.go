package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

var sample = "./sample.txt"
var recnum = 40

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

// Create the test database
func setup() {
	err, bs := NewBinSearch("sample.txt")
	defer bs.Close()
	if err != nil {
		panic(err)
	}

	f, err := os.Create(sample)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	//_, err = fmt.Fprintf(b, "Email addressess below, this first line is only for title!\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < recnum; i++ {
		_, err := fmt.Fprintf(b, "-%08dHello@World+\n", i)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	b.Flush()
}

// Clean up
func shutdown() {
	//os.Remove(sample)
}

func TestSecond(t *testing.T) {
	err, bs := NewBinSearch("sample.txt")
	defer bs.Close()
	if err != nil {
		panic(err)
	}

	rec := bs.ScanRecord(0)
	fmt.Println(rec)
	rec = bs.ScanRecord(8)
	fmt.Println(rec)

	n := bs.NextRecord(&rec)
	fmt.Println(n)

	if n.pos != 22 {
		t.Errorf("bs.nextRec(0) = %d; want 22, rec:%s\n", n.pos, n.con)
	}
}

func TestLast(t *testing.T) {
	err, bs := NewBinSearch("sample.txt")
	defer bs.Close()
	if err != nil {
		panic(err)
	}

	rec := bs.ScanRecord(bs.size)
	fmt.Println(rec)
	if rec.pos != bs.size {
		t.Errorf("bs.ScanRecord(bs.size) != bs.size\n")
	}

	rec = bs.ScanRecord(bs.size)
	last := bs.NextRecord(&rec)
	fmt.Println(last)
	if last.pos != bs.size {
		t.Errorf("bs.ScanRecord(bs.size) != bs.size\n")
	}

	rec = bs.ScanRecord(bs.size - 25)
	last = bs.NextRecord(&rec)
	last = bs.NextRecord(last)
	fmt.Println(last)
	if last.pos != bs.size {
		t.Errorf("bs.ScanRecord(bs.size) != bs.size\n")
	}
}

func TestBinSearchScanNext(t *testing.T) {
	err, bs := NewBinSearch("sample.txt")
	defer bs.Close()
	if err != nil {
		panic(err)
	}

	rec := bs.ScanRecord(0)
	fmt.Println(rec)
	rec = bs.ScanRecord(58)
	fmt.Println(rec)
	rec = bs.ScanRecord(59)
	fmt.Println(rec)
	rec = bs.ScanRecord(60)
	fmt.Println(rec)
	fmt.Println(bs.size)

	fmt.Println(rec)
}

func TestBinSearchHave(t *testing.T) {
	err, bs := NewBinSearch("sample.txt")
	defer bs.Close()
	if err != nil {
		panic(err)
	}

	sample, err := os.Open("check.txt")
	if err != nil {
		panic(err)
	}
	defer sample.Close()
	scanner := bufio.NewScanner(sample)
	for scanner.Scan() {
		valid := scanner.Text()
		fmt.Println(valid, bs.Have(valid)) // Println will add back the final '\n'

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
