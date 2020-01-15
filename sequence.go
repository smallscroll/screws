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
	reset()
}

//StartAndInitializeASequenceService 启动并初始化一个序列服务
func StartAndInitializeASequenceService(uniqueID uint) ISequence {
	s := &sequence{
		UniqueID:  uniqueID,
		BaseValue: time.Now().Unix(),
	}
	go s.reset()
	return s
}

//sequence 序列服务
type sequence struct {
	UniqueID  uint
	BaseValue int64
}

//Get 获取值
func (s *sequence) Get() int64 {
	vs := fmt.Sprintf("%d%d", s.UniqueID+100, atomic.AddInt64(&s.BaseValue, 1))
	v, err := strconv.ParseInt(vs, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

//reset 重置序列服务基准值
func (s *sequence) reset() {
	for {
		today := time.Now()
		tomorrow := today.Add(time.Hour * 24)
		tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
		<-time.After(tomorrow.Sub(today))
		atomic.StoreInt64(&s.BaseValue, time.Now().Unix())
	}
}
