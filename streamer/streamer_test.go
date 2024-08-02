package streamer

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNormalStreamer(t *testing.T) {
	s := NewStreamer(context.Background())
	dataSet := []string{"prehello", "preworld", "and", "iyou"}

	inCh := make(chan interface{})
	outCh := make(chan interface{})

	s.AddKnot(func(i interface{}) (interface{}, error) {
		rawData := i.(string)
		return strings.ToUpper(rawData), nil
	})

	s.AddKnot(func(i interface{}) (interface{}, error) {
		rawData := i.(string)
		return rawData + "tail", nil
	})

	s.AddKnot(func(i interface{}) (interface{}, error) {
		rawData := i.(string)
		return strings.TrimPrefix(rawData, "PRE"), nil
	})

	s.AddInChan(inCh)
	s.AddOutChan(outCh)
	s.Trigger()

	go func() {
		for _, v := range dataSet {
			select {
			case <-s.Aborted():
				t.Logf("Stremer aborted in advance, error=[%+v]", s.Error())
				return

			case inCh <- v:
			}
		}

		t.Log("Input finished. In channel closed. Waiting streamer to stop.")
		close(inCh)

		select {
		case <-s.Aborted():
			t.Logf("Streamer aborted , error=[%+v]", s.Error())
			return
		}

	}()

	for {
		select {
		case v, ok := <-outCh:
			if !ok {
				t.Logf("Output channel closed. Error=[%+v]", s.Error())
				time.Sleep(3 * time.Second)
				return
			}
			t.Log(v)
		}
	}
}

func TestAbnormalStreamer(t *testing.T) {
	s := NewStreamer(context.Background())
	dataSet := []string{"prehelloadyou", "preworld", "andyou", "iyou"}

	inCh := make(chan interface{})
	outCh := make(chan interface{})

	s.AddKnot(func(i interface{}) (interface{}, error) {
		rawData := i.(string)
		return strings.ToUpper(rawData), nil
	})

	s.AddKnot(func(i interface{}) (interface{}, error) {
		rawData := i.(string)
		if strings.Contains(rawData, "ANDYOU") {
			panic("test")
		}
		return rawData + "tail", nil
	})

	s.AddKnot(func(i interface{}) (interface{}, error) {
		rawData := i.(string)
		return strings.TrimPrefix(rawData, "PRE"), nil
	})

	s.AddInChan(inCh)
	s.AddOutChan(outCh)
	s.Trigger()

	go func() {
		for _, v := range dataSet {
			select {
			case <-s.Aborted():
				t.Logf("Streamer aborted in advance, error=[%+v].", s.Error())
				return

			case inCh <- v:
			}
		}

		t.Log("Input finished. In Channel closed. Waiting streamer to stop.")
		close(inCh)

		select {
		case <-s.Aborted():
			t.Logf("Streamer aborted , error=[%+v]", s.Error())
			return
		}
	}()

	for {
		select {
		case v, ok := <-outCh:
			if !ok {
				t.Logf("Output channel closed, Error=[%+v]", s.Error())
				time.Sleep(3 * time.Second)
				return
			}
			t.Log(v)
		}
	}

}
