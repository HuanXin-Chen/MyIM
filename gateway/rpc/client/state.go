package client

import (
	"context"
	"github.com/HuanXin-Chen/MyIM/common/config"
	"github.com/HuanXin-Chen/MyIM/common/prpc"
	"github.com/HuanXin-Chen/MyIM/state/rpc/service"
	"time"
)

var stateClient service.StateClient

func initStateClient() {
	pCli, err := prpc.NewPClient(config.GetStateServiceName())
	if err != nil {
		panic(err)
	}
	stateClient = service.NewStateClient(pCli.Conn())
}

// go 的 context 是工程化自定义的核心，后面有空研究一下

func CancelConn(ctx *context.Context, endpoint string, fd int32, playLoad []byte) error {
	rpcCtx, _ := context.WithTimeout(*ctx, 100*time.Millisecond)
	stateClient.CancelConn(rpcCtx, &service.StateRequest{
		Endpoint: endpoint,
		Fd:       fd,
		Data:     playLoad,
	})
	return nil
}

func SendMsg(ctx *context.Context, endpoint string, fd int32, playLoad []byte) error {
	rpcCtx, _ := context.WithTimeout(*ctx, 100*time.Millisecond)
	_, err := stateClient.SendMsg(rpcCtx, &service.StateRequest{
		Endpoint: endpoint,
		Fd:       fd,
		Data:     playLoad,
	})
	if err != nil {
		panic(err)
	}
	return nil
}
