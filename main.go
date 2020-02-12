package main

import (
	"bufio"
	"fmt"
	"os"
)

type BinSearch struct {
	fdb  *os.File // the "database" file pointer
	size int64
}

type Record struct {
	pos int64
	len int64
	con string
	eof bool
}

func NewBinSearch(path string) (err error, bs *BinSearch) {
	fdb, err := os.Open(path)
	if err != nil {
		return err, nil
	}

	fi, err := fdb.Stat()
	if err != nil {
		panic(err)
	}

	return nil, &BinSearch{fdb: fdb, size: fi.Size()}
}

func (bs *BinSearch) Close() {
	if bs.fdb != nil {
		bs.fdb.Close()
		bs.fdb = nil
	}
}

/*****************************************************************/
/*  Scan a record from the given position, it could be a halfone */
/*****************************************************************/
func (bs *BinSearch) ScanRecord(pos int64) (rec Record) {
	if pos >= bs.size {
		return Record{pos: bs.size, len: 0, con: "", eof: true}
	}

	bs.fdb.Seek(pos, 0)

	scanner := bufio.NewScanner(bs.fdb)
	if scanner.Scan() {
		return Record{pos: pos, len: int64(len(scanner.Text())), con: scanner.Text(), eof: false}
	}

	return Record{pos: bs.size, len: 0, con: "", eof: true}
}

/*****************************************************************/
/*  Scan the following record                                    */
/*****************************************************************/
func (bs *BinSearch) NextRecord(rec *Record) *Record {
	next := bs.ScanRecord(rec.pos + rec.len + 1)
	return &next
}

/*****************************************************************/
/*  Scan a complete record between begin and end                 */
/*****************************************************************/
func (bs *BinSearch) MidRecord(begin *Record, end *Record) (mid *Record) {
	midPos := (begin.pos + begin.len + end.pos) / 2
	midRec := bs.ScanRecord(midPos) // It could be a half record
	nextRec := bs.NextRecord(&midRec)

	switch {
	case begin.pos == end.pos: // the search intervall is zero
		return begin

	case nextRec.pos >= end.pos: // nextRec is out of the bound
		nextRec = bs.NextRecord(begin)
		return nextRec

	default:
		return nextRec
	}
}

func (bs *BinSearch) Have(query string) (found bool) {
	begin := bs.ScanRecord(0)     // the first line
	end := bs.ScanRecord(bs.size) // last pos

	if query < begin.con {
		//fmt.Println("search is below the first record")
		return false
	}

	if query == begin.con {
		//fmt.Println("search is on the first record")
		return true
	}

	for begin.pos != end.pos {
		mid := bs.MidRecord(&begin, &end)
		//fmt.Println(begin, mid, end)
		switch {
		case mid.con == query:
			return true

		case mid.con < query:
			begin = *mid
			continue

		case mid.con > query:
			end = *mid
			continue
		}
	}
	return false
}

func main() {
	err, bs := NewBinSearch("sample.txt")
	defer bs.Close()
	if err != nil {
		panic(err)
	}

	searchAddress := "-0000004hsasahjHello@World+"
	fmt.Printf("searching %s\n", searchAddress)
	fmt.Printf("found in the: %t\n", bs.Have(searchAddress))
	os.Exit(0)
}
