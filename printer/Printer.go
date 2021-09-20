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
	storage  []int                 // Хранилище сгенерированных случайных чисел.
	lim      int64                 // Максимально возможное значение числа.
	cmds     <-chan control.Cmd    // Команды, управляющие рабочим циклом принтера. Команды изменяют состояние принтера.
	feedback chan<- control.Signal // Обратная связь принтера (см. dispatcher).
	nums     <-chan int            // Сгенерированное число читается из nums.
	state    control.State         // Состояние принтера. Влияет на рабочий цикл.
	result   chan<- control.Signal
}

// New возвращает ссылку на новый экземпляр принтера случайных чисел.
func New(lim int64, cmds <-chan control.Cmd, feedback chan<- control.Signal, nums <-chan int, result chan<- control.Signal) *printer {
	if lim <= 0 {
		panic(fmt.Errorf("don't use value 0 of limit"))
	}
	return &printer{lim: lim,
		storage:  make([]int, 0),
		cmds:     cmds,
		feedback: feedback,
		nums:     nums,
		result:   result,
	}
}

// Run запускает основной рабочий цикл Printer.
func (p *printer) Run(wg *sync.WaitGroup) {
	for {
		select {
		case cmd := <-p.cmds:
			p.handleCommand(cmd)
			if cmd == control.Destroy {
				wg.Done()
				return
			}
		case n := <-p.nums:
			if p.state == control.Worked {
				if p.hasNum(n) == false && n != 0 {
					p.storage = append(p.storage, n)
					sort.IntSlice(p.storage).Sort()
					fmt.Printf("\rTotal: %v", p.storage)
					if len(p.storage) == int(p.lim) {
						fmt.Println()
						p.result <- control.Success
					}
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
		p.feedback <- control.Success
	case control.Stop:
		p.state = control.Waiting
		p.feedback <- control.Success
	case control.Destroy:
		p.state = control.Waiting
		p.feedback <- control.Success
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
