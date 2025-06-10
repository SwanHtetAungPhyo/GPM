package spinner

import (
	"fmt"
	"time"
)

type Spinner struct {
	done   chan struct{}
	active bool
}

func NewSpinner() *Spinner {
	return &Spinner{
		done: make(chan struct{}),
	}
}

func (s *Spinner) Start() {
	if s.active {
		return
	}
	s.active = true

	go func() {
		frames := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
		i := 0
		for {
			select {
			case <-s.done:
				fmt.Print("\r")
				return
			default:
				fmt.Printf("\r%c ", frames[i%len(frames)])
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()
}

func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	s.active = false
	s.done <- struct{}{}
	fmt.Print("\r")
}
