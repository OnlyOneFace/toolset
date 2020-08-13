// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/6/26

package util

import (
	"runtime"
	"sync"
	"time"
)

// 频繁的创建和关闭连接，对系统会造成很大负担
// 所以我们需要一个池子，里面事先创建好固定数量的连接资源，需要时就取，不需要就放回池中。
// 但是连接资源有一个特点，我们无法保证连接长时间会有效。
// 比如，网络原因，人为原因等都会导致连接失效。
// 所以我们设置一个超时时间，如果连接时间与当前时间相差超过超时时间，那么就关闭连接。

// 只要类型实现了Conn接口中的方法，就认为是一个连接资源类型
type Conn interface {
	Close() error
	Ping() error
}

// 工厂方法，用于创建连接资源
type Factory func() (Conn, error)

// 连接
type ConnWithTime struct {
	conn Conn
	// 连接时间
	time time.Time
}

// 连接池
type ConnPool struct {
	// 互斥锁，保证资源安全
	mu *sync.Mutex
	// 通道，保存所有连接资源
	connChan chan *ConnWithTime
	// 工厂方法，创建连接资源
	factory Factory
	// 判断池是否关闭
	closed bool
	// 连接超时时间
	connTimeOut time.Duration
	// 容量
	cap int
}

// 创建一个连接资源池
func NewConnPool(factory Factory, opts ...ConnPoolOption) *ConnPool {
	tempCap := runtime.NumCPU()
	cp := &ConnPool{
		mu:          new(sync.Mutex),
		connChan:    make(chan *ConnWithTime, tempCap), // 创建多余空间
		factory:     factory,
		closed:      false,
		connTimeOut: 5 * time.Second,
		cap:         tempCap,
	}
	for _, opt := range opts {
		opt.Apply(cp)
	}
	for i := 0; i < cp.cap; i++ {
		// 通过工厂方法创建连接资源
		conn, err := cp.factory()
		if err != nil {
			continue
		}
		// 将连接资源插入通道中
		cp.connChan <- &ConnWithTime{conn: conn, time: time.Now()}
	}
	return cp
}

// 获取连接资源
func (cp *ConnPool) Get() Conn {
	if cp.closed {
		return nil
	}
	for {
		select {
		// 从通道中获取连接资源
		case connValue, ok := <-cp.connChan:
			if !ok {
				return nil
			}
			// 判断连接中的时间，如果超时，或者心跳没有则关闭
			// 继续获取
			if time.Since(connValue.time) > cp.connTimeOut || connValue.conn.Ping() != nil {
				_ = connValue.conn.Close()
				continue
			}
			return connValue.conn
		default:
			// 如果无法从通道中获取资源，则重新创建一个资源返回
			conn, err := cp.factory()
			if err != nil {
				return nil
			}
			return conn
		}
	}
}

// 连接资源放回池中
func (cp *ConnPool) Put(conn Conn) {
	if cp.closed {
		return
	}
	if len(cp.connChan) >= cp.cap {
		// 如果无法加入,即连接池满了，则关闭连接
		_ = conn.Close()
	}
	// 向通道中加入连接资源
	cp.connChan <- &ConnWithTime{conn: conn, time: time.Now()}
}

// 关闭连接池
func (cp *ConnPool) Close() {
	if cp.closed {
		return
	}
	cp.mu.Lock()
	cp.closed = true
	// 关闭通道
	close(cp.connChan)
	// 循环关闭通道中的连接
	for conn := range cp.connChan {
		_ = conn.conn.Close()
	}
	cp.mu.Unlock()
}

// 线程的配置
type ConnPoolOption interface {
	Apply(*ConnPool)
}

type connPoolWithCap struct {
	cap int
}

func WithCap(capNum int) *connPoolWithCap {
	if capNum < 0 {
		capNum = runtime.NumCPU()
	}
	return &connPoolWithCap{cap: capNum}
}

func (c *connPoolWithCap) Apply(pool *ConnPool) {
	pool.cap = c.cap
	pool.connChan = make(chan *ConnWithTime, pool.cap)
}

type connPoolWithTimeOut struct {
	timeOut time.Duration
}

func WithTimeOut(t time.Duration) *connPoolWithTimeOut {
	if t < 0 {
		t = 5 * time.Second
	}
	return &connPoolWithTimeOut{timeOut: t}
}

func (c *connPoolWithTimeOut) Apply(pool *ConnPool) {
	pool.connTimeOut = c.timeOut
}
