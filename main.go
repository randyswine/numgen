package main

import (
	"flag"
	"fmt"
	"numgen/control"
	"numgen/dispatcher"
	"numgen/numbers"
	"numgen/printer"
	"os"
	"sync"
	"time"
)

var (
	flowCount int   // Количество потоков генерации случайных чисел.
	limit     int64 // Максимальное случайное число. Влияет на количество сгенерированых чисел (от 0 до limit).
	timeout   int64 // Таймаут генерации случайного числа.
)

func main() {
	var wg sync.WaitGroup
	defer func() {
		if r := recover(); r != nil {
			os.Exit(1)
		}
	}()
	flag.IntVar(&flowCount, "flowcount", 0, "Numbers of threads generating random numbers")
	flag.Int64Var(&limit, "limit", 0, "Maximum random number. Affects the number of generated numbers (from 0 to limit).")
	flag.Int64Var(&timeout, "timeout", 0, "Timeout for generating a random number.")
	flag.Parse()
	if flowCount == 0 || limit == 0 || timeout == 0 {
		flag.Usage()
		os.Exit(1)
	}
	initComponents(&wg)
	dispatcher.Dispatcher().StartAll()
	wg.Wait()
	os.Exit(0)
}

func initComponents(wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Component initalize error: %v\r\n", r)
		}
	}()
	// Инициализация диспетчера компонентов.
	dispatcher := dispatcher.Dispatcher()
	// Инициализация принтера.
	printerCmd := make(chan control.Cmd, 1)
	printerFeedback := make(chan control.Signal, 1)
	nums := make(chan int, flowCount)
	printer := printer.New(limit, printerCmd, printerFeedback, nums)
	// Передача служебных каналов принтера диспетчеру.
	dispatcher.AppendPrinter(printerCmd, printerFeedback)
	// Старт рабочего цикла принтера.
	wg.Add(1)
	go printer.Run(wg)
	// Вычисление таймаута.
	t := time.Duration(timeout * time.Hour.Milliseconds())
	// Инициализация генераторов.
	for i := 0; i < flowCount; i++ {
		generatorCmd := make(chan control.Cmd, 1)
		generatorFeedback := make(chan control.Signal, 1)
		generator := numbers.New(t, limit, generatorCmd, generatorFeedback, nums)
		dispatcher.AppendGenerator(generatorCmd, generatorFeedback)
		wg.Add(1)
		go generator.Run(wg)
	}
}
