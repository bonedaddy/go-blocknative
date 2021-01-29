package client

import "sync"

// MsgHistory is used to store a copy of all messages we send
// such that in the event of connection drops we can re-establish
// our state
type MsgHistory struct {
	mx     sync.RWMutex
	buffer []interface{}
}

// Push is used to push a message onto our buffer
func (mg *MsgHistory) Push(msg interface{}) {
	mg.mx.Lock()
	defer mg.mx.Unlock()
	mg.buffer = append(mg.buffer, msg)
}

// Pop is used to pop a message out of the buffer
func (mg *MsgHistory) Pop() interface{} {
	mg.mx.Lock()
	defer mg.mx.Unlock()
	if len(mg.buffer) == 0 {
		return nil
	}
	if len(mg.buffer) == 1 {
		item := mg.buffer[0]
		mg.buffer = nil
		return item
	}
	x, buff := mg.buffer[0], mg.buffer[1:]
	mg.buffer = buff
	return x
}

// PopAll returns all elements from the buffer, resetting the buffer
func (mg *MsgHistory) PopAll() []interface{} {
	mg.mx.Lock()
	defer mg.mx.Unlock()
	copied := make([]interface{}, len(mg.buffer))
	copy(copied, mg.buffer)
	mg.buffer = nil
	return copied
}

// Len returns the length of the msg history buffewr
func (mg *MsgHistory) Len() int {
	mg.mx.RLock()
	defer mg.mx.RUnlock()
	return len(mg.buffer)
}
