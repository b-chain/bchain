package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"log"
	"errors"
	"bchain.io/tool/simpleWallet/simple"
	"strconv"
	"github.com/shopspring/decimal"
	"bchain.io/common/types"
	"bchain.io/common/assert"
	"bchain.io/utils/crypto"
	"bchain.io/communication/p2p/discover"
)

func createAccount(c *cli.Context) error {
	if c.NArg() == 1 {
		pwd := c.Args().First()
		fmt.Println("password: ", pwd)
		simple.AccountCreate(pwd)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func createAccountWithPriKey(c *cli.Context) error {
	if c.NArg() == 2 {
		pwd := c.Args().Get(1)
		fmt.Println("password: ", pwd)

		key := c.Args().Get(0)
		fmt.Println("key: ", pwd)
		simple.AccountCreateByKey(key, pwd)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func createPrivateKey(c *cli.Context) error {
	if c.NArg() == 0 {
		simple.NewKey()
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func blockproductor_start(c *cli.Context) error {
	if c.NArg() == 1 {
		pwd := c.Args().First()
		fmt.Println("password: ", pwd)
		simple.Blockproductor_start(url, pwd)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_transfer(c *cli.Context) error {
	if c.NArg() == 5 {
		pwd := c.Args().Get(4)
		key, addr, err := simple.GetPriKey(pwd)
		if err != nil {
			fmt.Println(err)
			return err
		}
		toAddr := c.Args().Get(0)
		amountStr := c.Args().Get(1)
		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		amount = amount.Round(8)
		memoStr := c.Args().Get(2)

		txFeeStr := c.Args().Get(3)
		txfee, err := decimal.NewFromString(txFeeStr)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("from", addr.HexLower(), "to", toAddr, "amount(BC)", amount, "memo", memoStr, "txfee(C)", txfee)
		deci := decimal.NewFromFloat(100000000)
		amount = amount.Mul(deci)
		amountin, err := strconv.ParseInt(amount.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}
		nc := simple.GetAccountNonce(url, addr.HexLower())
		tx := simple.MakeBcTransaction(addr, key, nc, toAddr, uint64(amountin), uint64(txfeein), memoStr)
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_transferPriKey(c *cli.Context) error {
	if c.NArg() == 5 {
		key_str := c.Args().Get(4)
		key, err := crypto.HexToECDSA(key_str)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		addr := crypto.PubkeyToAddress(key.PublicKey)
		toAddr := c.Args().Get(0)
		amountStr := c.Args().Get(1)
		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		amount = amount.Round(8)
		memoStr := c.Args().Get(2)

		txFeeStr := c.Args().Get(3)
		txfee, err := decimal.NewFromString(txFeeStr)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("from", addr.HexLower(), "to", toAddr, "amount(BC)", amount, "memo", memoStr, "txfee(C)", txfee)
		deci := decimal.NewFromFloat(100000000)
		amount = amount.Mul(deci)
		amountin, err := strconv.ParseInt(amount.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}
		nc := simple.GetAccountNonce(url, addr.HexLower())
		tx := simple.MakeBcTransaction(addr, key, nc, toAddr, uint64(amountin), uint64(txfeein), memoStr)
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_getNonce(c *cli.Context) error {
	if c.NArg() == 1 {
		addr := c.Args().Get(0)
		simple.GetAccountNonce(url, addr)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_transferWithNonce(c *cli.Context) error {
	if c.NArg() == 6 {
		pwd := c.Args().Get(5)
		key, addr, err := simple.GetPriKey(pwd)
		if err != nil {
			fmt.Println(err)
			return err
		}
		nc := new(types.Uint64ForJson)
		ncStr := c.Args().Get(0)
		err = nc.UnmarshalText([]byte(ncStr))
		assert.AsserErr(err)
		toAddr := c.Args().Get(1)
		amountStr := c.Args().Get(2)
		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		amount = amount.Round(8)
		memoStr := c.Args().Get(3)

		txFeeStr := c.Args().Get(4)
		txfee, err := decimal.NewFromString(txFeeStr)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("from", addr.HexLower(), "nonce", ncStr, "to", toAddr, "amount(BC)", amount, "memo", memoStr, "txfee(C)", txfee)
		deci := decimal.NewFromFloat(100000000)
		amount = amount.Mul(deci)
		amountin, err := strconv.ParseInt(amount.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}
		tx := simple.MakeBcTransaction(addr, key, uint64(*nc), toAddr, uint64(amountin), uint64(txfeein), memoStr)
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_transferWithNonceByKey(c *cli.Context) error {
	if c.NArg() == 6 {
		key_str := c.Args().Get(5)
		key, err := crypto.HexToECDSA(key_str)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		addr := crypto.PubkeyToAddress(key.PublicKey)

		nc := new(types.Uint64ForJson)
		ncStr := c.Args().Get(0)
		err = nc.UnmarshalText([]byte(ncStr))
		assert.AsserErr(err)
		toAddr := c.Args().Get(1)
		amountStr := c.Args().Get(2)
		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		amount = amount.Round(8)
		memoStr := c.Args().Get(3)

		txFeeStr := c.Args().Get(4)
		txfee, err := decimal.NewFromString(txFeeStr)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("from", addr.HexLower(), "nonce", ncStr, "to", toAddr, "amount(BC)", amount, "memo", memoStr, "txfee(C)", txfee)
		deci := decimal.NewFromFloat(100000000)
		amount = amount.Mul(deci)
		amountin, err := strconv.ParseInt(amount.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}
		tx := simple.MakeBcTransaction(addr, key, uint64(*nc), toAddr, uint64(amountin), uint64(txfeein), memoStr)
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_balenceOf(c *cli.Context) error {
	if c.NArg() == 1 {
		simple.ActionCallBalenceofBc(url, c.Args().First())
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}


func bc_pledge(c *cli.Context) error {
	if c.NArg() == 3 {
		pwd := c.Args().Get(2)
		key, addr, err := simple.GetPriKey(pwd)
		if err != nil {
			fmt.Println(err)
			return err
		}
		amountStr := c.Args().Get(0)
		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		amount = amount.Round(8)

		txFeeStr := c.Args().Get(1)
		txfee, err := decimal.NewFromString(txFeeStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("from", addr.HexLower(), "amount(BC)", amount, "txfee(C)", txfee)
		deci := decimal.NewFromFloat(100000000)
		amount = amount.Mul(deci)
		amountin, err := strconv.ParseInt(amount.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}
		nc := simple.GetAccountNonce(url, addr.HexLower())
		tx := simple.MakeBcPledgeTransaction(addr, key, nc, uint64(amountin), uint64(txfeein))
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_pledgeByKey(c *cli.Context) error {
	if c.NArg() == 3 {
		key_str := c.Args().Get(2)
		key, err := crypto.HexToECDSA(key_str)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		addr := crypto.PubkeyToAddress(key.PublicKey)

		amountStr := c.Args().Get(0)
		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		amount = amount.Round(8)

		txFeeStr := c.Args().Get(1)
		txfee, err := decimal.NewFromString(txFeeStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("from", addr.HexLower(), "amount(BC)", amount, "txfee(C)", txfee)
		deci := decimal.NewFromFloat(100000000)
		amount = amount.Mul(deci)
		amountin, err := strconv.ParseInt(amount.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}
		nc := simple.GetAccountNonce(url, addr.HexLower())
		tx := simple.MakeBcPledgeTransaction(addr, key, nc, uint64(amountin), uint64(txfeein))
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_redeem(c *cli.Context) error {
	if c.NArg() == 3 {
		pwd := c.Args().Get(2)
		key, addr, err := simple.GetPriKey(pwd)
		if err != nil {
			fmt.Println(err)
			return err
		}
		amountStr := c.Args().Get(0)
		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		amount = amount.Round(8)

		txFeeStr := c.Args().Get(1)
		txfee, err := decimal.NewFromString(txFeeStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("from", addr.HexLower(), "amount:BC", amount, "txfee:C", txfeein)
		deci := decimal.NewFromFloat(100000000)
		amount = amount.Mul(deci)
		amountin, err := strconv.ParseInt(amount.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}
		nc := simple.GetAccountNonce(url, addr.HexLower())
		tx := simple.MakeBcRedeemTransaction(addr, key, nc, uint64(amountin), uint64(txfeein))
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_redeemByKey(c *cli.Context) error {
	if c.NArg() == 3 {
		key_str := c.Args().Get(2)
		key, err := crypto.HexToECDSA(key_str)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		addr := crypto.PubkeyToAddress(key.PublicKey)
		amountStr := c.Args().Get(0)
		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		amount = amount.Round(8)

		txFeeStr := c.Args().Get(1)
		txfee, err := decimal.NewFromString(txFeeStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("from", addr.HexLower(), "amount:BC", amount, "txfee:C", txfeein)
		deci := decimal.NewFromFloat(100000000)
		amount = amount.Mul(deci)
		amountin, err := strconv.ParseInt(amount.String(), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}
		nc := simple.GetAccountNonce(url, addr.HexLower())
		tx := simple.MakeBcRedeemTransaction(addr, key, nc, uint64(amountin), uint64(txfeein))
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_pledgeOf(c *cli.Context) error {
	if c.NArg() == 1 {
		simple.ActionCallPledgeofBc(url, c.Args().First())
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_getBlkByNum(c *cli.Context) error {
	if c.NArg() == 1 {
		simple.GetBlockByNumer(url, c.Args().First())
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bc_getBlkCertByNum(c *cli.Context) error {
	if c.NArg() == 1 {
		simple.GetBlocCertificateByNumer(url, c.Args().First())
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func checkErr(err error)  {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
func bc_publishContract(c *cli.Context) error {
	if c.NArg() == 3 {
		pwd := c.Args().Get(2)
		key, addr, err := simple.GetPriKey(pwd)
		checkErr(err)

		path := c.Args().Get(0)
		txFeeStr := c.Args().Get(1)
		txfee, err := decimal.NewFromString(txFeeStr)
		checkErr(err)

		fmt.Println("creator", addr.HexLower(), "code path", path, "txfee(C)", txfee)
		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		checkErr(err)

		nc := simple.GetAccountNonce(url, addr.HexLower())
		ncSys := simple.GetAccountNonce(url, simple.SystemContract)
		tx := simple.MakeSystemTransaction(key, nc, uint64(txfeein), path)
		contractAddr := crypto.CreateAddress(addr, ncSys)
		fmt.Println("contract address", contractAddr.HexLower())
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		checkErr(err)
		return nil
	}
}

func bc_publishContractBykey(c *cli.Context) error {
	if c.NArg() == 3 {
		key_str := c.Args().Get(2)
		key, err := crypto.HexToECDSA(key_str)
		checkErr(err)
		addr := crypto.PubkeyToAddress(key.PublicKey)

		path := c.Args().Get(0)
		txFeeStr := c.Args().Get(1)
		txfee, err := decimal.NewFromString(txFeeStr)
		checkErr(err)

		fmt.Println("creator", addr.HexLower(), "code path", path, "txfee(C)", txfee)
		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		checkErr(err)

		nc := simple.GetAccountNonce(url, addr.HexLower())
		ncSys := simple.GetAccountNonce(url, simple.SystemContract)
		tx := simple.MakeSystemTransaction(key, nc, uint64(txfeein), path)
		contractAddr := crypto.CreateAddress(addr, ncSys)
		fmt.Println("contract address", contractAddr.HexLower())
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		checkErr(err)
		return nil
	}
}

// big token part
// {toAddr} {amount} {symbol} {memo} {expiry} {contract addr} {txFee(C)} {password}
func bigToken_transfer(c *cli.Context) error {
	if c.NArg() == 8 {
		pwd := c.Args().Get(7)
		key, addr, err := simple.GetPriKey(pwd)
		checkErr(err)
		toAddr := c.Args().Get(0)
		amountStr := c.Args().Get(1)
		amount, err := decimal.NewFromString(amountStr)
		checkErr(err)

		symbolStr := c.Args().Get(2)
		if len(symbolStr)>64 {
			fmt.Println(symbolStr, "too long")
		}
		memoStr := c.Args().Get(3)
		if len(memoStr)>64 {
			fmt.Println(memoStr, "too long")
		}

		blkNumber :=simple.GetBlockNumer(url)


		expityStr := c.Args().Get(4)
		_, err = decimal.NewFromString(expityStr)
		checkErr(err)

		contractStr := c.Args().Get(5)
		if !types.IsHexAddress(contractStr) {
			fmt.Println(contractStr, "is not valide cntract address")
			os.Exit(-2)
		}
		conAddr := types.HexToAddress(contractStr)

		txFeeStr := c.Args().Get(6)
		txfee, err := decimal.NewFromString(txFeeStr)
		checkErr(err)

		fmt.Println("from", addr.HexLower(), "to", toAddr, "amount", amount, "symbol", symbolStr, "memo", memoStr)
		fmt.Println("txfee(C)", txfee, "contract address", contractStr, "block number", blkNumber, "expiry", expityStr)

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		checkErr(err)



		expiryIn, err := strconv.Atoi(expityStr)
		checkErr(err)

		nc := simple.GetAccountNonce(url, addr.HexLower())
		bt := simple.Big_token_tr{
			Key: key,
			ConAddr:conAddr,

			Nc:nc,
			TxFee:uint64(txfeein),

			To:toAddr,
			Amount: amount.String(),
			Symbol:symbolStr,
			Memo:memoStr,

			BlkNumber: blkNumber,
			Expiry: uint32(expiryIn),
		}
		tx := bt.MakeTransaction()
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bigToken_transferWithNonce(c *cli.Context) error {
	if c.NArg() == 9 {
		pwd := c.Args().Get(8)
		key, addr, err := simple.GetPriKey(pwd)
		checkErr(err)

		nc := new(types.Uint64ForJson)
		ncStr := c.Args().Get(0)
		err = nc.UnmarshalText([]byte(ncStr))
		checkErr(err)

		toAddr := c.Args().Get(1)
		amountStr := c.Args().Get(2)
		amount, err := decimal.NewFromString(amountStr)
		checkErr(err)

		symbolStr := c.Args().Get(3)
		if len(symbolStr)>64 {
			fmt.Println(symbolStr, "too long")
		}
		memoStr := c.Args().Get(4)
		if len(memoStr)>64 {
			fmt.Println(memoStr, "too long")
		}

		blkNumber :=simple.GetBlockNumer(url)

		expityStr := c.Args().Get(5)
		_, err = decimal.NewFromString(expityStr)
		checkErr(err)

		contractStr := c.Args().Get(6)
		if !types.IsHexAddress(contractStr) {
			fmt.Println(contractStr, "is not valide cntract address")
			os.Exit(-2)
		}
		conAddr := types.HexToAddress(contractStr)

		txFeeStr := c.Args().Get(7)
		txfee, err := decimal.NewFromString(txFeeStr)
		checkErr(err)

		fmt.Println("from", addr.HexLower(), "to", toAddr, "amount", amount, "symbol", symbolStr, "memo", memoStr)
		fmt.Println("txfee(C)", txfee, "contract address", contractStr, "block number", blkNumber, "expiry", expityStr)

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		checkErr(err)

		expiryIn, err := strconv.Atoi(expityStr)
		checkErr(err)

		bt := simple.Big_token_tr{
			Key: key,
			ConAddr:conAddr,

			Nc:uint64(*nc),
			TxFee:uint64(txfeein),

			To:toAddr,
			Amount: amount.String(),
			Symbol:symbolStr,
			Memo:memoStr,

			BlkNumber: blkNumber,
			Expiry: uint32(expiryIn),
		}
		tx := bt.MakeTransaction()
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

// bt-create {symbol} {name} {decimals} {supply} {is issue} {expiry} {contract addr} {txFee(C)} {password}
func bigToken_create(c *cli.Context) error {
	if c.NArg() == 9 {
		pwd := c.Args().Get(8)
		key, addr, err := simple.GetPriKey(pwd)
		checkErr(err)


		symbolStr := c.Args().Get(0)
		if len(symbolStr)>64 {
			fmt.Println(symbolStr, "too long")
		}

		NameStr := c.Args().Get(1)
		if len(symbolStr)>64 {
			fmt.Println(NameStr, "too long")
		}

		decimalsStr := c.Args().Get(2)
		_, err = decimal.NewFromString(decimalsStr)
		checkErr(err)
		if len(symbolStr)>64 {
			fmt.Println(NameStr, "too long")
		}

		supplyStr := c.Args().Get(3)
		_, err = decimal.NewFromString(supplyStr)
		checkErr(err)

		isIssue := c.Args().Get(4)
		_, err = decimal.NewFromString(isIssue)
		checkErr(err)

		blkNumber := simple.GetBlockNumer(url)

		expityStr := c.Args().Get(5)
		_, err = decimal.NewFromString(expityStr)
		checkErr(err)

		contractStr := c.Args().Get(6)
		if !types.IsHexAddress(contractStr) {
			fmt.Println(contractStr, "is not valide cntract address")
			os.Exit(-2)
		}
		conAddr := types.HexToAddress(contractStr)

		txFeeStr := c.Args().Get(7)
		txfee, err := decimal.NewFromString(txFeeStr)
		checkErr(err)

		fmt.Println("creator", addr.HexLower(), "symbol", symbolStr, "name", NameStr, "decimals", decimalsStr, "supply", supplyStr, "is issue", isIssue)
		fmt.Println("txfee(C)", txfee, "contract address", contractStr, "block number", blkNumber, "expiry", expityStr)

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		checkErr(err)

		expiryIn, err := strconv.Atoi(expityStr)
		checkErr(err)

		issueIn, err := strconv.Atoi(isIssue)
		checkErr(err)

		nc := simple.GetAccountNonce(url, addr.HexLower())
		bt := simple.Big_token_create{
			Key: key,
			ConAddr:conAddr,

			Nc:nc,
			TxFee:uint64(txfeein),

			Symbol:symbolStr,
			Name:NameStr,
			Decimals:decimalsStr,
			Supply:supplyStr,
			IsIssue:uint32(issueIn),

			BlkNumber: blkNumber,
			Expiry: uint32(expiryIn),
		}
		tx := bt.MakeTransaction()
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

// bt-issue {symbol} {amount} {memo} {expiry} {contract addr} {txFee(C)} {password}
func bigToken_issue(c *cli.Context) error {
	if c.NArg() == 7 {
		pwd := c.Args().Get(6)
		key, addr, err := simple.GetPriKey(pwd)
		checkErr(err)


		symbolStr := c.Args().Get(0)
		if len(symbolStr)>64 {
			fmt.Println(symbolStr, "too long")
		}

		amountStr := c.Args().Get(1)
		amount, err := decimal.NewFromString(amountStr)
		checkErr(err)

		memoStr := c.Args().Get(2)
		if len(memoStr)>64 {
			fmt.Println(memoStr, "too long")
		}


		blkNumber := simple.GetBlockNumer(url)


		expityStr := c.Args().Get(3)
		_, err = decimal.NewFromString(expityStr)
		checkErr(err)

		contractStr := c.Args().Get(4)
		if !types.IsHexAddress(contractStr) {
			fmt.Println(contractStr, "is not valide cntract address")
			os.Exit(-2)
		}
		conAddr := types.HexToAddress(contractStr)

		txFeeStr := c.Args().Get(5)
		txfee, err := decimal.NewFromString(txFeeStr)
		checkErr(err)

		fmt.Println("creator", addr.HexLower(), "symbol", symbolStr, "amount", amountStr)
		fmt.Println("txfee(C)", txfee, "contract address", contractStr, "block number", blkNumber, "expiry", expityStr)

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		checkErr(err)


		expiryIn, err := strconv.Atoi(expityStr)
		checkErr(err)

		nc := simple.GetAccountNonce(url, addr.HexLower())
		bt := simple.Big_token_issue{
			Key: key,
			ConAddr:conAddr,

			Nc:nc,
			TxFee:uint64(txfeein),

			Symbol:symbolStr,
			Amount:amount.String(),
			Memo:memoStr,

			BlkNumber: blkNumber,
			Expiry: uint32(expiryIn),
		}
		tx := bt.MakeTransaction()
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

///
func bigToken_transferByKey(c *cli.Context) error {
	if c.NArg() == 8 {
		key_str := c.Args().Get(7)
		key, err := crypto.HexToECDSA(key_str)
		checkErr(err)
		addr := crypto.PubkeyToAddress(key.PublicKey)


		toAddr := c.Args().Get(0)
		amountStr := c.Args().Get(1)
		amount, err := decimal.NewFromString(amountStr)
		checkErr(err)

		symbolStr := c.Args().Get(2)
		if len(symbolStr)>64 {
			fmt.Println(symbolStr, "too long")
		}
		memoStr := c.Args().Get(3)
		if len(memoStr)>64 {
			fmt.Println(memoStr, "too long")
		}

		blkNumber := simple.GetBlockNumer(url)

		expityStr := c.Args().Get(4)
		_, err = decimal.NewFromString(expityStr)
		checkErr(err)

		contractStr := c.Args().Get(5)
		if !types.IsHexAddress(contractStr) {
			fmt.Println(contractStr, "is not valide cntract address")
			os.Exit(-2)
		}
		conAddr := types.HexToAddress(contractStr)

		txFeeStr := c.Args().Get(6)
		txfee, err := decimal.NewFromString(txFeeStr)
		checkErr(err)

		fmt.Println("from", addr.HexLower(), "to", toAddr, "amount", amount, "symbol", symbolStr, "memo", memoStr)
		fmt.Println("txfee(C)", txfee, "contract address", contractStr, "block number", blkNumber, "expiry", expityStr)

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		checkErr(err)

		expiryIn, err := strconv.Atoi(expityStr)
		checkErr(err)

		nc := simple.GetAccountNonce(url, addr.HexLower())
		bt := simple.Big_token_tr{
			Key: key,
			ConAddr:conAddr,

			Nc:nc,
			TxFee:uint64(txfeein),

			To:toAddr,
			Amount: amount.String(),
			Symbol:symbolStr,
			Memo:memoStr,

			BlkNumber: blkNumber,
			Expiry: uint32(expiryIn),
		}
		tx := bt.MakeTransaction()
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bigToken_transferWithNonceByKey(c *cli.Context) error {
	if c.NArg() == 9 {
		key_str := c.Args().Get(8)
		key, err := crypto.HexToECDSA(key_str)
		checkErr(err)
		addr := crypto.PubkeyToAddress(key.PublicKey)

		nc := new(types.Uint64ForJson)
		ncStr := c.Args().Get(0)
		err = nc.UnmarshalText([]byte(ncStr))
		checkErr(err)

		toAddr := c.Args().Get(1)
		amountStr := c.Args().Get(2)
		amount, err := decimal.NewFromString(amountStr)
		checkErr(err)

		symbolStr := c.Args().Get(3)
		if len(symbolStr)>64 {
			fmt.Println(symbolStr, "too long")
		}
		memoStr := c.Args().Get(4)
		if len(memoStr)>64 {
			fmt.Println(memoStr, "too long")
		}

		blkNumber := simple.GetBlockNumer(url)

		expityStr := c.Args().Get(5)
		_, err = decimal.NewFromString(expityStr)
		checkErr(err)

		contractStr := c.Args().Get(6)
		if !types.IsHexAddress(contractStr) {
			fmt.Println(contractStr, "is not valide cntract address")
			os.Exit(-2)
		}
		conAddr := types.HexToAddress(contractStr)

		txFeeStr := c.Args().Get(7)
		txfee, err := decimal.NewFromString(txFeeStr)
		checkErr(err)

		fmt.Println("from", addr.HexLower(), "to", toAddr, "amount", amount, "symbol", symbolStr, "memo", memoStr)
		fmt.Println("txfee(C)", txfee, "contract address", contractStr, "block number", blkNumber, "expiry", expityStr)

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		checkErr(err)

		expiryIn, err := strconv.Atoi(expityStr)
		checkErr(err)

		bt := simple.Big_token_tr{
			Key: key,
			ConAddr:conAddr,

			Nc:uint64(*nc),
			TxFee:uint64(txfeein),

			To:toAddr,
			Amount: amount.String(),
			Symbol:symbolStr,
			Memo:memoStr,

			BlkNumber: blkNumber,
			Expiry: uint32(expiryIn),
		}
		tx := bt.MakeTransaction()
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

// bt-create {symbol} {name} {decimals} {supply} {is issue} {expiry} {contract addr} {txFee(C)} {password}
func bigToken_createByKey(c *cli.Context) error {
	if c.NArg() == 9 {
		key_str := c.Args().Get(8)
		key, err := crypto.HexToECDSA(key_str)
		checkErr(err)
		addr := crypto.PubkeyToAddress(key.PublicKey)


		symbolStr := c.Args().Get(0)
		if len(symbolStr)>64 {
			fmt.Println(symbolStr, "too long")
		}

		NameStr := c.Args().Get(1)
		if len(symbolStr)>64 {
			fmt.Println(NameStr, "too long")
		}

		decimalsStr := c.Args().Get(2)
		_, err = decimal.NewFromString(decimalsStr)
		checkErr(err)
		if len(symbolStr)>64 {
			fmt.Println(NameStr, "too long")
		}

		supplyStr := c.Args().Get(3)
		_, err = decimal.NewFromString(supplyStr)
		checkErr(err)

		isIssue := c.Args().Get(4)
		_, err = decimal.NewFromString(isIssue)
		checkErr(err)
		blkNumber := simple.GetBlockNumer(url)

		expityStr := c.Args().Get(5)
		_, err = decimal.NewFromString(expityStr)
		checkErr(err)

		contractStr := c.Args().Get(6)
		if !types.IsHexAddress(contractStr) {
			fmt.Println(contractStr, "is not valide cntract address")
			os.Exit(-2)
		}
		conAddr := types.HexToAddress(contractStr)

		txFeeStr := c.Args().Get(7)
		txfee, err := decimal.NewFromString(txFeeStr)
		checkErr(err)

		fmt.Println("creator", addr.HexLower(), "symbol", symbolStr, "name", NameStr, "decimals", decimalsStr, "supply", supplyStr, "is issue", isIssue)
		fmt.Println("txfee(C)", txfee, "contract address", contractStr, "block number", blkNumber, "expiry", expityStr)

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		checkErr(err)

		expiryIn, err := strconv.Atoi(expityStr)
		checkErr(err)

		issueIn, err := strconv.Atoi(isIssue)
		checkErr(err)

		nc := simple.GetAccountNonce(url, addr.HexLower())
		bt := simple.Big_token_create{
			Key: key,
			ConAddr:conAddr,

			Nc:nc,
			TxFee:uint64(txfeein),

			Symbol:symbolStr,
			Name:NameStr,
			Decimals:decimalsStr,
			Supply:supplyStr,
			IsIssue: uint32(issueIn),

			BlkNumber: blkNumber,
			Expiry: uint32(expiryIn),
		}
		tx := bt.MakeTransaction()
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

// bt-issue {symbol} {amount} {memo} {expiry} {contract addr} {txFee(C)} {password}
func bigToken_issueByKey(c *cli.Context) error {
	if c.NArg() == 7 {
		key_str := c.Args().Get(6)
		key, err := crypto.HexToECDSA(key_str)
		checkErr(err)
		addr := crypto.PubkeyToAddress(key.PublicKey)


		symbolStr := c.Args().Get(0)
		if len(symbolStr)>64 {
			fmt.Println(symbolStr, "too long")
		}

		amountStr := c.Args().Get(1)
		amount, err := decimal.NewFromString(amountStr)
		checkErr(err)

		memoStr := c.Args().Get(2)
		if len(memoStr)>64 {
			fmt.Println(memoStr, "too long")
		}

		blkNumber := simple.GetBlockNumer(url)

		expityStr := c.Args().Get(3)
		_, err = decimal.NewFromString(expityStr)
		checkErr(err)

		contractStr := c.Args().Get(4)
		if !types.IsHexAddress(contractStr) {
			fmt.Println(contractStr, "is not valide cntract address")
			os.Exit(-2)
		}
		conAddr := types.HexToAddress(contractStr)

		txFeeStr := c.Args().Get(5)
		txfee, err := decimal.NewFromString(txFeeStr)
		checkErr(err)

		fmt.Println("creator", addr.HexLower(), "symbol", symbolStr, "amount", amountStr)
		fmt.Println("txfee(C)", txfee, "contract address", contractStr, "block number", blkNumber, "expiry", expityStr)

		txfeein, err := strconv.ParseInt(txfee.String(), 10, 64)
		checkErr(err)


		expiryIn, err := strconv.Atoi(expityStr)
		checkErr(err)

		nc := simple.GetAccountNonce(url, addr.HexLower())
		bt := simple.Big_token_issue{
			Key: key,
			ConAddr:conAddr,

			Nc:nc,
			TxFee:uint64(txfeein),

			Symbol:symbolStr,
			Amount:amount.String(),
			Memo:memoStr,

			BlkNumber: blkNumber,
			Expiry: uint32(expiryIn),
		}
		tx := bt.MakeTransaction()
		simple.SendRawTransaction(url, tx)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bt_balanceOf(c *cli.Context) error {
	if c.NArg() == 3 {
		addr := c.Args().Get(2)
		contractAddr := types.HexToAddress(addr)
		simple.ActionCallBalanceOfBt(url, c.Args().First(), c.Args().Get(1), contractAddr)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bt_getSupply(c *cli.Context) error {
	if c.NArg() == 2 {
		addr := c.Args().Get(1)
		contractAddr := types.HexToAddress(addr)
		simple.ActionCallGetSupplyBt(url, c.Args().First(), contractAddr)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bt_getDecimals(c *cli.Context) error {
	if c.NArg() == 2 {
		addr := c.Args().Get(1)
		contractAddr := types.HexToAddress(addr)
		simple.ActionCallGetDecimalsBt(url, c.Args().First(), contractAddr)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func bt_getName(c *cli.Context) error {
	if c.NArg() == 2 {
		addr := c.Args().Get(1)
		contractAddr := types.HexToAddress(addr)
		simple.ActionCallGetNameBt(url, c.Args().First(), contractAddr)
		return nil
	} else {
		err := errors.New("args is not match")
		fmt.Println(err.Error())
		return err
	}
}

func printNodeId(c *cli.Context) error {
	keyfile := "./bchaind_node/nodekey"
	priv, err := crypto.LoadECDSA(keyfile)
	checkErr(err)
	nodeId := discover.PubkeyID(&priv.PublicKey)
	fmt.Println(nodeId)
	return nil
}

var url = "http://127.0.0.1:8989/"
func main() {
	app := cli.NewApp()
	app.Name = "b chain simple wallet tool"
	app.Version = "1.0.3"
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:        "url, u",
			Value:       "http://127.0.0.1:8989/",
			Usage:       "node url",
			Destination: &url,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "get-nodeId",
			Aliases: []string{"get-id"},
			Usage:   "get node id. get-nodeId",
			Action:  printNodeId,
		},
		{
			Name:    "account-create",
			Aliases: []string{"acc-c"},
			Usage:   "new a key store file with password. account-create {password}",
			Action:  createAccount,
		},
		{
			Name:    "account-import",
			Aliases: []string{"acc-i"},
			Usage:   "import private key, account-import {key password}",
			Action:  createAccountWithPriKey,
		},
		{
			Name:    "start-producer",
			Aliases: []string{"s-p"},
			Usage:   "start producing block. start-producer {password}",
			Action:  blockproductor_start,
		},
		{
			Name:    "bc-transfer",
			Aliases: []string{"bc-tr"},
			Usage:   "bchain BC transfer. bc-transfer {toAddr} {amount} {memo} {txFee(C)} {password}",
			Action:  bc_transfer,
		},
		{
			Name:    "bc-transferByKey",
			Aliases: []string{"bc-tr-k"},
			Usage:   "bchain BC transfer. bc-transferByKey {toAddr} {amount} {memo} {txFee(C)} {private key}",
			Action:  bc_transferPriKey,
		},
		{
			Name:    "bc-nonceOf",
			Aliases: []string{"bc-nof"},
			Usage:   "bchain get account nonce. bc-nonceOf {Addr}",
			Action:  bc_getNonce,
		},
		{
			Name:    "bc-transferWithNonce",
			Aliases: []string{"bc-trn"},
			Usage:   "bchain BC transfer with nonce. bc-transferWithNonce {account nonce} {toAddr} {amount} {memo} {txFee(C)} {password}",
			Action:  bc_transferWithNonce,
		},
		{
			Name:    "bc-transferWithNonceByKey",
			Aliases: []string{"bc-trn-k"},
			Usage:   "bchain BC transfer with nonce. bc-transferWithNonceByKey {account nonce} {toAddr} {amount} {memo} {txFee(C)} {private key}",
			Action:  bc_transferWithNonceByKey,
		},
		{
			Name:    "bc-balanceOf",
			Aliases: []string{"bc-of"},
			Usage:   "bchain BC balance of an addr. bc-balanceOf {Addr}",
			Action:  bc_balenceOf,
		},
		{
			Name:    "bc-pledge",
			Aliases: []string{"bc-pd"},
			Usage:   "pledge bchain BC to pledge pool. bc-pledge {amount} {txFee(C)} {password}",
			Action:  bc_pledge,
		},
		{
			Name:    "bc-pledgeByKey",
			Aliases: []string{"bc-pd-k"},
			Usage:   "pledge bchain BC to pledge pool. bc-pledgeByKey {amount} {txFee(C)} {private key}",
			Action:  bc_pledgeByKey,
		},
		{
			Name:    "bc-redeem",
			Aliases: []string{"bc-rd"},
			Usage:   "redeem bchain BC from pledge pool. bc-redeem {amount} {txFee(C)} {password}",
			Action:  bc_redeem,
		},
		{
			Name:    "bc-redeemByKey",
			Aliases: []string{"bc-rd-k"},
			Usage:   "redeem bchain BC from pledge pool. bc-redeemByKey {amount} {txFee(C)} {private key}",
			Action:  bc_redeemByKey,
		},
		{
			Name:    "bc-pledgeOf",
			Aliases: []string{"bc-pdof"},
			Usage:   "bchain BC pledge pool pledge of an addr. bc-pledgeOf {Addr}",
			Action:  bc_pledgeOf,
		},
		{
			Name:    "bc-getBlockByNumber",
			Aliases: []string{"bc-getBlkByBum"},
			Usage:   "get block by number. bc-getBlockByNumber {number}",
			Action:  bc_getBlkByNum,
		},
		{
			Name:    "bc-getBlockCertificateByNumber",
			Aliases: []string{"bc-getBlkCert"},
			Usage:   "get block Certificate by number. bc-getBlockCertificateByNumber {number}",
			Action:  bc_getBlkCertByNum,
		},
		{
			Name:    "bc-publishContract",
			Aliases: []string{"bc-publish"},
			Usage:   "publish b chain contract. bc-publishContract {path} {txFee(C)} {password}",
			Action:  bc_publishContract,
		},
		{
			Name:    "bc-publishContractByKey",
			Aliases: []string{"bc-publish-k"},
			Usage:   "publish b chain contract. bc-publishContractByKey {path} {txFee(C)} {private key}",
			Action:  bc_publishContractBykey,
		},
		{
			Name:    "bt-transfer",
			Aliases: []string{"bt-tr"},
			Usage:   "big token transfer. bt-transfer {toAddr} {amount} {symbol} {memo} {expiry} {contract addr} {txFee(C)} {password}",
			Action:  bigToken_transfer,
		},
		{
			Name:    "bt-transferWithNonce",
			Aliases: []string{"bt-trn"},
			Usage:   "big token transfer. bt-transferWithNonce {account nonce}{toAddr} {amount} {symbol} {memo} {expiry} {contract addr} {txFee(C)} {password}",
			Action:  bigToken_transferWithNonce,
		},
		{
			Name:    "bt-create",
			Aliases: []string{"bt-cr"},
			Usage:   "big token create token. bt-create {symbol} {name} {decimals} {supply} {is issue} {expiry} {contract addr} {txFee(C)} {password}",
			Action:  bigToken_create,
		},
		{
			Name:    "bt-issue",
			Aliases: []string{"bt-is"},
			Usage:   "big token issue token. bt-issue {symbol} {amount} {memo} {expiry} {contract addr} {txFee(C)} {password}",
			Action:  bigToken_issue,
		},
		{
			Name:    "bt-transferByKey",
			Aliases: []string{"bt-tr-k"},
			Usage:   "big token transfer. bt-transferByKey {toAddr} {amount} {symbol} {memo} {expiry} {contract addr} {txFee(C)} {password}",
			Action:  bigToken_transferByKey,
		},
		{
			Name:    "bt-transferWithNonceByKey",
			Aliases: []string{"bt-trn-k"},
			Usage:   "big token transfer. bt-transferWithNonceByKey {account nonce} {toAddr} {amount} {symbol} {memo} {expiry} {contract addr} {txFee(C)} {password}",
			Action:  bigToken_transferWithNonceByKey,
		},
		{
			Name:    "bt-createByKey",
			Aliases: []string{"bt-cr-k"},
			Usage:   "big token create token. bt-createByKey {symbol} {name} {decimals} {supply} {is issue} {expiry} {contract addr} {txFee(C)} {password}",
			Action:  bigToken_createByKey,
		},
		{
			Name:    "bt-issueByKey",
			Aliases: []string{"bt-is-k"},
			Usage:   "big token issue token. bt-issueByKey {symbol} {amount} {memo} {expiry} {contract addr} {txFee(C)} {password}",
			Action:  bigToken_issueByKey,
		},
		{
			Name:    "bt-balanceOf",
			Aliases: []string{"bt-of"},
			Usage:   "big token balance of an addr. bc-balanceOf {Addr} {symbol} {contract addr}",
			Action:  bt_balanceOf,
		},
		{
			Name:    "bt-getSupply",
			Aliases: []string{"bt-s"},
			Usage:   "big token Supply. bc-getSupply {symbol} {contract addr}",
			Action:  bt_getSupply,
		},
		{
			Name:    "bt-getDecimals",
			Aliases: []string{"bt-d"},
			Usage:   "big token decimals. bc-getDecimals {symbol} {contract addr}",
			Action:  bt_getDecimals,
		},
		{
			Name:    "bt-getName",
			Aliases: []string{"bt-n"},
			Usage:   "big token name. bc-getSupply {symbol} {contract addr}",
			Action:  bt_getName,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}