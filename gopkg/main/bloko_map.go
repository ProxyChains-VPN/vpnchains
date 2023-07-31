package main

import (
	"log"
	"sync"
	"time"
)

const MaxPacketsQueueSize = 16

type PacketsBuffer struct {
	sync.Mutex
	subs map[PacketOwner]chan *Packet
}

func NewPacketsBuffer() *PacketsBuffer {
	log.Println("creating packets buffer")
	return &PacketsBuffer{
		subs: make(map[PacketOwner]chan *Packet),
	}
}

func (buf *PacketsBuffer) PushPacket(key PacketOwner, value *Packet) {
	log.Println("push packet", key)
	buf.Lock()
	defer buf.Unlock()

	log.Println(len(buf.subs))
	if _, ok := buf.subs[key]; !ok {
		log.Println("l27 creating channel for", key)
		buf.subs[key] = make(chan *Packet, MaxPacketsQueueSize)
	}

	buf.subs[key] <- value
	log.Println("pushed packet", key)
}

func (buf *PacketsBuffer) WaitForPacket(key PacketOwner) *Packet {
	log.Println("wait for packet", key)
	buf.Lock()

	log.Println(len(buf.subs))
	if _, ok := buf.subs[key]; !ok {
		log.Println("l39 creating channel for", key)
		buf.subs[key] = make(chan *Packet, MaxPacketsQueueSize)
	}
	buf.Unlock()

	select {
	case packet := <-buf.subs[key]:
		log.Println("got packet", key, packet)
		return packet
	case <-time.After(time.Second * 2):
		return nil
	}
}
