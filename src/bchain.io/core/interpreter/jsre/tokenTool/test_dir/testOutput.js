
var ctxapi = ctxApi();


function getName(){
    return "ABC Token"
}

function getVersion(){
    return "V0.0.1"
}

;(function () {

    function strToLower(param){
        return param.toLowerCase()
    }

    var balancePrefix = "balance"
    var ratio = 10	// 1 bchain exchange ratio token

    function InitToken(){

        var symbol = "ABC"
        var tokenName = "ABC Token"
        var decimals = 18
        var totalSupply = new BigNumber(1e+10)
        var _totalSupply = totalSupply.times(1e+18)

        var owner = strToLower(ctxapi._creator)
        ctxapi.console.print(tokenName + "Owner:" +owner + "\n")

        ctxapi.db.set(owner + balancePrefix , _totalSupply.toString())
        ctxapi.db.set("symbol" , symbol)
        ctxapi.db.set("tokenName" , tokenName)
        ctxapi.db.set("decimals" , decimals)
        ctxapi.db.set("totalSupply" , totalSupply)
        ctxapi.console.print("=================Token init completed!!!!!!!!!!!!\n")
    }

    (function(){
        keyInit = "init_flag"
        initFlag = ctxapi.db.get(keyInit)
        if  (initFlag.length !== 0){
            return
        }
        InitToken()
        ctxapi.db.set(keyInit, "OK")
    })();

    this.balanceQuery = function(address){
        var _address = strToLower(address)
        ctxapi.console.print("================= ABC Query====================\n")
        var aBalanceString = ctxapi.db.get(_address + balancePrefix)
        var aBalance = new BigNumber(aBalanceString.length > 0 ? aBalanceString : 0)
        ctxapi.console.print(address+" balance= " + aBalance + "\n")
        ctxapi.result.setResult(aBalance.toString())
    },

   basicBalanceTransfer = function(from1,to1, amount){
       ctxapi.console.print("======================== ABC basicBalanceTransfer================================\n")
       if (amount <= 0){
           ctxapi.assert.assert(false , "ABC basicBalanceTransfer Amount <= 0")
       }
       var from = strToLower(from1)
       var to = strToLower(to1)

       if (from === to){
           ctxapi.assert.assert(false , "EthTransfer From == to:"+from)
       }

       ctxapi.console.print("---BalanceTransfer From:" + from + "\n")
       ctxapi.console.print("---BalanceTransfer To:" + to + "\n")

       // ctxapi.console.print("BalanceTransfer amount " + amount.toString()+"\n")
       var aBalanceHex = ctxapi.db.get(from + balancePrefix)
       var aBalance = new BigNumber(aBalanceHex.length > 0 ? aBalanceHex : 0)

       //ctxapi.console.print(from+"  before aBalance=" + aBalance +"\n")
       var bBalanceHex = ctxapi.db.get(to + balancePrefix)
       var bBalance = new BigNumber(bBalanceHex.length > 0 ? bBalanceHex : 0)

       //ctxapi.console.print(to+"  before bBalance=" + bBalance +"\n")

       if (aBalance.lt(amount)) {
           ctxapi.assert.assert(false , "balance is not enough")
       }

       aBalance = aBalance.minus(amount)
       bBalance = bBalance.plus(amount)

       aBalanceStr = aBalance.toString()
       bBalanceStr = bBalance.toString()
       //ctxapi.console.print(from+"  after aBalance=" + aBalance.toString()+"\n")
       //ctxapi.console.print(to+"  after bBalance=" + bBalance.toString()+"\n")
       ctxapi.db.set(from + balancePrefix , aBalanceStr)
       ctxapi.db.set(to + balancePrefix , bBalanceStr)
   },

    this.balanceTransfer = function(to1, amount){
        ctxapi.console.print("======================== ABC Transfer================================\n")
        basicBalanceTransfer(ctxapi._sender,to1,amount)
    },

    this.exchange = function(Amount){
        if (Amount <= 0){
            ctxapi.assert.assert(false , "ABC exchange Amount <= 0")
        }
        //if sender == creator,duihuan do a Recharge operation
        // exchange
        if (ctxapi._creator !== ctxapi._sender){

            ctxapi.console.print("======++++================ABC  Exchange================================\n")
            basicBalanceTransfer(ctxapi._creator , ctxapi._sender , Amount)
        }

        //Bchain exchange
        var bchainAmount = new BigNumber(Amount)
        bchainAmount = bchainAmount.dividedToIntegerBy(ratio)

        //the bchain receiver is contract
        to = strToLower(ctxapi._contract)
        var exchangeFunc = "balanceTransfer('" + to + "'," + bchainAmount.toString() + ")"
        ctxapi.call.call("0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B", exchangeFunc)
    },

    this.redeem = function(Amount){
        ctxapi.console.print("========================= ABC Redeem ================================\n")
        if (Amount <= 0){
            ctxapi.assert.assert(false , "ABC redeem Amount <= 0")
        }

        if (ctxapi._creator !== ctxapi._sender){
            //exchange
            this.balanceTransfer(ctxapi._creator , Amount)
        }

        //Bchain exchange
        var bchainAmount = new BigNumber(Amount)
        bchainAmount = bchainAmount.dividedToIntegerBy(ratio)
        //the bchain receiver is contract
        var to = strToLower(ctxapi._sender)
        var exchangeFunc = "balanceTransfer('" + to + "'," + bchainAmount.toString() + ")"
        ctxapi.call.innerCall("0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B", exchangeFunc)
    }

})(this);

