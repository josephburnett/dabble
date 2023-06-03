(module

  ;; memory
  (memory 0)
  (global $mem (mut i32) (i32.const 0))

  ;; testing compile
  (func $cons (param $lhs i32) (param $rhs i32) (result i32)
    (global.get $mem)
        (i32.store (i32.const 1))
    (global.set $mem (i32.const 1))
    i32.const 0))
