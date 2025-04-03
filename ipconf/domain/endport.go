package domain

import (
	"sync/atomic"
	"unsafe"
)

// 候选人状态维护

type Endport struct {
	IP          string       `json:"ip"`
	Port        string       `json:"port"`
	ActiveSorce float64      `json:"-"`
	StaticSorce float64      `json:"-"`
	Stats       *Stat        `json:"-"`
	window      *stateWindow `json:"-"`
}

func NewEndport(ip, port string) *Endport {
	ed := &Endport{
		IP:   ip,
		Port: port,
	}
	ed.window = newStateWindow()
	ed.Stats = ed.window.getStat()
	go func() { //每个状态用一个协程去维护信息，KB级别很轻量级
		for stat := range ed.window.statChan { // 鉴定statChan, 然后更新Stats
			ed.window.appendStat(stat)
			newStat := ed.window.getStat()
			atomic.SwapPointer((*unsafe.Pointer)((unsafe.Pointer)(ed.Stats)), unsafe.Pointer(newStat))
		}
	}()
	return ed
}

// 更新状态，传入通道
func (ed *Endport) UpdateStat(s *Stat) {
	ed.window.statChan <- s
}

// 返回数据源
func (ed *Endport) CalculateScore(ctx *IpConfConext) {
	// 如果 stats 字段是空的，则直接使用上一次计算的结果，此次不更新
	if ed.Stats != nil {
		ed.ActiveSorce = ed.Stats.CalculateActiveSorce()
		ed.StaticSorce = ed.Stats.CalculateStaticSorce()
	}
}
