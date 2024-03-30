package mux

import (
	"testing"
)

func makeMux[T any](nChans int) *Mux[T] {
	chans := []chan T{}
	for i := 0; i < nChans; i++ {
		chans = append(chans, make(chan T, 100))
	}
	return NewMux(chans)
}

func fillMux[T any](m *Mux[T], nOfVals int, val T) {
	for i := 0; i < nOfVals; i++ {
		m.Send(val)
	}
	m.CloseMux()
}
func TestSend(t *testing.T) {
	m := makeMux[int](2)
	m.Send(1)
	for _, c := range m.chans {
		select {
		case val := <-c:
			if val != 1 {
				t.Fatal("not equal")
			}
		default:
			t.Fatal("not sent to channel")
		}
	}
}

func TestSink(t *testing.T) {
	m := makeMux[int](5)
	c := m.NewSink()
	fillMux(m, 10, 1)

	cont := 0
	for v := range c {
		cont++
		if v != 1 {
			t.Fatal("wrong value")
		}
	}
	if cont != 10*5 {
		t.Fatal("Not enough messages")
	}
}

func FuzzSink(f *testing.F) {
	f.Add(uint(5), uint(10))
	f.Fuzz(func(t *testing.T, chans, fillN uint) {
		m := makeMux[int](int(chans))
		c := m.NewSink()
		go fillMux(m, int(fillN), 1)

		cont := 0
		for v := range c {
			cont++
			if v != 1 {
				t.Fatal("wrong value")
			}
		}
		if cont != int(chans*fillN) {
			t.Fatal("Not enough messages", cont)
		}
	})
}

func FuzzSend(f *testing.F) {
	f.Add(uint(2))
	f.Fuzz(func(t *testing.T, a uint) {
		m := makeMux[int](int(a))
		m.Send(1)
		for _, c := range m.chans {
			select {
			case val := <-c:
				if val != 1 {
					t.Fatal("not equal")
				}
			default:
				t.Fatal("not sent to channel")
			}
		}
	})
}
