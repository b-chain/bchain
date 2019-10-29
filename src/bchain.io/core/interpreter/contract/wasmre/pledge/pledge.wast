(module
  (type $FUNCSIG$vi (func (param i32)))
  (type $FUNCSIG$ii (func (param i32) (result i32)))
  (type $FUNCSIG$viii (func (param i32 i32 i32)))
  (type $FUNCSIG$iiii (func (param i32 i32 i32) (result i32)))
  (type $FUNCSIG$vii (func (param i32 i32)))
  (type $FUNCSIG$jj (func (param i64) (result i64)))
  (type $FUNCSIG$viiii (func (param i32 i32 i32 i32)))
  (import "env" "action_callWithPara" (func $action_callWithPara (param i32 i32 i32)))
  (import "env" "action_sender" (func $action_sender (param i32)))
  (import "env" "assert" (func $assert (param i32 i32)))
  (import "env" "contract_address" (func $contract_address (param i32)))
  (import "env" "contract_callWithPara" (func $contract_callWithPara (param i32 i32 i32)))
  (import "env" "db_get" (func $db_get (param i32 i32 i32) (result i32)))
  (import "env" "db_set" (func $db_set (param i32 i32 i32 i32)))
  (import "env" "getWeight" (func $getWeight (param i64) (result i64)))
  (import "env" "setResult" (func $setResult (param i32 i32)))
  (import "env" "str2lower" (func $str2lower (param i32)))
  (import "env" "strjoint" (func $strjoint (param i32 i32 i32) (result i32)))
  (import "env" "strlen" (func $strlen (param i32) (result i32)))
  (table 0 anyfunc)
  (memory $0 1)
  (data (i32.const 4) " A\00\00")
  (data (i32.const 16) "pledge\00")
  (data (i32.const 32) "transfer\00")
  (data (i32.const 48) "redeem\00")
  (data (i32.const 64) "0xb78f12Cb3924607A8BC6a66799e159E3459097e9\00")
  (data (i32.const 112) "reeem: pledge amount is not enough\00")
  (data (i32.const 160) "subTotalPledge amount >= tPledge\00")
  (data (i32.const 208) "subTotalWeight weight >= totalWeight\00")
  (data (i32.const 256) "_pledgeTotal\00")
  (data (i32.const 272) "_totalWeight\00")
  (export "memory" (memory $0))
  (export "_ZN10pledgePool6pledgeEy" (func $_ZN10pledgePool6pledgeEy))
  (export "_ZN10pledgePool6redeemEy" (func $_ZN10pledgePool6redeemEy))
  (export "_ZN10pledgePool8pledgeOfEPc" (func $_ZN10pledgePool8pledgeOfEPc))
  (export "_ZN10pledgePool11pledgeOfExtEPc" (func $_ZN10pledgePool11pledgeOfExtEPc))
  (export "pledge" (func $pledge))
  (export "redeem" (func $redeem))
  (export "pledgeOf" (func $pledgeOf))
  (export "pledgeOfExt" (func $pledgeOfExt))
  (func $_ZN10pledgePool6pledgeEy (param $0 i32) (param $1 i64)
    (local $2 i32)
    (local $3 i64)
    (local $4 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $4
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 272)
        )
      )
    )
    (i64.store offset=136
      (get_local $4)
      (get_local $1)
    )
    (call $action_sender
      (i32.add
        (get_local $4)
        (i32.const 80)
      )
    )
    (call $contract_address
      (i32.add
        (get_local $4)
        (i32.const 32)
      )
    )
    (i64.store
      (get_local $4)
      (i64.const 0)
    )
    (drop
      (call $db_get
        (i32.add
          (get_local $4)
          (i32.const 144)
        )
        (tee_local $2
          (call $strjoint
            (i32.const 16)
            (i32.add
              (get_local $4)
              (i32.const 80)
            )
            (i32.add
              (get_local $4)
              (i32.const 144)
            )
          )
        )
        (get_local $4)
      )
    )
    (set_local $3
      (call $getWeight
        (i64.load
          (get_local $4)
        )
      )
    )
    (i64.store
      (get_local $4)
      (i64.add
        (i64.load
          (get_local $4)
        )
        (get_local $1)
      )
    )
    (call $db_set
      (i32.add
        (get_local $4)
        (i32.const 144)
      )
      (get_local $2)
      (get_local $4)
      (i32.const 8)
    )
    (i64.store offset=16
      (get_local $4)
      (i64.const 0)
    )
    (drop
      (call $db_get
        (i32.const 256)
        (tee_local $2
          (call $strlen
            (i32.const 256)
          )
        )
        (i32.add
          (get_local $4)
          (i32.const 16)
        )
      )
    )
    (i64.store offset=16
      (get_local $4)
      (i64.add
        (i64.load offset=16
          (get_local $4)
        )
        (get_local $1)
      )
    )
    (call $db_set
      (i32.const 256)
      (get_local $2)
      (i32.add
        (get_local $4)
        (i32.const 16)
      )
      (i32.const 8)
    )
    (set_local $1
      (call $getWeight
        (i64.load
          (get_local $4)
        )
      )
    )
    (i64.store offset=16
      (get_local $4)
      (i64.const 0)
    )
    (drop
      (call $db_get
        (i32.const 272)
        (tee_local $2
          (call $strlen
            (i32.const 272)
          )
        )
        (i32.add
          (get_local $4)
          (i32.const 16)
        )
      )
    )
    (i64.store offset=16
      (get_local $4)
      (i64.add
        (i64.sub
          (get_local $1)
          (get_local $3)
        )
        (i64.load offset=16
          (get_local $4)
        )
      )
    )
    (call $db_set
      (i32.const 272)
      (get_local $2)
      (i32.add
        (get_local $4)
        (i32.const 16)
      )
      (i32.const 8)
    )
    (i32.store offset=144
      (get_local $4)
      (i32.const 4)
    )
    (i32.store offset=152
      (get_local $4)
      (i32.const 43)
    )
    (i32.store offset=148
      (get_local $4)
      (i32.add
        (get_local $4)
        (i32.const 32)
      )
    )
    (i32.store offset=24
      (get_local $4)
      (i32.const 8)
    )
    (i32.store offset=20
      (get_local $4)
      (i32.add
        (get_local $4)
        (i32.const 136)
      )
    )
    (i32.store offset=16
      (get_local $4)
      (i32.const 1)
    )
    (i32.store offset=4
      (get_local $4)
      (i32.const 16)
    )
    (i32.store
      (get_local $4)
      (i32.const 4)
    )
    (set_local $2
      (call $strlen
        (i32.const 16)
      )
    )
    (i32.store offset=12
      (get_local $4)
      (i32.const 0)
    )
    (i32.store offset=8
      (get_local $4)
      (i32.add
        (get_local $2)
        (i32.const 1)
      )
    )
    (i32.store offset=28
      (get_local $4)
      (get_local $4)
    )
    (i32.store offset=156
      (get_local $4)
      (i32.add
        (get_local $4)
        (i32.const 16)
      )
    )
    (call $action_callWithPara
      (i32.const 64)
      (i32.const 32)
      (i32.add
        (get_local $4)
        (i32.const 144)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $4)
        (i32.const 272)
      )
    )
  )
  (func $_ZN10pledgePool6redeemEy (param $0 i32) (param $1 i64)
    (local $2 i32)
    (local $3 i64)
    (local $4 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $4
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 224)
        )
      )
    )
    (i64.store offset=88
      (get_local $4)
      (get_local $1)
    )
    (call $action_sender
      (i32.add
        (get_local $4)
        (i32.const 32)
      )
    )
    (i64.store
      (get_local $4)
      (i64.const 0)
    )
    (drop
      (call $db_get
        (i32.add
          (get_local $4)
          (i32.const 96)
        )
        (tee_local $2
          (call $strjoint
            (i32.const 16)
            (i32.add
              (get_local $4)
              (i32.const 32)
            )
            (i32.add
              (get_local $4)
              (i32.const 96)
            )
          )
        )
        (get_local $4)
      )
    )
    (call $assert
      (i64.ge_u
        (i64.load
          (get_local $4)
        )
        (get_local $1)
      )
      (i32.const 112)
    )
    (set_local $3
      (call $getWeight
        (i64.load
          (get_local $4)
        )
      )
    )
    (i64.store
      (get_local $4)
      (i64.sub
        (i64.load
          (get_local $4)
        )
        (get_local $1)
      )
    )
    (call $db_set
      (i32.add
        (get_local $4)
        (i32.const 96)
      )
      (get_local $2)
      (get_local $4)
      (i32.const 8)
    )
    (i64.store offset=16
      (get_local $4)
      (i64.const 0)
    )
    (drop
      (call $db_get
        (i32.const 256)
        (tee_local $2
          (call $strlen
            (i32.const 256)
          )
        )
        (i32.add
          (get_local $4)
          (i32.const 16)
        )
      )
    )
    (call $assert
      (i64.gt_u
        (i64.load offset=16
          (get_local $4)
        )
        (get_local $1)
      )
      (i32.const 160)
    )
    (i64.store offset=16
      (get_local $4)
      (i64.sub
        (i64.load offset=16
          (get_local $4)
        )
        (get_local $1)
      )
    )
    (call $db_set
      (i32.const 256)
      (get_local $2)
      (i32.add
        (get_local $4)
        (i32.const 16)
      )
      (i32.const 8)
    )
    (set_local $1
      (call $getWeight
        (i64.load
          (get_local $4)
        )
      )
    )
    (i64.store offset=16
      (get_local $4)
      (i64.const 0)
    )
    (drop
      (call $db_get
        (i32.const 272)
        (tee_local $2
          (call $strlen
            (i32.const 272)
          )
        )
        (i32.add
          (get_local $4)
          (i32.const 16)
        )
      )
    )
    (call $assert
      (i64.gt_u
        (i64.load offset=16
          (get_local $4)
        )
        (tee_local $1
          (i64.sub
            (get_local $3)
            (get_local $1)
          )
        )
      )
      (i32.const 208)
    )
    (i64.store offset=16
      (get_local $4)
      (i64.sub
        (i64.load offset=16
          (get_local $4)
        )
        (get_local $1)
      )
    )
    (call $db_set
      (i32.const 272)
      (get_local $2)
      (i32.add
        (get_local $4)
        (i32.const 16)
      )
      (i32.const 8)
    )
    (i32.store offset=96
      (get_local $4)
      (i32.const 4)
    )
    (i32.store offset=104
      (get_local $4)
      (i32.const 43)
    )
    (i32.store offset=100
      (get_local $4)
      (i32.add
        (get_local $4)
        (i32.const 32)
      )
    )
    (i32.store offset=24
      (get_local $4)
      (i32.const 8)
    )
    (i32.store offset=20
      (get_local $4)
      (i32.add
        (get_local $4)
        (i32.const 88)
      )
    )
    (i32.store offset=16
      (get_local $4)
      (i32.const 1)
    )
    (i32.store
      (get_local $4)
      (i32.const 4)
    )
    (i32.store offset=4
      (get_local $4)
      (i32.const 48)
    )
    (set_local $2
      (call $strlen
        (i32.const 48)
      )
    )
    (i32.store offset=12
      (get_local $4)
      (i32.const 0)
    )
    (i32.store offset=8
      (get_local $4)
      (i32.add
        (get_local $2)
        (i32.const 1)
      )
    )
    (i32.store offset=28
      (get_local $4)
      (get_local $4)
    )
    (i32.store offset=108
      (get_local $4)
      (i32.add
        (get_local $4)
        (i32.const 16)
      )
    )
    (call $contract_callWithPara
      (i32.const 64)
      (i32.const 32)
      (i32.add
        (get_local $4)
        (i32.const 96)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $4)
        (i32.const 224)
      )
    )
  )
  (func $_ZN10pledgePool8pledgeOfEPc (param $0 i32) (param $1 i32)
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
    (call $str2lower
      (get_local $1)
    )
    (drop
      (call $db_get
        (get_local $2)
        (call $strjoint
          (i32.const 16)
          (get_local $1)
          (get_local $2)
        )
        (i32.add
          (get_local $2)
          (i32.const 136)
        )
      )
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
  (func $_ZN10pledgePool11pledgeOfExtEPc (param $0 i32) (param $1 i32)
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
    (call $str2lower
      (get_local $1)
    )
    (drop
      (call $db_get
        (get_local $2)
        (call $strjoint
          (i32.const 16)
          (get_local $1)
          (get_local $2)
        )
        (i32.add
          (get_local $2)
          (i32.const 136)
        )
      )
    )
    (call $setResult
      (i32.add
        (get_local $2)
        (i32.const 136)
      )
      (i32.const 8)
    )
    (i64.store offset=136
      (get_local $2)
      (i64.const 0)
    )
    (drop
      (call $db_get
        (i32.const 256)
        (call $strlen
          (i32.const 256)
        )
        (i32.add
          (get_local $2)
          (i32.const 136)
        )
      )
    )
    (call $setResult
      (i32.add
        (get_local $2)
        (i32.const 136)
      )
      (i32.const 8)
    )
    (i64.store offset=136
      (get_local $2)
      (i64.const 0)
    )
    (drop
      (call $db_get
        (i32.const 272)
        (call $strlen
          (i32.const 272)
        )
        (i32.add
          (get_local $2)
          (i32.const 136)
        )
      )
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
  (func $pledge (param $0 i64)
    (local $1 i32)
    (call $_ZN10pledgePool6pledgeEy
      (get_local $1)
      (get_local $0)
    )
  )
  (func $redeem (param $0 i64)
    (local $1 i32)
    (call $_ZN10pledgePool6redeemEy
      (get_local $1)
      (get_local $0)
    )
  )
  (func $pledgeOf (param $0 i32)
    (call $_ZN10pledgePool8pledgeOfEPc
      (get_local $0)
      (get_local $0)
    )
  )
  (func $pledgeOfExt (param $0 i32)
    (call $_ZN10pledgePool11pledgeOfExtEPc
      (get_local $0)
      (get_local $0)
    )
  )
)
