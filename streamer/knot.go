package streamer

import (
	"fmt"
)

type KnotFunc func(interface{}) (interface{}, error)
type Knot interface {
	Start()
}

type knot struct {
	fc       KnotFunc
	in       chan interface{}
	out      chan interface{}
	streamer *streamer
	index    int
}

func (k *knot) Start() {
	var err error

	defer func() {
		if r := recover(); r != nil {
			k.streamer.err = fmt.Errorf("%+v", r)
		} else if err != nil {
			k.streamer.err = err
		}

		if k.streamer.err != nil {
			k.streamer.cancelF()
		}

		if k.out != nil {
			close(k.out)
		}
		k.streamer.wg.Done()
	}()

	for {
		select {
		case <-k.streamer.ctx.Done():
			return

		case rawData, ok := <-k.in:
			if !ok {
				return
			}

			postData, e := k.fc(rawData)
			if e != nil {
				err = e
				return
			}

			if k.out != nil {
				// since the downstream knot is possible to be closed, here we should check again
				select {
				case <-k.streamer.ctx.Done():
					return
				case k.out <- postData:
				}
			}
		}
	}

}
