package api

import (
	"encoding/json"
	"log"
	"net/http"
	orderflow "orderFlow/libs/orderflow"
	"orderFlow/libs/shared"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// 设置Ping/Pong超时时间
	HandshakeTimeout: 20 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	// 允许所有子协议
	Subprotocols: []string{},
	// 允许所有请求头
	EnableCompression: true,
}

type ConnectionPool struct {
	connections map[*websocket.Conn]bool
	mu          sync.Mutex
}

var pool = &ConnectionPool{
	connections: make(map[*websocket.Conn]bool),
}

func (p *ConnectionPool) Add(conn *websocket.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.connections[conn] = true
}

func (p *ConnectionPool) Remove(conn *websocket.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.connections, conn)
}

// 广播消息给所有连接
func (p *ConnectionPool) Broadcast(message []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for conn := range p.connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("广播消息错误: %v", err)
			conn.Close()
			delete(p.connections, conn)
		}
	}
}

type orderFlowWs struct {
	Symbol   string `json:"symbol"`
	Interval string `json:"symbolInterval"`
	StepSize string `json:"stepSize"`
}

func SocketHandler(c *gin.Context) {
	log.Printf("收到WebSocket连接请求 远程地址: %s", c.Request.RemoteAddr)
	log.Printf("请求头: %v", c.Request.Header)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("升级WebSocket连接失败: %v", err)
		return
	}

	log.Printf("WebSocket连接成功建立 远程地址: %s", conn.RemoteAddr().String())

	// 将连接添加到连接池
	pool.Add(conn)
	defer func() {
		log.Printf("WebSocket连接关闭 远程地址: %s", conn.RemoteAddr().String())
		pool.Remove(conn)
		conn.Close()
	}()

	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(20 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(20 * time.Second))
		return nil
	})

	// 启动一个goroutine来发送ping消息
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					log.Printf("发送ping消息失败: %v", err)
					return
				}
			}
		}
	}()

	// 处理消息循环
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取消息错误: %v", err)
			}
			break
		}

		// 重置读取超时
		conn.SetReadDeadline(time.Now().Add(20 * time.Second))

		var orderFlowWs orderFlowWs
		err = json.Unmarshal(message, &orderFlowWs)
		if err != nil {
			log.Printf("解析消息错误: %v", err)
			continue
		}

		activeCandle := shared.GetActiveCandles(orderFlowWs.Symbol, orderFlowWs.Interval)
		if activeCandle == nil {
			log.Printf("没有找到活跃的Candle")
			continue
		}
		newCandle := *activeCandle
		var candles []orderflow.FootprintCandle
		candles = append(candles, newCandle)
		parseCandles(&candles, orderFlowWs.StepSize)
		jsonData, err := json.Marshal(candles)
		if err != nil {
			log.Printf("序列化数据错误: %v", err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			log.Printf("发送消息错误: %v", err)
			break
		}
	}
}
