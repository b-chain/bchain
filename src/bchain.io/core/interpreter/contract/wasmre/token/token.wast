(module
  (type $FUNCSIG$vi (func (param i32)))
  (type $FUNCSIG$ii (func (param i32) (result i32)))
  (type $FUNCSIG$vii (func (param i32 i32)))
  (type $FUNCSIG$iiii (func (param i32 i32 i32) (result i32)))
  (type $FUNCSIG$viiii (func (param i32 i32 i32 i32)))
  (import "env" "action_sender" (func $action_sender (param i32)))
  (import "env" "assert" (func $assert (param i32 i32)))
  (import "env" "db_emplace" (func $db_emplace (param i32 i32 i32 i32)))
  (import "env" "db_get" (func $db_get (param i32 i32 i32) (result i32)))
  (import "env" "db_set" (func $db_set (param i32 i32 i32 i32)))
  (import "env" "requireAuth" (func $requireAuth (param i32) (result i32)))
  (import "env" "setResult" (func $setResult (param i32 i32)))
  (import "env" "str2lower" (func $str2lower (param i32)))
  (import "env" "strjoint" (func $strjoint (param i32 i32 i32) (result i32)))
  (import "env" "strlen" (func $strlen (param i32) (result i32)))
  (table 0 anyfunc)
  (memory $0 1)
  (data (i32.const 4) "\a0A\00\00")
  (data (i32.const 16) "symbol len exceed\00")
  (data (i32.const 48) "name len exceed\00")
  (data (i32.const 64) "symbol \00")
  (data (i32.const 80) "name \00")
  (data (i32.const 96) "decimals \00")
  (data (i32.const 112) "supply \00")
  (data (i32.const 128) "total supply exceed\00")
  (data (i32.const 160) "total supply exceed!\00")
  (data (i32.const 192) "transfer: symbol is not exist!\00")
  (data (i32.const 224) "get sender token error\00")
  (data (i32.const 256) "insufficient token\00")
  (data (i32.const 288) "get supply db error\00")
  (data (i32.const 320) "get decimals db error\00")
  (data (i32.const 352) "get symbol db error\00")
  (data (i32.const 384) "get name db error\00")
  (export "memory" (memory $0))
  (export "_ZN5token6createEPcS0_iy" (func $_ZN5token6createEPcS0_iy))
  (export "_ZN5token8transferEPcyS0_S0_" (func $_ZN5token8transferEPcyS0_S0_))
  (export "_ZN5token9balanceOfEPcS0_" (func $_ZN5token9balanceOfEPcS0_))
  (export "_ZN5token9getSupplyEPc" (func $_ZN5token9getSupplyEPc))
  (export "_ZN5token11getDecimalsEPc" (func $_ZN5token11getDecimalsEPc))
  (export "_ZN5token9getSymbolEPc" (func $_ZN5token9getSymbolEPc))
  (export "_ZN5token7getNameEPc" (func $_ZN5token7getNameEPc))
  (export "create" (func $create))
  (export "transfer" (func $transfer))
  (export "balanceOf" (func $balanceOf))
  (export "getSupply" (func $getSupply))
  (export "getDecimals" (func $getDecimals))
  (export "getSymbol" (func $getSymbol))
  (export "getName" (func $getName))
  (func $_ZN5token6createEPcS0_iy (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i64)
    (local $5 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $5
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 208)
        )
      )
    )
    (i32.store offset=204
      (get_local $5)
      (get_local $3)
    )
    (i64.store offset=192
      (get_local $5)
      (get_local $4)
    )
    (call $action_sender
      (i32.add
        (get_local $5)
        (i32.const 144)
      )
    )
    (drop
      (call $requireAuth
        (i32.add
          (get_local $5)
          (i32.const 144)
        )
      )
    )
    (call $assert
      (i32.lt_s
        (call $strlen
          (get_local $1)
        )
        (i32.const 64)
      )
      (i32.const 16)
    )
    (call $assert
      (i32.lt_s
        (call $strlen
          (get_local $2)
        )
        (i32.const 64)
      )
      (i32.const 48)
    )
    (call $db_emplace
      (i32.add
        (get_local $5)
        (i32.const 16)
      )
      (call $strjoint
        (i32.const 64)
        (get_local $1)
        (i32.add
          (get_local $5)
          (i32.const 16)
        )
      )
      (get_local $1)
      (call $strlen
        (get_local $1)
      )
    )
    (call $db_emplace
      (i32.add
        (get_local $5)
        (i32.const 16)
      )
      (call $strjoint
        (i32.const 80)
        (get_local $1)
        (i32.add
          (get_local $5)
          (i32.const 16)
        )
      )
      (get_local $2)
      (call $strlen
        (get_local $2)
      )
    )
    (call $db_emplace
      (i32.add
        (get_local $5)
        (i32.const 16)
      )
      (call $strjoint
        (i32.const 96)
        (get_local $1)
        (i32.add
          (get_local $5)
          (i32.const 16)
        )
      )
      (i32.add
        (get_local $5)
        (i32.const 204)
      )
      (i32.const 4)
    )
    (call $db_emplace
      (i32.add
        (get_local $5)
        (i32.const 16)
      )
      (call $strjoint
        (i32.const 112)
        (get_local $1)
        (i32.add
          (get_local $5)
          (i32.const 16)
        )
      )
      (i32.add
        (get_local $5)
        (i32.const 192)
      )
      (i32.const 8)
    )
    (i64.store offset=8
      (get_local $5)
      (tee_local $4
        (i64.load offset=192
          (get_local $5)
        )
      )
    )
    (call $assert
      (i32.lt_u
        (i32.wrap/i64
          (i64.shr_u
            (get_local $4)
            (i64.const 56)
          )
        )
        (i32.const 25)
      )
      (i32.const 128)
    )
    (block $label$0
      (br_if $label$0
        (i32.lt_s
          (i32.load offset=204
            (get_local $5)
          )
          (i32.const 1)
        )
      )
      (i64.store offset=8
        (get_local $5)
        (tee_local $4
          (i64.mul
            (get_local $4)
            (i64.const 10)
          )
        )
      )
      (call $assert
        (i32.lt_u
          (i32.wrap/i64
            (i64.shr_u
              (get_local $4)
              (i64.const 56)
            )
          )
          (i32.const 25)
        )
        (i32.const 160)
      )
      (br_if $label$0
        (i32.lt_s
          (i32.load offset=204
            (get_local $5)
          )
          (i32.const 2)
        )
      )
      (set_local $2
        (i32.const 1)
      )
      (loop $label$1
        (i64.store offset=8
          (get_local $5)
          (tee_local $4
            (i64.mul
              (i64.load offset=8
                (get_local $5)
              )
              (i64.const 10)
            )
          )
        )
        (call $assert
          (i32.lt_u
            (i32.wrap/i64
              (i64.shr_u
                (get_local $4)
                (i64.const 56)
              )
            )
            (i32.const 25)
          )
          (i32.const 160)
        )
        (br_if $label$1
          (i32.lt_s
            (tee_local $2
              (i32.add
                (get_local $2)
                (i32.const 1)
              )
            )
            (i32.load offset=204
              (get_local $5)
            )
          )
        )
      )
    )
    (call $db_emplace
      (i32.add
        (get_local $5)
        (i32.const 16)
      )
      (call $strjoint
        (get_local $1)
        (i32.add
          (get_local $5)
          (i32.const 144)
        )
        (i32.add
          (get_local $5)
          (i32.const 16)
        )
      )
      (i32.add
        (get_local $5)
        (i32.const 8)
      )
      (i32.const 8)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $5)
        (i32.const 208)
      )
    )
  )
  (func $_ZN5token8transferEPcyS0_S0_ (param $0 i32) (param $1 i32) (param $2 i64) (param $3 i32) (param $4 i32)
    (local $5 i32)
    (local $6 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $6
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 320)
        )
      )
    )
    (call $action_sender
      (i32.add
        (get_local $6)
        (i32.const 272)
      )
    )
    (call $assert
      (i32.gt_s
        (call $db_get
          (i32.add
            (get_local $6)
            (i32.const 144)
          )
          (call $strjoint
            (i32.const 64)
            (get_local $3)
            (i32.add
              (get_local $6)
              (i32.const 144)
            )
          )
          (get_local $3)
        )
        (i32.const 0)
      )
      (i32.const 192)
    )
    (i64.store offset=136
      (get_local $6)
      (i64.const 0)
    )
    (call $assert
      (i32.eq
        (call $db_get
          (i32.add
            (get_local $6)
            (i32.const 144)
          )
          (tee_local $5
            (call $strjoint
              (get_local $3)
              (i32.add
                (get_local $6)
                (i32.const 272)
              )
              (i32.add
                (get_local $6)
                (i32.const 144)
              )
            )
          )
          (i32.add
            (get_local $6)
            (i32.const 136)
          )
        )
        (i32.const 8)
      )
      (i32.const 224)
    )
    (call $assert
      (i64.ge_u
        (i64.load offset=136
          (get_local $6)
        )
        (get_local $2)
      )
      (i32.const 256)
    )
    (i64.store offset=136
      (get_local $6)
      (i64.sub
        (i64.load offset=136
          (get_local $6)
        )
        (get_local $2)
      )
    )
    (call $db_set
      (i32.add
        (get_local $6)
        (i32.const 144)
      )
      (get_local $5)
      (i32.add
        (get_local $6)
        (i32.const 136)
      )
      (i32.const 8)
    )
    (i64.store offset=128
      (get_local $6)
      (i64.const 0)
    )
    (call $str2lower
      (get_local $1)
    )
    (drop
      (call $db_get
        (get_local $6)
        (tee_local $3
          (call $strjoint
            (get_local $3)
            (get_local $1)
            (get_local $6)
          )
        )
        (i32.add
          (get_local $6)
          (i32.const 128)
        )
      )
    )
    (i64.store offset=128
      (get_local $6)
      (i64.add
        (i64.load offset=128
          (get_local $6)
        )
        (get_local $2)
      )
    )
    (call $db_set
      (get_local $6)
      (get_local $3)
      (i32.add
        (get_local $6)
        (i32.const 128)
      )
      (i32.const 8)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $6)
        (i32.const 320)
      )
    )
  )
  (func $_ZN5token9balanceOfEPcS0_ (param $0 i32) (param $1 i32) (param $2 i32)
    (local $3 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $3
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 144)
        )
      )
    )
    (i64.store offset=136
      (get_local $3)
      (i64.const 0)
    )
    (call $str2lower
      (get_local $1)
    )
    (drop
      (call $db_get
        (get_local $3)
        (call $strjoint
          (get_local $2)
          (get_local $1)
          (get_local $3)
        )
        (i32.add
          (get_local $3)
          (i32.const 136)
        )
      )
    )
    (call $setResult
      (i32.add
        (get_local $3)
        (i32.const 136)
      )
      (i32.const 8)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $3)
        (i32.const 144)
      )
    )
  )
  (func $_ZN5token9getSupplyEPc (param $0 i32) (param $1 i32)
    (local $2 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $2
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 144)
        )
      )
    )
    (i64.store offset=136
      (get_local $2)
      (i64.const 0)
    )
    (call $assert
      (i32.eq
        (call $db_get
          (get_local $2)
          (call $strjoint
            (i32.const 112)
            (get_local $1)
            (get_local $2)
          )
          (i32.add
            (get_local $2)
            (i32.const 136)
          )
        )
        (i32.const 8)
      )
      (i32.const 288)
    )
    (call $setResult
      (i32.add
        (get_local $2)
        (i32.const 136)
      )
      (i32.const 8)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $2)
        (i32.const 144)
      )
    )
  )
  (func $_ZN5token11getDecimalsEPc (param $0 i32) (param $1 i32)
    (local $2 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $2
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 144)
        )
      )
    )
    (i32.store offset=140
      (get_local $2)
      (i32.const 0)
    )
    (call $assert
      (i32.eq
        (call $db_get
          (get_local $2)
          (call $strjoint
            (i32.const 96)
            (get_local $1)
            (get_local $2)
          )
          (i32.add
            (get_local $2)
            (i32.const 140)
          )
        )
        (i32.const 4)
      )
      (i32.const 320)
    )
    (call $setResult
      (i32.add
        (get_local $2)
        (i32.const 140)
      )
      (i32.const 4)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $2)
        (i32.const 144)
      )
    )
  )
  (func $_ZN5token9getSymbolEPc (param $0 i32) (param $1 i32)
    (local $2 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $2
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 192)
        )
      )
    )
    (call $assert
      (i32.lt_s
        (tee_local $1
          (call $db_get
            (get_local $2)
            (call $strjoint
              (i32.const 64)
              (get_local $1)
              (get_local $2)
            )
            (i32.add
              (get_local $2)
              (i32.const 128)
            )
          )
        )
        (i32.const 65)
      )
      (i32.const 352)
    )
    (call $setResult
      (i32.add
        (get_local $2)
        (i32.const 128)
      )
      (get_local $1)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $2)
        (i32.const 192)
      )
    )
  )
  (func $_ZN5token7getNameEPc (param $0 i32) (param $1 i32)
    (local $2 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $2
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 192)
        )
      )
    )
    (call $assert
      (i32.lt_s
        (tee_local $1
          (call $db_get
            (get_local $2)
            (call $strjoint
              (i32.const 80)
              (get_local $1)
              (get_local $2)
            )
            (i32.add
              (get_local $2)
              (i32.const 128)
            )
          )
        )
        (i32.const 65)
      )
      (i32.const 384)
    )
    (call $setResult
      (i32.add
        (get_local $2)
        (i32.const 128)
      )
      (get_local $1)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $2)
        (i32.const 192)
      )
    )
  )
  (func $create (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i64)
    (call $_ZN5token6createEPcS0_iy
      (get_local $0)
      (get_local $0)
      (get_local $1)
      (get_local $2)
      (get_local $3)
    )
  )
  (func $transfer (param $0 i32) (param $1 i64) (param $2 i32) (param $3 i32)
    (call $_ZN5token8transferEPcyS0_S0_
      (get_local $0)
      (get_local $0)
      (get_local $1)
      (get_local $2)
      (get_local $0)
    )
  )
  (func $balanceOf (param $0 i32) (param $1 i32)
    (call $_ZN5token9balanceOfEPcS0_
      (get_local $0)
      (get_local $0)
      (get_local $1)
    )
  )
  (func $getSupply (param $0 i32)
    (call $_ZN5token9getSupplyEPc
      (get_local $0)
      (get_local $0)
    )
  )
  (func $getDecimals (param $0 i32)
    (call $_ZN5token11getDecimalsEPc
      (get_local $0)
      (get_local $0)
    )
  )
  (func $getSymbol (param $0 i32)
    (call $_ZN5token9getSymbolEPc
      (get_local $0)
      (get_local $0)
    )
  )
  (func $getName (param $0 i32)
    (call $_ZN5token7getNameEPc
      (get_local $0)
      (get_local $0)
    )
  )
)
