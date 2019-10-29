var ctxapi = ctxApi();

;(function () {

    //var decimals = 18
    //var totalSupply = new BigNumber(1e+9)
    blockRewordAdjustNumber = new BigNumber(6250000)
    initReword = new BigNumber(80e+18)
    q = new BigNumber(0.5)

    function getRewordByNumber(number){
        ctxapi.console.print("===========getRewordByNumber===========" + number + "\n")
        idx = number.dividedToIntegerBy(blockRewordAdjustNumber)
        reword = (initReword.times(q.pow(idx))).toFixed(0)
        return new BigNumber(reword)
    }

    function strToLower(param){
        return param.toLowerCase()
    }

    this.balanceReward = function(from1 ){
        ctxapi.console.print("********* Reward Start *********\n")

        //call balanceReward before
        consensus_only_once = ctxapi.cache.get("flag_consensus_only_once")
        if (consensus_only_once.length != 0) {
            ctxapi.assert.assert(false , "called balanceReward before")
        }


        var from = strToLower(from1)

        blkNumber = new BigNumber(ctxapi._number)
        reword = getRewordByNumber(blkNumber)
        var rewardFunc = "balanceTransfer('" + from + "'," + reword.toString() + ")"
        ctxapi.call.call("0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B", rewardFunc)
        ctxapi.contract.emitEvent(from, true, reword.toString(), true)

        var queryFunc = "balanceQuery('" + from + "')"
        ctxapi.call.call("0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B", queryFunc)

        queryFunc = "balanceQuery('0x55fda7601ffa55f61b819642816460aa24883f7f')"
        ctxapi.call.call("0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B", queryFunc)
        if (consensus_only_once.length == 0) {
            ctxapi.console.print("======================111Reward Emplace================================\n")
            ctxapi.cache.emplace("flag_consensus_only_once" , "true")
        }
        ctxapi.console.print("********* Reward End *********\n")
    },

    this.balanceFee = function(amount){

        var check = new BigNumber(amount)
        if (check.lte(0)) {
            ctxapi.assert.assert(false , "balanceFee amount <= 0")
        }

        ctxapi.console.print("********* balanceFee Start *********\n")
        var miner = ctxapi._miner
        if (ctxapi._miner == ctxapi._sender) {
            ctxapi.console.print("============balanceFee Miner == Sender ================\n")
            return
        }
        ctxapi.console.print("============balanceFee Start================\n")
        ctxapi.console.print("======Miner>" + miner + "\n")
        var rewardFunc = "balanceTransfer('" + ctxapi._miner + "'," + amount.toString() + ")"
        ctxapi.call.call("0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B", rewardFunc)

        // basicBalanceTransfer(miner , amount)
        ctxapi.console.print("********* balanceFee End *********\n")
    }


})(this);

