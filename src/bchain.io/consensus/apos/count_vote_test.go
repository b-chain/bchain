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
// @File: count_vote_test.go
// @Date: 2018/07/19 14:15:19
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"bchain.io/common/types"
	"testing"
	"time"
)

func commitVote(s int, hash types.Hash) {
	logger.Info("Test commitVote [s hash]:", s, hash.String())
}

func hang() {
	logger.Info("Test hang")
}

func TestCvRun(t *testing.T) {
	cv := newCountVote(commitVote, hang, types.Hash{})
	go cv.run()
	time.Sleep(1 * time.Second)
	cv.startTimer(1)
	time.Sleep(500 * time.Second)
	cv.stopCh <- 1
}

func TestVoteSuccess(t *testing.T) {
	cv := newCountVote(commitVote, hang, types.Hash{})
	go cv.run()
	time.Sleep(1 * time.Second)
	cv.startTimer(int(Config().delayStep))
	hash := types.Hash{}
	hash[1] = 1

	Config().tStepThreshold = 50
	for i := 0; i < 52; i++ {
		ba := &ByzantineAgreementStar{hash, nil, &CredentialSign{Step: STEP_REDUCTION_1, votes: 1}}
		cv.sendMsg(ba)
	}

	time.Sleep(5 * time.Second)
	cv.stopCh <- 1
}

func TestVoteSuccess1(t *testing.T) {
	cv := newCountVote(commitVote, hang, types.Hash{})
	go cv.run()
	time.Sleep(1 * time.Second)
	cv.startTimer(int(Config().delayStep))
	hash := types.Hash{}
	hash[1] = 1

	Config().tStepThreshold = 50
	for i := 0; i < 52; i++ {
		ba := &ByzantineAgreementStar{hash, nil, &CredentialSign{Step: STEP_REDUCTION_2, votes: 1}}
		cv.sendMsg(ba)
	}

	time.Sleep(11 * time.Second)
	cv.stopCh <- 1
}

func TestVoteSuccess_bba(t *testing.T) {
	cv := newCountVote(commitVote, hang, types.Hash{})
	go cv.run()
	time.Sleep(1 * time.Second)
	cv.startTimer(int(Config().delayStep))
	hash := types.Hash{}
	hash[1] = 1

	Config().tStepThreshold = 50
	time.Sleep(2 * time.Second)
	for i := 0; i < 52; i++ {
		ba := &ByzantineAgreementStar{hash, nil, &CredentialSign{Step: 1, votes: 1}}
		cv.sendMsg(ba)
	}

	time.Sleep(30 * time.Second)
	cv.stopCh <- 1
}

func TestVoteSuccess_reduction_bba(t *testing.T) {
	cv := newCountVote(commitVote, hang, types.Hash{})
	go cv.run()
	time.Sleep(1 * time.Second)
	cv.startTimer(int(Config().delayStep))
	hash := types.Hash{}
	hash[1] = 1

	Config().tStepThreshold = 50
	for i := 0; i < 52; i++ {
		ba := &ByzantineAgreementStar{hash, nil, &CredentialSign{Step: STEP_REDUCTION_1, votes: 1}}
		cv.sendMsg(ba)
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < 52; i++ {
		ba := &ByzantineAgreementStar{hash, nil, &CredentialSign{Step: STEP_REDUCTION_2, votes: 1}}
		cv.sendMsg(ba)
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < 52; i++ {
		ba := &ByzantineAgreementStar{hash, nil, &CredentialSign{Step: 1, votes: 1}}
		cv.sendMsg(ba)
	}

	time.Sleep(30 * time.Second)
	cv.stop()
}

func TestVoteSuccess_reduction_bba_final(t *testing.T) {
	cv := newCountVote(commitVote, hang, types.Hash{})
	go cv.run()
	time.Sleep(1 * time.Second)
	cv.startTimer(int(Config().delayStep))
	hash := types.Hash{}
	hash[1] = 1

	Config().tStepThreshold = 50
	for i := 0; i < 52; i++ {
		ba := &ByzantineAgreementStar{hash, nil, &CredentialSign{Step: STEP_REDUCTION_1, votes: 1}}
		cv.sendMsg(ba)
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < 52; i++ {
		ba := &ByzantineAgreementStar{hash, nil, &CredentialSign{Step: STEP_REDUCTION_2, votes: 1}}
		cv.sendMsg(ba)
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < 52; i++ {
		ba := &ByzantineAgreementStar{hash, nil, &CredentialSign{Step: 1, votes: 1}}
		cv.sendMsg(ba)
	}
	time.Sleep(1 * time.Second)
	Config().tFinalThreshold = 50
	for i := 0; i < 52; i++ {
		ba := &ByzantineAgreementStar{hash, nil, &CredentialSign{Step: STEP_FINAL, votes: 1}}
		cv.sendMsg(ba)
	}

	time.Sleep(30 * time.Second)
	cv.stopCh <- 1
}
