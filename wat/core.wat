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

  ;; Atom - test if value is NOT a cons cell
  ;; Returns NUMBER(1) for true, nil for false
  (func $atom (export "atom") (param $val i64) (result i64)
    (if (result i64)
      (i32.eq (call $get_type (local.get $val)) (global.get $t_cons))
      (then (call $nil))
      (else (call $make_number (i32.const 1)))))

  ;; Eq - test equality of two values
  ;; Returns NUMBER(1) for true, nil for false
  (func $eq (export "eq") (param $a i64) (param $b i64) (result i64)
    (if (result i64)
      (i64.eq (local.get $a) (local.get $b))
      (then (call $make_number (i32.const 1)))
      (else (call $nil))))

  ;; ============================================================================
  ;; BINARY DATA
  ;; ============================================================================

  ;; Create BYTES1 value (1 byte)
  (func $make_bytes1 (export "make_bytes1") (param $byte i32) (result i64)
    (i64.or
      ;; Mask to 1 byte and place in low bits
      (i64.extend_i32_u (i32.and (local.get $byte) (i32.const 0xFF)))
      ;; Type tag (0x08) in high byte
      (i64.const 0x0800000000000000)))

  ;; Create BYTES2 value (2 bytes)
  (func $make_bytes2 (export "make_bytes2") (param $bytes i32) (result i64)
    (i64.or
      ;; Mask to 2 bytes and place in low bits
      (i64.extend_i32_u (i32.and (local.get $bytes) (i32.const 0xFFFF)))
      ;; Type tag (0x09) in high byte
      (i64.const 0x0900000000000000)))

  ;; Create BYTES3 value (3 bytes)
  (func $make_bytes3 (export "make_bytes3") (param $bytes i32) (result i64)
    (i64.or
      ;; Mask to 3 bytes and place in low bits
      (i64.extend_i32_u (i32.and (local.get $bytes) (i32.const 0xFFFFFF)))
      ;; Type tag (0x0A) in high byte
      (i64.const 0x0A00000000000000)))

  ;; Create BYTES4 value (4 bytes)
  (func $make_bytes4 (export "make_bytes4") (param $bytes i32) (result i64)
    (i64.or
      ;; All 4 bytes (no mask needed, i32 is already 4 bytes)
      (i64.extend_i32_u (local.get $bytes))
      ;; Type tag (0x0B) in high byte
      (i64.const 0x0B00000000000000)))

  ;; Get byte count from a BYTES* value (returns 1, 2, 3, or 4)
  (func $get_byte_count (export "get_byte_count") (param $val i64) (result i32)
    (local $type i32)
    (local.set $type (call $get_type (local.get $val)))
    ;; BYTES1=0x08, BYTES2=0x09, BYTES3=0x0A, BYTES4=0x0B
    ;; So count = (type - 7) for types 0x08-0x0B
    (i32.sub (local.get $type) (i32.const 7)))

  ;; ============================================================================
  ;; SYMBOLS
  ;; ============================================================================

  ;; Create a symbol from a binary chain (CONS of BYTES* cells)
  ;; Just re-tags a CONS pointer as SYMBOL
  (func $make_symbol (export "make_symbol") (param $binary_chain i64) (result i64)
    (i64.or
      ;; Keep the pointer (low 32 bits)
      (i64.and (local.get $binary_chain) (i64.const 0x00000000FFFFFFFF))
      ;; Set SYMBOL type tag (0x02)
      (i64.const 0x0200000000000000)))

  ;; Compare two binary chains byte-by-byte
  ;; Returns 1 if equal, 0 if not equal
  (func $binary_equal (param $chain1 i64) (param $chain2 i64) (result i32)
    (local $c1 i64)
    (local $c2 i64)
    (local $car1 i64)
    (local $car2 i64)
    (local $type1 i32)
    (local $type2 i32)

    (local.set $c1 (local.get $chain1))
    (local.set $c2 (local.get $chain2))

    (loop $compare
      ;; If both are nil, they're equal
      (if (i32.and (call $is_nil (local.get $c1)) (call $is_nil (local.get $c2)))
        (then (return (i32.const 1))))

      ;; If one is nil and the other isn't, not equal
      (if (call $is_nil (local.get $c1))
        (then (return (i32.const 0))))
      (if (call $is_nil (local.get $c2))
        (then (return (i32.const 0))))

      ;; Get car of each chain
      (local.set $car1 (call $car (local.get $c1)))
      (local.set $car2 (call $car (local.get $c2)))

      ;; Check if both are BYTES* types
      (local.set $type1 (call $get_type (local.get $car1)))
      (local.set $type2 (call $get_type (local.get $car2)))

      ;; Types must match and be BYTES1/2/3/4
      (if (i32.ne (local.get $type1) (local.get $type2))
        (then (return (i32.const 0))))

      ;; Values must match exactly (includes byte count via type tag)
      (if (i64.ne (local.get $car1) (local.get $car2))
        (then (return (i32.const 0))))

      ;; Move to next cells
      (local.set $c1 (call $cdr (local.get $c1)))
      (local.set $c2 (call $cdr (local.get $c2)))
      (br $compare))

    ;; Should never reach here
    (i32.const 0))

  ;; Compare two symbols for equality
  ;; Returns NUMBER(1) for true, nil for false
  (func $symbol_equal (export "symbol_equal") (param $sym1 i64) (param $sym2 i64) (result i64)
    (local $chain1 i64)
    (local $chain2 i64)

    ;; Convert SYMBOL pointers to CONS pointers to access binary chains
    (local.set $chain1
      (i64.or
        (i64.and (local.get $sym1) (i64.const 0x00000000FFFFFFFF))
        (i64.const 0x0300000000000000)))

    (local.set $chain2
      (i64.or
        (i64.and (local.get $sym2) (i64.const 0x00000000FFFFFFFF))
        (i64.const 0x0300000000000000)))

    (if (result i64) (call $binary_equal (local.get $chain1) (local.get $chain2))
      (then (call $make_number (i32.const 1)))
      (else (call $nil))))

  ;; ============================================================================
  ;; ERRORS
  ;; ============================================================================

  ;; Create an error from a binary chain (UTF-8 message)
  ;; Just re-tags a CONS pointer as ERROR
  (func $make_error (export "make_error") (param $message_chain i64) (result i64)
    (i64.or
      ;; Keep the pointer (low 32 bits)
      (i64.and (local.get $message_chain) (i64.const 0x00000000FFFFFFFF))
      ;; Set ERROR type tag (0x06)
      (i64.const 0x0600000000000000)))

  ;; Extract message from error (returns CONS pointer to binary chain)
  (func $error_message (export "error_message") (param $err i64) (result i64)
    (i64.or
      ;; Keep the pointer (low 32 bits)
      (i64.and (local.get $err) (i64.const 0x00000000FFFFFFFF))
      ;; Set CONS type tag (0x03)
      (i64.const 0x0300000000000000)))

  ;; ============================================================================
  ;; ENVIRONMENT OPERATIONS
  ;; ============================================================================

  ;; Lookup a symbol in an environment
  ;; Environment is an association list: ((sym . val) . ((sym . val) . ...))
  ;; Returns the value if found, or an error if not found
  (func $lookup (export "lookup") (param $sym i64) (param $env i64) (result i64)
    (local $e i64)
    (local $binding i64)
    (local $binding_sym i64)
    (local $binding_val i64)
    (local $is_equal i64)
    (local $err_chain i64)

    (local.set $e (local.get $env))

    (loop $search
      ;; If environment is nil, symbol not found
      (if (call $is_nil (local.get $e))
        (then
          ;; Create error "undefined symbol"
          (local.set $err_chain
            (call $cons
              (call $make_bytes4 (i32.const 0x65646E75))  ;; "unde" (little-endian)
              (call $cons
                (call $make_bytes4 (i32.const 0x656E6966))  ;; "fine" (little-endian)
                (call $cons
                  (call $make_bytes4 (i32.const 0x79732064))  ;; "d sy" (little-endian)
                  (call $cons
                    (call $make_bytes4 (i32.const 0x6C6F626D))  ;; "mbol" (little-endian)
                    (call $nil))))))
          (return (call $make_error (local.get $err_chain)))))

      ;; Get first binding (car env)
      (local.set $binding (call $car (local.get $e)))

      ;; Get symbol from binding (car binding)
      (local.set $binding_sym (call $car (local.get $binding)))

      ;; Check if symbols match
      (local.set $is_equal (call $symbol_equal (local.get $sym) (local.get $binding_sym)))

      (if (i32.eqz (call $is_nil (local.get $is_equal)))
        (then
          ;; Symbols match, return value (cdr binding)
          (return (call $cdr (local.get $binding)))))

      ;; Move to rest of environment
      (local.set $e (call $cdr (local.get $e)))
      (br $search))

    ;; Should never reach here
    (call $nil))

  ;; Extend environment with a new binding
  ;; Returns new environment: ((sym . val) . old_env)
  (func $extend (export "extend") (param $sym i64) (param $val i64) (param $env i64) (result i64)
    (local $binding i64)

    ;; Create binding: (sym . val)
    (local.set $binding (call $cons (local.get $sym) (local.get $val)))

    ;; Cons binding onto environment
    (call $cons (local.get $binding) (local.get $env)))

  ;; ============================================================================
  ;; EVALUATION ENGINE
  ;; ============================================================================

  ;; Evaluate arguments in a list
  (func $eval_args (param $args i64) (param $env i64) (result i64)
    (local $first i64)
    (local $rest i64)
    (local $evaled_first i64)
    (local $evaled_rest i64)

    ;; Base case: nil
    (if (call $is_nil (local.get $args))
      (then (return (call $nil))))

    ;; Evaluate first argument
    (local.set $first (call $car (local.get $args)))
    (local.set $evaled_first (call $eval (local.get $first) (local.get $env)))

    ;; Check for error
    (if (i32.eq (call $get_type (local.get $evaled_first)) (global.get $t_error))
      (then (return (local.get $evaled_first))))

    ;; Recursively evaluate rest
    (local.set $rest (call $cdr (local.get $args)))
    (local.set $evaled_rest (call $eval_args (local.get $rest) (local.get $env)))

    ;; Check for error in rest
    (if (i32.eq (call $get_type (local.get $evaled_rest)) (global.get $t_error))
      (then (return (local.get $evaled_rest))))

    ;; Cons evaluated arg onto result
    (call $cons (local.get $evaled_first) (local.get $evaled_rest)))

  ;; Apply a built-in function
  (func $apply_builtin (param $fn_id i32) (param $args i64) (result i64)
    (local $arg1 i64)
    (local $arg2 i64)
    (local $result i64)

    ;; Built-in ID 0: cons
    (if (i32.eq (local.get $fn_id) (i32.const 0))
      (then
        (local.set $arg1 (call $car (local.get $args)))
        (local.set $arg2 (call $car (call $cdr (local.get $args))))
        (return (call $cons (local.get $arg1) (local.get $arg2)))))

    ;; Built-in ID 1: car
    (if (i32.eq (local.get $fn_id) (i32.const 1))
      (then
        (local.set $arg1 (call $car (local.get $args)))
        (return (call $car (local.get $arg1)))))

    ;; Built-in ID 2: cdr
    (if (i32.eq (local.get $fn_id) (i32.const 2))
      (then
        (local.set $arg1 (call $car (local.get $args)))
        (return (call $cdr (local.get $arg1)))))

    ;; Built-in ID 3: atom
    (if (i32.eq (local.get $fn_id) (i32.const 3))
      (then
        (local.set $arg1 (call $car (local.get $args)))
        (return (call $atom (local.get $arg1)))))

    ;; Built-in ID 4: eq
    (if (i32.eq (local.get $fn_id) (i32.const 4))
      (then
        (local.set $arg1 (call $car (local.get $args)))
        (local.set $arg2 (call $car (call $cdr (local.get $args))))
        (return (call $eq (local.get $arg1) (local.get $arg2)))))

    ;; Unknown built-in
    (call $make_error
      (call $cons (call $make_bytes4 (i32.const 0x6E6B6E75)) ;; "unkn" (little-endian)
        (call $cons (call $make_bytes4 (i32.const 0x206E776F)) ;; "own " (little-endian)
          (call $cons (call $make_bytes4 (i32.const 0x6C697562)) ;; "buil" (little-endian)
            (call $cons (call $make_bytes3 (i32.const 0x6E6974)) ;; "tin" (little-endian)
              (call $nil)))))))

  ;; Special form: quote - return unevaluated
  (func $eval_quote (param $args i64) (result i64)
    ;; Return first argument unevaluated
    (call $car (local.get $args)))

  ;; Special form: if - conditional evaluation
  (func $eval_if (param $args i64) (param $env i64) (result i64)
    (local $cond i64)
    (local $then_expr i64)
    (local $else_expr i64)
    (local $cond_val i64)

    ;; Get condition, then, and else expressions
    (local.set $cond (call $car (local.get $args)))
    (local.set $then_expr (call $car (call $cdr (local.get $args))))
    (local.set $else_expr (call $car (call $cdr (call $cdr (local.get $args)))))

    ;; Evaluate condition
    (local.set $cond_val (call $eval (local.get $cond) (local.get $env)))

    ;; Check for error
    (if (i32.eq (call $get_type (local.get $cond_val)) (global.get $t_error))
      (then (return (local.get $cond_val))))

    ;; If condition is nil, evaluate else; otherwise evaluate then
    (if (result i64) (call $is_nil (local.get $cond_val))
      (then (call $eval (local.get $else_expr) (local.get $env)))
      (else (call $eval (local.get $then_expr) (local.get $env)))))

  ;; Special form: label - bind symbol and evaluate body
  (func $eval_label (param $args i64) (param $env i64) (result i64)
    (local $sym i64)
    (local $val_expr i64)
    (local $body i64)
    (local $val i64)
    (local $new_env i64)

    ;; Get symbol, value expression, and body
    (local.set $sym (call $car (local.get $args)))
    (local.set $val_expr (call $car (call $cdr (local.get $args))))
    (local.set $body (call $car (call $cdr (call $cdr (local.get $args)))))

    ;; Evaluate value
    (local.set $val (call $eval (local.get $val_expr) (local.get $env)))

    ;; Check for error
    (if (i32.eq (call $get_type (local.get $val)) (global.get $t_error))
      (then (return (local.get $val))))

    ;; Extend environment
    (local.set $new_env (call $extend (local.get $sym) (local.get $val) (local.get $env)))

    ;; Evaluate body in new environment
    (call $eval (local.get $body) (local.get $new_env)))

  ;; Special form: lambda - create function closure
  ;; Structure: (params . (body . env))
  (func $eval_lambda (param $args i64) (param $env i64) (result i64)
    (local $params i64)
    (local $body i64)
    (local $closure i64)

    ;; Get params and body
    (local.set $params (call $car (local.get $args)))
    (local.set $body (call $car (call $cdr (local.get $args))))

    ;; Create closure: (params . (body . env))
    (local.set $closure
      (call $cons
        (local.get $params)
        (call $cons
          (local.get $body)
          (local.get $env))))

    ;; Tag as LAMBDA
    (i64.or
      (i64.and (local.get $closure) (i64.const 0x00000000FFFFFFFF))
      (i64.const 0x0400000000000000)))

  ;; Special form: macro - create macro closure
  ;; Structure: (params . (body . env))
  (func $eval_macro (param $args i64) (param $env i64) (result i64)
    (local $params i64)
    (local $body i64)
    (local $closure i64)

    ;; Get params and body
    (local.set $params (call $car (local.get $args)))
    (local.set $body (call $car (call $cdr (local.get $args))))

    ;; Create closure: (params . (body . env))
    (local.set $closure
      (call $cons
        (local.get $params)
        (call $cons
          (local.get $body)
          (local.get $env))))

    ;; Tag as MACRO
    (i64.or
      (i64.and (local.get $closure) (i64.const 0x00000000FFFFFFFF))
      (i64.const 0x0500000000000000)))

  ;; Bind parameters to arguments
  ;; Returns new environment with params bound to args
  (func $bind_params (param $params i64) (param $args i64) (param $env i64) (result i64)
    (local $param i64)
    (local $arg i64)
    (local $rest_params i64)
    (local $rest_args i64)
    (local $new_env i64)

    ;; Base case: both nil
    (if (i32.and (call $is_nil (local.get $params)) (call $is_nil (local.get $args)))
      (then (return (local.get $env))))

    ;; Error: mismatched argument count
    (if (call $is_nil (local.get $params))
      (then
        (return (call $make_error
          (call $cons (call $make_bytes4 (i32.const 0x206F6F74)) ;; "too " (little-endian)
            (call $cons (call $make_bytes4 (i32.const 0x796E616D)) ;; "many" (little-endian)
              (call $cons (call $make_bytes4 (i32.const 0x67726120)) ;; " arg" (little-endian)
                (call $cons (call $make_bytes1 (i32.const 0x73)) ;; "s"
                  (call $nil)))))))))

    (if (call $is_nil (local.get $args))
      (then
        (return (call $make_error
          (call $cons (call $make_bytes4 (i32.const 0x206F6F74)) ;; "too " (little-endian)
            (call $cons (call $make_bytes4 (i32.const 0x20776566)) ;; "few " (little-endian)
              (call $cons (call $make_bytes4 (i32.const 0x73677261)) ;; "args" (little-endian)
                (call $nil))))))))

    ;; Get first param and arg
    (local.set $param (call $car (local.get $params)))
    (local.set $arg (call $car (local.get $args)))

    ;; Get rest
    (local.set $rest_params (call $cdr (local.get $params)))
    (local.set $rest_args (call $cdr (local.get $args)))

    ;; Extend environment with this binding
    (local.set $new_env (call $extend (local.get $param) (local.get $arg) (local.get $env)))

    ;; Recursively bind rest
    (call $bind_params (local.get $rest_params) (local.get $rest_args) (local.get $new_env)))

  ;; Apply lambda to arguments
  (func $apply_lambda (param $lambda i64) (param $args i64) (param $env i64) (result i64)
    (local $closure_ptr i64)
    (local $params i64)
    (local $body i64)
    (local $captured_env i64)
    (local $evaled_args i64)
    (local $new_env i64)

    ;; Convert LAMBDA to CONS to access closure
    (local.set $closure_ptr
      (i64.or
        (i64.and (local.get $lambda) (i64.const 0x00000000FFFFFFFF))
        (i64.const 0x0300000000000000)))

    ;; Extract closure components
    (local.set $params (call $car (local.get $closure_ptr)))
    (local.set $body (call $car (call $cdr (local.get $closure_ptr))))
    (local.set $captured_env (call $cdr (call $cdr (local.get $closure_ptr))))

    ;; Evaluate arguments
    (local.set $evaled_args (call $eval_args (local.get $args) (local.get $env)))

    ;; Check for error
    (if (i32.eq (call $get_type (local.get $evaled_args)) (global.get $t_error))
      (then (return (local.get $evaled_args))))

    ;; Bind parameters to values in captured environment
    (local.set $new_env
      (call $bind_params
        (local.get $params)
        (local.get $evaled_args)
        (local.get $captured_env)))

    ;; Check for error in binding
    (if (i32.eq (call $get_type (local.get $new_env)) (global.get $t_error))
      (then (return (local.get $new_env))))

    ;; Evaluate body in new environment
    (call $eval (local.get $body) (local.get $new_env)))

  ;; Apply macro to arguments
  (func $apply_macro (param $macro i64) (param $args i64) (param $env i64) (result i64)
    (local $closure_ptr i64)
    (local $params i64)
    (local $body i64)
    (local $captured_env i64)
    (local $new_env i64)
    (local $expansion i64)

    ;; Convert MACRO to CONS to access closure
    (local.set $closure_ptr
      (i64.or
        (i64.and (local.get $macro) (i64.const 0x00000000FFFFFFFF))
        (i64.const 0x0300000000000000)))

    ;; Extract closure components
    (local.set $params (call $car (local.get $closure_ptr)))
    (local.set $body (call $car (call $cdr (local.get $closure_ptr))))
    (local.set $captured_env (call $cdr (call $cdr (local.get $closure_ptr))))

    ;; Bind parameters to UNEVALUATED arguments in captured environment
    (local.set $new_env
      (call $bind_params
        (local.get $params)
        (local.get $args)  ;; Note: args NOT evaluated
        (local.get $captured_env)))

    ;; Check for error in binding
    (if (i32.eq (call $get_type (local.get $new_env)) (global.get $t_error))
      (then (return (local.get $new_env))))

    ;; Evaluate body to get expansion
    (local.set $expansion (call $eval (local.get $body) (local.get $new_env)))

    ;; Check for error
    (if (i32.eq (call $get_type (local.get $expansion)) (global.get $t_error))
      (then (return (local.get $expansion))))

    ;; Evaluate the expansion in original environment
    (call $eval (local.get $expansion) (local.get $env)))

  ;; Main evaluation function
  (func $eval (export "eval") (param $expr i64) (param $env i64) (result i64)
    (local $type i32)
    (local $op i64)
    (local $args i64)
    (local $op_val i64)
    (local $op_type i32)
    (local $evaled_args i64)
    (local $fn_id i32)
    (local $quote_sym i64)
    (local $if_sym i64)
    (local $label_sym i64)
    (local $lambda_sym i64)
    (local $macro_sym i64)
    (local $is_quote i64)
    (local $is_if i64)
    (local $is_label i64)
    (local $is_lambda i64)
    (local $is_macro i64)

    (local.set $type (call $get_type (local.get $expr)))

    ;; Self-evaluating types
    (if (i32.eq (local.get $type) (global.get $t_nil))
      (then (return (local.get $expr))))
    (if (i32.eq (local.get $type) (global.get $t_number))
      (then (return (local.get $expr))))
    (if (i32.eq (local.get $type) (global.get $t_error))
      (then (return (local.get $expr))))
    (if (i32.eq (local.get $type) (global.get $t_builtin))
      (then (return (local.get $expr))))
    ;; Lambda and Macro are self-evaluating (closures)
    (if (i32.eq (local.get $type) (global.get $t_lambda))
      (then (return (local.get $expr))))
    (if (i32.eq (local.get $type) (global.get $t_macro))
      (then (return (local.get $expr))))

    ;; Symbol - lookup in environment
    (if (i32.eq (local.get $type) (global.get $t_symbol))
      (then (return (call $lookup (local.get $expr) (local.get $env)))))

    ;; List - evaluate as function application
    (if (i32.eq (local.get $type) (global.get $t_cons))
      (then
        ;; Empty list evaluates to nil
        (if (call $is_nil (local.get $expr))
          (then (return (call $nil))))

        ;; Get operator and arguments
        (local.set $op (call $car (local.get $expr)))
        (local.set $args (call $cdr (local.get $expr)))

        ;; Check for special forms (before evaluating operator)
        (if (i32.eq (call $get_type (local.get $op)) (global.get $t_symbol))
          (then
            ;; Create symbols for special forms (little-endian byte order)
            ;; "quote" = "quot" (0x74756F71) + "e" (0x65)
            (local.set $quote_sym
              (call $make_symbol
                (call $cons (call $make_bytes4 (i32.const 0x746F7571)) ;; "quot"
                  (call $cons (call $make_bytes1 (i32.const 0x65)) ;; "e"
                    (call $nil)))))

            ;; "if" = 0x6669 (little-endian for "if")
            (local.set $if_sym
              (call $make_symbol
                (call $cons (call $make_bytes2 (i32.const 0x6669))
                  (call $nil))))

            ;; "label" = "labe" (0x6562616C) + "l" (0x6C)
            (local.set $label_sym
              (call $make_symbol
                (call $cons (call $make_bytes4 (i32.const 0x6562616C)) ;; "labe"
                  (call $cons (call $make_bytes1 (i32.const 0x6C)) ;; "l"
                    (call $nil)))))

            ;; Check if operator is "quote"
            (local.set $is_quote (call $symbol_equal (local.get $op) (local.get $quote_sym)))
            (if (i32.eqz (call $is_nil (local.get $is_quote)))
              (then (return (call $eval_quote (local.get $args)))))

            ;; Check if operator is "if"
            (local.set $is_if (call $symbol_equal (local.get $op) (local.get $if_sym)))
            (if (i32.eqz (call $is_nil (local.get $is_if)))
              (then (return (call $eval_if (local.get $args) (local.get $env)))))

            ;; Check if operator is "label"
            (local.set $is_label (call $symbol_equal (local.get $op) (local.get $label_sym)))
            (if (i32.eqz (call $is_nil (local.get $is_label)))
              (then (return (call $eval_label (local.get $args) (local.get $env)))))

            ;; "lambda" = "lamb" (0x626D616C) + "da" (0x6164)
            (local.set $lambda_sym
              (call $make_symbol
                (call $cons (call $make_bytes4 (i32.const 0x626D616C)) ;; "lamb"
                  (call $cons (call $make_bytes2 (i32.const 0x6164)) ;; "da"
                    (call $nil)))))

            ;; Check if operator is "lambda"
            (local.set $is_lambda (call $symbol_equal (local.get $op) (local.get $lambda_sym)))
            (if (i32.eqz (call $is_nil (local.get $is_lambda)))
              (then (return (call $eval_lambda (local.get $args) (local.get $env)))))

            ;; "macro" = "macr" (0x7263616D) + "o" (0x6F)
            (local.set $macro_sym
              (call $make_symbol
                (call $cons (call $make_bytes4 (i32.const 0x7263616D)) ;; "macr"
                  (call $cons (call $make_bytes1 (i32.const 0x6F)) ;; "o"
                    (call $nil)))))

            ;; Check if operator is "macro"
            (local.set $is_macro (call $symbol_equal (local.get $op) (local.get $macro_sym)))
            (if (i32.eqz (call $is_nil (local.get $is_macro)))
              (then (return (call $eval_macro (local.get $args) (local.get $env)))))))

        ;; Not a special form - evaluate operator
        (local.set $op_val (call $eval (local.get $op) (local.get $env)))

        ;; Check for error in operator
        (if (i32.eq (call $get_type (local.get $op_val)) (global.get $t_error))
          (then (return (local.get $op_val))))

        ;; Get operator type
        (local.set $op_type (call $get_type (local.get $op_val)))

        ;; Built-in function
        (if (i32.eq (local.get $op_type) (global.get $t_builtin))
          (then
            ;; Evaluate arguments
            (local.set $evaled_args (call $eval_args (local.get $args) (local.get $env)))

            ;; Check for error in arguments
            (if (i32.eq (call $get_type (local.get $evaled_args)) (global.get $t_error))
              (then (return (local.get $evaled_args))))

            ;; Apply built-in
            (local.set $fn_id (call $get_value (local.get $op_val)))
            (return (call $apply_builtin (local.get $fn_id) (local.get $evaled_args)))))

        ;; Lambda
        (if (i32.eq (local.get $op_type) (global.get $t_lambda))
          (then (return (call $apply_lambda (local.get $op_val) (local.get $args) (local.get $env)))))

        ;; Macro
        (if (i32.eq (local.get $op_type) (global.get $t_macro))
          (then (return (call $apply_macro (local.get $op_val) (local.get $args) (local.get $env)))))

        ;; Not a function
        (return (call $make_error
          (call $cons (call $make_bytes4 (i32.const 0x20746F6E)) ;; "not " (little-endian)
            (call $cons (call $make_bytes4 (i32.const 0x75662061)) ;; "a fu" (little-endian)
              (call $cons (call $make_bytes4 (i32.const 0x6974636E)) ;; "ncti" (little-endian)
                (call $cons (call $make_bytes2 (i32.const 0x6E6F)) ;; "on" (little-endian)
                  (call $nil)))))))))

    ;; Invalid expression type
    (call $make_error
      (call $cons (call $make_bytes4 (i32.const 0x61766E69)) ;; "inva" (little-endian)
        (call $cons (call $make_bytes4 (i32.const 0x2064696C)) ;; "lid " (little-endian)
          (call $cons (call $make_bytes4 (i32.const 0x72707865)) ;; "expr" (little-endian)
            (call $nil))))))

  ;; Create a built-in function value
  (func $make_builtin (export "make_builtin") (param $id i32) (result i64)
    (i64.or
      (i64.extend_i32_u (local.get $id))
      (i64.const 0x0700000000000000)))

)
