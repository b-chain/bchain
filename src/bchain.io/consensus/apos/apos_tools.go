package apos

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"math/big"
	"bchain.io/common/types"
	"bchain.io/consensus/apos/aposgensis"
	"bchain.io/core"
	"bchain.io/core/actioncontext"
	"bchain.io/core/blockchain/block"
	"bchain.io/core/interpreter"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/params"
	"bchain.io/utils/crypto"
	"bchain.io/utils/database"
	"bchain.io/utils/event"
	"reflect"
	"sync"
	"time"
	"bchain.io/core/interpreter/wasmre/para_paser"
	"encoding/json"
	"encoding/binary"
	"math"
)

type PriKeyHandler interface {
	GetBasePriKey(kind reflect.Type) *ecdsa.PrivateKey
}

type BlockChainHandler interface {
	CurrentBlock() *block.Block
	CurrentBlockNum() uint64
	GetNowBlockHash() types.Hash
	InsertChain(chain block.Blocks) (int, error)
	GetBlockByNumber(number uint64) *block.Block
	GetHeaderByNumber(number uint64) *block.Header
	StateAt(root types.Hash) (*state.StateDB, error)
	GetDb() database.IDatabase
	GetBlockByHash(hash types.Hash) *block.Block
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
	GetConsensusData(key types.Hash) []byte
	WriteConsensusData(key types.Hash, value []byte) error
	GetExtra(key []byte) []byte
	VerifyNextRoundBlock(block *block.Block) bool
}

type BlockProducerHandler interface {
	GetProducerNewBlock(data *block.ConsensusData, timeLimit int64) *block.Block
}

func generatePrivateKey() *ecdsa.PrivateKey {
	randBytes := make([]byte, 64)

	_, err := rand.Read(randBytes)
	if err != nil {
		panic("key generation: could not read from random source: " + err.Error())
	}
	reader := bytes.NewReader(randBytes)

	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), reader)
	if err != nil {
		panic("key generation: ecdsa.GenerateKey failed: " + err.Error())
	}
	return privateKeyECDSA
}

type bootWeightInfo struct {
	totalWeiht     uint64
	bootCommittees map[types.Address]uint64
}

type aposTools struct {
	lock              sync.RWMutex
	chanId            *big.Int
	basePriKey        *ecdsa.PrivateKey
	coinBase          types.Address
	tmpPriKeys        map[int]*ecdsa.PrivateKey
	blockChainHandler BlockChainHandler
	producerHandler   BlockProducerHandler
	bootLock          sync.RWMutex
	bootWeightInfo    bootWeightInfo
}

func newAposTools(chanId *big.Int, bcHandler BlockChainHandler, producerHander BlockProducerHandler) *aposTools {
	a := new(aposTools)
	//a.chanId = big.Int{}.Set(chanId)
	a.chanId = big.NewInt(0).Set(chanId)
	a.basePriKey = nil
	a.blockChainHandler = bcHandler
	a.producerHandler = producerHander
	a.loadBootCommittees()
	return a
}

func (this *aposTools) loadBootCommittees() {
	this.bootWeightInfo.bootCommittees = make(map[types.Address]uint64)
	data := this.blockChainHandler.GetExtra(aposgensis.BootCommitteeKey)
	if data == nil {
		logger.Warn("load boot Committees fail")
		return
	}
	var winfo aposgensis.WeightInfos
	byteBuf := bytes.NewBuffer(data)
	err := msgp.Decode(byteBuf, &winfo)
	if err != nil {
		logger.Error("WeightInfos.Decode err", err)
		return
	}

	for _, wi := range winfo {
		this.bootWeightInfo.bootCommittees[wi.AddrStr] = wi.Wt
		this.bootWeightInfo.totalWeiht += wi.Wt
	}
	fmt.Println("loadBootCommittees", this.bootWeightInfo)
}

func (this *aposTools) CreateTmpPriKey(step int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return
	if this.tmpPriKeys == nil {
		this.tmpPriKeys = make(map[int]*ecdsa.PrivateKey)
	}

	if _, ok := this.tmpPriKeys[step]; ok {
		return
	}

	tmpKey := generatePrivateKey()
	this.tmpPriKeys[step] = tmpKey
}

func (this *aposTools) SetPriKey(priKey *ecdsa.PrivateKey) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.basePriKey = priKey
}

func (this *aposTools) SetCoinBase(coinbase types.Address) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.coinBase = coinbase
}

func (this *aposTools) GetCoinBase() types.Address {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.coinBase
}

func (this *aposTools) Sig(pCs *CredentialSign) error {
	this.lock.RLock()
	defer this.lock.RUnlock()

	_, _, _, err := pCs.sign(this.basePriKey)
	return err
}

func (this *aposTools) SeedSig(pSd *SeedData) error {
	this.lock.RLock()
	defer this.lock.RUnlock()

	_, _, _, err := pSd.sign(this.basePriKey)
	return err
}

func (this *aposTools) Esig(pEphemeralSign *EphemeralSign) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	step := int(pEphemeralSign.GetStep())

	if true {
		pEphemeralSign.Signature.Init()
		_, _, _, err := pEphemeralSign.sign(this.basePriKey)
		return err
	} else {
		if pri, ok := this.tmpPriKeys[step]; ok {
			pEphemeralSign.Signature.Init()
			_, _, _, err := pEphemeralSign.sign(pri)
			return err
		}
	}

	return errors.New(fmt.Sprintf("Not Find TmpPrivKey About:%d", step))
}

func (this *aposTools) DelTmpKey(step int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return
	if _, ok := this.tmpPriKeys[step]; ok {
		delete(this.tmpPriKeys, step)
	}
}

func (this *aposTools) ClearTmpKeys() {
	this.lock.Lock()
	defer this.lock.Unlock()

	return
	this.tmpPriKeys = nil
}

func (this *aposTools) SigHash(hash types.Hash) []byte {
	this.lock.RLock()
	defer this.lock.RUnlock()

	sig, err := crypto.Sign(hash[:], this.basePriKey)
	if err != nil {
		logger.Error("aposTools SigErr:", err.Error())
		return nil
	}

	return sig
}

func (this *aposTools) ESigVerify(h types.Hash, sig []byte) error {
	return nil
}

func (this *aposTools) ESender(hash types.Hash, sig []byte) (types.Address, error) {
	return types.Address{}, nil
}

func (this *aposTools) GetLastQrSignature() []byte {
	blk := this.blockChainHandler.CurrentBlock()
	if blk == nil {
		return nil
	}
	return blk.H.Cdata.Para
}

func (this *aposTools) GetQrSignature(round uint64) []byte {
	blk := this.blockChainHandler.GetBlockByNumber(round)
	if blk == nil {
		return nil
	}
	return blk.H.Cdata.Para
}

func (this *aposTools) GetNowBlockNum() uint64 {
	return this.blockChainHandler.CurrentBlockNum()
}

func (this *aposTools) GetNextRound() int {
	return int(this.blockChainHandler.CurrentBlockNum() + 1)
}

func (this *aposTools) MakeEmptyBlock(data *block.ConsensusData) *block.Block {
	parent := this.blockChainHandler.CurrentBlock()

	header := block.CopyHeader(parent.H)
	header.ParentHash = parent.Hash()

	//r = r-1 + 1
	header.Number = types.NewBigInt(*big.NewInt(header.Number.IntVal.Int64() + 1))
	header.Cdata = *data

	header.Bloom = types.Bloom{}

	b := block.NewBlock(header, nil, nil)
	//use system private key to sign the block
	err := block.SignHeaderInner(b.H, block.NewBlockSigner(Config().GetChainId()), params.RewordPrikey)
	if err != nil {
		logger.Error("makeEmptyBlock error:", err)
		return nil
	}
	return b
}

func (this *aposTools) VerifyNextRoundBlock(block *block.Block) bool {
	return this.blockChainHandler.VerifyNextRoundBlock(block)
}

func (this *aposTools) GetNowBlockHash() types.Hash {
	return this.blockChainHandler.GetNowBlockHash()
}

func (this *aposTools) GetProducerNewBlock(data *block.ConsensusData, timeLimit int64) *block.Block {
	return this.producerHandler.GetProducerNewBlock(data, timeLimit)
}

func (this *aposTools) InsertChain(chain block.Blocks) (int, error) {
	return this.blockChainHandler.InsertChain(chain)
}

func (this *aposTools) GetCurrentBlock() *block.Block {
	return this.blockChainHandler.CurrentBlock()
}
func (this *aposTools) GetBlockByNum(num uint64) *block.Block {
	return this.blockChainHandler.GetBlockByNumber(uint64(num))
}
func (this *aposTools) GetBlockByHash(hash types.Hash) *block.Block {
	return this.blockChainHandler.GetBlockByHash(hash)
}
func (this *aposTools) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return this.blockChainHandler.SubscribeChainEvent(ch)
}

func (this *aposTools) GetBlockCertificate(blockHash types.Hash) BlockCertificate {
	msgpData := this.blockChainHandler.GetConsensusData(blockHash)
	if msgpData == nil || len(msgpData) == 0 {
		return nil
	}
	var blockCertificate BlockCertificate
	byteBuf := bytes.NewBuffer(msgpData)
	err := msgp.Decode(byteBuf, &blockCertificate)
	if err != nil {
		logger.Error("blockCertificate.Decode err", "hash", blockHash, "err", err)
		return nil
	}
	return blockCertificate
}

func (this *aposTools) WriteBlockCertificate(blk *block.Block, certificate BlockCertificate) error {
	var encData bytes.Buffer
	err := msgp.Encode(&encData, certificate)
	if err != nil {
		return err
	}

	key := blk.Hash()
	return this.blockChainHandler.WriteConsensusData(key, encData.Bytes())
}

func (this *aposTools) getWeight(r uint64, addr types.Address) (int64, int64) {
	this.bootLock.Lock()
	defer this.bootLock.Unlock()
	if weight, ok := this.bootWeightInfo.bootCommittees[addr]; ok {
		return int64(weight), int64(this.bootWeightInfo.totalWeiht)
	}
	return 0, int64(this.bootWeightInfo.totalWeiht)
}

//todo: should change this loop, just a test demo
//should delay when download
func (this *aposTools) delayWeight(pledgeNumber uint64) {
	var currentNumber uint64
	for {
		currentNumber = this.blockChainHandler.CurrentBlockNum()
		if pledgeNumber > currentNumber {
			time.Sleep(time.Millisecond)
		} else {
			break
		}
	}
}

//todo: just a demo, jsre is very slow alredy change to wasm speed is similer to local
const (
	configRound = 100
	pledgEffectiveRound = 997
)
func (this *aposTools) GetWeight(r uint64, addr types.Address) (int64, int64) {
	//return this.getWeight(r, addr)
	if r < configRound + pledgEffectiveRound{
		return this.getWeight(r, addr)
	}
	pledgeNumber := r - pledgEffectiveRound
	this.delayWeight(pledgeNumber)
	header := this.blockChainHandler.GetHeaderByNumber(pledgeNumber)
	if header == nil {
		logger.Error("GetWeight err : can not get header by number", pledgeNumber)
		return 0, 1
	}

	singner := block.NewBlockSigner(this.chanId)
	coinbase, err := singner.Sender(header)

	state, err := this.blockChainHandler.StateAt(header.StateRootHash)
	if err != nil {
		logger.Error("get state err", err)
		return 0, 1
	}

	act := &transaction.Action{}
	act.Contract = types.HexToAddress("0xFa58d9f83D1D86DF22435e67D5F7422337624737")
	arg := para_paser.Arg{para_paser.TypeAddress, []byte(addr.HexLower())}
	wp := &para_paser.WasmPara {
		FuncName: "pledgeOfExt",
		Args:     append([]para_paser.Arg{}, arg),
	}
	para, _ := json.Marshal(wp)
	act.Params = []byte(para)
	sender := types.Address{}

	tmpDb, _ := database.OpenMemDB()
	blkCtx := actioncontext.NewBlockContext(state, this.blockChainHandler.GetDb(), tmpDb, &header.Number.IntVal, coinbase)
	actCxt := actioncontext.NewContext(sender, act, blkCtx)
	if actCxt == nil {
		logger.Warn("contract query, new context return nil by contract ", act.Contract.Hex())
		return 0, 1
	}
	err = actCxt.Exec(interpreter.Singleton())
	if err != nil {
		return 0, 1
	}

	wOrignal := int64(binary.LittleEndian.Uint64(actCxt.ActionResult()[0]))
	W := int64(binary.LittleEndian.Uint64(actCxt.ActionResult()[2]))
	//w, _ := new(big.Int).SetString(string(actCxt.ActionResult()[0]), 10)
	//W, _ := new(big.Int).SetString(string(actCxt.ActionResult()[1]), 10)

	w := int64(0)
	if wOrignal >= 10000000000 {
		w = int64(math.Pow(float64(wOrignal),0.33))
	}

	return w, W
}
