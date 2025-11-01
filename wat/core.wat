(module

  ;; MEMORY LAYOUT
  ;;
  ;; All values are i64:
  ;; [TYPE_TAG|RESERVED|VALUE/POINTER]
  ;;  8 bits   24 bits  32 bits
  ;;
  ;; Type tags:
  ;; NIL     = 0x00
  ;; NUMBER  = 0x01
  ;; SYMBOL  = 0x02
  ;; CONS    = 0x03
  ;; LAMBDA  = 0x04
  ;; MACRO   = 0x05
  ;; ERROR   = 0x06
  ;; BUILTIN = 0x07
  ;; BYTES1  = 0x08
  ;; BYTES2  = 0x09
  ;; BYTES3  = 0x0A
  ;; BYTES4  = 0x0B

  ;; Type constants
  (global $t_nil i32 (i32.const 0x00))
  (global $t_number i32 (i32.const 0x01))
  (global $t_symbol i32 (i32.const 0x02))
  (global $t_cons i32 (i32.const 0x03))
  (global $t_lambda i32 (i32.const 0x04))
  (global $t_macro i32 (i32.const 0x05))
  (global $t_error i32 (i32.const 0x06))
  (global $t_builtin i32 (i32.const 0x07))
  (global $t_bytes1 i32 (i32.const 0x08))
  (global $t_bytes2 i32 (i32.const 0x09))
  (global $t_bytes3 i32 (i32.const 0x0A))
  (global $t_bytes4 i32 (i32.const 0x0B))

  ;; Memory (single heap starting at 0x0000)
  (memory 1)
  (global $heap_ptr (mut i32) (i32.const 0))

  ;; ============================================================================
  ;; BASIC VALUE HELPERS
  ;; ============================================================================

  ;; Return nil (0x0000000000000000)
  (func $nil (export "nil") (result i64)
    (i64.const 0))

  ;; Create a number from i32
  (func $make_number (export "make_number") (param $n i32) (result i64)
    (i64.or
      ;; Value in low 32 bits (sign-extended from i32)
      (i64.extend_i32_u (local.get $n))
      ;; Type tag (0x01) in high byte
      (i64.const 0x0100000000000000)))

  ;; Extract type tag from value (returns high byte)
  (func $get_type (export "get_type") (param $val i64) (result i32)
    (i32.wrap_i64
      (i64.shr_u (local.get $val) (i64.const 56))))

  ;; Extract value/pointer from value (returns low 32 bits)
  (func $get_value (export "get_value") (param $val i64) (result i32)
    (i32.wrap_i64 (local.get $val)))

  ;; Check if value is nil (returns 1 for true, 0 for false)
  (func $is_nil (export "is_nil") (param $val i64) (result i32)
    (i64.eqz (local.get $val)))

  ;; ============================================================================
  ;; MEMORY ALLOCATION
  ;; ============================================================================

  ;; Allocate 16-byte aligned cons cell (returns pointer)
  (func $alloc_cons (export "alloc_cons") (result i32)
    (local $ptr i32)
    ;; Get current heap pointer
    (local.set $ptr (global.get $heap_ptr))
    ;; Advance by 16 bytes
    (global.set $heap_ptr
      (i32.add (local.get $ptr) (i32.const 16)))
    ;; Return the allocated address
    (local.get $ptr))

  ;; ============================================================================
  ;; LIST OPERATIONS
  ;; ============================================================================

  ;; Cons - create new cons cell
  (func $cons (export "cons") (param $car i64) (param $cdr i64) (result i64)
    (local $ptr i32)
    (local $cell v128)

    ;; Allocate cell
    (local.set $ptr (call $alloc_cons))

    ;; Build v128 from car and cdr
    (local.set $cell
      (i64x2.replace_lane 1
        (i64x2.replace_lane 0
          (v128.const i64x2 0 0)
          (local.get $car))
        (local.get $cdr)))

    ;; Store cell
    (v128.store (local.get $ptr) (local.get $cell))

    ;; Return tagged pointer (CONS type = 0x03)
    (i64.or
      (i64.extend_i32_u (local.get $ptr))
      (i64.const 0x0300000000000000)))

  ;; Car - get first element
  (func $car (export "car") (param $cell i64) (result i64)
    (local $ptr i32)
    ;; Extract pointer from tagged value
    (local.set $ptr (call $get_value (local.get $cell)))
    ;; Load car (first i64 of cons cell)
    (i64.load (local.get $ptr)))

  ;; Cdr - get rest of list
  (func $cdr (export "cdr") (param $cell i64) (result i64)
    (local $ptr i32)
    ;; Extract pointer from tagged value
    (local.set $ptr (call $get_value (local.get $cell)))
    ;; Load cdr (second i64 of cons cell, offset by 8 bytes)
    (i64.load offset=8 (local.get $ptr)))

)
