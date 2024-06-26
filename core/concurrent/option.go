package concurrent

import "runtime"

const defaultMaxTaskChannelNum = 1000000

type Option func(*options) error

type options struct {
	maxTaskChannelNum     int
	customMinGoroutineNum int64
	customMaxGoroutineNum int64
}

func newDefaultConcurrentOptions() *options {
	defaultGoroutineNum := int32(float32(runtime.NumCPU()) + 1)
	return &options{
		maxTaskChannelNum:     defaultMaxTaskChannelNum,
		customMinGoroutineNum: int64(defaultGoroutineNum),
		customMaxGoroutineNum: int64(defaultGoroutineNum),
	}
}

func WithMaxTaskChannelNum(maxTaskChannelNum int) Option {
	return func(opt *options) error {
		opt.maxTaskChannelNum = maxTaskChannelNum
		return nil
	}
}

// WithGoroutineNum和WithCPUMulGoroutineNum二选一，指定数量范围
func WithGoroutineNum(minGoroutineNum int64, maxGoroutineNum int64) Option {
	return func(opt *options) error {
		opt.customMinGoroutineNum = minGoroutineNum
		opt.customMaxGoroutineNum = maxGoroutineNum
		return nil
	}
}

/*
cpuMul 表示cpu的倍数
建议:(1)cpu密集型 使用1  (2)i/o密集型使用2或者更高
Tips：WithGoroutineNum和WithCPUMulGoroutineNum二选一，指定CUP倍数的协程
*/
func WithCPUMulGoroutineNum(cupMul float32) Option {
	goroutineNum := int32(float32(runtime.NumCPU())*cupMul + 1)
	return func(opt *options) error {
		opt.customMinGoroutineNum = int64(goroutineNum)
		opt.customMaxGoroutineNum = int64(goroutineNum)
		return nil
	}
}
