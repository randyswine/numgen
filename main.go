package main

import (
	"numgen/control"
	"numgen/numbers"
	"numgen/printer"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	timer := time.NewTimer(10 * time.Second)
	pCmd := make(chan control.Cmd, 0)
	gCmd := make(chan control.Cmd, 0)
	nums := make(chan int, 10)
	p := printer.New(10, pCmd, nums)
	wg.Add(1)
	go p.Run(&wg)
	g := numbers.New(1, 10, gCmd, nums)
	wg.Add(1)
	go g.Run(&wg)
	pCmd <- control.Start
	gCmd <- control.Start
	<-timer.C
	gCmd <- control.Stop
	gCmd <- control.Destroy
	pCmd <- control.Stop
	pCmd <- control.Destroy
	wg.Wait()
	os.Exit(0)
}
