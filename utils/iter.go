package utils

import (
	"log"
	"fmt"
	"encoding/hex"
)

type iterator struct {
	mark, markEnd, bRemaining int
	srcIn []byte
	dBug bool
	dBugD []byte
}

func NewIterator(b []byte) iterator {
	var i iterator
	i.srcIn = b
	i.mark = 0
	i.markEnd = len(b)
	i.bRemaining = len(b)
	i.dBug = false
	return i
}

func (i *iterator) DeBug() {
	i.dBug = true
	return
}

func (i *iterator) deBug() {
	b := i.srcIn[i.mark:]
	xx := hex.Dump(b)
	log.Println("DeBug interator current state: Mark: ", i.mark, "End Mark: ", i.markEnd, "Bytes Remaining: ", i.bRemaining)
	fmt.Println(xx)
	i.dBug = false
	return
}

func (i *iterator) GetNextOne() (byte, bool) {
	if i.dBug {
		i.deBug()
	}
	if i.mark+1 <= i.markEnd {
		f := i.srcIn[i.mark]
		i.mark += 1
		i.bRemaining -= 1
		var ff []byte
		ff = append(ff, f)
		fmt.Println(hex.Dump(ff))
		return f, true
	} else if i.mark+1 > i.markEnd {
		log.Println("exceeds markEnd")
		return 0x00, false
		}
	return 0x00, false
}

func (i *iterator) GetNextN(it int) ([]byte, bool) {
	if i.dBug {
		i.deBug()
	}
	if i.mark+it < i.markEnd {
		markIn := i.mark
		i.mark += it
		i.bRemaining -= it	
		
		fmt.Println(hex.Dump(i.srcIn[markIn:i.mark]))
			
		return i.srcIn[markIn:i.mark], true
	} else {
		if i.mark+it > i.markEnd {
			log.Println("i exceeds markEnd")
			return nil, false
		} else if i.mark+it == i.markEnd {
			markIn := i.mark
			i.mark += it
			i.bRemaining -= it	
			
			fmt.Println(hex.Dump(i.srcIn[markIn:i.mark]))
					
			return i.srcIn[markIn:i.mark], true
		}
	}
	return nil, false
}