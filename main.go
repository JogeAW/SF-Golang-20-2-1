package main

import (
	"bufio"
	"fmt"
	"os"
	"pipeline/handlers"
	"pipeline/pipeline"
	"strconv"
	"strings"
	"time"
)

func dataSource() (<-chan int, <-chan bool) {
	channel := make(chan int)
	exit := make(chan bool)

	go func() {
		defer close(exit)

		scanner := bufio.NewScanner(os.Stdin)
		var value string
		for {
			scanner.Scan()
			value = scanner.Text()
			if strings.EqualFold(value, "exit") {
				fmt.Println("Finishing program")
				return
			}

			i, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Only integers!")
				continue
			}

			channel <- i
		}
	}()

	return channel, exit
}

const BufferLen = 15
const BufferDrainInterval = 10 * time.Second

func main() {

	handlers.BufferLen = BufferLen
	handlers.BufferDrainInterval = BufferDrainInterval

	consumer := func(exit <-chan bool, channel <-chan int) {
		for {
			select {
			case value := <-channel:
				fmt.Printf("Processed value: %d\n", value)
			case <-exit:
				return
			}
		}

	}
	source, exit := dataSource()
	Intpipeline := pipeline.NewPipeline(exit, handlers.FilterNegative, handlers.Multiple, handlers.Buffering)
	consumer(exit, Intpipeline.Run(source))
}
