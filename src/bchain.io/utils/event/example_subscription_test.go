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
// @File: example_subscription_test.go
// @Date: 2018/05/07 09:30:07
////////////////////////////////////////////////////////////////////////////////

package event_test

import (
	"fmt"

	"bchain.io/utils/event"
)

func ExampleNewSubscription() {
	// Create a subscription that sends 10 integers on ch.
	ch := make(chan int)
	sub := event.NewSubscription(func(quit <-chan struct{}) error {
		for i := 0; i < 10; i++ {
			select {
			case ch <- i:
			case <-quit:
				fmt.Println("unsubscribed")
				return nil
			}
		}
		return nil
	})

	// This is the consumer. It reads 5 integers, then aborts the subscription.
	// Note that Unsubscribe waits until the producer has shut down.
	for i := range ch {
		fmt.Println(i)
		if i == 4 {
			sub.Unsubscribe()
			break
		}
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// unsubscribed
}
