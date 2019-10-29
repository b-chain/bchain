
var ctxapi = ctxApi();

;(function () {
    function strToLower(param){
        return param.toLowerCase()
    }
    var balancePrefix = "balance"
    function BCHAINToken(){
        var symbol = "JOY"
        var tokenName = "BCHAIN Token"
        var decimals = 18
        var totalSupply = new BigNumber(1e+10)
        var _totalSupply = totalSupply.times(1e+18)
        
        var owner = strToLower(ctxapi._creator)
        ctxapi.console.print("---BCHAINToken Owner:" +owner + "\n")

        ctxapi.db.set(owner + balancePrefix , _totalSupply.toString())
        ctxapi.db.set("symbol" , symbol)
        ctxapi.db.set("tokenName" , tokenName)
        ctxapi.db.set("decimals" , decimals)
        ctxapi.db.set("totalSupply" , totalSupply)
        ctxapi.console.print("=================BCHAINToken init completed!!!!!!!!!!!!\n")
    }
    function construct(){
        keyInit = "init_flag"
        initFlag = ctxapi.db.get(keyInit)
        if  (initFlag.length != 0){
            return
        }
        BCHAINToken()
        ctxapi.db.set(keyInit, "OK")
    }
    construct()

    function basicBalanceTransfer(to1 , amount){
        ctxapi.console.print("Bchain Transfer sender:" + ctxapi._sender + "   Creator:" + ctxapi._creator + "\n")
        if (ctxapi._sender == ctxapi._creator){
            //call balanceReward before
            ret = ctxapi.cache.get("flag_consensus_only_once")
            if (ret.length != 0) {
                ctxapi.assert.assert(false , "invalid sender,code:"+ret)
            }
        }
        var from = strToLower(ctxapi._sender)
        var to = strToLower(to1)
        ctxapi.console.print("---BalanceTransfer From:" + from + "\n")
        ctxapi.console.print("---BalanceTransfer To:" + to + "\n")


        // ctxapi.console.print("BalanceTransfer amount " + amount.toString()+"\n")
        var aBalanceHex = ctxapi.db.get(from + balancePrefix)
        var aBalance = new BigNumber(0)
        if  (aBalanceHex.length != 0){
            aBalance = new BigNumber(aBalanceHex)
        }
        //ctxapi.console.print(from+"  before aBalance=" + aBalance +"\n")
        var bBalanceHex = ctxapi.db.get(to + balancePrefix)
        var bBalance = new BigNumber(0)
        if (bBalanceHex.length != 0) {
            bBalance = new BigNumber(bBalanceHex)
        }
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
        ctxapi.contract.emitEvent(from, true, to, true, amount.toString(), true)

    }

    this.balanceQuery = function(address){
        var _address = strToLower(address)
        ctxapi.console.print("=================balanceQuery====================\n")
        var aBalanceString = ctxapi.db.get(_address + balancePrefix)
        var aBalance = new BigNumber(0)
        if  (aBalanceString.length != 0){
            aBalance = new BigNumber(aBalanceString)
        }
        ctxapi.console.print(address+" balance= " + aBalance + "\n")
        ctxapi.result.setResult(aBalance.toString())
    },

    this.balanceTransfer = function(to1, amount){

        var check = new BigNumber(amount)
        if (check.lte(0)) {
            ctxapi.assert.assert(false , "balanceTransfer amount <= 0")
        }

        ctxapi.console.print("++++++ balanceTransfer Start +++++++\n")
        basicBalanceTransfer(to1 , amount)
        ctxapi.console.print("++++++ balanceTransfer End +++++++\n")
    }
})(this);

