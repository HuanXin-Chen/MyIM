package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/HuanXin-Chen/MyIM/common/prpc"
	"github.com/HuanXin-Chen/MyIM/gateway/rpc/client"
	"github.com/HuanXin-Chen/MyIM/gateway/rpc/service"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"

	"github.com/HuanXin-Chen/MyIM/common/config"
	"github.com/HuanXin-Chen/MyIM/common/tcp"
)

var cmdChannel chan *service.CmdContext

// RunMain 启动网关服务
func RunMain(path string) {
	config.Init(path)
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{Port: config.GetGatewayTCPServerPort()})
	if err != nil {
		log.Fatalf("StartTCPEPollServer err:%s", err.Error())
		panic(err)
	}
	initWorkPoll()
	initEpoll(ln, runProc)
	cmdChannel = make(chan *service.CmdContext, config.GetGatewayCmdChannelNum())
	s := prpc.NewPServer(
		prpc.WithServiceName(config.GetGatewayServiceName()),
		prpc.WithIP(config.GetGatewayServiceAddr()),
		prpc.WithPort(config.GetGatewayRPCServerPort()), prpc.WithWeight(config.GetGatewayRPCWeight()))
	fmt.Println(config.GetGatewayServiceName(), config.GetGatewayServiceAddr(), config.GetGatewayRPCServerPort(), config.GetGatewayRPCWeight())
	s.RegisterService(func(server *grpc.Server) {
		service.RegisterGatewayServer(server, &service.Service{CmdChannel: cmdChannel})
	})
	fmt.Println("-------------im gateway stated------------")
	// 启动rpc 客户端
	client.Init()
	// 启动 命令处理写协程
	go cmdHandler()
	// 启动 rpc server
	s.Start(context.TODO())
}

// 作为客户端的调用方式

func runProc(c *connection, ep *epoller) {
	ctx := context.Background()
	// step1: 读取一个完整的消息包
	dataBuf, err := tcp.ReadData(c.conn)
	if err != nil {
		// 如果读取conn时发现连接关闭，则直接端口连接
		// 通知 state 清理掉意外退出的 conn的状态信息
		if errors.Is(err, io.EOF) {
			ep.remove(c)
			client.CancelConn(&ctx, getEndpoint(), int32(c.fd), nil)
		}
		return
	}
	err = wPool.Submit(func() {
		// step2:交给 state server rpc 处理
		client.SendMsg(&ctx, getEndpoint(), int32(c.fd), dataBuf)
	})
	if err != nil {
		fmt.Errorf("runProc:err:%+v\n", err.Error())
	}
}

// 作为服务端的处理方式

func cmdHandler() {
	for cmd := range cmdChannel {
		// 异步提交到协池中完成发送任务
		switch cmd.Cmd {
		case service.DelConnCmd:
			wPool.Submit(func() { closeConn(cmd) })
		case service.PushCmd:
			wPool.Submit(func() { sendMsgByCmd(cmd) })
		default:
			panic("command undefined")
		}
	}
}

func closeConn(cmd *service.CmdContext) {
	if connPtr, ok := ep.tables.Load(cmd.FD); ok {
		conn, _ := connPtr.(*connection)
		conn.Close()
		ep.tables.Delete(cmd.FD)
	}
}

func sendMsgByCmd(cmd *service.CmdContext) {
	if connPtr, ok := ep.tables.Load(cmd.FD); ok {
		conn, _ := connPtr.(*connection)
		dp := tcp.DataPgk{
			Len:  uint32(len(cmd.Payload)),
			Data: cmd.Payload,
		}
		tcp.SendData(conn.conn, dp.Marshal())
	}
}

func getEndpoint() string {
	return fmt.Sprintf("%s:%d", config.GetGatewayServiceAddr(), config.GetGatewayRPCServerPort())
}
