package streamer

import (
	"context"
	"sync"
)

type Streamer interface {
	AddKnot(knotFunc KnotFunc)
	AddInChan(in chan interface{})
	AddOutChan(out chan interface{})
	AddReleaseFunc(fc func())
	Trigger()
	Aborted() <-chan struct{}
	Error() error
}

type streamer struct {
	ctx      context.Context
	cancelF  context.CancelFunc
	wg       *sync.WaitGroup
	releaseF func()
	in       chan interface{}
	out      chan interface{}
	finished chan struct{}
	err      error
	chain    []*knot
}

func NewStreamer(parent context.Context) Streamer {
	ctx, cancelF := context.WithCancel(parent)

	return &streamer{
		ctx:      ctx,
		cancelF:  cancelF,
		wg:       &sync.WaitGroup{},
		chain:    make([]*knot, 0),
		finished: make(chan struct{}),
	}
}

func (s *streamer) AddKnot(knotFunc KnotFunc) {
	k := &knot{
		fc:       knotFunc,
		streamer: s,
		index:    len(s.chain),
	}

	// connect previous know to this one
	if len(s.chain) != 0 {
		prevOut := make(chan interface{})
		s.chain[len(s.chain)-1].out = prevOut
		k.in = prevOut
	}

	s.chain = append(s.chain, k)
}

func (s *streamer) AddInChan(in chan interface{}) {
	if len(s.chain) == 0 {
		panic("AddInChan called before AddKnot")
	}
	s.in = in
	s.chain[0].in = in
}

func (s *streamer) AddOutChan(out chan interface{}) {
	if len(s.chain) == 0 {
		panic("AddOutChan called before AddKnot")
	}
	s.out = out
	s.chain[len(s.chain)-1].out = out
}

func (s *streamer) AddReleaseFunc(fc func()) {
	s.releaseF = fc
}

func (s *streamer) Trigger() {
	if s.in == nil {
		panic("Trigger called before AddInChan")
	}

	for _, k := range s.chain {
		_k := k
		go func() {
			_k.Start()
		}()
		s.wg.Add(1)
	}

	go func() {
		s.wg.Wait()
		if s.releaseF != nil {
			s.releaseF()
		}
		close(s.finished)
	}()
}

func (s *streamer) Aborted() <-chan struct{} {
	return s.finished
}

func (s *streamer) Error() error {
	return s.err
}
