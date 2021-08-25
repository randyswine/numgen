package numbers

import (
	"fmt"
	"math/rand"
	"numgen/control"
	"sync"
	"time"
)

// generator - объект, выполняющий генерацию случайного числа в диапазоне от 0 до lim.
// Сгенерированное число передаётся в канал out.
type generator struct {
	timeout    time.Duration      // Таймаут частоты генерации числа.
	lim        int64              // Максимально возможное значение числа.
	cmds       <-chan control.Cmd // Команды, управляющие рабочим циклом генератора.
	out        chan<- int         // Сгенерированное число передаётся в out.
	randomizer *rand.Rand
	state      control.State
	wg         *sync.WaitGroup
	ticker     *time.Ticker
}

func New(timeout time.Duration, lim int64, cmds <-chan control.Cmd, out chan<- int) *generator {
	if timeout == 0 {
		panic(fmt.Errorf("don't use value 0 of timeout"))
	}
	return &generator{
		timeout:    timeout,
		lim:        lim,
		cmds:       cmds,
		out:        out,
		state:      control.Waiting,
		randomizer: rand.New(rand.NewSource(lim)),
		ticker:     time.NewTicker(timeout * time.Second),
	}
}

func (g *generator) Run(wg *sync.WaitGroup) {
	if wg != nil {
		g.wg = wg
	}
	for {
		select {
		case cmd := <-g.cmds:
			g.handleCommand(cmd)
			if cmd == control.Destroy {
				return
			}
		case <-g.ticker.C:
			if g.state == control.Worked {
				g.makeNum()
			}
		}
	}
}

func (g *generator) handleCommand(cmd control.Cmd) {
	switch cmd {
	case control.Start:
		g.state = control.Worked
		g.ticker.Reset(g.timeout * time.Second)
	case control.Stop:
		g.ticker.Stop()
		g.state = control.Waiting
	case control.Destroy:
		g.ticker.Stop()
		g.state = control.Waiting
		if g.wg != nil {
			g.wg.Done()
		}
	}
}

func (g *generator) makeNum() {
	n := int(g.randomizer.Int63n(g.lim))
	g.out <- n
}
