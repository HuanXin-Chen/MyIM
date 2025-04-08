package plugin

import (
	"errors"
	"fmt"

	"github.com/HuanXin-Chen/MyIM/common/prpc/config"
	"github.com/HuanXin-Chen/MyIM/common/prpc/discov"
	"github.com/HuanXin-Chen/MyIM/common/prpc/discov/etcd"
)

// GetDiscovInstance 获取服务发现实例
func GetDiscovInstance() (discov.Discovery, error) {
	name := config.GetDiscovName()
	switch name {
	case "etcd":
		return etcd.NewETCDRegister(etcd.WithEndpoints(config.GetDiscovEndpoints()))
	}

	return nil, errors.New(fmt.Sprintf("not exist plugin:%s", name))
}
