(module
  (type $FUNCSIG$vi (func (param i32)))
  (type $FUNCSIG$ii (func (param i32) (result i32)))
  (type $FUNCSIG$vii (func (param i32 i32)))
  (type $FUNCSIG$iiii (func (param i32 i32 i32) (result i32)))
  (type $FUNCSIG$viiii (func (param i32 i32 i32 i32)))
  (type $FUNCSIG$j (func (result i64)))
  (import "env" "action_sender" (func $action_sender (param i32)))
  (import "env" "assert" (func $assert (param i32 i32)))
  (import "env" "big_add" (func $big_add (param i32 i32 i32) (result i32)))
  (import "env" "big_exp_safe" (func $big_exp_safe (param i32 i32 i32) (result i32)))
  (import "env" "big_mul_safe" (func $big_mul_safe (param i32 i32 i32) (result i32)))
  (import "env" "big_sub_safe" (func $big_sub_safe (param i32 i32 i32) (result i32)))
  (import "env" "block_number" (func $block_number (result i64)))
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
  (data (i32.const 4) "`A\00\00")
  (data (i32.const 16) "symbol len exceed\00")
  (data (i32.const 48) "name len exceed\00")
  (data (i32.const 64) "symbol \00")
  (data (i32.const 80) "isIssue \00")
  (data (i32.const 96) "name \00")
  (data (i32.const 112) "decimals \00")
  (data (i32.const 128) "supply \00")
  (data (i32.const 144) "10\00")
  (data (i32.const 160) "issue memo exceed\00")
  (data (i32.const 192) "issue: symbol is not exist!\00")
  (data (i32.const 224) "issue: token can not issue, illegal!\00")
  (data (i32.const 272) "transfer memo exceed\00")
  (data (i32.const 304) "transfer: symbol is not exist!\00")
  (data (i32.const 336) "action expired\00")
  (export "memory" (memory $0))
  (export "_ZN5token6createEPcS0_S0_S0_iyj" (func $_ZN5token6createEPcS0_S0_S0_iyj))
  (export "_ZN5token5issueEPcS0_S0_yj" (func $_ZN5token5issueEPcS0_S0_yj))
  (export "_ZN5token8transferEPcS0_S0_S0_yj" (func $_ZN5token8transferEPcS0_S0_S0_yj))
  (export "_ZN5token9balanceOfEPcS0_" (func $_ZN5token9balanceOfEPcS0_))
  (export "_ZN5token9getSupplyEPc" (func $_ZN5token9getSupplyEPc))
  (export "_ZN5token11getDecimalsEPc" (func $_ZN5token11getDecimalsEPc))
  (export "_ZN5token9getSymbolEPc" (func $_ZN5token9getSymbolEPc))
  (export "_ZN5token7getNameEPc" (func $_ZN5token7getNameEPc))
  (export "create" (func $create))
  (export "issue" (func $issue))
  (export "transfer" (func $transfer))
  (export "balanceOf" (func $balanceOf))
  (export "getSupply" (func $getSupply))
  (export "getDecimals" (func $getDecimals))
  (export "getSymbol" (func $getSymbol))
  (export "getName" (func $getName))
  (func $_ZN5token6createEPcS0_S0_S0_iyj (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i64) (param $7 i32)
    (local $8 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $8
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 320)
        )
      )
    )
    (i32.store offset=316
      (get_local $8)
      (get_local $5)
    )
    (call $_ZL17blkNumberValidateyj
      (get_local $6)
      (get_local $7)
    )
    (call $action_sender
      (i32.add
        (get_local $8)
        (i32.const 256)
      )
    )
    (drop
      (call $requireAuth
        (i32.add
          (get_local $8)
          (i32.const 256)
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
        (get_local $8)
        (i32.const 128)
      )
      (call $strjoint
        (i32.const 64)
        (get_local $1)
        (i32.add
          (get_local $8)
          (i32.const 128)
        )
      )
      (get_local $1)
      (call $strlen
        (get_local $1)
      )
    )
    (call $db_emplace
      (i32.add
        (get_local $8)
        (i32.const 128)
      )
      (call $strjoint
        (i32.const 80)
        (get_local $1)
        (i32.add
          (get_local $8)
          (i32.const 128)
        )
      )
      (i32.add
        (get_local $8)
        (i32.const 316)
      )
      (i32.const 4)
    )
    (call $db_emplace
      (i32.add
        (get_local $8)
        (i32.const 128)
      )
      (call $strjoint
        (i32.const 96)
        (get_local $1)
        (i32.add
          (get_local $8)
          (i32.const 128)
        )
      )
      (get_local $2)
      (call $strlen
        (get_local $2)
      )
    )
    (call $db_emplace
      (i32.add
        (get_local $8)
        (i32.const 128)
      )
      (call $strjoint
        (i32.const 112)
        (get_local $1)
        (i32.add
          (get_local $8)
          (i32.const 128)
        )
      )
      (get_local $3)
      (call $strlen
        (get_local $3)
      )
    )
    (call $db_emplace
      (i32.add
        (get_local $8)
        (i32.const 128)
      )
      (call $strjoint
        (i32.const 128)
        (get_local $1)
        (i32.add
          (get_local $8)
          (i32.const 128)
        )
      )
      (get_local $4)
      (call $strlen
        (get_local $4)
      )
    )
    (drop
      (call $big_exp_safe
        (i32.const 144)
        (get_local $3)
        (get_local $8)
      )
    )
    (set_local $4
      (call $big_mul_safe
        (get_local $8)
        (get_local $4)
        (get_local $8)
      )
    )
    (call $db_emplace
      (i32.add
        (get_local $8)
        (i32.const 128)
      )
      (call $strjoint
        (get_local $1)
        (i32.add
          (get_local $8)
          (i32.const 256)
        )
        (i32.add
          (get_local $8)
          (i32.const 128)
        )
      )
      (get_local $8)
      (get_local $4)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $8)
        (i32.const 320)
      )
    )
  )
  (func $_ZL17blkNumberValidateyj (param $0 i64) (param $1 i32)
    (local $2 i64)
    (local $3 i32)
    (set_local $3
      (i32.const 0)
    )
    (block $label$0
      (br_if $label$0
        (i64.lt_u
          (tee_local $2
            (call $block_number)
          )
          (get_local $0)
        )
      )
      (set_local $3
        (i64.le_u
          (get_local $2)
          (i64.add
            (i64.extend_u/i32
              (get_local $1)
            )
            (get_local $0)
          )
        )
      )
    )
    (call $assert
      (get_local $3)
      (i32.const 336)
    )
  )
  (func $_ZN5token5issueEPcS0_S0_yj (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i64) (param $5 i32)
    (local $6 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $6
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 704)
        )
      )
    )
    (call $assert
      (i32.lt_s
        (call $strlen
          (get_local $3)
        )
        (i32.const 64)
      )
      (i32.const 160)
    )
    (call $_ZL17blkNumberValidateyj
      (get_local $4)
      (get_local $5)
    )
    (call $action_sender
      (i32.add
        (get_local $6)
        (i32.const 656)
      )
    )
    (drop
      (call $requireAuth
        (i32.add
          (get_local $6)
          (i32.const 656)
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
      (i32.gt_s
        (call $db_get
          (i32.add
            (get_local $6)
            (i32.const 528)
          )
          (call $strjoint
            (i32.const 64)
            (get_local $1)
            (i32.add
              (get_local $6)
              (i32.const 528)
            )
          )
          (get_local $1)
        )
        (i32.const 0)
      )
      (i32.const 192)
    )
    (i32.store offset=524
      (get_local $6)
      (i32.const 0)
    )
    (drop
      (call $db_get
        (i32.add
          (get_local $6)
          (i32.const 528)
        )
        (call $strjoint
          (i32.const 80)
          (get_local $1)
          (i32.add
            (get_local $6)
            (i32.const 528)
          )
        )
        (i32.add
          (get_local $6)
          (i32.const 524)
        )
      )
    )
    (call $assert
      (i32.ne
        (i32.load offset=524
          (get_local $6)
        )
        (i32.const 0)
      )
      (i32.const 224)
    )
    (drop
      (call $db_get
        (i32.add
          (get_local $6)
          (i32.const 528)
        )
        (tee_local $5
          (call $strjoint
            (i32.const 128)
            (get_local $1)
            (i32.add
              (get_local $6)
              (i32.const 528)
            )
          )
        )
        (i32.add
          (get_local $6)
          (i32.const 384)
        )
      )
    )
    (call $db_set
      (i32.add
        (get_local $6)
        (i32.const 528)
      )
      (get_local $5)
      (i32.add
        (get_local $6)
        (i32.const 384)
      )
      (call $big_add
        (i32.add
          (get_local $6)
          (i32.const 384)
        )
        (get_local $2)
        (i32.add
          (get_local $6)
          (i32.const 384)
        )
      )
    )
    (drop
      (call $db_get
        (i32.add
          (get_local $6)
          (i32.const 528)
        )
        (call $strjoint
          (i32.const 112)
          (get_local $1)
          (i32.add
            (get_local $6)
            (i32.const 528)
          )
        )
        (i32.add
          (get_local $6)
          (i32.const 256)
        )
      )
    )
    (drop
      (call $big_exp_safe
        (i32.const 144)
        (i32.add
          (get_local $6)
          (i32.const 256)
        )
        (i32.add
          (get_local $6)
          (i32.const 128)
        )
      )
    )
    (drop
      (call $big_mul_safe
        (i32.add
          (get_local $6)
          (i32.const 128)
        )
        (get_local $2)
        (i32.add
          (get_local $6)
          (i32.const 128)
        )
      )
    )
    (i32.store8
      (i32.add
        (get_local $6)
        (call $db_get
          (i32.add
            (get_local $6)
            (i32.const 528)
          )
          (tee_local $1
            (call $strjoint
              (get_local $1)
              (i32.add
                (get_local $6)
                (i32.const 656)
              )
              (i32.add
                (get_local $6)
                (i32.const 528)
              )
            )
          )
          (get_local $6)
        )
      )
      (i32.const 0)
    )
    (call $db_set
      (i32.add
        (get_local $6)
        (i32.const 528)
      )
      (get_local $1)
      (get_local $6)
      (call $big_add
        (get_local $6)
        (i32.add
          (get_local $6)
          (i32.const 128)
        )
        (get_local $6)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $6)
        (i32.const 704)
      )
    )
  )
  (func $_ZN5token8transferEPcS0_S0_S0_yj (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i64) (param $6 i32)
    (local $7 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $7
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 560)
        )
      )
    )
    (call $assert
      (i32.lt_s
        (call $strlen
          (get_local $4)
        )
        (i32.const 64)
      )
      (i32.const 272)
    )
    (call $_ZL17blkNumberValidateyj
      (get_local $5)
      (get_local $6)
    )
    (call $action_sender
      (i32.add
        (get_local $7)
        (i32.const 512)
      )
    )
    (call $assert
      (i32.gt_s
        (call $db_get
          (i32.add
            (get_local $7)
            (i32.const 384)
          )
          (call $strjoint
            (i32.const 64)
            (get_local $3)
            (i32.add
              (get_local $7)
              (i32.const 384)
            )
          )
          (get_local $3)
        )
        (i32.const 0)
      )
      (i32.const 304)
    )
    (i32.store8
      (i32.add
        (i32.add
          (get_local $7)
          (i32.const 256)
        )
        (call $db_get
          (i32.add
            (get_local $7)
            (i32.const 384)
          )
          (tee_local $6
            (call $strjoint
              (get_local $3)
              (i32.add
                (get_local $7)
                (i32.const 512)
              )
              (i32.add
                (get_local $7)
                (i32.const 384)
              )
            )
          )
          (i32.add
            (get_local $7)
            (i32.const 256)
          )
        )
      )
      (i32.const 0)
    )
    (call $db_set
      (i32.add
        (get_local $7)
        (i32.const 384)
      )
      (get_local $6)
      (i32.add
        (get_local $7)
        (i32.const 256)
      )
      (call $big_sub_safe
        (i32.add
          (get_local $7)
          (i32.const 256)
        )
        (get_local $2)
        (i32.add
          (get_local $7)
          (i32.const 256)
        )
      )
    )
    (call $str2lower
      (get_local $1)
    )
    (i32.store8
      (i32.add
        (i32.add
          (get_local $7)
          (i32.const 128)
        )
        (call $db_get
          (get_local $7)
          (tee_local $3
            (call $strjoint
              (get_local $3)
              (get_local $1)
              (get_local $7)
            )
          )
          (i32.add
            (get_local $7)
            (i32.const 128)
          )
        )
      )
      (i32.const 0)
    )
    (call $db_set
      (get_local $7)
      (get_local $3)
      (i32.add
        (get_local $7)
        (i32.const 128)
      )
      (call $big_add
        (i32.add
          (get_local $7)
          (i32.const 128)
        )
        (get_local $2)
        (i32.add
          (get_local $7)
          (i32.const 128)
        )
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $7)
        (i32.const 560)
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
          (i32.const 256)
        )
      )
    )
    (call $str2lower
      (get_local $1)
    )
    (block $label$0
      (block $label$1
        (br_if $label$1
          (i32.eqz
            (tee_local $1
              (call $db_get
                (get_local $3)
                (call $strjoint
                  (get_local $2)
                  (get_local $1)
                  (get_local $3)
                )
                (i32.add
                  (get_local $3)
                  (i32.const 128)
                )
              )
            )
          )
        )
        (call $setResult
          (i32.add
            (get_local $3)
            (i32.const 128)
          )
          (get_local $1)
        )
        (br $label$0)
      )
      (i32.store8 offset=128
        (get_local $3)
        (i32.const 48)
      )
      (call $setResult
        (i32.add
          (get_local $3)
          (i32.const 128)
        )
        (i32.const 1)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $3)
        (i32.const 256)
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
          (i32.const 256)
        )
      )
    )
    (block $label$0
      (block $label$1
        (br_if $label$1
          (i32.eqz
            (tee_local $1
              (call $db_get
                (get_local $2)
                (call $strjoint
                  (i32.const 128)
                  (get_local $1)
                  (get_local $2)
                )
                (i32.add
                  (get_local $2)
                  (i32.const 128)
                )
              )
            )
          )
        )
        (call $setResult
          (i32.add
            (get_local $2)
            (i32.const 128)
          )
          (get_local $1)
        )
        (br $label$0)
      )
      (i32.store8 offset=128
        (get_local $2)
        (i32.const 48)
      )
      (call $setResult
        (i32.add
          (get_local $2)
          (i32.const 128)
        )
        (i32.const 1)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $2)
        (i32.const 256)
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
          (i32.const 256)
        )
      )
    )
    (block $label$0
      (block $label$1
        (br_if $label$1
          (i32.eqz
            (tee_local $1
              (call $db_get
                (get_local $2)
                (call $strjoint
                  (i32.const 112)
                  (get_local $1)
                  (get_local $2)
                )
                (i32.add
                  (get_local $2)
                  (i32.const 128)
                )
              )
            )
          )
        )
        (call $setResult
          (i32.add
            (get_local $2)
            (i32.const 128)
          )
          (get_local $1)
        )
        (br $label$0)
      )
      (i32.store8 offset=128
        (get_local $2)
        (i32.const 0)
      )
      (call $setResult
        (i32.add
          (get_local $2)
          (i32.const 128)
        )
        (i32.const 1)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $2)
        (i32.const 256)
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
          (i32.const 256)
        )
      )
    )
    (block $label$0
      (block $label$1
        (br_if $label$1
          (i32.eqz
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
          )
        )
        (call $setResult
          (i32.add
            (get_local $2)
            (i32.const 128)
          )
          (get_local $1)
        )
        (br $label$0)
      )
      (i32.store8 offset=128
        (get_local $2)
        (i32.const 0)
      )
      (call $setResult
        (i32.add
          (get_local $2)
          (i32.const 128)
        )
        (i32.const 1)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $2)
        (i32.const 256)
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
          (i32.const 256)
        )
      )
    )
    (block $label$0
      (block $label$1
        (br_if $label$1
          (i32.eqz
            (tee_local $1
              (call $db_get
                (get_local $2)
                (call $strjoint
                  (i32.const 96)
                  (get_local $1)
                  (get_local $2)
                )
                (i32.add
                  (get_local $2)
                  (i32.const 128)
                )
              )
            )
          )
        )
        (call $setResult
          (i32.add
            (get_local $2)
            (i32.const 128)
          )
          (get_local $1)
        )
        (br $label$0)
      )
      (i32.store8 offset=128
        (get_local $2)
        (i32.const 0)
      )
      (call $setResult
        (i32.add
          (get_local $2)
          (i32.const 128)
        )
        (i32.const 1)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $2)
        (i32.const 256)
      )
    )
  )
  (func $create (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i64) (param $6 i32)
    (call $_ZN5token6createEPcS0_S0_S0_iyj
      (get_local $0)
      (get_local $0)
      (get_local $1)
      (get_local $2)
      (get_local $3)
      (get_local $4)
      (get_local $5)
      (get_local $6)
    )
  )
  (func $issue (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i64) (param $4 i32)
    (call $_ZN5token5issueEPcS0_S0_yj
      (get_local $0)
      (get_local $0)
      (get_local $1)
      (get_local $2)
      (get_local $3)
      (get_local $4)
    )
  )
  (func $transfer (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i64) (param $5 i32)
    (call $_ZN5token8transferEPcS0_S0_S0_yj
      (get_local $0)
      (get_local $0)
      (get_local $1)
      (get_local $2)
      (get_local $3)
      (get_local $4)
      (get_local $5)
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
