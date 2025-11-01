# Dabble WebAssembly Implementation Design

## Philosophy

This WebAssembly Text (WAT) implementation of Dabble prioritizes simplicity and clarity over efficiency. The goal is to create a complete, understandable Lisp implementation that can serve as both a learning tool and a working interpreter.

Key principles:
- Everything is a cons cell (including binary data, symbols, and errors)
- Uniform value representation (all values are i64)
- Single heap with linear allocation (no special memory regions)
- Maximum code clarity over performance

## Memory Layout

### Single Heap Model

```
Memory starts at address 0x0000 and grows linearly:

[0x0000] Start of heap - cons cells allocated here
   |
   v
[heap_ptr] Next free location (16-byte aligned)
   |
   v
[...] Unused memory
```

All memory allocation happens on a single heap with a bump allocator. There are no special regions - symbols, cons cells, and all data share the same space.

### Value Representation (i64)

Every value is a 64-bit integer with this layout:

```
Bit Position:  [63-56]    [55-48]    [47-32]    [31-0]
Content:       TYPE_TAG   RESERVED   RESERVED   VALUE/POINTER
               (8 bits)   (8 bits)   (16 bits)  (32 bits)
```

The high byte contains the type tag, the low 32 bits contain the value or memory pointer.

### Type System

| Type    | Tag  | Description                              | Value Field Contains        |
|---------|------|------------------------------------------|-----------------------------|
| NIL     | 0x00 | Empty list / false                      | Always 0                    |
| NUMBER  | 0x01 | 32-bit signed integer                   | The integer value           |
| SYMBOL  | 0x02 | Symbol (stored as binary data)          | Pointer to binary chain     |
| CONS    | 0x03 | Regular cons cell                        | Pointer to cell             |
| LAMBDA  | 0x04 | Function closure                         | Pointer to cell             |
| MACRO   | 0x05 | Macro closure                            | Pointer to cell             |
| ERROR   | 0x06 | Error with UTF-8 message                 | Pointer to binary chain     |
| BUILTIN | 0x07 | Built-in function                        | Function ID                 |
| BYTES1  | 0x08 | Binary data cell with 1 byte            | 1 byte in bits 0-7          |
| BYTES2  | 0x09 | Binary data cell with 2 bytes           | 2 bytes in bits 0-15        |
| BYTES3  | 0x0A | Binary data cell with 3 bytes           | 3 bytes in bits 0-23        |
| BYTES4  | 0x0B | Binary data cell with 4 bytes           | 4 bytes in bits 0-31        |

## Data Structure Layouts

### Cons Cells (16 bytes, always 16-byte aligned)

```
Address | +0 | +1 | +2 | +3 | +4 | +5 | +6 | +7 | +8 | +9 | +10| +11| +12| +13| +14| +15|
--------|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|
Content | [         car (i64)         ] | [         cdr (i64)         ]
        | [type|reserved|value/ptr    ] | [type|reserved|value/ptr    ]
```

Cons cells are stored as v128 vectors for atomic load/store operations.

### Binary Data as Linked Lists

Binary data (including symbols and error messages) is stored as linked lists of cons cells, where each cell contains 1-4 bytes:

```
"Hello" (5 bytes) stored as:

Address 0x0010: [BYTES4|"Hell"] [CONS|0x0020]  <- First cell
Address 0x0020: [BYTES1|"o"]    [NIL]          <- Second cell

Complete structure:
Cell 1: car = 0x0B00000048656C6C  (BYTES4 with "Hell")
        cdr = 0x0300000000000020  (CONS pointer to next)
Cell 2: car = 0x080000000000006F  (BYTES1 with "o")
        cdr = 0x0000000000000000  (NIL)
```

### Symbol Representation

Symbols are stored exactly like binary data - as chains of BYTES1/2/3/4 cells:

```
Symbol 'foo stored as:

Address 0x0030: [BYTES3|"foo"] [NIL]

The symbol value is: 0x0200000000000030 (SYMBOL type pointing to the binary chain)
```

Symbol equality requires traversing both binary chains and comparing bytes.

### Error Representation

Errors are also binary chains, just with a different type tag:

```
Error "File not found" stored as:

Address 0x0040: [BYTES4|"File"] [CONS|0x0050]
Address 0x0050: [BYTES4|" not"] [CONS|0x0060]
Address 0x0060: [BYTES4|" fou"] [CONS|0x0070]
Address 0x0070: [BYTES2|"nd"]   [NIL]

The error value is: 0x0600000000000040 (ERROR type pointing to the binary chain)
```

### Lambda/Macro Structure

Functions and macros are cons cells with the structure:

```
(params . (body . env))

Memory layout for (lambda (x y) (+ x y)) with empty environment:

Address 0x0080: [CONS|0x0090]   [CONS|0x00B0]    <- Lambda cell
Address 0x0090: [SYMBOL|x_ptr]  [CONS|0x00A0]    <- Params list
Address 0x00A0: [SYMBOL|y_ptr]  [NIL]            <- Rest of params
Address 0x00B0: [CONS|body_ptr] [NIL]            <- Body and env

The lambda value is: 0x0400000000000080 (LAMBDA type pointing to the structure)
```

### Environment Structure

Environments are association lists mapping symbols to values:

```
Environment ((x . 42) (y . nil)):

Address 0x0100: [CONS|0x0110]   [CONS|0x0130]    <- Env list
Address 0x0110: [SYMBOL|x_ptr]  [NUMBER|42]      <- First binding
Address 0x0130: [CONS|0x0140]   [NIL]            <- Rest of env
Address 0x0140: [SYMBOL|y_ptr]  [NIL]            <- Second binding
```

## Memory Allocation

### Heap Allocator

```wat
;; Global heap pointer (starts at address 0)
(global $heap_ptr (mut i32) (i32.const 0))

;; Allocate 16-byte aligned cons cell
(func $alloc_cons (result i32)
  (local $ptr i32)
  ;; Get current heap pointer
  (local.set $ptr (global.get $heap_ptr))
  ;; Advance by 16 bytes
  (global.set $heap_ptr
    (i32.add (local.get $ptr) (i32.const 16)))
  ;; Return the allocated address
  (local.get $ptr))
```

## Symbol Comparison

Without interning, symbol comparison requires byte-by-byte comparison:

```wat
(func $symbol_equal (param $sym1 i64) (param $sym2 i64) (result i32)
  ;; Extract pointers to binary chains
  (local $ptr1 i32)
  (local $ptr2 i32)
  (local.set $ptr1 (call $get_value (local.get $sym1)))
  (local.set $ptr2 (call $get_value (local.get $sym2)))

  ;; Compare the binary chains byte by byte
  (call $binary_equal (local.get $ptr1) (local.get $ptr2)))
```

## Example Memory Snapshot

Here's what memory might look like after evaluating `(cons 'foo (error "bad"))`:

```
Address | Data                          | Description
--------|-------------------------------|----------------------------------
0x0000  | [0x0B00666F6F] [0x00000000] | "foo" as BYTES3 cell
0x0010  | [0x0B00626164] [0x00000000] | "bad" as BYTES3 cell
0x0020  | [0x0200000000] [0x0600000010] | Main cons: (SYMBOL . ERROR)

Result value: 0x0300000000000020 (CONS pointer to address 0x0020)
```

Breaking down the main cons cell at 0x0020:
- car: 0x0200000000000000 = SYMBOL pointing to "foo" at 0x0000
- cdr: 0x0600000000000010 = ERROR pointing to "bad" at 0x0010

## Advantages of This Memory Model

1. **Simplicity** - One heap, one allocation strategy
2. **Uniformity** - Everything is cons cells
3. **Flexibility** - Symbols can be created dynamically
4. **Debugging** - Easy to dump and inspect memory
5. **No special cases** - All data follows the same patterns

## Disadvantages (Acceptable for a Toy)

1. **Symbol comparison is O(n)** - Must compare bytes
2. **No memory reuse** - Linear allocation only
3. **Memory fragmentation** - No compaction
4. **Larger memory usage** - Symbols stored as full strings

## Implementation Order

1. **Core memory**: Heap allocator, cons cell allocation
2. **Basic types**: nil, numbers, type checking
3. **Binary data**: BYTES1/2/3/4 cell creation and chaining
4. **Symbols**: Binary chains with SYMBOL tag
5. **Lists**: cons, car, cdr operations using v128
6. **Errors**: Binary chains with ERROR tag
7. **Evaluation**: Environment lookup, function application
8. **Built-ins**: Core functions implementation
9. **Special forms**: if, quote, lambda, macro
10. **Reader/Printer**: Text to values and back

This design achieves maximum simplicity by having only one type of memory allocation (16-byte cons cells) and representing everything as chains of these cells.