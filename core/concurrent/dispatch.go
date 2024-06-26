package concurrent

import "sync"

type dispatch struct {
	option *options

	queueIdChannel chan int64
	workerQueue    chan task
	tasks          chan task
	idle           bool
	workerNum      int32
	cbChannel      chan func(error)

	mapTaskQueueSession map[int64]*queue.Deque[task]

	waitWorker   sync.WaitGroup
	waitDispatch sync.WaitGroup
}

func (d *dispatch) open(opts *options, tasks chan task, cbChannel chan func(error)) {
	d.option = opts
	d.tasks = tasks
	d.cbChannel = cbChannel
	go d.run()
}

func (d *dispatch) run() {

}
