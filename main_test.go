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
	f, err := os.Create(sample)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	_, err = fmt.Fprintf(b, "Email addressess below, this first line is only for title!\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < recnum; i++ {
		_, err := fmt.Fprintf(b, "%08dHello@World\n", i)
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

func TestBinSearch(t *testing.T) {
	err, bs := NewBinSearch("sample.txt")
	defer bs.close()
	if err != nil {
		panic(err)
	}

	_, pos, rec, _ := bs.ScanNext(0)

	if pos != 59 {
		t.Errorf("bs.nextRec(0) = %d; want 59, rec:%s\n", pos, rec)
	}

	_, nPos, nRec, _ := bs.ScanNext(pos)

	if pos == nPos {
		t.Errorf("bs.nextRec(%d) = %d; want different than %d, rec:%s\n", pos, nPos, pos, nRec)
	}
}

func TestBinSearchScanNext(t *testing.T) {
	err, bs := NewBinSearch("sample.txt")
	defer bs.close()
	if err != nil {
		panic(err)
	}

	err, pos, rec, eof := bs.ScanNext(0)
	fmt.Println(err, pos, rec, eof, len(rec))
	err, pos, rec, eof = bs.ScanNext(58)
	fmt.Println(err, pos, rec, eof, len(rec))
	err, pos, rec, eof = bs.ScanNext(59)
	fmt.Println(err, pos, rec, eof, len(rec))
	err, pos, rec, eof = bs.ScanNext(60)
	fmt.Println(err, pos, rec, eof, len(rec))
	fmt.Println(bs.size)

	err, pos, rec, eof = bs.ScanNext(859)
	fmt.Println(err, pos, rec, eof)
}

func TestBinSearchHave(t *testing.T) {
	err, bs := NewBinSearch("sample.txt")
	defer bs.close()
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
