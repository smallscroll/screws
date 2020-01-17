package screws

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

//ISequence 序列服务接口
type ISequence interface {
	Get() int64
}

//StartAndInitializeASequenceService 启动并初始化一个序列服务
func StartAndInitializeASequenceService() ISequence {
	s := &sequence{
		BaseValue: time.Now().Unix(),
		LastValue: 0,
	}
	go func() {
		c := time.Tick(time.Second * 60)
		for range c {
			atomic.StoreInt64(&s.BaseValue, time.Now().Unix())
			atomic.StoreInt64(&s.LastValue, 0)
		}
	}()
	return s
}

//sequence 序列服务
type sequence struct {
	BaseValue int64
	LastValue int64
}

//Get 获取值
func (s *sequence) Get() int64 {
	vs := fmt.Sprintf("%d%09d", s.BaseValue, atomic.AddInt64(&s.LastValue, 1))
	v, err := strconv.ParseInt(vs, 10, 64)
	if err != nil {
		return 0
	}
	return v
}
