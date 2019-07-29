# bchain
bchain aims to build up a democratic ecosystem.

## Supported platforms
bchain currently supports the following platforms:
1. Linux
2. Arm
3. Darwin

windows platform node is in developing, will coming soon...

## Bchaind

bchaind is the bchain node program.

## Bchain-cli

bchain-cli is a simple node control program.

```
root@test245:/home/bc# ./bchain-cli
NAME:
   b chain simple wallet tool - A new cli application

USAGE:
   bchain-cli [global options] command [command options][arguments...]

VERSION:
   1.0.1

COMMANDS:
 account-create, acc-c         new a key store file with password. account-create {password}
 start-producer, s-p           start producing block. start-producer {password}
 bc-transfer, bc-tr            bchain BC transfer. bc-transfer {toAddr} {amount} {memo} {txFee(C)} {password}
 bc-nonceOf, bc-nof            bchain get account nonce. bc-nonceOf {Addr}
 bc-transferWithNonce, bc-trn  bchain BC transfer with nonce. bc-transferWithNonce {account nonce} {toAddr} {amount} {memo} {txFee(C)} {password}
 bc-balanceOf, bc-of           bchain BC balance of an addr. bc-balanceOf {Addr}
 bc-pledge, bc-pd              pledge bchain BC to pledge pool. bc-pledge {amount} {txFee(C)} {password}
 bc-redeem, bc-rd              redeem bchain BC from pledge pool. bc-redeem {amount} {txFee(C)} {password}
 bc-pledgeOf, bc-pdof          bchain BC pledge pool pledge of an addr. bc-pledgeOf {Addr}
 help, h                       Shows a list of commands or help for one 

GLOBAL OPTIONS:
   --url value, -u value  node url (default: "http://127.0.0.1:8989/")
   --help, -h             show help
   --version, -v          print the version
```



## How to config node

1. download program bchaind and bchain-cli

2. download boot config file bootCommittee.json and node.toml

3. create config directories

   ```shell
   root@test245:mkdir config -p
   ```

   move bootCommittee.json  and node.toml to config directories.

   ​

## Run a node

1. create an account

   ```shell
   root@test245:/home/bc# ./bchain-cli acc-c 6789212
      password:  6789212
      WellCome to CreateAccount
      Print NewAccount Address:0x74b22c8abe8423cce73d14f3609907327208c4ac,   Url:keystore:///home/bc/keystore/UTC--2019-07-28T13-23-22.456113068Z--74b22c8abe8423cce73d14f3609907327208c4ac
      After Address:0x74b22c8abe8423cce73d14f3609907327208c4ac,   Url:keystore:///home/bc/keystore/UTC--2019-07-28T13-23-22.456113068Z--74b22c8abe8423cce73d14f3609907327208c4ac
   ```

2. run bchaind

   ```shell
   root@test245:/home/bc# nohup ./bchaind >/dev/null &
   ```

3. start producing block(before this, need to pledge BC...)

   ```shell
   root@test245:/home/bc# ./bchain-cli s-p 6789212
   password:  6789212
   {"jsonrpc":"2.0","id":"1","result":null}
   ```

##  Pledge BC

1. pledge BC 

   ```shell
   root@test245:/home/bc# ./bchain-cli bc-pd 50 100 6789212
   0x74b22C8aBE8423CcE73D14f3609907327208c4Ac
   from 0x74b22c8abe8423cce73d14f3609907327208c4ac amount(BC) 50 txfee(C) 100
   {"jsonrpc":"2.0","id":"1","result":"0x0"}

   nonceOf 0x74b22c8abe8423cce73d14f3609907327208c4ac is 0x0
   {"jsonrpc":"2.0","id":"1","result":"0x85fe99c37031219c7b13d67c6709efbc33eeae564eafdff8418b6e6496cfc342"}
   ```

2. query pledge amount

   ```shell
   root@test245:/home/bc# ./bchain-cli bc-pdof 0x74b22c8abe8423cce73d14f3609907327208c4ac
   {"jsonrpc":"2.0","id":"1","result":["0x00f2052a01000000","0x00b1676e46000000","0xf53e000000000000"]}

   pledgeOf 0x74b22c8abe8423cce73d14f3609907327208c4ac is 50 BC
   pledge pool total is 3025 BC
   ```

   ​