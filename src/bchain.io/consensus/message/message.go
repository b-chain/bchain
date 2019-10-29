////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The bchain-go Authors.
//
// The bchain-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: message.go
// @Date: 2018/06/13 09:34:13
////////////////////////////////////////////////////////////////////////////////

package message

import (
	"bchain.io/utils/event"
	"sync"
)

// the channel of consensus' message
type messageChannel struct {
	data chan interface{}
	stop chan struct{}
}

// each message must implement this interface
// external interface
type Message interface {
	Send() error // send message
	Close()      // close message processing
}

// each message must implement this interface
// to handle message
type Handler interface {
	Handle(h Handleable)
}

// each message must implement this interface
// to handle data and stop
type Handleable interface {
	DataHandle(data interface{})
	StopHandle()
}

// private message struct
// The basic structure and interface of the message are implemented and could be inherited
type MsgPriv struct {
	channel messageChannel
}

func NewMsgPriv() *MsgPriv {
	msg := MsgPriv{
		channel: messageChannel{
			data: make(chan interface{}),
			stop: make(chan struct{}),
		},
	}
	return &msg
}

func (msg MsgPriv) Send() error {
	msg.channel.data <- msg
	return nil
}

func (msg MsgPriv) Close() {
	close(msg.channel.stop)
}

func (msg MsgPriv) Handle(h Handleable) {
	for {
		select {
		case data := <-msg.channel.data:
			h.DataHandle(data)
		case <-msg.channel.stop:
			h.StopHandle()
			return
		}
	}
}

func isHandler(msg interface{}) bool {
	_, ok := msg.(Handler)
	return ok
}

func getHandler(msg interface{}) Handler {
	hd, ok := msg.(Handler)
	if !ok {
		panic("not a Handler")
	}
	return hd
}

func isHandleable(msg interface{}) bool {
	_, ok := msg.(Handleable)
	return ok
}

func getHandleable(msg interface{}) Handleable {
	handle, ok := msg.(Handleable)
	if !ok {
		panic("not a Handleable")
	}
	return handle
}

// TODO:
type msgcore struct {
	scope event.SubscriptionScope
}

// about msgcore singleton
var (
	instance *msgcore
	once     sync.Once
)

// get the msgcore singleton
func Msgcore() *msgcore {
	once.Do(func() {
		instance = &msgcore{}
	})
	return instance
}

// go routine, handle msg
func (mc msgcore) Handle(msg interface{}) {
	Handler := getHandler(msg)
	h := getHandleable(msg)
	go Handler.Handle(h)
}

// Track starts tracking a subscription of consensus
func (mc msgcore) SubTrack(s event.Subscription) event.Subscription {
	return mc.scope.Track(s)
}

// Close calls Unsubscribe on all tracked subscriptions and prevents further additions to
// the tracked set. Calls to Track after Close return nil.
func (mc msgcore) SubClose() {
	go mc.scope.Close()
}
