(module
  (type $FUNCSIG$viiii (func (param i32 i32 i32 i32)))
  (type $FUNCSIG$vi (func (param i32)))
  (type $FUNCSIG$iiii (func (param i32 i32 i32) (result i32)))
  (type $FUNCSIG$vii (func (param i32 i32)))
  (type $FUNCSIG$ii (func (param i32) (result i32)))
  (type $FUNCSIG$viii (func (param i32 i32 i32)))
  (type $FUNCSIG$j (func (result i64)))
  (import "env" "action_sender" (func $action_sender (param i32)))
  (import "env" "assert" (func $assert (param i32 i32)))
  (import "env" "block_number" (func $block_number (result i64)))
  (import "env" "db_emplace" (func $db_emplace (param i32 i32 i32 i32)))
  (import "env" "db_get" (func $db_get (param i32 i32 i32) (result i32)))
  (import "env" "db_set" (func $db_set (param i32 i32 i32 i32)))
  (import "env" "emitEvent" (func $emitEvent (param i32 i32)))
  (import "env" "isHexAddress" (func $isHexAddress (param i32) (result i32)))
  (import "env" "log" (func $log (param i32) (result i32)))
  (import "env" "memcmp" (func $memcmp (param i32 i32 i32) (result i32)))
  (import "env" "memcpy" (func $memcpy (param i32 i32 i32) (result i32)))
  (import "env" "memset" (func $memset (param i32 i32 i32) (result i32)))
  (import "env" "sha1" (func $sha1 (param i32 i32 i32)))
  (import "env" "sha256" (func $sha256 (param i32 i32 i32)))
  (import "env" "sha512" (func $sha512 (param i32 i32 i32)))
  (table 0 anyfunc)
  (memory $0 1)
  (data (i32.const 4) "@A\00\00")
  (data (i32.const 16) "symol\00")
  (data (i32.const 32) "name\00")
  (data (i32.const 48) "error\n\00")
  (data (i32.const 64) "xxb\00")
  (data (i32.const 80) "xxx\00")
  (data (i32.const 96) "hello world!\n\00")
  (data (i32.const 112) "0x2e68b0583021d78c122f719fc82036529a90571d\00")
  (data (i32.const 160) "0x2e68b0583021d78c122f719fc82036529a90571d is a valid hex address\n\00")
  (data (i32.const 240) "test assert\00")
  (data (i32.const 256) "abcdefghi\00")
  (data (i32.const 272) "sha1 ret: \n\00")
  (data (i32.const 288) "sha256 ret: \n\00")
  (data (i32.const 304) "sha512 ret: \n\00")
  (export "memory" (memory $0))
  (export "_Z6strlenPc" (func $_Z6strlenPc))
  (export "_ZN5token6createEPcS0_ii" (func $_ZN5token6createEPcS0_ii))
  (export "_ZN5token7transerEPci" (func $_ZN5token7transerEPci))
  (export "_ZN5token9balenceOfEPc" (func $_ZN5token9balenceOfEPc))
  (export "create" (func $create))
  (export "transer" (func $transer))
  (export "balenceOf" (func $balenceOf))
  (export "test1" (func $test1))
  (export "testAssert1" (func $testAssert1))
  (export "testAssert2" (func $testAssert2))
  (export "testCryotoApi" (func $testCryotoApi))
  (export "testStack" (func $testStack))
  (export "varPara" (func $varPara))
  (export "testVarPara" (func $testVarPara))
  (export "testEvent" (func $testEvent))
  (export "testBlockNumber" (func $testBlockNumber))
  (export "memTest" (func $memTest))
  (export "loopForever" (func $loopForever))
  (func $_Z6strlenPc (param $0 i32) (result i32)
    (local $1 i32)
    (local $2 i32)
    (local $3 i32)
    (block $label$0
      (br_if $label$0
        (i32.eqz
          (i32.load8_u
            (get_local $0)
          )
        )
      )
      (set_local $1
        (i32.add
          (get_local $0)
          (i32.const 1)
        )
      )
      (set_local $0
        (i32.const 0)
      )
      (loop $label$1
        (set_local $2
          (i32.add
            (get_local $1)
            (get_local $0)
          )
        )
        (set_local $0
          (tee_local $3
            (i32.add
              (get_local $0)
              (i32.const 1)
            )
          )
        )
        (br_if $label$1
          (i32.load8_u
            (get_local $2)
          )
        )
      )
      (return
        (get_local $3)
      )
    )
    (i32.const 0)
  )
  (func $_ZN5token6createEPcS0_ii (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    (local $7 i32)
    (local $8 i32)
    (local $9 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $9
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 64)
        )
      )
    )
    (set_local $8
      (i32.const 0)
    )
    (set_local $7
      (i32.const 0)
    )
    (block $label$0
      (br_if $label$0
        (i32.eqz
          (i32.load8_u
            (get_local $0)
          )
        )
      )
      (set_local $4
        (i32.add
          (get_local $0)
          (i32.const 1)
        )
      )
      (set_local $6
        (i32.const 0)
      )
      (loop $label$1
        (set_local $5
          (i32.add
            (get_local $4)
            (get_local $6)
          )
        )
        (set_local $6
          (tee_local $7
            (i32.add
              (get_local $6)
              (i32.const 1)
            )
          )
        )
        (br_if $label$1
          (i32.load8_u
            (get_local $5)
          )
        )
      )
    )
    (call $db_emplace
      (i32.const 16)
      (i32.const 5)
      (get_local $0)
      (get_local $7)
    )
    (block $label$2
      (br_if $label$2
        (i32.eqz
          (i32.load8_u
            (get_local $1)
          )
        )
      )
      (set_local $7
        (i32.add
          (get_local $1)
          (i32.const 1)
        )
      )
      (set_local $6
        (i32.const 0)
      )
      (loop $label$3
        (set_local $5
          (i32.add
            (get_local $7)
            (get_local $6)
          )
        )
        (set_local $6
          (tee_local $8
            (i32.add
              (get_local $6)
              (i32.const 1)
            )
          )
        )
        (br_if $label$3
          (i32.load8_u
            (get_local $5)
          )
        )
      )
    )
    (call $db_emplace
      (i32.const 32)
      (i32.const 4)
      (get_local $1)
      (get_local $8)
    )
    (i32.store offset=60
      (get_local $9)
      (get_local $3)
    )
    (block $label$4
      (br_if $label$4
        (i32.lt_s
          (get_local $2)
          (i32.const 1)
        )
      )
      (loop $label$5
        (set_local $3
          (i32.mul
            (get_local $3)
            (i32.const 10)
          )
        )
        (br_if $label$5
          (tee_local $2
            (i32.add
              (get_local $2)
              (i32.const -1)
            )
          )
        )
      )
      (i32.store offset=60
        (get_local $9)
        (get_local $3)
      )
    )
    (call $action_sender
      (get_local $9)
    )
    (block $label$6
      (block $label$7
        (br_if $label$7
          (i32.eqz
            (i32.load8_u
              (get_local $9)
            )
          )
        )
        (set_local $5
          (i32.or
            (get_local $9)
            (i32.const 1)
          )
        )
        (set_local $6
          (i32.const 0)
        )
        (loop $label$8
          (set_local $3
            (i32.add
              (get_local $5)
              (get_local $6)
            )
          )
          (set_local $6
            (tee_local $2
              (i32.add
                (get_local $6)
                (i32.const 1)
              )
            )
          )
          (br_if $label$8
            (i32.load8_u
              (get_local $3)
            )
          )
          (br $label$6)
        )
      )
      (set_local $2
        (i32.const 0)
      )
    )
    (call $db_emplace
      (get_local $9)
      (get_local $2)
      (i32.add
        (get_local $9)
        (i32.const 60)
      )
      (i32.const 4)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $9)
        (i32.const 64)
      )
    )
  )
  (func $_ZN5token7transerEPci (param $0 i32) (param $1 i32) (param $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    (local $7 i32)
    (local $8 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $8
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 64)
        )
      )
    )
    (call $action_sender
      (i32.add
        (get_local $8)
        (i32.const 16)
      )
    )
    (set_local $7
      (i32.const 0)
    )
    (i32.store offset=12
      (get_local $8)
      (i32.const 0)
    )
    (i32.store offset=8
      (get_local $8)
      (i32.const 0)
    )
    (set_local $6
      (i32.const 0)
    )
    (block $label$0
      (br_if $label$0
        (i32.eqz
          (i32.load8_u offset=16
            (get_local $8)
          )
        )
      )
      (set_local $3
        (i32.or
          (i32.add
            (get_local $8)
            (i32.const 16)
          )
          (i32.const 1)
        )
      )
      (set_local $5
        (i32.const 0)
      )
      (loop $label$1
        (set_local $4
          (i32.add
            (get_local $3)
            (get_local $5)
          )
        )
        (set_local $5
          (tee_local $6
            (i32.add
              (get_local $5)
              (i32.const 1)
            )
          )
        )
        (br_if $label$1
          (i32.load8_u
            (get_local $4)
          )
        )
      )
    )
    (call $assert
      (i32.eq
        (call $db_get
          (i32.add
            (get_local $8)
            (i32.const 16)
          )
          (get_local $6)
          (i32.add
            (get_local $8)
            (i32.const 12)
          )
        )
        (i32.const 4)
      )
      (i32.const 48)
    )
    (block $label$2
      (br_if $label$2
        (i32.eqz
          (i32.load8_u
            (get_local $1)
          )
        )
      )
      (set_local $6
        (i32.add
          (get_local $1)
          (i32.const 1)
        )
      )
      (set_local $5
        (i32.const 0)
      )
      (loop $label$3
        (set_local $4
          (i32.add
            (get_local $6)
            (get_local $5)
          )
        )
        (set_local $5
          (tee_local $7
            (i32.add
              (get_local $5)
              (i32.const 1)
            )
          )
        )
        (br_if $label$3
          (i32.load8_u
            (get_local $4)
          )
        )
      )
    )
    (drop
      (call $db_get
        (get_local $1)
        (get_local $7)
        (i32.add
          (get_local $8)
          (i32.const 8)
        )
      )
    )
    (i32.store offset=12
      (get_local $8)
      (i32.sub
        (i32.load offset=12
          (get_local $8)
        )
        (get_local $2)
      )
    )
    (i32.store offset=8
      (get_local $8)
      (i32.add
        (i32.load offset=8
          (get_local $8)
        )
        (get_local $2)
      )
    )
    (set_local $7
      (i32.const 0)
    )
    (set_local $6
      (i32.const 0)
    )
    (block $label$4
      (br_if $label$4
        (i32.eqz
          (i32.load8_u offset=16
            (get_local $8)
          )
        )
      )
      (set_local $3
        (i32.or
          (i32.add
            (get_local $8)
            (i32.const 16)
          )
          (i32.const 1)
        )
      )
      (set_local $5
        (i32.const 0)
      )
      (loop $label$5
        (set_local $4
          (i32.add
            (get_local $3)
            (get_local $5)
          )
        )
        (set_local $5
          (tee_local $6
            (i32.add
              (get_local $5)
              (i32.const 1)
            )
          )
        )
        (br_if $label$5
          (i32.load8_u
            (get_local $4)
          )
        )
      )
    )
    (call $db_set
      (i32.add
        (get_local $8)
        (i32.const 16)
      )
      (get_local $6)
      (i32.add
        (get_local $8)
        (i32.const 12)
      )
      (i32.const 4)
    )
    (block $label$6
      (br_if $label$6
        (i32.eqz
          (i32.load8_u
            (get_local $1)
          )
        )
      )
      (set_local $6
        (i32.add
          (get_local $1)
          (i32.const 1)
        )
      )
      (set_local $5
        (i32.const 0)
      )
      (loop $label$7
        (set_local $4
          (i32.add
            (get_local $6)
            (get_local $5)
          )
        )
        (set_local $5
          (tee_local $7
            (i32.add
              (get_local $5)
              (i32.const 1)
            )
          )
        )
        (br_if $label$7
          (i32.load8_u
            (get_local $4)
          )
        )
      )
    )
    (call $db_set
      (get_local $1)
      (get_local $7)
      (i32.add
        (get_local $8)
        (i32.const 8)
      )
      (i32.const 4)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $8)
        (i32.const 64)
      )
    )
  )
  (func $_ZN5token9balenceOfEPc (param $0 i32) (param $1 i32) (result i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $6
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 16)
        )
      )
    )
    (set_local $5
      (i32.const 0)
    )
    (i32.store offset=12
      (get_local $6)
      (i32.const 0)
    )
    (block $label$0
      (br_if $label$0
        (i32.eqz
          (i32.load8_u
            (get_local $1)
          )
        )
      )
      (set_local $2
        (i32.add
          (get_local $1)
          (i32.const 1)
        )
      )
      (set_local $4
        (i32.const 0)
      )
      (loop $label$1
        (set_local $3
          (i32.add
            (get_local $2)
            (get_local $4)
          )
        )
        (set_local $4
          (tee_local $5
            (i32.add
              (get_local $4)
              (i32.const 1)
            )
          )
        )
        (br_if $label$1
          (i32.load8_u
            (get_local $3)
          )
        )
      )
    )
    (drop
      (call $db_get
        (get_local $1)
        (get_local $5)
        (i32.add
          (get_local $6)
          (i32.const 12)
        )
      )
    )
    (set_local $4
      (i32.load offset=12
        (get_local $6)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $6)
        (i32.const 16)
      )
    )
    (get_local $4)
  )
  (func $create
    (call $_ZN5token6createEPcS0_ii
      (i32.const 64)
      (i32.const 80)
      (i32.const 2)
      (i32.const 1000)
    )
  )
  (func $transer (param $0 i32) (param $1 i32)
    (call $_ZN5token7transerEPci
      (get_local $0)
      (get_local $0)
      (get_local $1)
    )
  )
  (func $balenceOf (param $0 i32) (result i32)
    (local $1 i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $5
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 16)
        )
      )
    )
    (set_local $4
      (i32.const 0)
    )
    (i32.store offset=12
      (get_local $5)
      (i32.const 0)
    )
    (block $label$0
      (br_if $label$0
        (i32.eqz
          (i32.load8_u
            (get_local $0)
          )
        )
      )
      (set_local $1
        (i32.add
          (get_local $0)
          (i32.const 1)
        )
      )
      (set_local $3
        (i32.const 0)
      )
      (loop $label$1
        (set_local $2
          (i32.add
            (get_local $1)
            (get_local $3)
          )
        )
        (set_local $3
          (tee_local $4
            (i32.add
              (get_local $3)
              (i32.const 1)
            )
          )
        )
        (br_if $label$1
          (i32.load8_u
            (get_local $2)
          )
        )
      )
    )
    (drop
      (call $db_get
        (get_local $0)
        (get_local $4)
        (i32.add
          (get_local $5)
          (i32.const 12)
        )
      )
    )
    (set_local $3
      (i32.load offset=12
        (get_local $5)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $5)
        (i32.const 16)
      )
    )
    (get_local $3)
  )
  (func $test1 (param $0 i32) (result i32)
    (local $1 i32)
    (drop
      (call $log
        (i32.const 96)
      )
    )
    (set_local $1
      (i32.const 3)
    )
    (block $label$0
      (br_if $label$0
        (i32.eqz
          (call $isHexAddress
            (i32.const 112)
          )
        )
      )
      (drop
        (call $log
          (i32.const 160)
        )
      )
      (set_local $1
        (i32.const 1)
      )
    )
    (i32.add
      (get_local $1)
      (get_local $0)
    )
  )
  (func $testAssert1 (result i32)
    (call $assert
      (i32.const 1)
      (i32.const 240)
    )
    (i32.const 0)
  )
  (func $testAssert2
    (call $assert
      (i32.const 0)
      (i32.const 240)
    )
  )
  (func $testCryotoApi
    (local $0 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $0
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 272)
        )
      )
    )
    (call $sha1
      (i32.const 256)
      (i32.const 9)
      (i32.add
        (get_local $0)
        (i32.const 224)
      )
    )
    (drop
      (call $log
        (i32.const 272)
      )
    )
    (i32.store16 offset=266
      (get_local $0)
      (i32.const 10)
    )
    (drop
      (call $log
        (i32.add
          (get_local $0)
          (i32.const 224)
        )
      )
    )
    (call $sha256
      (i32.const 256)
      (i32.const 9)
      (i32.add
        (get_local $0)
        (i32.const 144)
      )
    )
    (drop
      (call $log
        (i32.const 288)
      )
    )
    (i32.store16 offset=208
      (get_local $0)
      (i32.const 10)
    )
    (drop
      (call $log
        (i32.add
          (get_local $0)
          (i32.const 144)
        )
      )
    )
    (call $sha512
      (i32.const 256)
      (i32.const 9)
      (get_local $0)
    )
    (drop
      (call $log
        (i32.const 304)
      )
    )
    (i32.store16 offset=128
      (get_local $0)
      (i32.const 10)
    )
    (drop
      (call $log
        (get_local $0)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $0)
        (i32.const 272)
      )
    )
  )
  (func $testStack
    (local $0 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $0
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 128)
        )
      )
    )
    (call $sha1
      (i32.const 256)
      (i32.const 9)
      (tee_local $0
        (call $memset
          (get_local $0)
          (i32.const 0)
          (i32.const 44)
        )
      )
    )
    (drop
      (call $log
        (i32.const 272)
      )
    )
    (i32.store16 offset=42
      (get_local $0)
      (i32.const 10)
    )
    (drop
      (call $log
        (get_local $0)
      )
    )
    (call $sha256
      (i32.const 256)
      (i32.const 9)
      (i32.add
        (get_local $0)
        (i32.const 48)
      )
    )
    (drop
      (call $log
        (i32.const 288)
      )
    )
    (i32.store16 offset=112
      (get_local $0)
      (i32.const 10)
    )
    (drop
      (call $log
        (i32.add
          (get_local $0)
          (i32.const 48)
        )
      )
    )
    (drop
      (call $log
        (i32.const 272)
      )
    )
    (drop
      (call $log
        (get_local $0)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $0)
        (i32.const 128)
      )
    )
  )
  (func $varPara (param $0 i32) (param $1 i32) (result i32)
    (i32.const 1)
  )
  (func $testVarPara (result i32)
    (unreachable)
  )
  (func $testEvent
    (local $0 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $0
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 80)
        )
      )
    )
    (i32.store offset=72
      (get_local $0)
      (i32.const 0)
    )
    (i32.store8 offset=40
      (get_local $0)
      (i32.const 1)
    )
    (i32.store8
      (get_local $0)
      (i32.const 2)
    )
    (i32.store offset=32
      (get_local $0)
      (i32.add
        (get_local $0)
        (i32.const 40)
      )
    )
    (call $emitEvent
      (get_local $0)
      (i32.const 0)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $0)
        (i32.const 80)
      )
    )
  )
  (func $testBlockNumber (result i64)
    (i64.add
      (call $block_number)
      (i64.const 1)
    )
  )
  (func $memTest (result i32)
    (local $0 i32)
    (local $1 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $1
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 48)
        )
      )
    )
    (drop
      (call $memset
        (get_local $1)
        (i32.const 115)
        (i32.const 43)
      )
    )
    (i32.store16 offset=42
      (get_local $1)
      (i32.const 10)
    )
    (drop
      (call $log
        (get_local $1)
      )
    )
    (drop
      (call $memcpy
        (get_local $1)
        (i32.const 256)
        (i32.const 9)
      )
    )
    (drop
      (call $log
        (get_local $1)
      )
    )
    (set_local $0
      (call $memcmp
        (i32.const 256)
        (get_local $1)
        (i32.const 9)
      )
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $1)
        (i32.const 48)
      )
    )
    (get_local $0)
  )
  (func $loopForever
    (loop $label$0
      (br $label$0)
    )
  )
)
