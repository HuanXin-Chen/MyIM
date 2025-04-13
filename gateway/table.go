package gateway

import "sync"

var tables table

// 维护映射表,did->连接,便于快速push

type table struct {
	did2conn sync.Map
}

func InitTables() {
	tables = table{
		did2conn: sync.Map{},
	}
}
