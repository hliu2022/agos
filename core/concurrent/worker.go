package concurrent

import (
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	"runtime"
	"sync"
)

type task struct {
	queueId int64
	fn      func() bool
	cb      func(err error)
}

type worker struct {
	*dispatch
}

func (w *worker) start(waitGroup *sync.WaitGroup, t *task, d *dispatch) {
	w.dispatch = d
	d.workerNum += 1
	waitGroup.Add(1)
	go w.run(waitGroup, *t)
}

func (w *worker) run(waitGroup *sync.WaitGroup, t task) {
	defer waitGroup.Done()

	w.exec(&t)
	for {
		select {
		case tw := <-w.workerQueue:
			if tw.isExistTask() {
				//exit goroutine
				//log.Info("worker goroutine exit")
				return
			}
			w.exec(&tw)
		}
	}
}

func (w *worker) exec(t *task) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			errString := fmt.Sprint(r)

			cb := t.cb
			t.cb = func(err error) {
				cb(errors.New(errString))
			}
			log.Dump(string(buf[:l]), log.String("error", errString))
			w.endCallFun(true, t)
		}
	}()

	w.endCallFun(t.fn(), t)
}
