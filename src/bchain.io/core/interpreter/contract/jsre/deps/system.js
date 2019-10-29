
var ctxapi = ctxApi();

;(function () {

    this.createContract = function(code){
        ctxapi.console.print("****** CreateContract Start ******")
        creator = ctxapi._sender
        contractAddress = ctxapi.contract.create(creator, code)
        ctxapi.console.print("ctxapi.contract.create:" + contractAddress + "\n")
        ctxapi.contract.emitEvent("createContract", true, contractAddress, false, contractAddress, true)
        ctxapi.result.setResult(contractAddress)
        ctxapi.console.print("****** CreateContract End ******")
    }
})(this);

