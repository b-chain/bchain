(module
  (type $FUNCSIG$vi (func (param i32)))
  (type $FUNCSIG$ii (func (param i32) (result i32)))
  (type $FUNCSIG$iiii (func (param i32 i32 i32) (result i32)))
  (type $FUNCSIG$viiii (func (param i32 i32 i32 i32)))
  (type $FUNCSIG$vii (func (param i32 i32)))
  (import "env" "action_sender" (func $action_sender (param i32)))
  (import "env" "db_get" (func $db_get (param i32 i32 i32) (result i32)))
  (import "env" "db_set" (func $db_set (param i32 i32 i32 i32)))
  (import "env" "isHexAddress" (func $isHexAddress (param i32) (result i32)))
  (import "env" "requireAuth" (func $requireAuth (param i32) (result i32)))
  (import "env" "setResult" (func $setResult (param i32 i32)))
  (import "env" "str2lower" (func $str2lower (param i32)))
  (import "env" "strjoint" (func $strjoint (param i32 i32 i32) (result i32)))
  (table 0 anyfunc)
  (memory $0 1)
  (data (i32.const 4) " @\00\00")
  (data (i32.const 16) "acc\00")
  (export "memory" (memory $0))
  (export "_ZN7account3setEPcS0_i" (func $_ZN7account3setEPcS0_i))
  (export "_ZN7account3getEPc" (func $_ZN7account3getEPc))
  (export "set" (func $set))
  (export "get" (func $get))
  (func $_ZN7account3setEPcS0_i (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
    (local $4 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $4
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 176)
        )
      )
    )
    (call $action_sender
      (i32.add
        (get_local $4)
        (i32.const 128)
      )
    )
    (drop
      (call $requireAuth
        (i32.add
          (get_local $4)
          (i32.const 128)
        )
      )
    )
    (drop
      (call $isHexAddress
        (get_local $1)
      )
    )
    (call $str2lower
      (get_local $1)
    )
    (call $db_set
      (get_local $4)
      (call $strjoint
        (i32.const 16)
        (get_local $1)
        (get_local $4)
      )
      (get_local $2)
      (get_local $3)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $4)
        (i32.const 176)
      )
    )
  )
  (func $_ZN7account3getEPc (param $0 i32) (param $1 i32)
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
    (call $str2lower
      (get_local $1)
    )
    (call $setResult
      (i32.add
        (get_local $2)
        (i32.const 128)
      )
      (call $db_get
        (get_local $2)
        (call $strjoint
          (i32.const 16)
          (get_local $1)
          (get_local $2)
        )
        (i32.add
          (get_local $2)
          (i32.const 128)
        )
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
  (func $set (param $0 i32) (param $1 i32) (param $2 i32)
    (call $_ZN7account3setEPcS0_i
      (get_local $0)
      (get_local $0)
      (get_local $1)
      (get_local $2)
    )
  )
  (func $get (param $0 i32)
    (call $_ZN7account3getEPc
      (get_local $0)
      (get_local $0)
    )
  )
)
