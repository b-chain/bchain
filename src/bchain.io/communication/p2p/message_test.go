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
// @File: message_test.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package p2p

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"io"
	"runtime"
	"strings"
	"testing"
	"time"
)

func ExampleMsgPipe() {
	rw1, rw2 := MsgPipe()
	go func() {
		var num msgp.Number
		num.AsUint(0x0000)
		Send(rw1, 8, &num)
		num.AsUint(0x0101)
		Send(rw1, 5, &num)
		rw1.Close()
	}()

	for {
		msg, err := rw2.ReadMsg()
		if err != nil {
			break
		}
		var data msgp.Number
		msg.Decode(&data)
		a, _ := data.Int()
		fmt.Printf("msg: %d, %04x\n", msg.Code, a)
	}
	// Output:
	// msg: 8, 0000
	// msg: 5, 0101
}

func TestMsgPipeUnblockWrite(t *testing.T) {
loop:
	for i := 0; i < 100; i++ {
		rw1, rw2 := MsgPipe()
		done := make(chan struct{})
		go func() {
			if err := SendItems(rw1, 1); err == nil {
				t.Error("EncodeMsg returned nil error")
			} else if err != ErrPipeClosed {
				t.Errorf("EncodeMsg returned wrong error: got %v, want %v", err, ErrPipeClosed)
			}
			close(done)
		}()

		// this call should ensure that EncodeMsg is waiting to
		// deliver sometimes. if this isn't done, Close is likely to
		// be executed before EncodeMsg starts and then we won't test
		// all the cases.
		runtime.Gosched()

		rw2.Close()
		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
			t.Errorf("write didn't unblock")
			break loop
		}
	}
}

// This test should panic if concurrent close isn't implemented correctly.
func TestMsgPipeConcurrentClose(t *testing.T) {
	rw1, _ := MsgPipe()
	for i := 0; i < 10; i++ {
		go rw1.Close()
	}
}

func TestEOFSignal(t *testing.T) {
	rb := make([]byte, 10)

	// empty reader
	eof := make(chan struct{}, 1)
	sig := &eofSignal{new(bytes.Buffer), 0, eof}
	if n, err := sig.Read(rb); n != 0 || err != io.EOF {
		t.Errorf("Read returned unexpected values: (%v, %v)", n, err)
	}
	select {
	case <-eof:
	default:
		t.Error("EOF chan not signaled")
	}

	// count before error
	eof = make(chan struct{}, 1)
	sig = &eofSignal{bytes.NewBufferString("aaaaaaaa"), 4, eof}
	if n, err := sig.Read(rb); n != 4 || err != nil {
		t.Errorf("Read returned unexpected values: (%v, %v)", n, err)
	}
	select {
	case <-eof:
	default:
		t.Error("EOF chan not signaled")
	}

	// error before count
	eof = make(chan struct{}, 1)
	sig = &eofSignal{bytes.NewBufferString("aaaa"), 999, eof}
	if n, err := sig.Read(rb); n != 4 || err != nil {
		t.Errorf("Read returned unexpected values: (%v, %v)", n, err)
	}
	if n, err := sig.Read(rb); n != 0 || err != io.EOF {
		t.Errorf("Read returned unexpected values: (%v, %v)", n, err)
	}
	select {
	case <-eof:
	default:
		t.Error("EOF chan not signaled")
	}

	// no signal if neither occurs
	eof = make(chan struct{}, 1)
	sig = &eofSignal{bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaa"), 999, eof}
	if n, err := sig.Read(rb); n != 10 || err != nil {
		t.Errorf("Read returned unexpected values: (%v, %v)", n, err)
	}
	select {
	case <-eof:
		t.Error("unexpected EOF signal")
	default:
	}
}

func unhex(str string) []byte {
	r := strings.NewReplacer("\t", "", " ", "", "\n", "")
	b, err := hex.DecodeString(r.Replace(str))
	if err != nil {
		panic(fmt.Sprintf("invalid hex string: %q", str))
	}
	return b
}
