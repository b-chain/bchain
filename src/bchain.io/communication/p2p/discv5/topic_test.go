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
// @File: topic_test.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package discv5

import (
	"encoding/binary"
	"testing"
	"time"

	"bchain.io/common/mclock"
	"bchain.io/common/types"
)

func TestTopicRadius(t *testing.T) {
	now := mclock.Now()
	topic := Topic("qwerty")
	rad := newTopicRadius(topic)
	targetRad := (^uint64(0)) / 100

	waitFn := func(addr types.Hash) time.Duration {
		prefix := binary.BigEndian.Uint64(addr[0:8])
		dist := prefix ^ rad.topicHashPrefix
		relDist := float64(dist) / float64(targetRad)
		relTime := (1 - relDist/2) * 2
		if relTime < 0 {
			relTime = 0
		}
		return time.Duration(float64(targetWaitTime) * relTime)
	}

	bcnt := 0
	cnt := 0
	var sum float64
	for cnt < 100 {
		addr := rad.nextTarget(false).target
		wait := waitFn(addr)
		ticket := &ticket{
			topics:  []Topic{topic},
			regTime: []mclock.AbsTime{mclock.AbsTime(wait)},
			node:    &Node{nodeNetGuts: nodeNetGuts{sha: addr}},
		}
		rad.adjustWithTicket(now, addr, ticketRef{ticket, 0})
		if rad.radius != maxRadius {
			cnt++
			sum += float64(rad.radius)
		} else {
			bcnt++
			if bcnt > 500 {
				t.Errorf("Radius did not converge in 500 iterations")
			}
		}
	}
	avgRel := sum / float64(cnt) / float64(targetRad)
	if avgRel > 1.05 || avgRel < 0.95 {
		t.Errorf("Average/target ratio is too far from 1 (%v)", avgRel)
	}
}
