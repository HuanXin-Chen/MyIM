package domain

const (
	windowSize = 5
)

// 滑动窗口维护最近5个状态

type stateWindow struct {
	stateQueue []*Stat    // 维护状态窗口
	statChan   chan *Stat // 事件传递通道
	sumStat    *Stat      // 状态总和，缓存加快计算
	idx        int64      //下一个位置
}

func newStateWindow() *stateWindow {
	return &stateWindow{
		stateQueue: make([]*Stat, windowSize),
		statChan:   make(chan *Stat),
		sumStat:    &Stat{},
	}
}

func (sw *stateWindow) getStat() *Stat {
	res := sw.sumStat.Clone()
	res.Avg(windowSize)
	return res
}

func (sw *stateWindow) appendStat(s *Stat) {
	// 减去即将被删除的state
	sw.sumStat.Sub(sw.stateQueue[sw.idx%windowSize])
	// 更新最新的stat
	sw.stateQueue[sw.idx%windowSize] = s
	// 计算最新的窗口和
	sw.sumStat.Add(s)
	sw.idx++
}
