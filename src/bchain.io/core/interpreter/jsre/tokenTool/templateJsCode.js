
var ctxapi = ctxApi();


function getName(){
    return "template_tokenName"
}

function getVersion(){
    return "V0.0.1"
}

;(function () {

    function strToLower(param){
        return param.toLowerCase()
    }

    var balancePrefix = "balance"
    var ratio = template_Ratio	// 1 bchain exchange ratio token

    function InitToken(){

        var symbol = "template_Symbol"
        var tokenName = "template_tokenName"
        var decimals = template_Decimals
        var totalSupply = new BigNumber(template_totalSupply)
        var _totalSupply = totalSupply.times(1e+template_Decimals)

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
        ctxapi.console.print("================= template_Symbol Query====================\n")
        var aBalanceString = ctxapi.db.get(_address + balancePrefix)
        var aBalance = new BigNumber(aBalanceString.length > 0 ? aBalanceString : 0)
        ctxapi.console.print(address+" balance= " + aBalance + "\n")
        ctxapi.result.setResult(aBalance.toString())
    },

   basicBalanceTransfer = function(from1,to1, amount){
       ctxapi.console.print("======================== template_Symbol basicBalanceTransfer================================\n")
       if (amount <= 0){
           ctxapi.assert.assert(false , "template_Symbol basicBalanceTransfer Amount <= 0")
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
        ctxapi.console.print("======================== template_Symbol Transfer================================\n")
        basicBalanceTransfer(ctxapi._sender,to1,amount)
    },

    this.exchange = function(Amount){
        if (Amount <= 0){
            ctxapi.assert.assert(false , "template_Symbol exchange Amount <= 0")
        }
        //if sender == creator,duihuan do a Recharge operation
        // exchange
        if (ctxapi._creator !== ctxapi._sender){

            ctxapi.console.print("======++++================template_Symbol  Exchange================================\n")
            basicBalanceTransfer(ctxapi._creator , ctxapi._sender , Amount)
        }

        //Bchain exchange
        var bchainAmount = new BigNumber(Amount)
        bchainAmount = bchainAmount.dividedToIntegerBy(ratio)

        //the bchain receiver is contract
        to = strToLower(ctxapi._contract)
        var exchangeFunc = "balanceTransfer('" + to + "'," + bchainAmount.toString() + ")"
        ctxapi.call.call("BchainContractAddress", exchangeFunc)
    },
    //redeem just a private function now
    redeem = function(Amount){
        ctxapi.console.print("========================= template_Symbol Redeem ================================\n")
        if (Amount <= 0){
            ctxapi.assert.assert(false , "template_Symbol redeem Amount <= 0")
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
        ctxapi.call.innerCall("BchainContractAddress", exchangeFunc)
    }

})(this);

