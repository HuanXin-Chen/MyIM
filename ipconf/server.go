package ipconf

import (
	"github.com/HuanXin-Chen/MyIM/common/config"
	"github.com/HuanXin-Chen/MyIM/ipconf/domain"
	"github.com/HuanXin-Chen/MyIM/ipconf/source"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RunMain 启动web容器
func RunMain(path string) {
	config.Init(path)
	source.Init() //数据源要优先启动
	domain.Init() // 初始化调度层
	s := server.Default(server.WithHostPorts(":6789"))
	s.GET("/ip/list", GetIpInfoList)
	s.Spin()
}
