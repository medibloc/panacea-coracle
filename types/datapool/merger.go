package datapool

import (
	"fmt"
	"sync"
)

type Merger struct {
	channels []<-chan []byte
}

func NewMerger() *Merger {
	return &Merger{}
}

func (m *Merger) Add(newChan <-chan []byte) {
	m.channels = append(m.channels, newChan)
}

func (m *Merger) Merge(errPipeline chan error) <-chan []byte {
	var wg sync.WaitGroup
	out := make(chan []byte, len(m.channels))

	output := func(c <-chan []byte) {
		defer func() {
			wg.Done()
		}()

		select {
		case err := <-errPipeline:
			fmt.Println("error3")
			errPipeline <- err
			return
		default:
			for n := range c {
				out <- n
			}
		}
	}

	for _, c := range m.channels {
		wg.Add(1)
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
