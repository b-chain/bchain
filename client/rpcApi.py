# -*- coding: UTF-8 -*-

import argparse
import sys
import urllib.parse
import urllib.request
import urllib.error
import json
import os
import time
import _thread
import random
import base64

url = "http://127.0.0.1:8989/"
def print_error(file, lineno, err):
    print("[" + file + ":" + str(lineno) + "]", "Err:", err)


def post(url, values):
    headers = {
        'User-Agent':'bchain-test',
        "Content-Type": "application/json",
        "Connection": "keep-alive"
    }

    try:
        data = json.dumps(values).encode(encoding='utf-8')

        request = urllib.request.Request(url, data, headers)
        respose = urllib.request.urlopen(request).read().decode('utf-8')

    except urllib.error.HTTPError as e:
        print_error(__file__, sys._getframe().f_lineno, e)
        exit(1)

    except urllib.error.URLError as e:
        print_error(__file__, sys._getframe().f_lineno, e)
        exit(1)
    except TypeError as e:
        print_error(__file__, sys._getframe().f_lineno, e)

    return respose

def commonRet(ret):
    try:
        print(json.loads(ret)['result'])
        return json.loads(ret)['result']
    except KeyError as e:
        print(ret)
        return ret
    
def allRet(ret):
    print(ret)
def jsonRet(ret):
    try:
        print(json.dumps(json.loads(ret)['result'], indent=4))
        return json.loads(ret)['result']
    except KeyError as e:
        print(ret)
        return ret 
#PublicbchainAPI
def protocolVersion():
    print("rpc: bchain_protocolVersion")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_protocolVersion",
        "params": [],
        'id': '1'
     }
    return post(url, values)

def protocolVersionRet(ret):
    #print(ret)
    print(json.loads(ret)['result'])
    
def syncing():
    print("rpc: bchain_syncing")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_syncing",
        "params": [],
        'id': '1'
     }
    return post(url, values)

#NewPublicTxPoolAPI
def content():
    print("rpc: txpool_content")
    values = {
        "jsonrpc": "2.0",
        "method": "txpool_content",
        "params": [],
        'id': '1'
     }
    return post(url, values)

def status():
    print("rpc: txpool_status")
    values = {
        "jsonrpc": "2.0",
        "method": "txpool_status",
        "params": [],
        'id': '1'
     }
    return post(url, values)

def inspect():
    print("rpc: txpool_inspect")
    values = {
        "jsonrpc": "2.0",
        "method": "txpool_inspect",
        "params": [],
        'id': '1'
     }
    return post(url, values)

#PublicAccountAPI
def accounts():
    print("rpc: bchain_accounts")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_accounts",
        "params": [],
        'id': '1'
     }
    return post(url, values)

#PrivateAccountAPI
def listAccounts():
    print("rpc: personal_listAccounts")
    values = {
        "jsonrpc": "2.0",
        "method": "personal_listAccounts",
        "params": [],
        'id': '1'
     }
    return post(url, values)
def listWallets():
    print("rpc: personal_listWallets")
    values = {
        "jsonrpc": "2.0",
        "method": "personal_listWallets",
        "params": [],
        'id': '1'
     }
    return post(url, values)
def newAccount(password):
    print("rpc: personal_newAccount")
    values = {
        "jsonrpc": "2.0",
        "method": "personal_newAccount",
        "params": [password],
        'id': '1'
        }
    return post(url, values)

def importRawKey(prikey, password):
    print("rpc: personal_importRawKey")
    values = {
        "jsonrpc": "2.0",
        "method": "personal_importRawKey",
        "params": [prikey, password],
        'id': '1'
        }
    return post(url, values)
def unlockAccount(addr, password, duration):
    print("rpc: personal_unlockAccount")
    values = {
        "jsonrpc": "2.0",
        "method": "personal_unlockAccount",
        "params": [addr, password, duration],
        'id': '1'
    }
    return post(url, values)

def lockAccount(addr):
    print("rpc: personal_lockAccount")
    values = {
        "jsonrpc": "2.0",
        "method": "personal_lockAccount",
        "params": [addr],
        'id': '1'
    }
    return post(url, values)

def action(contract_address, params):
    return {"address": contract_address, "params": "0x"+params.encode(encoding='utf-8').hex()}
def actions():
    return []
def actions_append(actions, action):
    actions.append(action)
    return actions
def sendTransaction(fromAddess,password, acts):
    print("rpc: personal_sendTransaction")
    values = {
        "jsonrpc": "2.0",
        "method": "personal_sendTransaction",
        "params": [{"from":fromAddess,"actions":acts}, password],
        'id': '1'
    }
    return post(url, values)

def sign(data, addr,password):
    print("rpc: personal_sign")
    dataHex = "0x" + data.hex()
    values = {
        "jsonrpc": "2.0",
        "method": "personal_sign",
        "params": [dataHex, addr, password],
        'id': '1'
    }
    return post(url, values)

def signRet(ret):
    try:
        rlt = json.loads(ret)['result']
        print(rlt)
        bytes_rlt = bytes.fromhex(rlt[2:])
        print(bytes_rlt)
        #print(bytes_rlt.hex())
        return bytes_rlt
    except:
        print(ret)
        return ret
    
def ecRecover(data, sig):
    print("rpc: personal_ecRecover")
    dataHex = "0x" + data.hex()
    sigHex = "0x" + sig.hex()
    values = {
        "jsonrpc": "2.0",
        "method": "personal_ecRecover",
        "params": [dataHex, sigHex],
        'id': '1'
    }
    return post(url, values)


#PublicBlockChainAPI
def blockNumber():
    print("rpc: bchain_blockNumber")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_blockNumber",
        "params": [],
        'id': '1'
    }
    return post(url, values)

def getStatInfoByNumber(number):
    print("rpc: bchain_getStatInfoByNumber")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getStatInfoByNumber",
        "params": [number],
        'id': '1'
    }
    return post(url, values)

def getBlockByNumber(number, fullTx):
    print("rpc: bchain_getBlockByNumber")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getBlockByNumber",
        "params": [number, fullTx],
        'id': '1'
    }
    return post(url, values)


def getBlockByHash(blkHash, fullTx):
    print("rpc: bchain_getBlockByHash")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getBlockByHash",
        "params": [blkHash, fullTx],
        'id': '1'
    }
    return post(url, values)

def getCode(addr, number):
    print("rpc: bchain_getCode")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getCode",
        "params": [addr, number],
        'id': '1'
    }
    return post(url, values)

def getCodeRet(ret):
    rlt = json.loads(ret)['result']
    bytes_rlt = bytes.fromhex(rlt[2:])
    print(bytes_rlt)
    return bytes_rlt

def getStorageAt(addr, keyHash, number):
    print("rpc: bchain_getStorageAt")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getStorageAt",
        "params": [addr, keyHash, number],
        'id': '1'
    }
    return post(url, values)

def actionCall(action, number):
    print("rpc: bchain_actionCall")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_actionCall",
        "params": [action, number],
        'id': '1'
    }
    return post(url, values)

#PublicTransactionPoolAPI
def getBlockTransactionCountByNumber(number):
    print("rpc: bchain_getBlockTransactionCountByNumber")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getBlockTransactionCountByNumber",
        "params": [number],
        'id': '1'
    }
    return post(url, values)

def getBlockTransactionCountByHash(blkHash):
    print("rpc: bchain_getBlockTransactionCountByHash")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getBlockTransactionCountByHash",
        "params": [blkHash],
        'id': '1'
    }
    return post(url, values)

def getTransactionByBlockNumberAndIndex(number, index):
    print("rpc: bchain_getTransactionByBlockNumberAndIndex")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getTransactionByBlockNumberAndIndex",
        "params": [number, index],
        'id': '1'
    }
    return post(url, values)

def getTransactionByBlockHashAndIndex(blkHash, index):
    print("rpc: bchain_getTransactionByBlockHashAndIndex")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getTransactionByBlockHashAndIndex",
        "params": [blkHash, index],
        'id': '1'
    }
    return post(url, values)

def getRawTransactionByBlockNumberAndIndex(number, index):
    print("rpc: bchain_getRawTransactionByBlockNumberAndIndex")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getRawTransactionByBlockNumberAndIndex",
        "params": [number, index],
        'id': '1'
    }
    return post(url, values)

def getRawTransactionByBlockHashAndIndex(blkHash, index):
    print("rpc: bchain_getRawTransactionByBlockHashAndIndex")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getRawTransactionByBlockHashAndIndex",
        "params": [blkHash, index],
        'id': '1'
    }
    return post(url, values)

def getAccountNonce(addr, number):
    print("rpc: bchain_getAccountNonce")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getAccountNonce",
        "params": [addr, number],
        'id': '1'
    }
    return post(url, values)

def getTransactionByHash(txHash):
    print("rpc: bchain_getTransactionByHash")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getTransactionByHash",
        "params": [txHash],
        'id': '1'
    }
    return post(url, values)
	
def getTransactionByAddress(txAddr, nonce):
    print("rpc: getTransactionByAddress")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getTransactionByAddress",
        "params": [txAddr, nonce],
        'id': '1'
    }
    return post(url, values)

def getRawTransactionByHash(txHash):
    print("rpc: bchain_getRawTransactionByHash")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getRawTransactionByHash",
        "params": [txHash],
        'id': '1'
    }
    return post(url, values)

def getTransactionReceipt(txHash):
    print("rpc: bchain_getTransactionReceipt")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getTransactionReceipt",
        "params": [txHash],
        'id': '1'
    }
    return post(url, values)

def getCoinBaseLogByBlockNumber(nunmber):
    print("rpc: bchain_getCoinBaseLogByBlockNumber")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_getCoinBaseLogByBlockNumber",
        "params": [nunmber],
        'id': '1'
    }
    return post(url, values)

def sendTransaction_nopassword(fromAddess, acts):
    print("rpc: bchain_sendTransaction")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_sendTransaction",
        "params": [{"from":fromAddess,"actions":acts}],
        'id': '1'
    }
    return post(url, values)

def sign_nopassword(addr, data):
    print("rpc: bchain_sign")
    dataHex = "0x" + data.hex()
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_sign",
        "params": [addr, dataHex],
        'id': '1'
    }
    return post(url, values)

def signTransaction_nopassword(fromAddess, acts):
    print("rpc: bchain_signTransaction")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_signTransaction",
        "params": [{"from":fromAddess,"actions":acts}],
        'id': '1'
    }
    return post(url, values)

def signTransactionWithNonce_nopassword(fromAddess, nonce, acts):
    print("rpc: bchain_signTransaction")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_signTransaction",
        "params": [{"from":fromAddess,"actions":acts, "nonce": nonce}],
        'id': '1'
    }
    return post(url, values)

def sendRawTransaction(msgpSignedTx):
    print("rpc: bchain_sendRawTransaction")
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_sendRawTransaction",
        "params": [msgpSignedTx],
        'id': '1'
    }
    return post(url, values)

	
def blockproductor_start(url):
    values = {
        "jsonrpc": "2.0",
        "method": "blockproducer_start",
        "params": [10, pass_coninbase],
        'id': '22'
    }
    return post(url, values)

pass_coninbase = "123"
pass_other = "123"

def start_blockproductor():
    url1 = "http://192.168.2.245:7981/"
    url2 = "http://localhost:7981/"
    url3 = "http://localhost:7982/"
    url4 = "http://localhost:7983/"
    blockproductor_start(url1)
    print('start blockproductor!!!')
    blockproductor_start(url2)
    print('start blockproductor!!!')
    blockproductor_start(url3)
    print('start blockproductor!!!')
    blockproductor_start(url4)
    print('start blockproductor!!!')

def exampleCallApi():
    ret = protocolVersion()
    protocolVersionRet(ret)
    
    ret = syncing()
    commonRet(ret)

    ret = content()
    commonRet(ret)

    ret = status()
    commonRet(ret)

    ret = inspect()
    commonRet(ret)
    
    ret = accounts()
    commonRet(ret)

    ret = listAccounts()
    commonRet(ret)

    ret = listWallets()
    commonRet(ret)

    #ret = newAccount("123")
    #commonRet(ret)

    ret = importRawKey("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292", "123" )
    allRet(ret)

    ret = unlockAccount("0x55fda7601ffa55f61b819642816460aa24883f7f", "123", 100)
    commonRet(ret)

    ret = lockAccount("0x55fda7601ffa55f61b819642816460aa24883f7f")
    commonRet(ret)
    
    #act = action("0xBdA0610666Fb9Bc52Dcf30bca4a42288c853A94B", "pledge(99);")
    act = action("0xBdA0610666Fb9Bc52Dcf30bca4a42288c853A94B", "redeem(99);")
    acts = actions()
    acts = actions_append(acts, act)
    ret = sendTransaction("0x2e68b0583021d78c122f719fc82036529a90571d", "123", acts)
    commonRet(ret)

    ret = sign("12333".encode(encoding='utf-8'), "0x2e68b0583021d78c122f719fc82036529a90571d", "123")
    sig = signRet(ret)

    ret = ecRecover("12333".encode(encoding='utf-8'), sig)
    commonRet(ret)

    ret = blockNumber()
    commonRet(ret)

    ret = getStatInfoByNumber("latest")
    commonRet(ret)

    ret = getBlockByNumber("latest", True)
    jsonRet(ret)

    ret = getBlockByHash("0xa09eea0cb17474f81b0941f3ea4f78e6a62e2f5a33bb2416b38d6e58092b11d3", True)
    jsonRet(ret)

    ret = getCode("0xBdA0610666Fb9Bc52Dcf30bca4a42288c853A94B", "latest")
    getCodeRet(ret)

    ret = getStorageAt("0xBdA0610666Fb9Bc52Dcf30bca4a42288c853A94B", "0xa09eea0cb17474f81b0941f3ea4f78e6a62e2f5a33bb2416b38d6e58092b11d3", "latest")
    allRet(ret)

    act = action("0xBdA0610666Fb9Bc52Dcf30bca4a42288c853A94B", 'pledgeOf("0x2e68b0583021d78c122f719fc82036529a90571d");')
    ret = actionCall(act, "latest")
    commonRet(ret)

    ret = getBlockTransactionCountByNumber("0xe43")
    commonRet(ret)

    ret = getBlockTransactionCountByHash("0xa09eea0cb17474f81b0941f3ea4f78e6a62e2f5a33bb2416b38d6e58092b11d3")
    commonRet(ret)

    ret = getTransactionByBlockNumberAndIndex("0xe43", "0x0")
    jsonRet(ret)

    ret = getTransactionByBlockHashAndIndex("0xc8702b6c4ac37e864c15e142a132c4f5f6e191dec7ba8f854ed6aa88635b0b2f", "0x0")
    jsonRet(ret)

    ret = getRawTransactionByBlockNumberAndIndex("0xe43", "0x0")
    jsonRet(ret)

    ret = getRawTransactionByBlockHashAndIndex("0xc8702b6c4ac37e864c15e142a132c4f5f6e191dec7ba8f854ed6aa88635b0b2f", "0x0")
    jsonRet(ret)

    ret = getAccountNonce("0x2e68b0583021d78c122f719fc82036529a90571d", "latest")
    commonRet(ret)

    ret = getTransactionByHash("0x0cf080e5630679c9adf5ce3d6aab4a02b0f9119f705c4be3d5562e6f8850de00")
    jsonRet(ret)

    ret = getRawTransactionByHash("0x0cf080e5630679c9adf5ce3d6aab4a02b0f9119f705c4be3d5562e6f8850de00")
    jsonRet(ret)

    ret = getTransactionReceipt("0x0cf080e5630679c9adf5ce3d6aab4a02b0f9119f705c4be3d5562e6f8850de00")
    jsonRet(ret)

    ret = getCoinBaseLogByBlockNumber("latest")
    jsonRet(ret)

    ret = unlockAccount("0x2e68b0583021d78c122f719fc82036529a90571d", "123", 100)
    commonRet(ret)
    
    act = action("0xBdA0610666Fb9Bc52Dcf30bca4a42288c853A94B", "pledge(10);")
    acts = actions()
    acts = actions_append(acts, act)
    ret = sendTransaction_nopassword("0x2e68b0583021d78c122f719fc82036529a90571d", acts)
    commonRet(ret)

    ret = unlockAccount("0x2e68b0583021d78c122f719fc82036529a90571d", "123", 100)
    commonRet(ret)
    ret = sign_nopassword("0x2e68b0583021d78c122f719fc82036529a90571d", "12333".encode(encoding='utf-8'))
    signRet(ret)

    ret = getAccountNonce("0x2e68b0583021d78c122f719fc82036529a90571d", "latest")
    nonce = commonRet(ret)
    
    ret = unlockAccount("0x2e68b0583021d78c122f719fc82036529a90571d", "123", 100)
    commonRet(ret)
    act = action("0xBdA0610666Fb9Bc52Dcf30bca4a42288c853A94B", "pledge(10);")
    acts = actions()
    acts = actions_append(acts, act)
    ret = signTransaction_nopassword("0x2e68b0583021d78c122f719fc82036529a90571d", acts)
    msgpSignTx = jsonRet(ret)['raw']

    ret = sendRawTransaction(msgpSignTx)
    commonRet(ret)
		
def main():
    parser = argparse.ArgumentParser()

    parser.add_argument("-u", "--url",
                        required=False,
                        nargs=1,
                        type=str,
                        default=["http://127.0.0.1:8989/"],
                        dest="base_url",
                        metavar="URL",
                        help="the url of peer (default: %(default)s)")


    args = parser.parse_args()
    print(args)
    url = args.base_url[0]

    os.system("pause")

if __name__ == "__main__":
    main()
