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

	scanner := bufio.NewScanner(bs.fdb)
	if scanner.Scan() == false {
		return Record{pos: bs.size, len: 0, con: "", eof: true}
	}

	return Record{pos: pos, len: int64(len(scanner.Text())), con: scanner.Text(), eof: false}
}

/*****************************************************************/
/*  Scan the following record                                    */
/*****************************************************************/
func (bs *BinSearch) NextRecord(rec *Record) *Record {
	next := bs.ScanRecord(rec.pos + rec.len)
	return &next
}

/*****************************************************************/
/*  Scan a complete record between begin and end                 */
/*****************************************************************/
func (bs *BinSearch) MidRecord(begin *Record, end *Record) (mid *Record) {
	midPos := (begin.pos + begin.len + end.pos) / 2
	midRec := bs.ScanRecord(midPos)
	nextRec := bs.NextRecord(&midRec)

	switch {
	case nextRec.pos != end.pos:
		return nextRec

	case begin.pos == end.pos:
		return begin

	default:
		return bs.NextRecord(begin)
	}
}

func (bs *BinSearch) Have(query string) (found bool) {
	begin := bs.ScanRecord(0)     // the first line
	end := bs.ScanRecord(bs.size) // last pos

	if query < begin.con {
		return false
	}

	if query == begin.con {
		return true
	}

	for begin.pos != end.pos {
		mid := bs.MidRecord(&begin, &end)
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

	searchAddress := "p000000020Hello@World"
	fmt.Printf("searching %s\n", searchAddress)
	fmt.Printf("found in the: %t\n", bs.Have(searchAddress))
	os.Exit(0)
}
