package main

import (
	"bufio"
	"fmt"
	"os"
)

type binSearch struct {
	fdb  *os.File // the "database" file pointer
	size int64
}

type record struct {
	pos int64
	len int64
	con string
}

func NewBinSearch(path string) (err error, bs *binSearch) {
	fdb, err := os.Open(path)
	if err != nil {
		return err, nil
	}

	fi, err := fdb.Stat()
	if err != nil {
		panic(err)
	}

	return nil, &binSearch{fdb: fdb, size: fi.Size()}
}

func (bs *binSearch) close() {
	if bs.fdb != nil {
		bs.fdb.Close()
		bs.fdb = nil
	}
}

func (bs *binSearch) ScanNext(pos int64) (err error, rpos int64, address string, eof bool) {
	bs.fdb.Seek(pos, 0)

	scanner := bufio.NewScanner(bs.fdb)
	if scanner.Scan() == false {
		return nil, pos, scanner.Text(), true
	}
	rpos = pos + int64(len(scanner.Text()))

	if scanner.Scan() == false {
		return nil, rpos, scanner.Text(), true
	}
	return nil, rpos, scanner.Text(), false
}

func (bs *binSearch) Have(query string) (found bool) {
	beg := int64(0)
	end := bs.size
	pos := (beg + end) / 2
	eof := false

	var err error
	var rec string

	for {
		err, pos, rec, eof = bs.ScanNext(pos)
		if err != nil {
			panic(err)
		}
		switch {
		case rec == query:
			return true

		case eof:
			return false

		case rec < query:
			beg = pos
			pos = (beg + end) / 2
			continue

		case rec > query:
			end = pos
			pos = (beg + end) / 2
			continue

		case beg == end:
			return false
		}
	}
}

func main() {
	err, bs := NewBinSearch("sample.txt")
	defer bs.close()
	if err != nil {
		panic(err)
	}

	searchAddress := "p000000020Hello@World"
	fmt.Printf("searching %s\n", searchAddress)
	fmt.Printf("found in the: %t\n", bs.Have(searchAddress))
	os.Exit(0)
}
