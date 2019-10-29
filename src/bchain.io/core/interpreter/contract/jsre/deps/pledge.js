
var ctxapi = ctxApi();

function getName(){
    return "pledge"
}

function getVersion(){
    return "V0.0.1"
}

;(function () {
    var pledgePrefix = "pledge"
    var pledgeTotal = "_pledgeTotal"

    function strToLower(param){
        return param.toLowerCase()
    }

    function addTotalPledge(amount){
        var totalStr = ctxapi.db.get(pledgeTotal)
        var pTotal = new BigNumber(totalStr.length > 0 ? totalStr : 0)
        pTotal = pTotal.plus(amount)
        ctxapi.db.set(pledgeTotal, pTotal.toString())
    }

    function subTotalPledge(amount){
        var totalStr = ctxapi.db.get(pledgeTotal)
        var pTotal = new BigNumber(totalStr.length > 0 ? totalStr : 0)
        if (pTotal.lt(amount)) {
            ctxapi.assert.assert(false , "total pledge balance is not enough")
        }
        pTotal = pTotal.minus(amount)
        ctxapi.db.set(pledgeTotal, pTotal.toString())
    }

    function pledge(fromAddress, amount) {
        var from = strToLower(fromAddress)
        var pledgeStr = ctxapi.db.get(from + pledgePrefix)
        var pledge = new BigNumber(pledgeStr.length > 0 ? pledgeStr : 0)

        pledge = pledge.plus(amount)
        ctxapi.db.set(from + pledgePrefix, pledge.toString())
        addTotalPledge(amount)
    }

    function redeem(fromAddress, amount) {
        var from = strToLower(fromAddress)
        var pledgeStr = ctxapi.db.get(from + pledgePrefix)
        var pledge = new BigNumber(pledgeStr.length > 0 ? pledgeStr : 0)
        if (pledge.lt(amount)) {
            ctxapi.assert.assert(false , "pledge balance is not enough")
        }

        pledge = pledge.minus(amount)
        ctxapi.db.set(from + pledgePrefix , pledge.toString())
        subTotalPledge(amount)
    }

    this.pledge = function(tokens){
        var amount = new BigNumber(tokens)
        if (amount.lte(0)) {
            ctxapi.assert.assert(false , "pledge amount <= 0")
        }
        ctxapi.console.print("-----------pledge from:" + ctxapi._sender + " amount " + amount + "\n")
        pledge(ctxapi._sender, amount)

        // bchain pledge
        to = strToLower(ctxapi._contract)
        transferFunc = "balanceTransfer('" + to + "'," + (amount.times(1e+18)).toString() + ")"
        ctxapi.call.call("0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B", transferFunc)
    },

    this.redeem = function(tokens){
        var amount = new BigNumber(tokens)
        if (amount.lte(0)) {
            ctxapi.assert.assert(false , "redeem amount <= 0")
        }
        ctxapi.console.print("-----------redeem from:" + ctxapi._sender + " amount " + amount + "\n")
        redeem(ctxapi._sender, amount)

        //bchain redeem
        to = strToLower(ctxapi._sender)
        transferFunc = "balanceTransfer('" + to + "'," + (amount.times(1e+18)).toString() + ")"
        ctxapi.call.innerCall("0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B", transferFunc)
    },

    this.pledgeOf = function(address){
        var _address = strToLower(address)
        var pledgeStr = ctxapi.db.get(_address + pledgePrefix)
        var pledge = new BigNumber(pledgeStr.length > 0 ? pledgeStr : 0)
        ctxapi.console.print(address+" pledge = " + pledge + "\n")
        ctxapi.result.setResult(pledge.toString(10))
    },

    this.totalPledge = function(){
        var pledgeStr = ctxapi.db.get(pledgeTotal)
        var pledge = new BigNumber(pledgeStr.length > 0 ? pledgeStr : 0)
        ctxapi.console.print("total pledge = " + pledge + "\n")
        ctxapi.result.setResult(pledge.toString(10))
    },

    this.pledgeOfExt = function(address){
        var _address = strToLower(address)
        var pledgeStr = ctxapi.db.get(_address + pledgePrefix)
        var pledge = new BigNumber(pledgeStr.length > 0 ? pledgeStr : 0)
        ctxapi.console.print(address+" pledge = " + pledge + "\n")
        ctxapi.result.setResult(pledge.toString(10))

        var pTotalStr = ctxapi.db.get(pledgeTotal)
        var pTotal = new BigNumber(pTotalStr.length > 0 ? pTotalStr : 0)
        ctxapi.console.print("total pledge = " + pTotal + "\n")
        ctxapi.result.setResult(pTotal.toString(10))
    }
})(this);
