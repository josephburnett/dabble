(module

  ;; memory
  (memory 0)
  (global $next_cell_ptr (mut i32) (i32.const 0))
  (func $peek_cell (result i32)
    (global.get $next_cell_ptr))
  (func $next_cell (result i32)
    (global.get $next_cell_ptr)
    (global.set $next_cell_ptr (i32.add (global.get $next_cell_ptr) (i32.const 1))))

  ;; cons
  (func $cons (param $car i32) (param $cdr i32) (result i32)
    (call $peek_cell)
    (i32.store (call $next_cell) (local.get $car))
    (i32.store (call $next_cell) (local.get $cdr)))

  ;; car
  (func $car (param $cell i32) (result i32)
    (i32.load (local.get $cell)))

  ;; cdr
  (func $cdr (param $cell i32) (result i32)
    (i32.load (i32.add (local.get $cell) (i32.const 1)))))