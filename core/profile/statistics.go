package profile

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/idealeak/goserver.v3/core/basic"
	"github.com/idealeak/goserver.v3/core/logger"
	"github.com/idealeak/goserver.v3/core/utils"
)

const (
	TIME_ELEMENT_ACTION = iota
	TIME_ELEMENT_TASK
	TIME_ELEMENT_TIMER
	TIME_ELEMENT_NET
	TIME_ELEMENT_COMMAND
	TIME_ELEMENT_MODULE
	TIME_ELEMENT_JOB
)

var TimeStatisticMgr = &timeStatisticMgr{
	elements: make(map[string]*TimeElement),
}

type TimeElement struct {
	Name        string
	ElementType int
	Times       int64
	TotalTick   int64
	MaxTick     int64
	MinTick     int64
}

type timeStatisticMgr struct {
	elements map[string]*TimeElement
	l        sync.RWMutex
}

func (this *timeStatisticMgr) WatchStart(name string, elementype int) basic.IStatsWatch {
	tw := newTimeWatcher(name, elementype)
	return tw
}

func (this *timeStatisticMgr) addStatistic(name string, elementype int, d int64) {
	this.l.RLock()
	if te, exist := this.elements[name]; exist {
		this.l.RUnlock()
		times := atomic.AddInt64(&te.Times, 1)
		total := atomic.AddInt64(&te.TotalTick, d)
		if d > te.MaxTick {
			atomic.StoreInt64(&te.MaxTick, d)
			if Config.SlowMS > 0 && d >= int64(Config.SlowMS)*int64(time.Millisecond) {
				logger.Logger.Warnf("###slow timespan name: %s  take:%s avg used:%s", strings.ToLower(te.Name), utils.ToS(time.Duration(d)), utils.ToS(time.Duration(total/times)))
			}
		}
		if d < te.MinTick {
			atomic.StoreInt64(&te.MinTick, d)
		}

	} else {
		this.l.RUnlock()
		te := &TimeElement{
			Name:        name,
			ElementType: elementype,
			Times:       1,
			TotalTick:   d,
			MaxTick:     d,
			MinTick:     d,
		}
		this.l.Lock()
		this.elements[name] = te
		this.l.Unlock()
	}
}

func (this *timeStatisticMgr) GetStats() map[string]TimeElement {
	elements := make(map[string]TimeElement)
	this.l.RLock()
	for k, v := range this.elements {
		te := *v
		te.TotalTick /= int64(time.Millisecond)
		te.MinTick /= int64(time.Millisecond)
		te.MaxTick /= int64(time.Millisecond)
		elements[k] = te
	}
	this.l.RUnlock()
	return elements
}

func (this *timeStatisticMgr) dump(w io.Writer) {
	elements := make(map[string]*TimeElement)
	this.l.RLock()
	for k, v := range this.elements {
		elements[k] = v
	}
	this.l.RUnlock()
	fmt.Fprintf(w, "| % -30s| % -10s | % -16s | % -16s | % -16s | % -16s |\n", "name", "times", "used", "max used", "min used", "avg used")
	for k, v := range elements {
		fmt.Fprintf(w, "| % -30s| % -10d | % -16s | % -16s | % -16s | % -16s |\n", strings.ToLower(k), v.Times, utils.ToS(time.Duration(v.TotalTick)), utils.ToS(time.Duration(v.MaxTick)), utils.ToS(time.Duration(v.MinTick)), utils.ToS(time.Duration(int64(v.TotalTick)/v.Times)))
	}
}

func GetStats() map[string]TimeElement {
	return TimeStatisticMgr.GetStats()
}

func init() {
	basic.StatsWatchMgr = TimeStatisticMgr
}
