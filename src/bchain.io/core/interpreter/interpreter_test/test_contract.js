//var apijs = require('context_api') // not support

var ctxapi = ctxApi();


function PrintHello(){
    ctxapi.console.print("================================PrintHello======================================\n");
}

;(function () {


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
        }
})(this);