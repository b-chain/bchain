import websocket
import json
import hashlib
try:
    import _thread
except ImportError:
    import _thread as thread
import time

def jsonRet(ret):
    try:
        print(json.dumps(json.loads(ret), indent=4))
    except KeyError as e:
        print(ret)
def on_message(ws, message):
    jsonRet(message)

def on_error(ws, error):
    print(error)

def on_close(ws):
    print("### closed ###")

def on_open(ws):
    print("### web socket opened ###")
    
##### filter Api
def newPendingTransactions():
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_subscribe",
        "params": ["newPendingTransactions"],
        'id': '1'
     }
    data = json.dumps(values).encode(encoding='utf-8')
    ws.send(data)
    
def newHeads():
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_subscribe",
        "params": ["newHeads"],
        'id': '1'
     }
    data = json.dumps(values).encode(encoding='utf-8')
    ws.send(data)

def logs(fromBlk, toBlk, addr, topics):
    values = {
        "jsonrpc": "2.0",
        "method": "bchain_subscribe",
        "params": ["logs",{"FromBlock":fromBlk, "ToBlock":toBlk, "Addresses": addr, "Topics":topics}],
        'id': '1'
     }
    data = json.dumps(values).encode(encoding='utf-8')
    ws.send(data)
    


if __name__ == "__main__":
    ws = websocket.WebSocketApp("ws://127.0.0.1:9999",on_message = on_message,on_error = on_error,on_close = on_close)
    ws.on_open = on_open
    #xx.run_forever()
    _thread.start_new_thread(ws.run_forever, ())
    data = hashlib.sha256("80000000000000000000".encode('utf-8')).hexdigest()
    print(data)

