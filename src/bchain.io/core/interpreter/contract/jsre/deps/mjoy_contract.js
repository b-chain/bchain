

//var apijs = require('context_api') // not support

var ctxapi = ctxApi();


function PrintHello(){
    ctxapi.console.print("================================PrintHello======================================\n");
}

;(function () {
    //var decimals = 18
    //var totalSupply = new BigNumber(1e+9)
    blockRewordAdjustNumber = new BigNumber(6250000)
    initReword = new BigNumber(80e+18)
    q = new BigNumber(0.5)

    //private function
    function PrintPrivateInfo(){
        ctxapi.console.print("===========Print Private Info===========\n")
    }

    function getRewordByNumber(number){
        ctxapi.console.print("===========getRewordByNumber===========" + number + "\n")
        idx = number.dividedToIntegerBy(blockRewordAdjustNumber)
        reword = (initReword.times(q.pow(idx))).toFixed(0)
        return new BigNumber(reword)
    }

    function strToLower(param){
        return param.toLowerCase()
    }

    this.testLower = function(param1) {
        var param = strToLower(param1)
        ctxapi.console.print(param + "========================\n");
    },

        this.test = function(param) {
            ctxapi.console.print(param)
            ctxapi.console.print("\n")
            ctxapi.console.print("============================================================================\n");
            ctxapi.console.print("ctxapi._contract: " + ctxapi._contract + "\n");
            ctxapi.console.print("ctxapi._sender:   " + ctxapi._sender + "\n");
            ctxapi.console.print("============================================================================\n");
        },
        this.crypto = function(param) {
            ctxapi.console.print("raw param: " + param)
            ctxapi.console.print("ctxapi.cypto.sha1: " + ctxapi.crypto.sha1(param) + "\n")
            ctxapi.console.print("ctxapi.cypto.sha256: " + ctxapi.crypto.sha256(param) + "\n")
            ctxapi.console.print("ctxapi.cypto.sha512: " + ctxapi.crypto.sha512(param) + "\n")
            var msg="0x41b1a0649752af1b28b3dc29a1556eee781e4a4c3a1f7f53f90fa834de098c4d";
            var sig="0xd155e94305af7e07dd8c32873e5c03cb95c9e05960ef85be9c07f671da58c73718c19adc397a211aa9e87e519e2038c5a3b658618db335f74f800b8e0cfeef4401";
            ctxapi.console.print("msg data: " + msg + "\n")
            ctxapi.console.print("sig data: " + sig + "\n")
            ctxapi.console.print("ctxapi.cypto.recover pubkey:0x" + ctxapi.crypto.recover(msg, sig) + "\n")
        },
        this.auth = function(param) {
            ctxapi.console.print("======================================================\n")
            ctxapi.console.print("Auth raw param:" + param + "\n")
            var addr="970e8128ab834e8eac17ab8e3812f010678c970e";
            ctxapi.console.print("ctxapi.auth.isHexAddress: " + ctxapi.auth.isHexAddress(addr) + "\n")
            ctxapi.console.print("ctxapi.auth.requireAuth: " + ctxapi.auth.requireAuth(addr) + "\n")
            ctxapi.console.print("ctxapi.auth.isAccount: " + ctxapi.auth.isAccount(addr) + "\n")
            ctxapi.console.print("ctxapi.auth.isContract: " + ctxapi.auth.isContract(addr) + "\n")
        },
        this.contractCreate = function(creator, code) {
            ctxapi.console.print("======================================================\n")
            ctxapi.console.print("Contract Create Test: \n interpreter:" + interpreter + "\n creator:" + creator +"\n code:" + code +"\n")
            ctxapi.console.print("ctxapi.contract.create:" + ctxapi.contract.create(creator, code) + "\n")
        },
        this.balanceQuery = function(from1){
            var from = strToLower(from1)
            ctxapi.console.print("======++++================44444Queryr================================\n")
            var balancePrefix = "balance"
            var aBalanceString = ctxapi.db.get(from + balancePrefix)

            var aBalance = new BigNumber(0)
            if  (aBalanceString.length != 0){
                aBalance = new BigNumber(aBalanceString)
            }

            ctxapi.console.print(from+" balance= " + aBalance + "\n")
            ctxapi.result.setResult(aBalance.toString())
        },
        this.balanceTransfer = function(from1 , to1 , amount){
            var from = strToLower(from1)
            var to = strToLower(to1)
            ctxapi.console.print("======++++================222BalanceTransfer================================\n")
            // ctxapi.console.print("BalanceTransfer amount " + amount.toString()+"\n")
            var balancePrefix = "balance"
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
        },

        this.balanceReward = function(from1 ){
            var from = strToLower(from1)
            //ctxapi.console.print("======================111Reward================================\n")
            var balancePrefix = "balance"
            var aBalanceString = ctxapi.db.get(from + balancePrefix)
            var aBalance = new BigNumber(0)
            if  (aBalanceString.length != 0){
                aBalance = new BigNumber(aBalanceString)
            }
            //ctxapi.console.print(from+"  before aBalance=" + aBalance+"\n")
            blkNumber = new BigNumber(ctxapi._number)
            reword = getRewordByNumber(blkNumber)
            aBalance = aBalance.plus(reword)
            aBalanceStr = aBalance.toString()

            ctxapi.db.set(from + balancePrefix, aBalanceStr)

        },
        this.createContract = function(code){
            //ctxapi.console.print("======================createContract===========================\n")
            creator = ctxapi._sender
            contractAddress = ctxapi.contract.create(creator, code)
            ctxapi.console.print("ctxapi.contract.create:" + contractAddress + "\n")
            ctxapi.result.setResult(contractAddress)
        },
        this.callContract = function(contractAddress, para){
            ctxapi.console.print("======================callContract===========================\n")
            ctxapi.call.call(contractAddress, para)
        },
        this.pledgeAdd = function(from1) {
            var fro = strToLower(from1)
            ctxapi.console.print("======================333pledge================================\n")
            var pledgePrefix = "pledge"
            var aPledgeString = ctxapi.db.get(from + pledgePrefix)
            var aPledge = new BigNumber(0)
            if  (aPledgeString.length != 0){
                aPledge = new BigNumber(aPledgeString)
            }
            ctxapi.console.print(from + "  before aPledge= " + aPledge.toString()+ "\n")

            aPledge = aPledge.plus(1)
            aPledgeStr = aPledge.toString()

            ctxapi.db.set(from + pledgePrefix, aPledgeStr)
        },
        this.getRewardByBlockNum = function(blockNum) {
            ctxapi.console.print("========================getRewardByBlockNum==============================\n")
            totalSupply  = new BigNumber(1e+9)
            precision  = new BigNumber(1e+8)
            q = 0.5 //ratio
            a1BlockNumPerStep = new BigNumber(6250000)

            SnAllReward = totalSupply.times(precision)
            initReward = SnAllReward.times(1 - q)
            initReward = initReward.div(a1BlockNumPerStep)
            blk = new BigNumber(blockNum)
            idx = blk.dividedToIntegerBy(a1BlockNumPerStep)
            currentReward = initReward.dividedToIntegerBy(new BigNumber(2).pow(idx))
            ctxapi.console.print("current block num is " + blockNum + " , reward is:" + currentReward + " TOKEN\n")

            return currentReward

        }

})(this);

