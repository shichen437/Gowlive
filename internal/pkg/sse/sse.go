package sse

import (
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

var (
	clients            = make(map[*ghttp.Request]clientInfo)
	channelConnections = make(map[string]int)
	clientsMu          sync.Mutex
	autoSendTickers    = make(map[string]*time.Ticker)
	autoSendMu         sync.Mutex
	autoMonitorStarted bool
)

type clientInfo struct {
	channel string
	ch      chan string
}

type ChannelHandler interface {
	OnConnect(channel string)
	OnDisconnect(channel string)
}

var channelHandlers = make(map[string]ChannelHandler)

func RegisterChannelHandler(channel string, handler ChannelHandler) {
	channelHandlers[channel] = handler
}

func HandleSSE(r *ghttp.Request) {
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache, no-transform")
	r.Response.Header().Set("Connection", "keep-alive")
	r.Response.Header().Set("X-Accel-Buffering", "no")
	r.Response.Header().Set("Access-Control-Allow-Origin", "*")
	channel := r.Get("channel").String()
	if channel == "" {
		channel = "default"
	}

	clientChan := make(chan string)

	clientsMu.Lock()
	clients[r] = clientInfo{
		channel: channel,
		ch:      clientChan,
	}
	channelConnections[channel]++
	if handler, exists := channelHandlers[channel]; exists && channelConnections[channel] == 1 {
		go handler.OnConnect(channel)
	}
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		channelConnections[channel]--
		if handler, exists := channelHandlers[channel]; exists && channelConnections[channel] <= 0 {
			go handler.OnDisconnect(channel)
			delete(channelConnections, channel)
		}
		delete(clients, r)
		clientsMu.Unlock()
		close(clientChan)
	}()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-clientChan:
			event := fmt.Sprintf("data: %s\n\n", msg)
			r.Response.Write(event)
			r.Response.Flush()
		}
	}
}

func BroadcastMessage(channel, message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	if channelConnections[channel] <= 0 {
		return
	}

	for _, info := range clients {
		if info.channel == channel {
			select {
			case info.ch <- message:
			default:
			}
		}
	}
}

func GetChannelConnections(channel string) int {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	return channelConnections[channel]
}

func HasChannelConnections(channel string) bool {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	return channelConnections[channel] > 0
}

func StartAutoSend(channel string, interval time.Duration, fn func()) {
	autoSendMu.Lock()
	defer autoSendMu.Unlock()

	if ticker, exists := autoSendTickers[channel]; exists {
		ticker.Stop()
	}

	ticker := time.NewTicker(interval)
	autoSendTickers[channel] = ticker

	go func() {
		for range ticker.C {
			if HasChannelConnections(channel) {
				fn()
			}
		}
	}()
}

func StopAutoSend(channel string) {
	autoSendMu.Lock()
	defer autoSendMu.Unlock()

	if ticker, exists := autoSendTickers[channel]; exists {
		ticker.Stop()
		delete(autoSendTickers, channel)
	}
}
