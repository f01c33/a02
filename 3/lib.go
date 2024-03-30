package mux

type Mux[T any] struct {
	chans []chan T
}

func NewMux[T any](chans []chan T) *Mux[T] {
	return &Mux[T]{
		chans: chans,
	}
}

func (m *Mux[T]) CloseMux() {
	for _, c := range m.chans {
		close(c)
	}
}

func (m *Mux[T]) Send(t T) {
	for _, c := range m.chans {
		// go func() { c <- t }()
		c <- t
	}
}

func (m *Mux[T]) NewSink() chan T {
	out := make(chan T, 100)

	go func() {
		for {
			closed := 0
			for _, c := range m.chans {
				select {
				case v, ok := <-c:
					if !ok {
						closed++
					} else {
						out <- v
					}
				default:
				}
			}
			if closed == len(m.chans) {
				close(out)
				return
			}
		}
	}()
	return out
}
