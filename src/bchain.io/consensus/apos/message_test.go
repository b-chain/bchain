package apos
/*
import (
	"fmt"
	"math/big"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/utils/crypto"
	"testing"
	"time"
)

//func (s *Signature) init() {
//	s.R = new(types.BigInt)
//	s.S = new(types.BigInt)
//	s.V = new(types.BigInt)
//}
// End condition 0 for message bp bba
func TestBba_EndCondition0(t *testing.T) {
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	Config().prLeader = 10000000000
	Config().prVerifier = 10000000000
	Config().maxPotLeaders = big.NewInt(3)
	Config().maxPotVerifiers = big.NewInt(4)
	an := newAllNodeManager()
	verifierCnt := an.initTestCommonNew(0)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX, "Verifier Cnt:", verifierCnt, COLOR_SHORT_RESET)
	Config().maxPotLeaders = big.NewInt(3)
	Config().maxPotVerifiers = big.NewInt(4)

	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 1
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}
	msgcs := NewMsgCredential(cs)
	msgcs.Send()

	bp := newBlockProposal()
	bp.Credential = cs
	bp.Block = an.actualNode.makeEmptyBlockForTest(bp.Credential)
	fmt.Println(bp.Block)
	hash := bp.Block.Hash()

	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = bp.Credential.Step
	bp.Esig.val = hash.Bytes()
	bp.Esig.Signature.init()
	if _, _, _, err := bp.Esig.sign(priKey); err != nil {
		fmt.Println("2222", err)
	}

	msgbp := NewMsgBlockProposal(bp)
	msgbp.Send()

	msg_css := []*msgCredentialSig{}
	msg_bbas := []*msgBinaryByzantineAgreement{}

	for i := 1; i <= 4; i++ {
		//time.Sleep(1 * time.Second)
		priKey := generatePrivateKey()
		cs := &CredentialSign{}
		cs.Round = 100
		cs.Step = 4 + 3
		cs.Signature.init()
		if _, _, _, err := cs.sign(priKey); err != nil {
			fmt.Println("333", err)
			return
		}
		msgcs := NewMsgCredential(cs)
		msg_css = append(msg_css, msgcs)
		bba := newBinaryByzantineAgreement()

		bba.Credential = cs
		bba.B = 0
		bba.Hash = hash
		//b
		bba.EsigB.round = bba.Credential.Round
		bba.EsigB.step = bba.Credential.Step
		bba.EsigB.val = big.NewInt(int64(bba.B)).Bytes()
		bba.EsigB.Signature.init()
		bba.EsigB.sign(priKey)

		//hash
		bba.EsigV.round = bba.Credential.Round
		bba.EsigV.step = bba.Credential.Step
		bba.EsigV.val = hash.Bytes()
		bba.EsigV.Signature.init()
		bba.EsigV.sign(priKey)

		msgBba := NewMsgBinaryByzantineAgreement(bba)
		msg_bbas = append(msg_bbas, msgBba)
	}

	for _, mcs := range msg_css {
		mcs.Send()
	}

	for _, mbba := range msg_bbas {
		time.Sleep(1 * time.Second)
		mbba.Send()
	}

	select {
	case <-an.actualNode.StopCh():
	}
}

func TestBba_EndConditionM3(t *testing.T) {
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	Config().prLeader = 10000000000
	Config().prVerifier = 10000000000
	Config().maxPotLeaders = big.NewInt(3)
	Config().maxPotVerifiers = big.NewInt(4)
	an := newAllNodeManager()
	verifierCnt := an.initTestCommonNew(0)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX, "Verifier Cnt:", verifierCnt, COLOR_SHORT_RESET)
	Config().maxPotLeaders = big.NewInt(3)
	Config().maxPotVerifiers = big.NewInt(2)

	msg_css := []*msgCredentialSig{}
	msg_bbas := []*msgBinaryByzantineAgreement{}

	for i := 1; i <= 2; i++ {
		//time.Sleep(1 * time.Second)
		priKey := generatePrivateKey()
		cs := &CredentialSign{}
		cs.Round = 100
		cs.Step = 12 + 3
		cs.Signature.init()
		if _, _, _, err := cs.sign(priKey); err != nil {
			fmt.Println("333", err)
			return
		}
		msgcs := NewMsgCredential(cs)
		msg_css = append(msg_css, msgcs)
		bba := newBinaryByzantineAgreement()

		bba.Credential = cs
		bba.B = 1
		hash := an.actualNode.roundCtx.getEmptyBlockHash()
		bba.Hash = hash
		//b
		bba.EsigB.round = bba.Credential.Round
		bba.EsigB.step = bba.Credential.Step
		bba.EsigB.val = big.NewInt(int64(bba.B)).Bytes()
		bba.EsigB.Signature.init()
		bba.EsigB.sign(priKey)

		//hash
		bba.EsigV.round = bba.Credential.Round
		bba.EsigV.step = bba.Credential.Step
		bba.EsigV.val = hash.Bytes()
		bba.EsigV.Signature.init()
		bba.EsigV.sign(priKey)

		msgBba := NewMsgBinaryByzantineAgreement(bba)
		msg_bbas = append(msg_bbas, msgBba)
	}

	for _, mcs := range msg_css {
		mcs.Send()
	}

	for _, mbba := range msg_bbas {
		time.Sleep(1 * time.Second)
		mbba.Send()
	}

	select {
	case <-an.actualNode.StopCh():
	}
}

func TestEvent1(t *testing.T) {
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	Config().prLeader = 10000000000
	Config().prVerifier = 10000000000
	Config().maxPotLeaders = big.NewInt(3)
	Config().maxPotVerifiers = big.NewInt(4)

	actualNode := NewApos(MsgTransfer(), newOutCommonTools())

	csCh := make(chan CsEvent, 1000)
	csSub := MsgTransfer().SubscribeCsEvent(csCh)
	fcs := func() {
		for {
			select {
			case event := <-csCh:
				logger.Info("Test: receive cs message", event.Cs.Round, event.Cs.Step)
				// Err() channel will be closed when unsubscribing.
			case <-csSub.Err():
				logger.Info("Test :receive cs stop message")
				return
			}
		}
	}

	bbach := make(chan BbaEvent, 1000)
	bbaSub := MsgTransfer().SubscribeBbaEvent(bbach)
	fbba := func() {
		for {
			select {
			case event := <-bbach:
				logger.Info("Test: receive bba message", event.Bba.Credential.Round, event.Bba.Credential.Step)
				// Err() channel will be closed when unsubscribing.
			case <-bbaSub.Err():
				logger.Info("Test :receive bba stop message")
				return
			}
		}
	}

	go actualNode.Run()
	go fcs()
	go fbba()

	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 1
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}
	msgcs := NewMsgCredential(cs)
	msgcs.Send()

	select {
	case <-actualNode.StopCh():
	}
}

func TestCs_validate_success(t *testing.T) {
	Config().prVerifier = 10000000000
	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 2
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}
	msgcs := NewMsgCredential(cs)
	msgcs.Send()
	time.Sleep(2 * time.Second)
}

//credential has no right to verify
func TestCs_validate_fail_1(t *testing.T) {
	Config().prVerifier = 1
	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 2
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}
	msgcs := NewMsgCredential(cs)
	msgcs.Send()
	time.Sleep(2 * time.Second)
}

//verify CredentialSig fail: invalid chain id for signer
func TestCs_validate_fail_2(t *testing.T) {
	Config().prVerifier = 10000000000
	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 2
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}
	cs.V.IntVal.Add(&cs.V.IntVal, common.Big2)
	msgcs := NewMsgCredential(cs)
	msgcs.Send()
	time.Sleep(2 * time.Second)
}

func TestCs_sava(t *testing.T) {
	Config().blockDelay = 20
	Config().verifyDelay = 10
	Config().maxBBASteps = 12
	Config().maxPotLeaders = big.NewInt(3)
	an := newAllNodeManager()
	verifierCnt := an.initTestCommonNew(0)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX, "Verifier Cnt:", verifierCnt, COLOR_SHORT_RESET)

	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 1
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	msgcs := NewMsgCredential(cs)
	msgcs.Send()
	msgcs.Send()

	for i := 1; i <= 10; i++ {
		time.Sleep(1 * time.Second)
		priKey := generatePrivateKey()
		cs := &CredentialSign{}
		cs.Round = 100
		cs.Step = 1
		cs.Signature.init()
		if _, _, _, err := cs.sign(priKey); err != nil {
			fmt.Println("111", err)
			return
		}

		msgcs := NewMsgCredential(cs)
		msgcs.Send()
		msgcs.Send()

	}

	select {
	case <-an.actualNode.StopCh():
	}
}

func TestBp_validate_success(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 1
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	apos := NewApos(nil, newOutCommonTools())
	bp := newBlockProposal()
	bp.Credential = cs
	bp.Block = apos.makeEmptyBlockForTest(bp.Credential)
	hash := bp.Block.Hash()

	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = bp.Credential.Step
	bp.Esig.val = hash.Bytes()
	bp.Esig.Signature.init()
	if _, _, _, err := bp.Esig.sign(priKey); err != nil {
		fmt.Println("2222", err)
	}

	msgbp := NewMsgBlockProposal(bp)
	msgbp.Send()
	time.Sleep(2 * time.Second)
}

//credential has no right to verify
func TestBp_validate_fail_1(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 1
	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 1
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	apos := NewApos(nil, newOutCommonTools())
	bp := newBlockProposal()
	bp.Credential = cs
	bp.Block = apos.makeEmptyBlockForTest(bp.Credential)
	hash := bp.Block.Hash()

	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = bp.Credential.Step
	bp.Esig.val = hash.Bytes()
	bp.Esig.Signature.init()
	if _, _, _, err := bp.Esig.sign(priKey); err != nil {
		fmt.Println("2222", err)
	}

	msgbp := NewMsgBlockProposal(bp)
	msgbp.Send()
	time.Sleep(2 * time.Second)
}

//Block Proposal step is not 1
func TestBp_validate_fail_2(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 1
	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 2
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	apos := NewApos(nil, newOutCommonTools())
	bp := newBlockProposal()
	bp.Credential = cs
	bp.Block = apos.makeEmptyBlockForTest(bp.Credential)
	hash := bp.Block.Hash()

	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = bp.Credential.Step
	bp.Esig.val = hash.Bytes()
	bp.Esig.Signature.init()
	if _, _, _, err := bp.Esig.sign(priKey); err != nil {
		fmt.Println("2222", err)
	}

	msgbp := NewMsgBlockProposal(bp)
	msgbp.Send()
	time.Sleep(2 * time.Second)
}

func TestBb_sava(t *testing.T) {
	Config().blockDelay = 20
	Config().verifyDelay = 10
	Config().maxBBASteps = 12
	Config().maxPotLeaders = big.NewInt(3)
	an := newAllNodeManager()
	verifierCnt := an.initTestCommonNew(0)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX, "Verifier Cnt:", verifierCnt, COLOR_SHORT_RESET)

	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 1
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	msgcs := NewMsgCredential(cs)
	msgcs.Send()

	bp := newBlockProposal()
	bp.Credential = cs
	bp.Block = an.actualNode.makeEmptyBlockForTest(bp.Credential)
	//fmt.Println(bp.Block)
	hash := bp.Block.Hash()

	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = bp.Credential.Step
	bp.Esig.val = hash.Bytes()
	bp.Esig.Signature.init()
	if _, _, _, err := bp.Esig.sign(priKey); err != nil {
		fmt.Println("2222", err)
	}

	//an.SendDataPackToActualNode(m1)
	msgbp := NewMsgBlockProposal(bp)
	msgbp.Send()

	for i := 1; i <= 10; i++ {
		time.Sleep(1 * time.Second)
		priKey := generatePrivateKey()
		cs := &CredentialSign{}
		cs.Round = 100
		cs.Step = 1
		cs.Signature.init()
		if _, _, _, err := cs.sign(priKey); err != nil {
			fmt.Println("111", err)
			return
		}

		msgcs := NewMsgCredential(cs)
		msgcs.Send()

		bp := newBlockProposal()
		bp.Credential = cs
		bp.Block = an.actualNode.makeEmptyBlockForTest(bp.Credential)
		//fmt.Println(bp.Block)
		hash := bp.Block.Hash()

		bp.Esig.round = bp.Credential.Round
		bp.Esig.step = bp.Credential.Step
		bp.Esig.val = hash.Bytes()
		bp.Esig.Signature.init()
		if _, _, _, err := bp.Esig.sign(priKey); err != nil {
			fmt.Println("2222", err)
		}

		//an.SendDataPackToActualNode(m1)
		msgbp := NewMsgBlockProposal(bp)
		msgbp.Send()
		msgbp.Send()
	}

	select {
	case <-an.actualNode.StopCh():
	}
}

//BP verify ephemeral signature fail: invalid chain id for signer
func TestBp_validate_fail_3(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 1
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	apos := NewApos(nil, newOutCommonTools())
	bp := newBlockProposal()
	bp.Credential = cs
	bp.Block = apos.makeEmptyBlockForTest(bp.Credential)
	hash := bp.Block.Hash()

	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = bp.Credential.Step
	bp.Esig.val = hash.Bytes()
	bp.Esig.Signature.init()
	if _, _, _, err := bp.Esig.sign(priKey); err != nil {
		fmt.Println("2222", err)
	}

	bp.Esig.Signature.V.IntVal.Add(&bp.Esig.Signature.V.IntVal, common.Big2)

	msgbp := NewMsgBlockProposal(bp)
	msgbp.Send()
	time.Sleep(2 * time.Second)
}

//sender's address is not equal in Credential and Ephemeral signature
func TestBp_validate_fail_4(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	Orignaddress := crypto.PubkeyToAddress(priKey.PublicKey)
	logger.Debug("Orignaddress", Orignaddress.Hex())
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 1
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	apos := NewApos(nil, newOutCommonTools())
	bp := newBlockProposal()
	bp.Credential = cs
	bp.Block = apos.makeEmptyBlockForTest(bp.Credential)
	hash := bp.Block.B_header.Hash()

	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = bp.Credential.Step
	bp.Esig.val = hash.Bytes()
	bp.Esig.Signature.init()
	if _, _, _, err := bp.Esig.sign(priKey); err != nil {
		fmt.Println("2222", err)
	}

	bp.Block.B_header.Time.IntVal.Add(&bp.Block.B_header.Time.IntVal, common.Big2)

	msgbp := NewMsgBlockProposal(bp)
	msgbp.Send()
	time.Sleep(2 * time.Second)
}

func TestGc_validate_success(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	Orignaddress := crypto.PubkeyToAddress(priKey.PublicKey)
	logger.Debug("Orignaddress", Orignaddress.Hex())
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 2
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	gc := newGradedConsensus()
	gc.Credential = cs
	hash := types.Hash{}
	//hash[1] = 1

	gc.Esig.round = gc.Credential.Round
	gc.Esig.step = gc.Credential.Step
	gc.Esig.val = hash.Bytes()
	gc.Esig.Signature.init()
	if _, _, _, err := gc.Esig.sign(priKey); err != nil {
		fmt.Println("2222", err)
	}

	msgGc := NewMsgGradedConsensus(gc)
	msgGc.Send()
	time.Sleep(2 * time.Second)
}

//message GradedConsensus validate error: Graded Consensus step is not 2 or 3
func TestGc_validate_fail_1(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 4
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	gc := newGradedConsensus()
	gc.Credential = cs
	hash := types.Hash{}
	hash[1] = 1

	gc.Esig.round = gc.Credential.Round
	gc.Esig.step = gc.Credential.Step
	gc.Esig.val = hash.Bytes()
	gc.Esig.Signature.init()
	if _, _, _, err := gc.Esig.sign(priKey); err != nil {
		fmt.Println("2222", err)
	}

	msgGc := NewMsgGradedConsensus(gc)
	msgGc.Send()
	time.Sleep(2 * time.Second)
}

//sender address is not equal
func TestGc_validate_fail_2(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	Orignaddress := crypto.PubkeyToAddress(priKey.PublicKey)
	logger.Debug("Orignaddress", Orignaddress.Hex())
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 2
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	gc := newGradedConsensus()
	gc.Credential = cs
	hash := types.Hash{}
	hash[1] = 1

	gc.Esig.round = gc.Credential.Round
	gc.Esig.step = gc.Credential.Step
	gc.Esig.val = hash.Bytes()
	gc.Esig.Signature.init()
	if _, _, _, err := gc.Esig.sign(priKey); err != nil {
		fmt.Println("2222", err)
	}

	msgGc := NewMsgGradedConsensus(gc)
	msgGc.Send()
	time.Sleep(2 * time.Second)
}

func TestBba_validate_success(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	Orignaddress := crypto.PubkeyToAddress(priKey.PublicKey)
	logger.Debug("Orignaddress", Orignaddress.Hex())
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 4 + 3
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	bba := newBinaryByzantineAgreement()

	bba.Credential = cs
	bba.B = 0
	bba.Hash = types.Hash{}
	bba.Hash[1] = 1
	//b
	bba.EsigB.round = bba.Credential.Round
	bba.EsigB.step = bba.Credential.Step
	bba.EsigB.val = big.NewInt(int64(bba.B)).Bytes()
	bba.EsigB.Signature.init()
	bba.EsigB.sign(priKey)

	//hash
	bba.EsigV.round = bba.Credential.Round
	bba.EsigV.step = bba.Credential.Step
	bba.EsigV.val = bba.Hash.Bytes()
	bba.EsigV.Signature.init()
	bba.EsigV.sign(priKey)

	msgBba := NewMsgBinaryByzantineAgreement(bba)
	msgBba.Send()
	time.Sleep(2 * time.Second)
}

//step is not right
func TestBba_validate_fail_1(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	Orignaddress := crypto.PubkeyToAddress(priKey.PublicKey)
	logger.Debug("Orignaddress", Orignaddress.Hex())
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 3
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	bba := newBinaryByzantineAgreement()

	bba.Credential = cs
	bba.B = 0
	bba.Hash = types.Hash{}
	bba.Hash[1] = 1
	//b
	bba.EsigB.round = bba.Credential.Round
	bba.EsigB.step = bba.Credential.Step
	bba.EsigB.val = big.NewInt(int64(bba.B)).Bytes()
	bba.EsigB.Signature.init()
	bba.EsigB.sign(priKey)

	//hash
	bba.EsigV.round = bba.Credential.Round
	bba.EsigV.step = bba.Credential.Step
	bba.EsigV.val = bba.Hash.Bytes()
	bba.EsigV.Signature.init()
	bba.EsigV.sign(priKey)

	msgBba := NewMsgBinaryByzantineAgreement(bba)
	msgBba.Send()
	time.Sleep(2 * time.Second)
}

//B value 2 is not right in apos protocal
func TestBba_validate_fail_2(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	Orignaddress := crypto.PubkeyToAddress(priKey.PublicKey)
	logger.Debug("Orignaddress", Orignaddress.Hex())
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 7
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	bba := newBinaryByzantineAgreement()

	bba.Credential = cs
	bba.B = 2
	bba.Hash = types.Hash{}
	bba.Hash[1] = 1
	//b
	bba.EsigB.round = bba.Credential.Round
	bba.EsigB.step = bba.Credential.Step
	bba.EsigB.val = big.NewInt(int64(bba.B)).Bytes()
	bba.EsigB.Signature.init()
	bba.EsigB.sign(priKey)

	//hash
	bba.EsigV.round = bba.Credential.Round
	bba.EsigV.step = bba.Credential.Step
	bba.EsigV.val = bba.Hash.Bytes()
	bba.EsigV.Signature.init()
	bba.EsigV.sign(priKey)

	msgBba := NewMsgBinaryByzantineAgreement(bba)
	msgBba.Send()
	time.Sleep(2 * time.Second)
}

//sender's address is not equal in Credential and B Ephemeral signature
func TestBba_validate_fail_3(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	Orignaddress := crypto.PubkeyToAddress(priKey.PublicKey)
	logger.Debug("Orignaddress", Orignaddress.Hex())
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 7
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	bba := newBinaryByzantineAgreement()

	bba.Credential = cs
	bba.B = 0
	bba.Hash = types.Hash{}
	bba.Hash[1] = 1
	//b
	bba.EsigB.round = bba.Credential.Round
	bba.EsigB.step = bba.Credential.Step
	bba.EsigB.val = big.NewInt(int64(bba.B)).Bytes()
	bba.EsigB.Signature.init()
	bba.EsigB.sign(priKey)

	bba.B = 1

	//hash
	bba.EsigV.round = bba.Credential.Round
	bba.EsigV.step = bba.Credential.Step
	bba.EsigV.val = bba.Hash.Bytes()
	bba.EsigV.Signature.init()
	bba.EsigV.sign(priKey)

	msgBba := NewMsgBinaryByzantineAgreement(bba)
	msgBba.Send()
	time.Sleep(2 * time.Second)
}

//sender's address is not equal in Credential and V Ephemeral
func TestBba_validate_fail_4(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	Orignaddress := crypto.PubkeyToAddress(priKey.PublicKey)
	logger.Debug("Orignaddress", Orignaddress.Hex())
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 7
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	bba := newBinaryByzantineAgreement()

	bba.Credential = cs
	bba.B = 0
	bba.Hash = types.Hash{}
	bba.Hash[1] = 1
	//b
	bba.EsigB.round = bba.Credential.Round
	bba.EsigB.step = bba.Credential.Step
	bba.EsigB.val = big.NewInt(int64(bba.B)).Bytes()
	bba.EsigB.Signature.init()
	bba.EsigB.sign(priKey)

	//hash
	bba.EsigV.round = bba.Credential.Round
	bba.EsigV.step = bba.Credential.Step
	bba.EsigV.val = bba.Hash.Bytes()
	bba.EsigV.Signature.init()
	bba.EsigV.sign(priKey)
	bba.Hash[1] = 2

	msgBba := NewMsgBinaryByzantineAgreement(bba)
	msgBba.Send()
	time.Sleep(2 * time.Second)
}

//bba m + 3 step message'b is not equal 1
func TestBba_validate_max(t *testing.T) {
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	priKey := generatePrivateKey()
	Orignaddress := crypto.PubkeyToAddress(priKey.PublicKey)
	logger.Debug("Orignaddress", Orignaddress.Hex())
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 180 + 3
	cs.Signature.init()
	if _, _, _, err := cs.sign(priKey); err != nil {
		fmt.Println("111", err)
		return
	}

	bba := newBinaryByzantineAgreement()

	bba.Credential = cs
	bba.B = 0
	bba.Hash = types.Hash{}
	bba.Hash[1] = 1
	//b
	bba.EsigB.round = bba.Credential.Round
	bba.EsigB.step = bba.Credential.Step
	bba.EsigB.val = big.NewInt(int64(bba.B)).Bytes()
	bba.EsigB.Signature.init()
	bba.EsigB.sign(priKey)

	//hash
	bba.EsigV.round = bba.Credential.Round
	bba.EsigV.step = bba.Credential.Step
	bba.EsigV.val = bba.Hash.Bytes()
	bba.EsigV.Signature.init()
	bba.EsigV.sign(priKey)

	msgBba := NewMsgBinaryByzantineAgreement(bba)
	msgBba.Send()
	time.Sleep(2 * time.Second)
}

func TestBba_e(t *testing.T) {
	a:= STEP_REDUCTION_1
	fmt.Println(a)
}
*/