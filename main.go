package main

import (
	"flag"
	"fmt"
	"numgen/control"
	"numgen/dispatcher"
	"numgen/numbers"
	"numgen/printer"
	"os"
	"time"
)

var (
	flowCount int   // Количество потоков генерации случайных чисел.
	limit     int64 // Максимальное случайное число. Влияет на количество сгенерированых чисел (от 0 до limit).
	timeout   int64 // Таймаут генерации случайного числа.
)

func main() {
	flag.IntVar(&flowCount, "flowcount", 1, "Numbers of threads generating random numbers")
	flag.Int64Var(&limit, "limit", 10, "Maximum random number. Affects the number of generated numbers (from 0 to limit).")
	flag.Int64Var(&timeout, "flowcount", 1, "Timeout (second) for generating a random number.")
	defer func() {
		if r := recover(); r != nil {
			os.Exit(1)
		}
	}()
	initComponents()
	dispatcher.Dispatcher().StartAll()
	os.Exit(0)
}

func initComponents() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Component initalize error: %v\r\n", r)
		}
	}()
	// Инициализация диспетчера компонентов.
	dispatcher := dispatcher.Dispatcher()
	// Инициализация принтера.
	printerCmd := make(chan control.Cmd, 0)
	printerFeedback := make(chan control.Signal, 0)
	nums := make(chan int, 0)
	printer := printer.New(limit, printerCmd, printerFeedback, nums)
	// Передача служебных каналов принтера диспетчеру.
	dispatcher.AppendPrinter(printerCmd, printerFeedback)
	// Старт рабочего цикла принтера.
	printer.Run()
	// Вычисление таймаута.
	t := time.Duration(timeout * int64(time.Second))
	// Инициализация генераторов.
	for i := 0; i < flowCount; i++ {
		generatorCmd := make(chan control.Cmd, 0)
		generatorFeedback := make(chan control.Signal, 0)
		generator := numbers.New(t, limit, generatorCmd, generatorFeedback, nums)
		dispatcher.AppendGenerator(generatorCmd, generatorFeedback)
		generator.Run()
	}
}
