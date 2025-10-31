# Dabble

A toy Lisp interpreter with the core implemented in multiple languages. The goal is to implement the absolute minimum primitives in each host language, then build the rest of the language in Lisp itself.

## Implementations

- **Go** (`/go/`) - Most complete with REPL and automatic bootstrap loading
- **C** (`/c/`) - Feature-complete reference implementation
- **Lua** (`/lua/`) - Clean implementation with try/throw extensions
- **WebAssembly** (`/wat/`) - Experimental stub using SIMD vectors

## Core Primitives

Each implementation provides ~15 primitives:
- List operations: `cons`, `car`, `cdr`, `atom`
- Control flow: `if`/`cond`, `quote`, `eq`
- Functions: `label`, `lambda`, `macro`, `apply`, `recur`
- Error handling: `error`

## Bootstrap Library

Higher-level functions written in Lisp itself (`/src/core/`):
- `let` - Variable binding
- `and`, `or`, `not` - Logical operations
- `list` - List construction

## Examples

```lisp
; Define a function
(label last
  (lambda (x)
    (cond
      (eq () (cdr x)) (car x)
      1 (recur (cdr x)))))

(last (1 2 3 4))  ; Returns 4

; Macros with hygiene
(label m
  (macro (x xs)
    (cond
      (eq x 'y) 1
      (eq x 'z) 2)))

(m 'z)  ; Returns 2

; Recursive functions
(label fact
  (lambda (n)
    (cond
      (eq n 0) 1
      1 (cons n (fact (cons (car n) (cdr n)))))))
```

## Running

### Go (with REPL)
```bash
cd go && go run main.go
```

### C (file evaluation)
```bash
cd c && make
./dabble filename.lisp
```

## Testing

```bash
make test
```

Or watch for changes:
```bash
while true; do inotifywait -e close_write -r . 2>/dev/null ; make test; clear; ./test; done
```