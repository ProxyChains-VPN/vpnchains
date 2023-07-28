package main

import "sync"

func NewMap() *Map {
	return &Map{
		m:    make(map[packetOwner]packet),
		subs: make(map[packetOwner][]chan packet),
	}
}

type Map struct {
	sync.Mutex

	m    map[packetOwner]packet
	subs map[packetOwner][]chan packet
}

func (m *Map) Set(key packetOwner, value packet) {
	m.Lock()
	defer m.Unlock()

	m.m[key] = value

	// Send the new value to all waiting subscribers of the key
	for _, sub := range m.subs[key] {
		sub <- value
	}
	delete(m.subs, key)
}

func (m *Map) Wait(key packetOwner) packet {
	m.Lock()
	// Unlock cannot be deferred so we can unblock Set() while waiting

	value, ok := m.m[key]
	if ok {
		delete(m.m, key)
		m.Unlock()
		return value
	}

	// if there is no value yet, subscribe to any new values for this key
	ch := make(chan packet)
	m.subs[key] = append(m.subs[key], ch)
	m.Unlock()

	return <-ch
}
