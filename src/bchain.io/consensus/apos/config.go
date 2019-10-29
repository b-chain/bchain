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
// @File: types.go
// @Date: 2018/06/12 11:01:51
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"fmt"
	"math/big"
	"bchain.io/common"
	"bchain.io/params"
	"sync"
)

var (
	decimal         = big.NewInt(10)
	honestPercision = big.NewInt(100)
	maxUint256      = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
)

//go:generate gencodec -type config -field-override configMarshaling -out gen_config.go

//some system param(apos system param) for step goroutine.
type config struct {
	lookback         int    `json:"lookback"`             // lookback val, r - k
	prPrecision      uint64 `json:"precision"`            // the precision
	prLeader         uint64 `json:"probability-leader"`   // the probability of Leaders
	prVerifier       uint64 `json:"probability-verifier"` // the probability of Verifiers
	maxBBASteps      int    `json:"max-steps"`            // the max number of BBA steps
	maxNodesPerRound int    `json:"max-nodes-per-round"`  // the max number of nodes per round
	prH              uint64 `json:"probability-honest"`   // the probability of honest
	blockDelay       int    `json:"block-delay"`          // time A, sec
	verifyDelay      int    `json:"verify-delay"`         // time λ, sec

	prP             *big.Int `json:"-"` // 10 ^ prPrecision
	maxPotLeaders   *big.Int `json:"-"` // the max number of potential leaders
	maxPotVerifiers *big.Int `json:"-"` // the max number of potential verifiers

	// chain info
	chainId    *big.Int `json:"-"`
	chainIdMul *big.Int `json:"-"`

	//new struct
	R               uint  `json:"r"`               // seed refresh interval (# of rounds)
	tProposer       int64 `json:"tProposer"`       // expected # of block proposers
	tStep           int64 `json:"tStep"`           // expected # of committee members
	tStepThreshold  int64 `json:"tStepThreshold"`  // threshold # of τstep for BA⋆
	tFinal          int64 `json:"tFinal"`          // expected # of final committee members
	tFinalThreshold int64 `json:"tFinalThreshold"` // threshold # of τfinal for BA⋆
	maxStep         uint  `json:"maxStep"`         // maximum number of steps in BinaryBA⋆
	delayPriority   uint  `json:"delayPriority"`   // time to gossip sortition proofs
	delayStep       uint  `json:"delayStep"`       // timeout for receiving a block
	delayBlock      uint  `json:"delayBlock"`      // timeout for BA⋆ step
	delayStepVar    uint  `json:"delayStepVar"`    // estimate of BA⋆ completion time variance

	sp            SortitionPriority
	delayRecovery uint `json:"delayRecovery"` // timeout for recovery protocol
}

func (c *config) setDefault() {
	c.lookback = 100
	c.prPrecision = 10
	c.prLeader = 1000000000   // 0.1
	c.prVerifier = 5000000000 // 0.5
	c.maxBBASteps = 180
	c.maxNodesPerRound = 10
	c.maxPotLeaders = big.NewInt(3)
	c.maxPotLeaders = big.NewInt(4)
	c.prH = 67
	c.blockDelay = 10
	c.verifyDelay = 5
	c.chainId = big.NewInt(int64(params.DefaultChainId))
	c.chainIdMul = new(big.Int).Mul(c.chainId, common.Big2)

	//new struct
	c.R = 1000
	c.tProposer = 26
	c.tStep = 2000
	c.tStepThreshold = c.tStep * 685 / 1000
	c.tFinal = 10000
	c.tFinalThreshold = c.tFinal * 74 / 100
	c.maxStep = 150
	c.delayPriority = 5
	c.delayStep = 5
	c.delayBlock = 60
	c.delayStepVar = 5

	c.sp = new(binomialDistribution)
	c.delayRecovery = 10
}

func (c *config) setSpecialConfig() {
	//set config
	c.R = 1000
	c.tProposer = 26
	c.tStep = 100
	c.tStepThreshold = 69
	c.tFinal = 200
	c.tFinalThreshold = 148
	c.maxStep = 60
	c.delayPriority = 5
	c.delayStep = 5
	c.delayBlock = 20
	c.delayStepVar = 5
}

// about msgcore singleton
var (
	instance *config
	once     sync.Once
)

// get the msgcore singleton
func Config() *config {
	once.Do(func() {
		instance = &config{}
		instance.setDefault()
		instance.setSpecialConfig()
		instance.Verify()
		//instance.verifier()
		instance.chain()
		fmt.Println(instance)
	})

	return instance
}

func (c *config) GetChainId() *big.Int {
	return c.chainId
}

func (c *config) precision() *big.Int {
	if c.prP == nil {
		c.prP = new(big.Int).Exp(decimal, big.NewInt(0).SetUint64(c.prPrecision), big.NewInt(0))
	}
	return c.prP
}

func (c *config) verifier() (uint64, uint64, uint64, *big.Int, *big.Int) {
	if c.maxPotLeaders == nil {
		c.maxPotLeaders = big.NewInt(int64(c.maxNodesPerRound))
		c.maxPotLeaders.Mul(c.maxPotLeaders, big.NewInt(0).SetUint64(c.prLeader))
		c.maxPotLeaders.Div(c.maxPotLeaders, c.precision())
	}

	if c.maxPotVerifiers == nil {
		c.maxPotVerifiers = big.NewInt(int64(c.maxNodesPerRound))
		c.maxPotVerifiers.Mul(c.maxPotVerifiers, big.NewInt(0).SetUint64(c.prVerifier))
		c.maxPotVerifiers.Div(c.maxPotVerifiers, c.precision())
	}

	return c.prPrecision, c.prLeader, c.prVerifier, c.maxPotLeaders, c.maxPotVerifiers
}

func (c *config) chain() (chainId *big.Int, chainIdMul *big.Int) {
	if c.chainId != nil {
		c.chainIdMul = new(big.Int).Mul(c.chainId, common.Big2)
		return c.chainId, c.chainIdMul
	}
	return chainId, chainIdMul
}

func (c *config) Verify() {
	if c.lookback <= 0 {
		panic(fmt.Errorf("lookback <= 0 \n"))
	}

	if c.maxBBASteps <= 0 {
		panic(fmt.Errorf("maxBBASteps <= 0 \n"))
	}

	if c.maxNodesPerRound <= 0 {
		panic(fmt.Errorf("maxNodesPerRound <= 0 \n"))
	}

	if c.blockDelay <= 0 {
		panic(fmt.Errorf("blockDelay <= 0 \n"))
	}

	if c.verifyDelay <= 0 {
		panic(fmt.Errorf("verifyDelay <= 0 \n"))
	}

	if c.prH == 0 {
		panic(fmt.Errorf("prH == 0 \n"))
	}
	if big.NewInt(0).SetUint64(c.prH).Cmp(honestPercision) > 0 {
		panic(fmt.Errorf("prH > 100 \n"))
	}

	if c.precision().Cmp(maxUint256) > 0 {
		panic(fmt.Errorf("PrLeader > precision \n"))
	}

	if c.precision().Cmp(big.NewInt(0).SetUint64(c.prLeader)) < 0 {
		panic(fmt.Errorf("prLeader < precision \n"))
	}

	if c.precision().Cmp(big.NewInt(0).SetUint64(c.prVerifier)) < 0 {
		panic(fmt.Errorf("prVerifier < Precision \n"))
	}
}
