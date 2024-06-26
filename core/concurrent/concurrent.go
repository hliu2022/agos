package concurrent

type IConcurrent interface {
	InitWithOptions(opts ...Option) error
	AsyncDoByQueue(queueId int64, fn func() bool, cb func(err error))
	AsyncDo(fn func() bool, cb func(err error))
}

type Concurrent struct {
	dispatch
	option *options

	tasks     chan task
	cbChannel chan func(error)
}

func (c *Concurrent) InitWithOptions(opts ...Option) error {
	c.option = newDefaultConcurrentOptions()
	for _, opt := range opts {
		if err := opt(c.option); err != nil {
			return err
		}
	}
	c.tasks = make(chan task, c.option.maxTaskChannelNum)
	c.cbChannel = make(chan func(error), c.option.maxTaskChannelNum)
	c.dispatch.open(c.option, c.tasks, c.cbChannel)
	return nil
}

func (c *Concurrent) AsyncDo(fn func() bool, cb func(err error)) {
	c.AsyncDoByQueue(0, fn, cb)
}

func (c *Concurrent) AsyncDoByQueue(queueId int64, fn func() bool, cb func(err error)) {

}
