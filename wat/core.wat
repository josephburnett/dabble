(module

  ;; MEMORY

  ;; cell    = v128: (t_value, t_value)
  ;; t_value = i64: (value, type)
  ;; type    = i32: integer | float | letter | byte | error | *cell
  ;; value   = i32


  (global $t_integer i32 (i32.const 1))
  (global $t_float i32 (i32.const 2))
  (global $t_letter i32 (i32.const 3))
  (global $t_byte i32 (i32.const 4))
  (global $t_error i32 (i32.const 5))
  (global $t_cell i32 (i32.const 6))
  (global $nil i64 (i64.const 6))

  (memory 0)
  (global $next_cell_ptr (mut i32) (i32.const 0))

  (func $peek_cell (result i32)
    (global.get $next_cell_ptr))
  (func $next_cell (result i32)
    (global.get $next_cell_ptr)
    (global.set $next_cell_ptr (i32.add (global.get $next_cell_ptr) (i32.const 16))))

  ;; CORE FUNCTIONS

  ;; cons
  (func $cons (param $car i64) (param $cdr i64) (result i64)
    (call $peek_cell)           ;; address of new cons cell
    (v128.const i64x2 0 0)      ;; pack car and cdr t_values into a new cell
    (local.get $car)
    (i64x2.replace_lane 0)
    (local.get $cdr)
    (i64x2.replace_lane 1)
    (v128.store)                ;; store the new cell
    (call $next_cell)           ;; address of new cons cell
    (i64.extend_i32_u)
    (i64.shl (i64.const 32))    ;; set value
    (global.get $t_cell)        ;; set type *cell
    (i64.extend_i32_u)
    (i64.or))                   ;; combine value and type

  ;; car
  ;; (func $car (param $cell i64) (result i64)
  ;;   (local.get $cell)
  ;;   (i32x2.extract_lane 0)
    

  ;; cdr
  (func $cdr (param $cell i32) (result i32)
    (i32.load (i32.add (local.get $cell) (i32.const 1))))

  ;; atom
  ;;(func $atom (param $cell i32) (result i32)
  ;;  (local.get 

)