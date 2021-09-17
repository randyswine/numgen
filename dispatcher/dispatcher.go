package dispatcher

import (
	"fmt"
	"numgen/control"
)

// component предназначен для удобного хранения служебных канолов генераторов и принтера.
type component struct {
	control  chan<- control.Cmd
	feedback <-chan control.Signal
}

var instance *dispatcher

// dispatcher управляет рабочими циклами генераторов случайных чисел и принтера.
type dispatcher struct {
	generators []component // Коллекция служебных каналов генераторов случайных чисел.
	printer    component   // Служебные каналы принтера.
}

// New возвращает экземпляр диспетчера компонетов.
func Dispatcher() *dispatcher {
	if instance == nil {
		instance = &dispatcher{}
	}
	return instance
}

// AppendGenerator добавляет служебные каналы генератора в коллекцию.
func (d *dispatcher) AppendGenerator(control chan<- control.Cmd, feedback <-chan control.Signal) {
	d.generators = append(d.generators, component{control: control, feedback: feedback})
}

// AppendPrinter добавляет служебные каналы принтера.
func (d *dispatcher) AppendPrinter(control chan<- control.Cmd, feedback <-chan control.Signal) {
	d.printer = component{control: control, feedback: feedback}
}

// StartAll запускает полезное действие генераторов случайных чисел и принтера, и дожидается обратной связи
// об успешном запуске.
func (d *dispatcher) StartAll() {
	var result control.Signal
	for _, generator := range d.generators {
		generator.control <- control.Start
		result = <-generator.feedback
		if result != control.Success {
			panic(fmt.Errorf("ошибка старта генератора случайных чисел"))
		}
	}
}

// StopAll останавливает выполение полезного действия всех генераторов, принтера, и дожидается обратной связи
// об успешной остановке каждого компонента.
func (d *dispatcher) StopAll() {
	var result control.Signal
	for _, generator := range d.generators {
		generator.control <- control.Stop
		result = <-generator.feedback
		if result != control.Success {
			panic(fmt.Errorf("ошибка остановки генератора случайных чисел"))
		}
	}
	d.printer.control <- control.Stop
	result = <-d.printer.feedback
	if result != control.Success {
		panic(fmt.Errorf("ошибка остановки принтера случайных чисел"))
	}
}

// DestroyAll останавливает рабочие циклы всех генераторов, принтера, и дожидается обратной связи
// об успешной остановке каждого компонента.
func (d *dispatcher) DestroyAll() {
	var result control.Signal
	for _, generator := range d.generators {
		generator.control <- control.Destroy
		result = <-generator.feedback
		if result != control.Success {
			panic(fmt.Errorf("ошибка завершения рабочего цикла генератора случайных чисел"))
		}
	}
	d.printer.control <- control.Stop
	result = <-d.printer.feedback
	if result != control.Success {
		panic(fmt.Errorf("ошибка завершения рабочего цикла принтера случайных чисел"))
	}
}
