package screws

import (
	"sync"

	"github.com/gorilla/websocket"
)

//IWSManager 连接管理器接口
type IWSManager interface {
	Start()
	Broadcast(message []byte)
	Notice(message []byte, client IWSClient)
	GetClients() []string
	addClient(client *wsClient)
	addr() *wsManager
}

//IWSClient 客户端接口
type IWSClient interface {
	Reading()
	Writing()
	addr() *wsClient
}

// NewWSManager 初始化连接管理器
func NewWSManager() IWSManager {
	return &wsManager{
		BroadcastChan: make(chan []byte),
		AddChan:       make(chan *wsClient, 1000),
		DeleteChan:    make(chan *wsClient),
		Clients:       &sync.Map{},
	}
}

//NewWSClient 初始化客户端(连接ID，连接，管理器)
func NewWSClient(id string, conn *websocket.Conn, manager IWSManager) IWSClient {
	wsClient := &wsClient{
		ID:        id,
		Socket:    conn,
		SendChan:  make(chan []byte),
		Manager:   manager.addr(),
		AddedChan: make(chan bool),
	}
	manager.addClient(wsClient)
	<-wsClient.AddedChan
	return wsClient
}

//wsManager  管理器
type wsManager struct {
	Clients       *sync.Map
	BroadcastChan chan []byte
	AddChan       chan *wsClient
	DeleteChan    chan *wsClient
	Lock          *sync.RWMutex
}

//wsClient 客户端
type wsClient struct {
	ID        string
	Socket    *websocket.Conn
	SendChan  chan []byte
	Manager   *wsManager
	AddedChan chan bool
}

// Start 启动连接管理器
func (m *wsManager) Start() {
	for {
		select {
		case conn := <-m.AddChan:
			m.Clients.Store(conn, true)
			conn.AddedChan <- true
		case conn := <-m.DeleteChan:
			if _, ok := m.Clients.Load(conn); ok {
				close(conn.SendChan)
				m.Clients.Delete(conn)
			}
		case message := <-m.BroadcastChan:
			m.Clients.Range(func(k, v interface{}) bool {
				k.(*wsClient).SendChan <- message
				return true
			})
		}
	}
}

//Add 添加客户端到连接管理器
func (m *wsManager) addClient(client *wsClient) {
	m.AddChan <- client
}

// Broadcast 消息广播
func (m *wsManager) Broadcast(message []byte) {
	m.Clients.Range(func(k, v interface{}) bool {
		k.(*wsClient).SendChan <- message
		return true
	})

}

// Notice 消息通知
func (m *wsManager) Notice(message []byte, client IWSClient) {
	if _, ok := m.Clients.Load(client.addr()); ok {
		client.addr().SendChan <- message
	}
}

// GetClients 查询全部连接
func (m *wsManager) GetClients() []string {
	clients := []string{}
	m.Clients.Range(func(k, v interface{}) bool {
		clients = append(clients, k.(*wsClient).ID)
		return true
	})
	return clients
}

func (m *wsManager) addr() *wsManager {
	return m
}

func (c *wsClient) addr() *wsClient {
	return c
}

// Reading 启动读服务
func (c *wsClient) Reading() {
	defer func() {
		c.Manager.DeleteChan <- c
		c.Socket.Close()
	}()
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			c.Manager.DeleteChan <- c
			c.Socket.Close()
			break
		}
		c.Manager.BroadcastChan <- message
	}
}

// Writing 启动写服务
func (c *wsClient) Writing() {
	defer func() {
		c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.SendChan:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
