package state

import (
	"context"
	"fmt"
	"github.com/HuanXin-Chen/MyIM/common/config"
	"github.com/HuanXin-Chen/MyIM/common/prpc"
	"github.com/HuanXin-Chen/MyIM/state/rpc/client"
	"github.com/HuanXin-Chen/MyIM/state/rpc/service"
	"google.golang.org/grpc"
)

var cmdChannel chan *service.CmdContext

func RunMain(path string) {
	config.Init(path)
	cmdChannel = make(chan *service.CmdContext, config.GetSateCmdChannelNum())

	s := prpc.NewPServer(
		prpc.WithServiceName(config.GetStateServiceName()),
		prpc.WithIP(config.GetSateServiceAddr()),
		prpc.WithPort(config.GetSateServerPort()), prpc.WithWeight(config.GetSateRPCWeight()))
	fmt.Println(config.GetStateServiceName(), config.GetSateServiceAddr(), config.GetSateServerPort(), config.GetSateRPCWeight())
	s.RegisterService(func(server *grpc.Server) {
		service.RegisterStateServer(server, &service.Service{CmdChannel: cmdChannel})
	})
	fmt.Println("-------------im state stated------------")
	// 初始化RPC 客户端
	client.Init()
	// 启动 命令处理写协程
	go cmdHandler()
	// 启动 rpc server
	s.Start(context.TODO())
}

func cmdHandler() {
	for cmd := range cmdChannel {
		switch cmd.Cmd {
		case service.CancelConnCmd:
			fmt.Printf("cancelconn endpoint:%s, fd:%d, data:%+v", cmd.Endpoint, cmd.FD, cmd.Playload)
		case service.SendMsgCmd:
			fmt.Println("cmdHandler", string(cmd.Playload))
			client.Push(cmd.Ctx, int32(cmd.FD), cmd.Playload)
		}
	}
}
