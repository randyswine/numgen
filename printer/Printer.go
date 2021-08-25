package printer

import (
	"fmt"
	"numgen/control"
	"sort"
	"sync"
)

type printer struct {
	storage []int
	lim     int64
	cmds    <-chan control.Cmd
	nums    <-chan int
	state   control.State
}

func New(lim int64, cmds <-chan control.Cmd, nums <-chan int) *printer {
	if lim <= 0 {
		panic(fmt.Errorf("don't use value 0 of limit"))
	}
	return &printer{lim: lim,
		storage: make([]int, 0),
		cmds:    cmds,
		nums:    nums,
	}
}

func (p *printer) Run(wg *sync.WaitGroup) {
	for {
		select {
		case cmd := <-p.cmds:
			switch cmd {
			case control.Start:
				p.state = control.Worked
			case control.Stop:
				p.state = control.Waiting
			case control.Destroy:
				wg.Done()
				fmt.Println()
				return
			}
		case n := <-p.nums:
			if p.state == control.Worked {
				if p.hasNum(n) == false {
					p.storage = append(p.storage, n)
					sort.IntSlice(p.storage).Sort()
					fmt.Printf("\rTotal: %v", p.storage)
				}
			}
		}
	}
}

func (p *printer) hasNum(n int) bool {
	for _, val := range p.storage {
		if n == val {
			return true
		}
	}
	return false
}
