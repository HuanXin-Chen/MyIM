package gateway

import (
	"fmt"
	"github.com/HuanXin-Chen/MyIM/common/config"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
)

// 全局对象
var ep *ePool    // epoll池
var tcpNum int32 // 当前服务允许接入的最大tcp连接数

type ePool struct {
	eChan  chan *connection // 连接通道
	tables sync.Map         // 关系映射
	eSize  int              // 轮询器数量大小
	done   chan struct{}    // 资源回收，用于控制协程

	ln *net.TCPListener                 // 简单器
	f  func(c *connection, ep *epoller) // 回调处理函数
}

func initEpoll(ln *net.TCPListener, f func(c *connection, ep *epoller)) {
	setLimit()               // 打开最大进程数限制，默认是1024
	ep = newEPool(ln, f)     // 新建epoll池
	ep.createAcceptProcess() // 创建监听协程
	ep.startEPool()          // 启动
}

func newEPool(ln *net.TCPListener, cb func(c *connection, ep *epoller)) *ePool {
	return &ePool{
		eChan:  make(chan *connection, config.GetGatewayEpollerChanNum()),
		done:   make(chan struct{}),
		eSize:  config.GetGatewayEpollerNum(),
		tables: sync.Map{},
		ln:     ln,
		f:      cb,
	}
}

// 创建一个专门处理 accept 事件的协程，与当前cpu的核数对应，能够发挥最大功效
func (e *ePool) createAcceptProcess() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				conn, e := e.ln.AcceptTCP()
				// 限流熔断
				if !checkTcp() {
					_ = conn.Close()
					continue
				}
				setTcpConifg(conn)
				if e != nil {
					if ne, ok := e.(net.Error); ok && ne.Temporary() {
						fmt.Errorf("accept temp err: %v", ne)
						continue
					}
					fmt.Errorf("accept err: %v", e)
				}
				c := connection{
					conn: conn,
					fd:   socketFD(conn),
				}
				ep.addTask(&c)
			}
		}()
	}
}

// 启动轮询器处理池
func (e *ePool) startEPool() {
	for i := 0; i < e.eSize; i++ {
		go e.startEProc()
	}
}

// 轮询器池 处理器
func (e *ePool) startEProc() {
	ep, err := newEpoller()
	if err != nil {
		panic(err)
	}
	// 监听连接创建事件
	go func() {
		for {
			select {
			case <-e.done: // 退出时优雅关闭
				return
			case conn := <-e.eChan:
				addTcpNum()
				fmt.Printf("tcpNum:%d\n", tcpNum)
				if err := ep.add(conn); err != nil {
					fmt.Printf("failed to add connection %v\n", err)
					conn.Close() //登录未成功直接关闭连接
					continue
				}
				fmt.Printf("EpollerPool new connection[%v] tcpSize:%d\n", conn.RemoteAddr(), tcpNum)
			}
		}
	}()
	// 轮询器在这里轮询等待, 当有wait发生时则调用回调函数去处理
	for {
		select {
		case <-e.done: // 退出时优雅关闭
			return
		default:
			connections, err := ep.wait(200) // 200ms 一次轮询避免 忙轮询

			if err != nil && err != syscall.EINTR {
				fmt.Printf("failed to epoll wait %v\n", err)
				continue
			}
			for _, conn := range connections {
				if conn == nil {
					break
				}
				e.f(conn, ep)
			}
		}
	}
}

func (e *ePool) addTask(c *connection) {
	e.eChan <- c
}

type epoller struct {
	fd int
}

func newEpoller() (*epoller, error) {
	fd, err := unix.EpollCreate1(0) // 一个红黑树根
	if err != nil {
		return nil, err
	}
	return &epoller{
		fd: fd,
	}, nil
}

func (e *epoller) add(conn *connection) error {
	fd := conn.fd
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.EPOLLIN | unix.EPOLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}
	ep.tables.Store(fd, conn)
	return nil
}

func (e *epoller) remove(c *connection) error {
	subTcpNum()
	fd := c.fd
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	ep.tables.Delete(fd)
	return nil
}

func (e *epoller) wait(msec int) ([]*connection, error) {
	events := make([]unix.EpollEvent, config.GetGatewayEpollWaitQueueSize())
	n, err := unix.EpollWait(e.fd, events, msec)
	if err != nil {
		return nil, err
	}
	var connections []*connection
	for i := 0; i < n; i++ {
		if conn, ok := ep.tables.Load(int(events[i].Fd)); ok {
			connections = append(connections, conn.(*connection))
		}
	}
	return connections, err
}

func socketFD(conn *net.TCPConn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(*conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

// 设置go 进程打开文件数的限制
func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	log.Printf("set cur limit: %d", rLimit.Cur)
}

func getTcpNum() int32 {
	return atomic.LoadInt32(&tcpNum)
}

func addTcpNum() {
	atomic.AddInt32(&tcpNum, 1)
}

func subTcpNum() {
	atomic.AddInt32(&tcpNum, -1)
}

func checkTcp() bool {
	num := getTcpNum()
	maxTcpNum := config.GetGatewayMaxTcpNum()
	return num <= maxTcpNum
}

func setTcpConifg(c *net.TCPConn) {
	_ = c.SetKeepAlive(true)
}
