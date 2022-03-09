package handlers

import (
	"pipeline/ringbuffer"
	"time"
)

var BufferLen int = 10
var BufferDrainInterval = 10 * time.Second

func FilterNegative(exit <-chan bool, v <-chan int) <-chan int {
	positiveChannel := make(chan int)
	go func() {
		for {
			select {
			case value := <-v:
				if value > 0 {
					select {
					case positiveChannel <- value:
					case <-exit:
						return
					}
				}
			case <-exit:
				return
			}
		}
	}()

	return positiveChannel
}

func Multiple(exit <-chan bool, v <-chan int) <-chan int {
	multipleChan := make(chan int)
	go func() {
		for {
			select {
			case value := <-v:
				if value != 0 && value%3 == 0 {
					select {
					case multipleChan <- value:
					case <-exit:
						return
					}
				}
			case <-exit:
				return
			}
		}
	}()

	return multipleChan
}

func Buffering(exit <-chan bool, v <-chan int) <-chan int {
	bufferChan := make(chan int)
	buffer := ringbuffer.NewRingBuffer(BufferLen)

	go func() {
		for {
			select {
			case value := <-v:
				buffer.Append(value)
			case <-exit:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(BufferDrainInterval):
				bufferData := buffer.Read()
				if bufferData != nil {
					for _, value := range bufferData {
						select {
						case bufferChan <- value:
						case <-exit:
							return
						}
					}
				}
			case <-exit:
				return
			}
		}
	}()

	return bufferChan
}
