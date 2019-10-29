package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"log"
	"bchain.io/common/types"
	"bchain.io/tool/tps_test/tpsTest"
	"os"
	"os/signal"
	"time"
)

func main() {
	key, addr := tpsTest.GetPriKey()
	nc := tpsTest.GetAccountNonce()
	fmt.Println(nc)

	stopChan := make(chan interface{})
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt)
		defer signal.Stop(sigc)
		<-sigc
		log.Println("Got interrupt,shutting down...")
		close(stopChan)
	}()
	tps := tpsTest.GetConfig().Tps
	cnt := 0
	timer := time.NewTicker(time.Second)
	for stop := false; !stop; {
		select {
		case <-timer.C:
			fmt.Println("count", cnt)
			senTx(tps, addr, key, nc)
			nc += uint64(tps)
			cnt++
		case <-stopChan:
			stop = true
			break
		}
	}
	timer.Stop()
	fmt.Println("exit!")

}

func sendRawTransaction(addr types.Address, key *ecdsa.PrivateKey, nc uint64) {
	tx := tpsTest.MakeTransaction(addr, key, nc)
	var encData bytes.Buffer
	err := msgp.Encode(&encData, tx)
	if err != nil {
		panic(err)
	}
	tpsTest.TxPost("bchain_sendRawTransaction", types.BytesForJson(encData.Bytes()))
	//fmt.Println(string(rsp))
}

func senTx(c int, addr types.Address, key *ecdsa.PrivateKey, nc uint64) {
	tnc := nc
	for i := 0; i < c; i++ {
		sendRawTransaction(addr, key, tnc)
		tnc++
	}
}
