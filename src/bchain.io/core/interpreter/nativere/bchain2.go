package nativere

import (
	"bchain.io/common/assert"
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"bchain.io/core/transaction"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"math/big"
)

const (
	DecaimasBase  = uint64(100000000)
	PledgeLimit   = uint64(10 * DecaimasBase)
	PooolLimit    = uint64(5000 * DecaimasBase)
	ProducerLimit = uint64(500 * DecaimasBase)
)

type pledge struct {
	ctx *actioncontext.Context
}

func (p *pledge) run(para []byte) {
	np := &NativePara{}
	err := json.Unmarshal(para, np)
	assert.AsserErr(err)
	paraLen := len(np.Args)
	ErrLen := "args length err"

	switch np.FuncName {
	case "rewords":
		fmt.Println("rewords")
		assert.AssertEx(paraLen == 0, ErrLen)
		p.rewords()
	case "transfer":
		fmt.Println("transfer")
		assert.AssertEx(paraLen == 3, ErrLen)
		p.transfer(np.Args[0], np.Args[1], np.Args[2])
	case "transferFee":
		fmt.Println("transfer")
		assert.AssertEx(paraLen == 1, ErrLen)
		p.transferFee(np.Args[0])
	case "pledge":
		fmt.Println("pledge")
		assert.AssertEx(paraLen == 2, ErrLen)
		p.pledge(np.Args[0], np.Args[1])
	case "redeem":
		assert.AssertEx(paraLen == 1, ErrLen)
		fmt.Println("redeem")
		p.redeem(np.Args[0])
	case "proxy":
		assert.AssertEx(paraLen == 1, ErrLen)
		fmt.Println("proxy")
		p.proxy(np.Args[0])
	case "cancelProxy":
		assert.AssertEx(paraLen == 0, ErrLen)
		fmt.Println("cancelProxy")
		p.cancelProxy()
	case "makeProducer":
		assert.AssertEx(paraLen == 1, ErrLen)
		fmt.Println("makeProducer")
		p.makeProducer(np.Args[0])
	case "balanceOf":
		assert.AssertEx(paraLen == 1, ErrLen)
		fmt.Println("balanceOf")
		p.balanceOf(np.Args[0])
	case "pledgeOf":
		assert.AssertEx(paraLen == 1, ErrLen)
		fmt.Println("pledgeOf")
		p.pledgeOf(np.Args[0])
	case "pledgeOfExt":
		assert.AssertEx(paraLen == 1, ErrLen)
		fmt.Println("pledgeOfExt")
		p.pledgeOfExt(np.Args[0])
	default:
		assert.AssertEx(false, "func name is invalid")
	}
}

func (p *pledge) transfer(to string, amount string, memo string) {
	assert.AssertEx(types.IsHexAddress(to), "to addr invalid")
	assert.AssertEx(len(memo) < 64, "memo length invalid")
	amount_int, err := strconv.ParseInt(amount, 10, 64)
	assert.AsserErr(err)
	amount_val := uint64(amount_int)
	assert.AssertEx(amount_val > 0, "amount is invalid")

	toAddr := types.HexToAddress(to)
	from := p.sender()
	assert.AssertEx(from != toAddr, "to address is self")

	fromVal := p.dbGet(from)
	fromVal_uint := binary.LittleEndian.Uint64(fromVal)

	toVal := p.dbGet(toAddr)
	toVal_uint := binary.LittleEndian.Uint64(toVal)

	assert.AssertEx(fromVal_uint >= amount_val, "insufficient balance")
	fromVal_uint -= amount_val
	toVal_uint += amount_val
	newFrom := make([]byte, 8)
	newTo := make([]byte, 8)
	binary.LittleEndian.PutUint64(newFrom, fromVal_uint)
	binary.LittleEndian.PutUint64(newTo, toVal_uint)
	p.dbSet(from, newFrom)
	p.dbSet(toAddr, newTo)
}

func (p *pledge) transferFee(amount string) {
	amount_int, err := strconv.ParseInt(amount, 10, 64)
	assert.AsserErr(err)
	amount_val := uint64(amount_int)
	assert.AssertEx(amount_val > 0, "amount is invalid")

	producer := p.producer()
	from := p.sender()
	if producer == from {
		return
	}

	fromVal := p.dbGet(from)
	fromVal_uint := binary.LittleEndian.Uint64(fromVal)

	toVal := p.dbGet(producer)
	toVal_uint := binary.LittleEndian.Uint64(toVal)

	assert.AssertEx(fromVal_uint >= amount_val, "insufficient balance")
	fromVal_uint -= amount_val
	toVal_uint += amount_val
	newFrom := make([]byte, 8)
	newTo := make([]byte, 8)
	binary.LittleEndian.PutUint64(newFrom, fromVal_uint)
	binary.LittleEndian.PutUint64(newTo, toVal_uint)
	p.dbSet(from, newFrom)
	p.dbSet(producer, newTo)

}

func (p *pledge) requireRewordAuth() {
	memDbApi := actioncontext.BlockMemDbApi{}
	memDbApi.SetCtx(p.ctx)
	memDbApi.Emplace([]byte("producer"), []byte("produced"))
}
func (p *pledge) getCurRewordsNumber() uint64 {
	actApi := actioncontext.DatabaseApi{}
	actApi.SetCtx(p.ctx)
	val := actApi.Get([]byte("_rNumber"))
	number := uint64(0)
	if val != nil {
		number = binary.LittleEndian.Uint64(val)
	}
	number++
	valInc := make([]byte, 8)
	binary.LittleEndian.PutUint64(valInc, number)
	actApi.Set([]byte("_rNumber"), valInc)
	return number
}
func (p *pledge) getRewordsValue(number uint64) uint64 {
	if number == 1 {
		return 2100000 * DecaimasBase
	}
	a0 := uint64(75600000)
	period := uint64(12500000)

	idx := number / period
	a := a0 >> idx
	return a
}
func (p *pledge) addToken(addr types.Address, amount uint64) {
	fromVal := p.dbGet(addr)
	fromVal_uint := binary.LittleEndian.Uint64(fromVal)
	fromVal_uint += amount

	newVal := make([]byte, 8)
	binary.LittleEndian.PutUint64(newVal, fromVal_uint)
	p.dbSet(addr, newVal)
}

func (p *pledge) getTokenAmount(total uint64, numerator, denominator uint64) uint64 {
	totalB := big.NewInt(int64(total))
	numeratorB := big.NewInt(int64(numerator))
	denominatorB := big.NewInt(int64(denominator))

	a := new(big.Int).Mul(totalB, numeratorB)
	b := new(big.Int).Div(a, denominatorB)
	return b.Uint64()
}

func (p *pledge) rewordsPd(headIn types.Address, totalRw uint64, totalAmount uint64)  {
	head := headIn
	for {
		if head == nilAddr {
			break
		}
		hData := p.getPledgeData(head)
		pdRewords := p.getTokenAmount(totalRw, hData.Amount, totalAmount)
		p.addToken(head, pdRewords)
		head = hData.Next
	}
}

func (p *pledge) getRewordsTotalAmount(headIn types.Address, totalAmount uint64) uint64 {
	middle := totalAmount/2
	head := headIn
	ret := uint64(0)
	for {
		if head == nilAddr {
			break
		}
		hData := p.getPoolData(head)
		if hData.Amount >= middle {
			ret += hData.Amount
		}
		head = hData.Next
	}
	return ret
}

func (p *pledge) rewordsPool(headIn types.Address, totalRw uint64, totalAmount uint64) {
	middle := totalAmount/2
	head := headIn
	newTAmount := p.getRewordsTotalAmount(headIn, totalAmount)
	for {
		if head == nilAddr {
			break
		}
		hData := p.getPoolData(head)

		if hData.Amount >= middle {
			pdHeadData := p.getPledgeData(head)
			poolReword := p.getTokenAmount(totalRw, pdHeadData.Amount, newTAmount)
			poolOwnerRewords := p.getTokenAmount(poolReword, 10, 100)
			p.addToken(head, poolOwnerRewords)
			p.rewordsPd(hData.ChildHead, poolReword - poolOwnerRewords, pdHeadData.Amount)
		}
		head = hData.Next
	}
}
var nilAddr = types.Address{}
func (p *pledge) rewords() {
	p.requireRewordAuth()
	producer := p.producer()
	number := p.getCurRewordsNumber()
	rewordsVal := p.getRewordsValue(number)

	pData := p.getProducerData(producer)
	if pData == nil {
		p.addToken(producer, rewordsVal)
		return
	} else {
		producerRewords := p.getTokenAmount(rewordsVal, 10, 100)
		p.addToken(producer, producerRewords)
	}

	leftRewords := p.getTokenAmount(rewordsVal, 90, 100)
	head := pData.ChildHead
	if head == nilAddr {
		p.addToken(producer, leftRewords)
		return
	}
	p.rewordsPool(head, leftRewords, pData.Amount)
}

func (p *pledge) pledge(amount string, beneficiary string) {
	assert.AssertEx(types.IsHexAddress(beneficiary), "to addr invalid")
	amount_int, err := strconv.ParseInt(amount, 10, 64)
	amount_val := uint64(amount_int)
	assert.AsserErr(err)
	assert.AssertEx(amount_val >= PledgeLimit, "amount is lower than limit")

	parent := types.HexToAddress(beneficiary)

	from := p.sender()
	// find
	selfData := p.getPledgeData(from)
	part1init := false
	if selfData == nil {
		part1init = true
		selfData = &PledgeData{
			Parent: parent,
			Amount: amount_val,
		}
	} else {
		if selfData.Amount < PledgeLimit {
			part1init = true
		}
		selfData.Amount += amount_val
	}
	assert.AssertEx(parent == selfData.Parent, "beneficiary is not match last")

	parentData := p.getPoolData(parent)
	if parentData == nil {
		parentData = &PoolData{
			Amount: amount_val,
		}
	} else {
		parentData.Amount += amount_val
	}

	if part1init == true {
		//insert chain
		if parentData.ChildHead == nilAddr {
			parentData.ChildHead = from
		} else {
			lastHead := p.getPledgeData(parentData.ChildHead)
			lastHead.Prev = from
			p.savePledgeData(parentData.ChildHead, lastHead)
			selfData.Next = parentData.ChildHead
			parentData.ChildHead = from
		}
	}
	p.savePledgeData(from, selfData)
	p.savePoolData(parent, parentData)

	producerData := p.getProducerData(parentData.Producer)
	if producerData != nil {
		oriWeight := p.getWeight(producerData.Amount)
		producerData.Amount += amount_val
		newWeight := p.getWeight(producerData.Amount)
		p.saveProducerData(parentData.Producer, producerData)
		p.addTotalPledgeWeight(newWeight - oriWeight)
	}

	p.addTotalPledgeVal(amount_val)
	p.transfer(p.conAddr().HexLower(), amount, "pledge")
}

func (p *pledge) redeem(amount string) {
	amount_int, err := strconv.ParseInt(amount, 10, 64)
	amount_val := uint64(amount_int)
	assert.AsserErr(err)
	assert.AssertEx(amount_val >= 0, "amount is lower than limit")

	from := p.sender()
	// find
	selfData := p.getPledgeData(from)
	assert.AssertEx(selfData != nil, "not pledge yet!")
	assert.AssertEx(amount_val < selfData.Amount, "insufficient pledge")

	popChain := false
	if selfData.Amount >= PledgeLimit {
		popChain = true
	}
	selfData.Amount -= amount_val

	parentData := p.getPoolData(selfData.Parent)
	assert.AssertEx(parentData != nil, "parentData nil")

	popProducerChain := false
	if parentData.Amount >= PooolLimit {
		popProducerChain = true
	}
	parentData.Amount -= amount_val

	if selfData.Amount < PledgeLimit && popChain {
		prev := p.getPledgeData(selfData.Prev)
		next := p.getPledgeData(selfData.Next)
		if prev != nil {
			prev.Next = selfData.Next
			p.savePledgeData(selfData.Prev, prev)
		}
		if next != nil {
			next.Prev = selfData.Prev
			p.savePledgeData(selfData.Next, next)
		}
		if parentData.ChildHead == from {
			parentData.ChildHead = selfData.Next
		}
		selfData.Next = nilAddr
		selfData.Prev = nilAddr
	}

	producerData := p.getProducerData(parentData.Producer)
	if producerData != nil {
		oriWeight := p.getWeight(producerData.Amount)

		producerData.Amount -= amount_val
		newWeight := p.getWeight(producerData.Amount)

		if parentData.Amount < PooolLimit && popProducerChain {
			prev := p.getPoolData(parentData.Prev)
			next := p.getPoolData(parentData.Next)
			if prev != nil {
				prev.Next = parentData.Next
				p.savePoolData(parentData.Prev, prev)
			}
			if next != nil {
				next.Prev = parentData.Prev
				p.savePoolData(parentData.Next, next)
			}
			if producerData.ChildHead == selfData.Parent {
				producerData.ChildHead = parentData.Next
			}
			parentData.Next = nilAddr
			parentData.Prev = nilAddr
		}

		p.saveProducerData(parentData.Producer, producerData)
		p.subTotalPledgeWeight(oriWeight - newWeight)
	}

	p.savePledgeData(from, selfData)
	p.savePoolData(selfData.Parent, parentData)

	p.subTotalPledgeVal(amount_val)
	act := &transaction.Action{p.conAddr(), []byte{}}
	ctx := actioncontext.NewContext(p.conAddr(), act, p.ctx.GetBlockContext())
	conP := pledge{ctx}
	conP.transfer(p.conAddr().HexLower(), amount, "pledge")
}

func (p *pledge) proxy(beneficiary string) {
	assert.AssertEx(types.IsHexAddress(beneficiary), "to addr invalid")
	producer := types.HexToAddress(beneficiary)
	nilAddress := types.Address{}
	from := p.sender()
	// find
	selfData := p.getPoolData(from)
	assert.AssertEx(selfData != nil, "not have pd data")
	assert.AssertEx(selfData.Amount >= PooolLimit, "not enough pd balance")
	assert.AssertEx(selfData.Producer == nilAddress, "already proxy")
	assert.AssertEx(selfData.Producer != producer, "duplicate proxy")
	selfData.Producer = producer

	producerData := p.getProducerData(producer)
	assert.AssertEx(producerData != nil, "producer invalid 0")
	assert.AssertEx(producerData.ProducerCert >= ProducerLimit, "producer invalid")

	oriWeight := p.getWeight(producerData.Amount)
	producerData.Amount += selfData.Amount
	newWeight := p.getWeight(producerData.Amount)

	//insert chain
	if producerData.ChildHead == nilAddress {
		producerData.ChildHead = from
	} else {
		lastHead := p.getPoolData(producerData.ChildHead)
		lastHead.Prev = from
		p.savePoolData(producerData.ChildHead, lastHead)
		selfData.Next = producerData.ChildHead
		producerData.ChildHead = from
	}

	p.savePoolData(from, selfData)
	p.saveProducerData(producer, producerData)
	p.addTotalPledgeWeight(newWeight - oriWeight)
	p.addTotalPledgeVal(selfData.Amount)
}

func (p *pledge) makeProducer(amount string) {
	amount_int, err := strconv.ParseInt(amount, 10, 64)
	amount_val := uint64(amount_int)
	assert.AsserErr(err)
	assert.AssertEx(amount_val >= ProducerLimit, "makeProducer amount is lower than 500")

	from := p.sender()
	// find
	selfData := p.getProducerData(from)
	if selfData == nil {
		selfData = &ProducerData{
			ProducerCert: amount_val,
		}
	} else {
		selfData.ProducerCert += amount_val
	}
	p.saveProducerData(from, selfData)
	p.transfer(p.conAddr().HexLower(), amount, "makeProducer")
}

func (p *pledge) cancelProxy() {
	nilAddr := types.Address{}
	from := p.sender()
	// find
	selfData := p.getPoolData(from)
	assert.AssertEx(selfData != nil, "not pledge data yet!")
	assert.AssertEx(selfData.Producer != nilAddr, "not have proxy yet")

	popChain := false
	if selfData.Amount >= PooolLimit {
		popChain = true
	}

	producer := selfData.Producer

	producerData := p.getProducerData(producer)
	assert.AssertEx(producerData != nil, "producerData nil")
	oriWeight := p.getWeight(producerData.Amount)
	producerData.Amount -= selfData.Amount
	newWeight := p.getWeight(producerData.Amount)
	p.subTotalPledgeWeight(oriWeight - newWeight)
	p.subTotalPledgeVal(selfData.Amount)

	if popChain {
		prev := p.getPoolData(selfData.Prev)
		next := p.getPoolData(selfData.Next)
		if prev != nil {
			prev.Next = selfData.Next
			p.savePoolData(selfData.Prev, prev)
		}
		if next != nil {
			next.Prev = selfData.Prev
			p.savePoolData(selfData.Next, next)
		}
		if producerData.ChildHead == from {
			producerData.ChildHead = selfData.Next
		}
		selfData.Next = nilAddr
		selfData.Prev = nilAddr
	}
	selfData.Producer = nilAddr
	p.savePoolData(from, selfData)
	p.saveProducerData(producer, producerData)
}

func (p *pledge) balanceOf(addr string) {
	assert.AssertEx(types.IsHexAddress(addr), "to addr invalid")
	toAddr := types.HexToAddress(addr)
	val := p.dbGet(toAddr)
	valU := binary.LittleEndian.Uint64(val)
	ret := strconv.FormatUint(valU, 10)
	p.setResult(ret)
}

func (p *pledge) pledgeOf(addr string) {
	assert.AssertEx(types.IsHexAddress(addr), "to addr invalid")
	toAddr := types.HexToAddress(addr)
	data := p.getPledgeData(toAddr)
	val := uint64(0)
	if data != nil {
		val = data.Amount
	}
	ret := strconv.FormatUint(val, 10)
	p.setResult(ret)
}

func (p *pledge) pledgeOfExt(addr string) {
	assert.AssertEx(types.IsHexAddress(addr), "to addr invalid")
	toAddr := types.HexToAddress(addr)
	val := p.dbGet(toAddr)
	valU := binary.LittleEndian.Uint64(val)
	ret := strconv.FormatUint(valU, 10)
	p.setResult(ret)

	valPd := uint64(0)
	valPool := uint64(0)
	valProduce := uint64(0)
	data := p.getPledgeData(toAddr)
	if data != nil {
		valPd = data.Amount
	}
	data1 := p.getPoolData(toAddr)
	if data1 != nil {
		valPool = data1.Amount
	}
	data2 := p.getProducerData(toAddr)
	if data2 != nil {
		valProduce = data2.Amount
	}
	ret = strconv.FormatUint(valPd, 10)
	p.setResult(ret)
	ret = strconv.FormatUint(valPool, 10)
	p.setResult(ret)
	ret = strconv.FormatUint(valProduce, 10)
	p.setResult(ret)

	tt := p.getTotalAmount()
	retTt := strconv.FormatUint(tt, 10)
	p.setResult(retTt)

	tw := p.getTotalWeight()
	retTw := strconv.FormatUint(tw, 10)
	p.setResult(retTw)
}

func printDumpDataPd(amount uint64, pdData *PledgeData)  {
	dd := &DumpData{
		Amount: amount,
		PledgeData: pdData,
	}
	jsonData, _ := json.Marshal(dd)
	fmt.Println(string(jsonData))
}

func printDumpDataPool(amount uint64, pdData *PoolData)  {
	dd := &DumpData{
		Amount: amount,
		PoolData: pdData,
	}
	jsonData, _ := json.Marshal(dd)
	fmt.Println(string(jsonData))
}

func (p *pledge) dumpPool(addr string) {
	assert.AssertEx(types.IsHexAddress(addr), "to addr invalid")
	address := types.HexToAddress(addr)
	selfData := p.getPoolData(address)
	if selfData == nil {
		return
	}
	fmt.Println(addr, selfData.Amount, selfData.ChildHead.HexLower())
	head := selfData.ChildHead
	for {
		if head == nilAddr {
			break
		}
		hData := p.getPledgeData(head)
		val := p.dbGet(head)
		valX := binary.LittleEndian.Uint64(val)
		printDumpDataPd(valX, hData)
		head = hData.Next
	}
}

func (p *pledge) dumpProducer(addr string) {
	assert.AssertEx(types.IsHexAddress(addr), "to addr invalid")
	address := types.HexToAddress(addr)
	selfData := p.getProducerData(address)
	if selfData == nil {
		return
	}
	fmt.Println(addr, selfData.Amount, selfData.ChildHead.HexLower())
	head := selfData.ChildHead
	for {
		if head == nilAddr {
			break
		}
		hData := p.getPoolData(head)
		val := p.dbGet(head)
		valX := binary.LittleEndian.Uint64(val)
		printDumpDataPool(valX, hData)
		head = hData.Next
	}
}

