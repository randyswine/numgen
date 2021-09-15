package printer

import (
	"fmt"
	"numgen/control"
	"sort"
	"sync"
)

// printer агрегирует все уникальные поступающие от генераторов случайные числа,
// и выводит их в консоль.
type printer struct {
	storage []int
	lim     int64
	cmds    <-chan control.Cmd
	nums    <-chan int
	state   control.State
}

// New возвращает ссылку на новый экземпляр принтера случайных чисел.
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

// Run запускает основной рабочий цикл Printer.
func (p *printer) Run(wg *sync.WaitGroup) {
	for {
		select {
		case cmd := <-p.cmds:
			p.handleCommand(cmd)
			if cmd == control.Destroy {
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

// handleCommand выполняет обработку поступающих команд и управляет тикером генерации чисел.
func (p *printer) handleCommand(cmd control.Cmd) {
	switch cmd {
	case control.Start:
		p.state = control.Worked
	case control.Stop:
		p.state = control.Waiting

	case control.Destroy:
		p.state = control.Waiting
	}
}

// hasNum возвращает true, если входящее от генератора число уже было получено ранее.
func (p *printer) hasNum(n int) bool {
	for _, val := range p.storage {
		if n == val {
			return true
		}
	}
	return false
}
