(module
  (type $FUNCSIG$vii (func (param i32 i32)))
  (type $FUNCSIG$ii (func (param i32) (result i32)))
  (type $FUNCSIG$vi (func (param i32)))
  (type $FUNCSIG$iiii (func (param i32 i32 i32) (result i32)))
  (type $FUNCSIG$viiii (func (param i32 i32 i32 i32)))
  (type $FUNCSIG$v (func))
  (import "env" "action_sender" (func $action_sender (param i32)))
  (import "env" "assert" (func $assert (param i32 i32)))
  (import "env" "block_producer" (func $block_producer (param i32)))
  (import "env" "db_get" (func $db_get (param i32 i32 i32) (result i32)))
  (import "env" "db_set" (func $db_set (param i32 i32 i32 i32)))
  (import "env" "requireRewordAuth" (func $requireRewordAuth))
  (import "env" "setResult" (func $setResult (param i32 i32)))
  (import "env" "str2lower" (func $str2lower (param i32)))
  (import "env" "strjoint" (func $strjoint (param i32 i32 i32) (result i32)))
  (import "env" "strlen" (func $strlen (param i32) (result i32)))
  (table 0 anyfunc)
  (memory $0 1)
  (data (i32.const 4) "0G\00\00")
  (data (i32.const 16) "memo exceed\00")
  (data (i32.const 32) "transfer get sender BC error\00")
  (data (i32.const 64) "transfer insufficient BC\00")
  (data (i32.const 96) "transferFee fee is 0\00")
  (data (i32.const 128) "transferFee get sender BC error\00")
  (data (i32.const 160) "transferFee insufficient BC\00")
  (data (i32.const 188) "\08\00\00\00")
  (data (i32.const 192) "BC\00")
  (data (i32.const 208) "_rNumber\00")
  (data (i32.const 224) "\00\c2\eb\0b\00\00\00\00\80\a7\e5\0b\00\00\00\00 \90\df\0b\00\00\00\00\de{\d9\0b\00\00\00\00\b9j\d3\0b\00\00\00\00\af\\\cd\0b\00\00\00\00\bfQ\c7\0b\00\00\00\00\e7I\c1\0b\00\00\00\00&E\bb\0b\00\00\00\00yC\b5\0b\00\00\00\00\e0D\af\0b\00\00\00\00XI\a9\0b\00\00\00\00\e1P\a3\0b\00\00\00\00x[\9d\0b\00\00\00\00\1ci\97\0b\00\00\00\00\ccy\91\0b\00\00\00\00\85\8d\8b\0b\00\00\00\00G\a4\85\0b\00\00\00\00\10\be\7f\0b\00\00\00\00\de\day\0b\00\00\00\00\af\fas\0b\00\00\00\00\83\1dn\0b\00\00\00\00WCh\0b\00\00\00\00*lb\0b\00\00\00\00\fb\97\\\0b\00\00\00\00\c8\c6V\0b\00\00\00\00\8f\f8P\0b\00\00\00\00O-K\0b\00\00\00\00\07eE\0b\00\00\00\00\b5\9f?\0b\00\00\00\00W\dd9\0b\00\00\00\00\ec\1d4\0b\00\00\00\00ra.\0b\00\00\00\00\e8\a7(\0b\00\00\00\00L\f1\"\0b\00\00\00\00\9d=\1d\0b\00\00\00\00\da\8c\17\0b\00\00\00\00\00\df\11\0b\00\00\00\00\0f4\0c\0b\00\00\00\00\05\8c\06\0b\00\00\00\00\e0\e6\00\0b\00\00\00\00\9fD\fb\n\00\00\00\00@\a5\f5\n\00\00\00\00\c2\08\f0\n\00\00\00\00$o\ea\n\00\00\00\00d\d8\e4\n\00\00\00\00\80D\df\n\00\00\00\00x\b3\d9\n\00\00\00\00I%\d4\n\00\00\00\00\f2\99\ce\n\00\00\00\00r\11\c9\n\00\00\00\00\c7\8b\c3\n\00\00\00\00\f0\08\be\n\00\00\00\00\eb\88\b8\n\00\00\00\00\b7\0b\b3\n\00\00\00\00S\91\ad\n\00\00\00\00\bd\19\a8\n\00\00\00\00\f3\a4\a2\n\00\00\00\00\f52\9d\n\00\00\00\00\c0\c3\97\n\00\00\00\00TW\92\n\00\00\00\00\ae\ed\8c\n\00\00\00\00\ce\86\87\n\00\00\00\00\b2\"\82\n\00\00\00\00Y\c1|\n\00\00\00\00\c1bw\n\00\00\00\00\e9\06r\n\00\00\00\00\cf\adl\n\00\00\00\00rWg\n\00\00\00\00\d0\03b\n\00\00\00\00\e9\b2\\\n\00\00\00\00\badW\n\00\00\00\00C\19R\n\00\00\00\00\82\d0L\n\00\00\00\00u\8aG\n\00\00\00\00\1cGB\n\00\00\00\00u\06=\n\00\00\00\00~\c87\n\00\00\00\006\8d2\n\00\00\00\00\9cT-\n\00\00\00\00\ae\1e(\n\00\00\00\00k\eb\"\n\00\00\00\00\d2\ba\1d\n\00\00\00\00\e1\8c\18\n\00\00\00\00\97a\13\n\00\00\00\00\f28\0e\n\00\00\00\00\f2\12\t\n\00\00\00\00\94\ef\03\n\00\00\00\00\d8\ce\fe\t\00\00\00\00\bc\b0\f9\t\00\00\00\00?\95\f4\t\00\00\00\00_|\ef\t\00\00\00\00\1bf\ea\t\00\00\00\00rR\e5\t\00\00\00\00bA\e0\t\00\00\00\00\eb2\db\t\00\00\00\00\n\'\d6\t\00\00\00\00\bf\1d\d1\t\00\00\00\00\08\17\cc\t\00\00\00\00\e4\12\c7\t\00\00\00\00Q\11\c2\t\00\00\00\00N\12\bd\t\00\00\00\00\da\15\b8\t\00\00\00\00\f4\1b\b3\t\00\00\00\00\9a$\ae\t\00\00\00\00\cb/\a9\t\00\00\00\00\85=\a4\t\00\00\00\00\c8M\9f\t\00\00\00\00\92`\9a\t\00\00\00\00\e2u\95\t\00\00\00\00\b6\8d\90\t\00\00\00\00\0d\a8\8b\t\00\00\00\00\e6\c4\86\t\00\00\00\00@\e4\81\t\00\00\00\00\19\06}\t\00\00\00\00p*x\t\00\00\00\00DQs\t\00\00\00\00\93zn\t\00\00\00\00\\\a6i\t\00\00\00\00\9e\d4d\t\00\00\00\00X\05`\t\00\00\00\00\888[\t\00\00\00\00.nV\t\00\00\00\00H\a6Q\t\00\00\00\00\d4\e0L\t\00\00\00\00\d2\1dH\t\00\00\00\00@]C\t\00\00\00\00\1d\9f>\t\00\00\00\00g\e39\t\00\00\00\00\1e*5\t\00\00\00\00@s0\t\00\00\00\00\cc\be+\t\00\00\00\00\c0\0c\'\t\00\00\00\00\1c]\"\t\00\00\00\00\de\af\1d\t\00\00\00\00\05\05\19\t\00\00\00\00\90\\\14\t\00\00\00\00}\b6\0f\t\00\00\00\00\cc\12\0b\t\00\00\00\00{q\06\t\00\00\00\00\89\d2\01\t\00\00\00\00\f45\fd\08\00\00\00\00\bc\9b\f8\08\00\00\00\00\df\03\f4\08\00\00\00\00\\n\ef\08\00\00\00\002\db\ea\08\00\00\00\00_J\e6\08\00\00\00\00\e3\bb\e1\08\00\00\00\00\bc/\dd\08\00\00\00\00\e9\a5\d8\08\00\00\00\00i\1e\d4\08\00\00\00\00;\99\cf\08\00\00\00\00]\16\cb\08\00\00\00\00\ce\95\c6\08\00\00\00\00\8e\17\c2\08\00\00\00\00\9b\9b\bd\08\00\00\00\00\f3!\b9\08\00\00\00\00\96\aa\b4\08\00\00\00\00\825\b0\08\00\00\00\00\b6\c2\ab\08\00\00\00\002R\a7\08\00\00\00\00\f3\e3\a2\08\00\00\00\00\f9w\9e\08\00\00\00\00C\0e\9a\08\00\00\00\00\cf\a6\95\08\00\00\00\00\9cA\91\08\00\00\00\00\aa\de\8c\08\00\00\00\00\f6}\88\08\00\00\00\00\80\1f\84\08\00\00\00\00G\c3\7f\08\00\00\00\00Ii{\08\00\00\00\00\86\11w\08\00\00\00\00\fc\bbr\08\00\00\00\00\aahn\08\00\00\00\00\8f\17j\08\00\00\00\00\aa\c8e\08\00\00\00\00\fa{a\08\00\00\00\00}1]\08\00\00\00\003\e9X\08\00\00\00\00\1a\a3T\08\00\00\00\001_P\08\00\00\00\00w\1dL\08\00\00\00\00\eb\ddG\08\00\00\00\00\8c\a0C\08\00\00\00\00Ye?\08\00\00\00\00P,;\08\00\00\00\00q\f56\08\00\00\00\00\ba\c02\08\00\00\00\00*\8e.\08\00\00\00\00\c0]*\08\00\00\00\00|/&\08\00\00\00\00[\03\"\08\00\00\00\00]\d9\1d\08\00\00\00\00\81\b1\19\08\00\00\00\00\c6\8b\15\08\00\00\00\00*h\11\08\00\00\00\00\adF\0d\08\00\00\00\00M\'\t\08\00\00\00\00\n\n\05\08\00\00\00\00\e2\ee\00\08\00\00\00\00")
  (data (i32.const 1824) "B Chain Token\00")
  (export "memory" (memory $0))
  (export "_ZN6bchain8transferEPcyS0_" (func $_ZN6bchain8transferEPcyS0_))
  (export "_ZN6bchain6rewordEv" (func $_ZN6bchain6rewordEv))
  (export "_ZN6bchain11transferFeeEy" (func $_ZN6bchain11transferFeeEy))
  (export "_ZN6bchain9balenceOfEPc" (func $_ZN6bchain9balenceOfEPc))
  (export "_ZN6bchain9getSupplyEv" (func $_ZN6bchain9getSupplyEv))
  (export "_ZN6bchain11getDecimalsEv" (func $_ZN6bchain11getDecimalsEv))
  (export "_ZN6bchain9getSymbolEv" (func $_ZN6bchain9getSymbolEv))
  (export "_ZN6bchain7getNameEv" (func $_ZN6bchain7getNameEv))
  (export "transfer" (func $transfer))
  (export "reword" (func $reword))
  (export "transferFee" (func $transferFee))
  (export "balenceOf" (func $balenceOf))
  (export "getSupply" (func $getSupply))
  (export "getDecimals" (func $getDecimals))
  (export "getSymbol" (func $getSymbol))
  (export "getName" (func $getName))
  (func $_ZN6bchain8transferEPcyS0_ (param $0 i32) (param $1 i32) (param $2 i64) (param $3 i32)
    (local $4 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $4
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 336)
        )
      )
    )
    (call $assert
      (i32.lt_s
        (call $strlen
          (get_local $3)
        )
        (i32.const 32)
      )
      (i32.const 16)
    )
    (call $action_sender
      (i32.add
        (get_local $4)
        (i32.const 288)
      )
    )
    (i64.store offset=152
      (get_local $4)
      (i64.const 0)
    )
    (call $assert
      (i32.eq
        (call $db_get
          (i32.add
            (get_local $4)
            (i32.const 160)
          )
          (tee_local $3
            (call $strjoint
              (i32.const 192)
              (i32.add
                (get_local $4)
                (i32.const 288)
              )
              (i32.add
                (get_local $4)
                (i32.const 160)
              )
            )
          )
          (i32.add
            (get_local $4)
            (i32.const 152)
          )
        )
        (i32.const 8)
      )
      (i32.const 32)
    )
    (call $assert
      (i64.ge_u
        (i64.load offset=152
          (get_local $4)
        )
        (get_local $2)
      )
      (i32.const 64)
    )
    (i64.store offset=152
      (get_local $4)
      (i64.sub
        (i64.load offset=152
          (get_local $4)
        )
        (get_local $2)
      )
    )
    (call $db_set
      (i32.add
        (get_local $4)
        (i32.const 160)
      )
      (get_local $3)
      (i32.add
        (get_local $4)
        (i32.const 152)
      )
      (i32.const 8)
    )
    (i64.store offset=8
      (get_local $4)
      (i64.const 0)
    )
    (call $str2lower
      (get_local $1)
    )
    (drop
      (call $db_get
        (i32.add
          (get_local $4)
          (i32.const 16)
        )
        (tee_local $1
          (call $strjoint
            (i32.const 192)
            (get_local $1)
            (i32.add
              (get_local $4)
              (i32.const 16)
            )
          )
        )
        (i32.add
          (get_local $4)
          (i32.const 8)
        )
      )
    )
    (i64.store offset=8
      (get_local $4)
      (i64.add
        (i64.load offset=8
          (get_local $4)
        )
        (get_local $2)
      )
    )
    (call $db_set
      (i32.add
        (get_local $4)
        (i32.const 16)
      )
      (get_local $1)
      (i32.add
        (get_local $4)
        (i32.const 8)
      )
      (i32.const 8)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $4)
        (i32.const 336)
      )
    )
  )
  (func $_ZN6bchain6rewordEv (param $0 i32)
    (local $1 i64)
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
          (i32.const 192)
        )
      )
    )
    (call $requireRewordAuth)
    (call $block_producer
      (i32.add
        (get_local $4)
        (i32.const 144)
      )
    )
    (i64.store offset=16
      (get_local $4)
      (i64.const 0)
    )
    (drop
      (call $db_get
        (i32.const 208)
        (tee_local $2
          (call $strlen
            (i32.const 208)
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
        (i64.const 1)
      )
    )
    (call $db_set
      (i32.const 208)
      (get_local $2)
      (i32.add
        (get_local $4)
        (i32.const 16)
      )
      (i32.const 8)
    )
    (set_local $1
      (i64.load offset=16
        (get_local $4)
      )
    )
    (i64.store offset=8
      (get_local $4)
      (i64.const 0)
    )
    (drop
      (call $db_get
        (i32.add
          (get_local $4)
          (i32.const 16)
        )
        (tee_local $2
          (call $strjoint
            (i32.const 192)
            (i32.add
              (get_local $4)
              (i32.const 144)
            )
            (i32.add
              (get_local $4)
              (i32.const 16)
            )
          )
        )
        (i32.add
          (get_local $4)
          (i32.const 8)
        )
      )
    )
    (set_local $3
      (i64.const 1750000000000000)
    )
    (block $label$0
      (br_if $label$0
        (i64.eq
          (get_local $1)
          (i64.const 1)
        )
      )
      (set_local $3
        (i64.const 0)
      )
      (br_if $label$0
        (i64.gt_u
          (get_local $1)
          (i64.const 49999999)
        )
      )
      (set_local $3
        (i64.load
          (i32.add
            (i32.shl
              (i32.wrap/i64
                (i64.div_u
                  (get_local $1)
                  (i64.const 250000)
                )
              )
              (i32.const 3)
            )
            (i32.const 224)
          )
        )
      )
    )
    (i64.store offset=8
      (get_local $4)
      (i64.add
        (i64.load offset=8
          (get_local $4)
        )
        (get_local $3)
      )
    )
    (call $db_set
      (i32.add
        (get_local $4)
        (i32.const 16)
      )
      (get_local $2)
      (i32.add
        (get_local $4)
        (i32.const 8)
      )
      (i32.const 8)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $4)
        (i32.const 192)
      )
    )
  )
  (func $_ZN6bchain11transferFeeEy (param $0 i32) (param $1 i64)
    (local $2 i32)
    (local $3 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $3
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 384)
        )
      )
    )
    (call $assert
      (i64.ne
        (get_local $1)
        (i64.const 0)
      )
      (i32.const 96)
    )
    (call $block_producer
      (i32.add
        (get_local $3)
        (i32.const 336)
      )
    )
    (call $action_sender
      (i32.add
        (get_local $3)
        (i32.const 288)
      )
    )
    (i64.store offset=280
      (get_local $3)
      (i64.const 0)
    )
    (call $assert
      (i32.eq
        (call $db_get
          (i32.add
            (get_local $3)
            (i32.const 144)
          )
          (tee_local $2
            (call $strjoint
              (i32.const 192)
              (i32.add
                (get_local $3)
                (i32.const 288)
              )
              (i32.add
                (get_local $3)
                (i32.const 144)
              )
            )
          )
          (i32.add
            (get_local $3)
            (i32.const 280)
          )
        )
        (i32.const 8)
      )
      (i32.const 128)
    )
    (call $assert
      (i64.ge_u
        (i64.load offset=280
          (get_local $3)
        )
        (get_local $1)
      )
      (i32.const 160)
    )
    (i64.store offset=280
      (get_local $3)
      (i64.sub
        (i64.load offset=280
          (get_local $3)
        )
        (get_local $1)
      )
    )
    (call $db_set
      (i32.add
        (get_local $3)
        (i32.const 144)
      )
      (get_local $2)
      (i32.add
        (get_local $3)
        (i32.const 280)
      )
      (i32.const 8)
    )
    (i64.store offset=8
      (get_local $3)
      (i64.const 0)
    )
    (drop
      (call $db_get
        (i32.add
          (get_local $3)
          (i32.const 16)
        )
        (tee_local $2
          (call $strjoint
            (i32.const 192)
            (i32.add
              (get_local $3)
              (i32.const 336)
            )
            (i32.add
              (get_local $3)
              (i32.const 16)
            )
          )
        )
        (i32.add
          (get_local $3)
          (i32.const 8)
        )
      )
    )
    (i64.store offset=8
      (get_local $3)
      (i64.add
        (i64.load offset=8
          (get_local $3)
        )
        (get_local $1)
      )
    )
    (call $db_set
      (i32.add
        (get_local $3)
        (i32.const 16)
      )
      (get_local $2)
      (i32.add
        (get_local $3)
        (i32.const 8)
      )
      (i32.const 8)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $3)
        (i32.const 384)
      )
    )
  )
  (func $_ZN6bchain9balenceOfEPc (param $0 i32) (param $1 i32)
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
          (i32.const 192)
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
  (func $_ZN6bchain9getSupplyEv (param $0 i32)
    (local $1 i32)
    (local $2 i64)
    (local $3 i32)
    (local $4 i32)
    (i32.store offset=4
      (i32.const 0)
      (tee_local $4
        (i32.sub
          (i32.load offset=4
            (i32.const 0)
          )
          (i32.const 16)
        )
      )
    )
    (set_local $2
      (i64.const 100000000)
    )
    (i64.store offset=8
      (get_local $4)
      (i64.const 100000000)
    )
    (set_local $3
      (i32.const 0)
    )
    (block $label$0
      (br_if $label$0
        (i32.eqz
          (tee_local $1
            (i32.load offset=188
              (i32.const 0)
            )
          )
        )
      )
      (loop $label$1
        (set_local $2
          (i64.mul
            (get_local $2)
            (i64.const 10)
          )
        )
        (br_if $label$1
          (i32.lt_u
            (tee_local $3
              (i32.add
                (get_local $3)
                (i32.const 1)
              )
            )
            (get_local $1)
          )
        )
      )
      (i64.store offset=8
        (get_local $4)
        (get_local $2)
      )
    )
    (call $setResult
      (i32.add
        (get_local $4)
        (i32.const 8)
      )
      (i32.const 8)
    )
    (i32.store offset=4
      (i32.const 0)
      (i32.add
        (get_local $4)
        (i32.const 16)
      )
    )
  )
  (func $_ZN6bchain11getDecimalsEv (param $0 i32)
    (call $setResult
      (i32.const 188)
      (i32.const 4)
    )
  )
  (func $_ZN6bchain9getSymbolEv (param $0 i32)
    (call $setResult
      (i32.const 192)
      (call $strlen
        (i32.const 192)
      )
    )
  )
  (func $_ZN6bchain7getNameEv (param $0 i32)
    (call $setResult
      (i32.const 1824)
      (call $strlen
        (i32.const 1824)
      )
    )
  )
  (func $transfer (param $0 i32) (param $1 i64) (param $2 i32)
    (call $_ZN6bchain8transferEPcyS0_
      (get_local $0)
      (get_local $0)
      (get_local $1)
      (get_local $2)
    )
  )
  (func $reword
    (local $0 i32)
    (call $_ZN6bchain6rewordEv
      (get_local $0)
    )
  )
  (func $transferFee (param $0 i64)
    (local $1 i32)
    (call $_ZN6bchain11transferFeeEy
      (get_local $1)
      (get_local $0)
    )
  )
  (func $balenceOf (param $0 i32)
    (call $_ZN6bchain9balenceOfEPc
      (get_local $0)
      (get_local $0)
    )
  )
  (func $getSupply
    (local $0 i32)
    (call $_ZN6bchain9getSupplyEv
      (get_local $0)
    )
  )
  (func $getDecimals
    (call $setResult
      (i32.const 188)
      (i32.const 4)
    )
  )
  (func $getSymbol
    (call $setResult
      (i32.const 192)
      (call $strlen
        (i32.const 192)
      )
    )
  )
  (func $getName
    (call $setResult
      (i32.const 1824)
      (call $strlen
        (i32.const 1824)
      )
    )
  )
)
