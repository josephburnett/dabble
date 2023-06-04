(module

  ;; MEMORY

  ;; cell    = (v128: t_value, t_value)
  ;; t_value = (i64: type, value)
  ;; type    = (i32) integer | float | letter | byte | error | *cell
  ;; value   = (i32)
  ;; integer = 1
  ;; float   = 2
  ;; letter  = 4
  ;; byte    = 8
  ;; error   = 16
  ;; *cell   = 32

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
    (call $peek_cell)           ;; address of new cons cell for store operation
    (v128.const i64x2 0 0)      ;; pack car and cdr t_values into a new cell
    (local.get $car)
    (i64x2.replace_lane 0)
    (local.get $cdr)
    (i64x2.replace_lane 1)
    (v128.store)                ;; store the new cell
    (v128.const i32x4 0 0 0 0)  ;; pack type and value into a t_value
    (i32.const 32)              ;; type *cell
    (i32x4.replace_lane 0)
    (call $next_cell)           ;; value is address of new cons cell
    (i32x4.replace_lane 1)
    (i64x2.extract_lane 0))

  ;; car
  (func $car (param $cell i32) (result i32)
    (i32.load (local.get $cell)))

  ;; cdr
  (func $cdr (param $cell i32) (result i32)
    (i32.load (i32.add (local.get $cell) (i32.const 1))))

  ;; atom
  ;;(func $atom (param $cell i32) (result i32)
  ;;  (local.get 

)