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

func (bs *binSearch) nextRec(pos int64) (err error, rpos int64, address string) {
	bs.fdb.Seek(pos, 0)
	scanner := bufio.NewScanner(bs.fdb)
	if scanner.Scan() == false {
		panic("Can't read it")
	}
	rpos = pos + int64(len(scanner.Text()))
	fmt.Printf("read rec %s\n", scanner.Text())

	if scanner.Scan() == false {
		panic("Can't read it")
	}
	fmt.Printf("Real pos: %d, line: %s", rpos, scanner.Text())
	return nil, rpos, scanner.Text()
}

func (bs *binSearch) have(query string) (found bool) {
	beg := int64(0)
	end := bs.size
	pos := (beg + end) / 2

	for {
		err, pos, rec := bs.nextRec(pos)
		if err != nil {
			panic(err)
		}
		switch {
		case rec == query:
			return true

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

	searchAddress := "0009999821Hello@World"
	fmt.Printf("searching %s", searchAddress)
	fmt.Printf("found in the %s", bs.have(searchAddress))
	os.Exit(0)
}
