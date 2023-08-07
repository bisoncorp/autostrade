package utils

type UnboundedChan[T any] struct {
	in, out chan *T
	close   chan struct{}
	queue   []*T
}

func NewUnboundedChan[T any]() *UnboundedChan[T] {
	u := &UnboundedChan[T]{
		in:    make(chan *T, 16),
		out:   make(chan *T, 16),
		close: make(chan struct{}),
	}

	go u.processing()
	return u
}

func (u *UnboundedChan[T]) In() chan<- *T {
	return u.in
}

func (u *UnboundedChan[T]) Out() <-chan *T {
	return u.out
}

func (u *UnboundedChan[T]) Close() {
	u.close <- struct{}{}
}

func (u *UnboundedChan[T]) processing() {
	u.queue = make([]*T, 0, 1<<10)
	for {
		select {
		case e, ok := <-u.in:
			if !ok {
				panic("utils: misuse of unbounded channel, In() was closed")
			}
			u.queue = append(u.queue, e)
		case <-u.close:
			u.closed()
			return
		}
		for len(u.queue) > 0 {
			select {
			case u.out <- u.queue[0]:
				u.queue[0] = nil
				u.queue = u.queue[1:]
			case e, ok := <-u.in:
				if !ok {
					panic("utils: misuse of unbounded channel, In() was closed")
				}
				u.queue = append(u.queue, e)
			case <-u.close:
				u.closed()
				return
			}
		}
		if cap(u.queue) < 1<<5 {
			u.queue = make([]*T, 0, 1<<10)
		}
	}
}

func (u *UnboundedChan[T]) closed() {
	close(u.in)
	for e := range u.in {
		u.queue = append(u.queue, e)
	}
	for len(u.queue) > 0 {
		select {
		case u.out <- u.queue[0]:
			u.queue[0] = nil // de-reference earlier to help GC
			u.queue = u.queue[1:]
		default:
		}
	}
	close(u.out)
	close(u.close)
}
