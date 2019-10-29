(module
  (type $FUNCSIG$vi (func (param i32)))
  (type $FUNCSIG$viii (func (param i32 i32 i32)))
  (type $FUNCSIG$ii (func (param i32) (result i32)))
  (type $FUNCSIG$vii (func (param i32 i32)))
  (import "env" "action_sender" (func $action_sender (param i32)))
  (import "env" "contract_create" (func $contract_create (param i32 i32 i32)))
  (import "env" "emitEvent" (func $emitEvent (param i32 i32)))
  (import "env" "setResult" (func $setResult (param i32 i32)))
  (import "env" "strlen" (func $strlen (param i32) (result i32)))
  (table 0 anyfunc)
  (memory $0 1)
  (data (i32.const 4) " @\00\00")
  (data (i32.const 16) "createContract\00")
  (export "memory" (memory $0))
  (export "_ZN6system14createContractEPc" (func $_ZN6system14createContractEPc))
  (export "createContract" (func $createContract))
  (func $_ZN6system14createContractEPc (param $0 i32) (param $1 i32)
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
    (call $action_sender
      (i32.add
        (get_local $2)
        (i32.const 48)
      )
    )
    (call $contract_create
      (i32.add
        (get_local $2)
        (i32.const 48)
      )
      (get_local $1)
      (i32.add
        (get_local $2)
        (i32.const 96)
      )
    )
    (set_local $1
      (call $strlen
        (i32.add
          (get_local $2)
          (i32.const 96)
        )
      )
    )
    (i32.store offset=32
      (get_local $2)
      (i32.const 16)
    )
    (i32.store offset=36
      (get_local $2)
      (call $strlen
        (i32.const 16)
      )
    )
    (i32.store offset=24
      (get_local $2)
      (i32.const 0)
    )
    (i32.store offset=20
      (get_local $2)
      (get_local $1)
    )
    (i32.store offset=16
      (get_local $2)
      (i32.add
        (get_local $2)
        (i32.const 96)
      )
    )
    (i32.store offset=40
      (get_local $2)
      (i32.add
        (get_local $2)
        (i32.const 16)
      )
    )
    (i32.store offset=8
      (get_local $2)
      (i32.const 0)
    )
    (i32.store offset=4
      (get_local $2)
      (get_local $1)
    )
    (i32.store
      (get_local $2)
      (i32.add
        (get_local $2)
        (i32.const 96)
      )
    )
    (call $emitEvent
      (i32.add
        (get_local $2)
        (i32.const 32)
      )
      (get_local $2)
    )
    (call $setResult
      (i32.add
        (get_local $2)
        (i32.const 96)
      )
      (get_local $1)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $2)
        (i32.const 144)
      )
    )
  )
  (func $createContract (param $0 i32)
    (call $_ZN6system14createContractEPc
      (get_local $0)
      (get_local $0)
    )
  )
)
