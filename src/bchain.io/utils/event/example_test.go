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
// @File: example_test.go
// @Date: 2018/05/07 09:30:07
////////////////////////////////////////////////////////////////////////////////

package event

import "fmt"

func ExampleTypeMux() {
	type someEvent struct{ I int }
	type otherEvent struct{ S string }
	type yetAnotherEvent struct{ X, Y int }

	var mux TypeMux

	// Start a subscriber.
	done := make(chan struct{})
	sub := mux.Subscribe(someEvent{}, otherEvent{})
	go func() {
		for event := range sub.Chan() {
			fmt.Printf("Received: %#v\n", event.Data)
		}
		fmt.Println("done")
		close(done)
	}()

	// Post some events.
	mux.Post(someEvent{5})
	mux.Post(yetAnotherEvent{X: 3, Y: 4})
	mux.Post(someEvent{6})
	mux.Post(otherEvent{"whoa"})

	// Stop closes all subscription channels.
	// The subscriber goroutine will print "done"
	// and exit.
	mux.Stop()

	// Wait for subscriber to return.
	<-done

	// Output:
	// Received: event.someEvent{I:5}
	// Received: event.someEvent{I:6}
	// Received: event.otherEvent{S:"whoa"}
	// done
}
