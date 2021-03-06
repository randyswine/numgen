package numbers

import (
	"fmt"
	"math/rand"
	"numgen/control"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// generator - объект, выполняющий генерацию случайного числа в диапазоне от 0 до lim.
// Сгенерированное число передаётся в канал out.
type generator struct {
	token      string                // Уникальный идентификатор генератора.
	timeout    time.Duration         // Таймаут частоты генерации числа.
	lim        int64                 // Максимально возможное значение числа.
	cmds       <-chan control.Cmd    // Команды, управляющие рабочим циклом генератора. Команды изменяют состояние генератора.
	feedback   chan<- control.Signal // Обратная связь генератора (см. dispatcher).
	out        chan<- int            // Сгенерированное число передаётся в out.
	randomizer *rand.Rand            // Генератор случайных чисел.
	state      control.State         // Состояние генератора. Влияет на рабочий цикл.
	ticker     *time.Ticker          // Тикер генерации числа.
}

// New возвращает ссылку на новый экземпляр генератора случайных чисел.
func New(timeout time.Duration, lim int64, cmds <-chan control.Cmd, feedback chan<- control.Signal, out chan<- int) *generator {
	if timeout == 0 {
		panic(fmt.Errorf("don't use value 0 of timeout"))
	}
	return &generator{
		token:      uuid.NewV4().String(),
		timeout:    timeout,
		lim:        lim,
		cmds:       cmds,
		feedback:   feedback,
		out:        out,
		state:      control.Waiting,
		randomizer: rand.New(rand.NewSource(time.Now().UnixNano())),
		ticker:     time.NewTicker(timeout * time.Second),
	}
}

//	Run запускает рабочий цикл генератора случайных чисел. Каждую итерацию цикла ожидается
//	получение команды или сигнала тикера для генерации числа. Тикер может прислать сигнал
//	только в том случае, если состояние генератора - control.Worked. Состояние изменяется командами, поступающими
//	по каналу cmds. Рабочий цикл прерывается после получения control.Destroy.
func (g *generator) Run(wg *sync.WaitGroup) {
	for {
		select {
		case cmd := <-g.cmds:
			// Обработка команды, если она получена.
			g.handleCommand(cmd)
			if cmd == control.Destroy {
				wg.Done()
				return
			}
		case <-g.ticker.C:
			// Генерация числа в состоянии worked.
			if g.state == control.Worked {
				g.makeNum()
			}
		}
	}
}

// handleCommand выполняет обработку поступающих команд и управляет тикером генерации чисел.
func (g *generator) handleCommand(cmd control.Cmd) {
	switch cmd {
	case control.Start:
		g.state = control.Worked
		g.ticker.Reset(g.timeout)
		g.feedback <- control.Success
	case control.Stop:
		g.ticker.Stop()
		g.state = control.Waiting
		g.feedback <- control.Success
	case control.Destroy:
		g.ticker.Stop()
		g.state = control.Waiting
		g.feedback <- control.Success
	}
}

func (g *generator) makeNum() {
	n := int(g.randomizer.Int63n(g.lim + 1))
	//fmt.Printf("\rGenerator %s maked number: %d", g.token, n)
	g.out <- n
}
