#!/usr/bin/env node

const fs = require('fs');
const { execSync } = require('child_process');

// Compile WAT to WASM using wat2wasm (from wabt)
function compileWat(watPath) {
  const wasmPath = watPath.replace('.wat', '.wasm');
  try {
    execSync(`wat2wasm ${watPath} -o ${wasmPath}`, { encoding: 'utf8' });
    return wasmPath;
  } catch (e) {
    console.error('Failed to compile WAT:', e.message);
    process.exit(1);
  }
}

// Load and instantiate WASM module
async function loadWasm(wasmPath) {
  const wasmBuffer = fs.readFileSync(wasmPath);
  const wasmModule = await WebAssembly.compile(wasmBuffer);
  const instance = await WebAssembly.instantiate(wasmModule);
  return instance.exports;
}

// Test utilities
let passCount = 0;
let failCount = 0;

function assert(condition, message) {
  if (condition) {
    passCount++;
    console.log(`✓ ${message}`);
  } else {
    failCount++;
    console.error(`✗ ${message}`);
  }
}

function assertEquals(actual, expected, message) {
  if (actual === expected) {
    passCount++;
    console.log(`✓ ${message}`);
  } else {
    failCount++;
    console.error(`✗ ${message}`);
    console.error(`  Expected: ${expected} (0x${expected.toString(16)})`);
    console.error(`  Actual:   ${actual} (0x${actual.toString(16)})`);
  }
}

// Extract type tag from i64 value
function getType(val) {
  return Number(val >> 56n) & 0xFF;
}

// Extract value/pointer from i64
function getValue(val) {
  return Number(val & 0xFFFFFFFFn);
}

// Main test runner
async function runTests() {
  console.log('Compiling core.wat...\n');
  const wasmPath = compileWat('./core.wat');

  console.log('Loading WASM module...\n');
  const exports = await loadWasm(wasmPath);

  console.log('Running tests...\n');

  // Test basic value helpers
  console.log('--- Basic Value Helpers ---');

  // Test nil
  const nilVal = exports.nil();
  assertEquals(nilVal, 0n, 'nil() returns 0');

  // Test make_number
  const num42 = exports.make_number(42);
  assertEquals(getType(num42), 0x01, 'make_number(42) has NUMBER type');
  assertEquals(getValue(num42), 42, 'make_number(42) has value 42');

  const numNeg = exports.make_number(-5);
  assertEquals(getType(numNeg), 0x01, 'make_number(-5) has NUMBER type');
  assertEquals(getValue(numNeg), 0xFFFFFFFB, 'make_number(-5) has value -5 as u32');

  // Test get_type
  assertEquals(exports.get_type(0n), 0, 'get_type(nil) returns 0');
  assertEquals(exports.get_type(num42), 0x01, 'get_type(NUMBER) returns 0x01');

  // Test get_value
  assertEquals(exports.get_value(num42), 42, 'get_value(NUMBER 42) returns 42');

  // Test is_nil
  assertEquals(exports.is_nil(0n), 1, 'is_nil(0) returns true');
  assertEquals(exports.is_nil(num42), 0, 'is_nil(NUMBER) returns false');

  // Test allocator
  console.log('\n--- Memory Allocation ---');
  const ptr1 = exports.alloc_cons();
  const ptr2 = exports.alloc_cons();
  const ptr3 = exports.alloc_cons();
  assertEquals(ptr1, 0, 'First allocation at address 0');
  assertEquals(ptr2, 16, 'Second allocation at address 16');
  assertEquals(ptr3, 32, 'Third allocation at address 32');

  // Test cons/car/cdr
  console.log('\n--- List Operations ---');

  // Simple cons cell: (42 . nil)
  const cell1 = exports.cons(num42, nilVal);
  assertEquals(getType(cell1), 0x03, 'cons returns CONS type');
  assertEquals(getValue(cell1), 48, 'cons pointer points to allocated memory');

  const car1 = exports.car(cell1);
  assertEquals(car1, num42, 'car returns the first element');

  const cdr1 = exports.cdr(cell1);
  assertEquals(cdr1, nilVal, 'cdr returns the second element');

  // Nested cons: (1 . (2 . (3 . nil)))
  const num1 = exports.make_number(1);
  const num2 = exports.make_number(2);
  const num3 = exports.make_number(3);

  const list = exports.cons(num1,
                exports.cons(num2,
                  exports.cons(num3, nilVal)));

  assertEquals(getType(list), 0x03, 'nested cons has CONS type');

  // Traverse the list
  const first = exports.car(list);
  assertEquals(getValue(first), 1, 'first element is 1');

  const rest1 = exports.cdr(list);
  const second = exports.car(rest1);
  assertEquals(getValue(second), 2, 'second element is 2');

  const rest2 = exports.cdr(rest1);
  const third = exports.car(rest2);
  assertEquals(getValue(third), 3, 'third element is 3');

  const rest3 = exports.cdr(rest2);
  assertEquals(rest3, 0n, 'end of list is nil');

  // Test atom
  console.log('\n--- Atom Predicate ---');
  const atomNil = exports.atom(nilVal);
  assert(getValue(atomNil) === 1, 'atom(nil) returns true');

  const atomNum = exports.atom(num42);
  assert(getValue(atomNum) === 1, 'atom(NUMBER) returns true');

  const atomCons = exports.atom(cell1);
  assertEquals(atomCons, 0n, 'atom(CONS) returns false (nil)');

  // Test eq
  console.log('\n--- Equality ---');
  const eqNils = exports.eq(nilVal, nilVal);
  assert(getValue(eqNils) === 1, 'eq(nil, nil) returns true');

  const eqNums = exports.eq(num42, num42);
  assert(getValue(eqNums) === 1, 'eq(42, 42) returns true');

  const eqDiff = exports.eq(num42, num1);
  assertEquals(eqDiff, 0n, 'eq(42, 1) returns false');

  const eqCons = exports.eq(cell1, cell1);
  assert(getValue(eqCons) === 1, 'eq(same cons, same cons) returns true');

  const cell2 = exports.cons(num42, nilVal);
  const eqDiffCons = exports.eq(cell1, cell2);
  assertEquals(eqDiffCons, 0n, 'eq(different cons, different cons) returns false');

  // Test binary data
  console.log('\n--- Binary Data ---');

  // BYTES1: single byte
  const byte1 = exports.make_bytes1(0x41);  // 'A'
  assertEquals(getType(byte1), 0x08, 'make_bytes1 has BYTES1 type');
  assertEquals(getValue(byte1), 0x41, 'make_bytes1 stores 1 byte');
  assertEquals(exports.get_byte_count(byte1), 1, 'BYTES1 has count 1');

  // BYTES2: two bytes
  const byte2 = exports.make_bytes2(0x4142);  // 'AB'
  assertEquals(getType(byte2), 0x09, 'make_bytes2 has BYTES2 type');
  assertEquals(getValue(byte2), 0x4142, 'make_bytes2 stores 2 bytes');
  assertEquals(exports.get_byte_count(byte2), 2, 'BYTES2 has count 2');

  // BYTES3: three bytes
  const byte3 = exports.make_bytes3(0x414243);  // 'ABC'
  assertEquals(getType(byte3), 0x0A, 'make_bytes3 has BYTES3 type');
  assertEquals(getValue(byte3), 0x414243, 'make_bytes3 stores 3 bytes');
  assertEquals(exports.get_byte_count(byte3), 3, 'BYTES3 has count 3');

  // BYTES4: four bytes
  const byte4 = exports.make_bytes4(0x41424344);  // 'ABCD'
  assertEquals(getType(byte4), 0x0B, 'make_bytes4 has BYTES4 type');
  assertEquals(getValue(byte4), 0x41424344, 'make_bytes4 stores 4 bytes');
  assertEquals(exports.get_byte_count(byte4), 4, 'BYTES4 has count 4');

  // Binary data in cons cells (building blocks for symbols/errors)
  const binaryChain = exports.cons(byte4,
                        exports.cons(byte1, nilVal));
  assertEquals(getType(binaryChain), 0x03, 'binary chain is CONS type');
  const firstBytes = exports.car(binaryChain);
  assertEquals(getType(firstBytes), 0x0B, 'car of chain is BYTES4');
  assertEquals(getValue(firstBytes), 0x41424344, 'first chunk has correct bytes');

  // Test symbols
  console.log('\n--- Symbols ---');

  // Create symbol 'foo (3 bytes)
  const fooBytes = exports.make_bytes3(0x666F6F);  // "foo"
  const fooChain = exports.cons(fooBytes, nilVal);
  const symFoo = exports.make_symbol(fooChain);
  assertEquals(getType(symFoo), 0x02, 'make_symbol creates SYMBOL type');
  assertEquals(getValue(symFoo), getValue(fooChain), 'symbol points to binary chain');

  // Create another symbol 'foo
  const fooChain2 = exports.cons(exports.make_bytes3(0x666F6F), nilVal);
  const symFoo2 = exports.make_symbol(fooChain2);

  // Create symbol 'bar (3 bytes)
  const barBytes = exports.make_bytes3(0x626172);  // "bar"
  const barChain = exports.cons(barBytes, nilVal);
  const symBar = exports.make_symbol(barChain);

  // Test symbol equality
  const eqFooFoo = exports.symbol_equal(symFoo, symFoo2);
  assert(getValue(eqFooFoo) === 1, 'symbol_equal("foo", "foo") returns true');

  const eqFooBar = exports.symbol_equal(symFoo, symBar);
  assertEquals(eqFooBar, 0n, 'symbol_equal("foo", "bar") returns false');

  // Multi-byte symbol: 'hello (5 bytes = BYTES4 + BYTES1)
  const helloChain = exports.cons(
    exports.make_bytes4(0x68656C6C),  // "hell"
    exports.cons(
      exports.make_bytes1(0x6F),       // "o"
      nilVal));
  const symHello = exports.make_symbol(helloChain);
  assertEquals(getType(symHello), 0x02, 'multi-byte symbol has SYMBOL type');

  const helloChain2 = exports.cons(
    exports.make_bytes4(0x68656C6C),
    exports.cons(exports.make_bytes1(0x6F), nilVal));
  const symHello2 = exports.make_symbol(helloChain2);

  const eqHello = exports.symbol_equal(symHello, symHello2);
  assert(getValue(eqHello) === 1, 'symbol_equal("hello", "hello") returns true');

  // Test errors
  console.log('\n--- Errors ---');

  // Create error with message "bad"
  const badChain = exports.cons(exports.make_bytes3(0x626164), nilVal);
  const errBad = exports.make_error(badChain);
  assertEquals(getType(errBad), 0x06, 'make_error creates ERROR type');

  // Extract message from error
  const msgChain = exports.error_message(errBad);
  assertEquals(getType(msgChain), 0x03, 'error_message returns CONS chain');
  const msgBytes = exports.car(msgChain);
  assertEquals(getValue(msgBytes), 0x626164, 'error message contains "bad"');

  // Test environment operations
  console.log('\n--- Environment Operations ---');

  // Create empty environment
  let env = nilVal;

  // Extend with x = 42
  env = exports.extend(symFoo, num42, env);
  assertEquals(getType(env), 0x03, 'extend returns CONS');

  // Lookup x
  const lookedUp = exports.lookup(symFoo, env);
  assertEquals(lookedUp, num42, 'lookup finds value in environment');

  // Extend with y = 1
  env = exports.extend(symBar, num1, env);

  // Lookup both
  const lookupFoo = exports.lookup(symFoo, env);
  assertEquals(lookupFoo, num42, 'lookup finds first binding');

  const lookupBar = exports.lookup(symBar, env);
  assertEquals(lookupBar, num1, 'lookup finds second binding');

  // Lookup non-existent symbol
  const symBaz = exports.make_symbol(exports.cons(exports.make_bytes3(0x62617A), nilVal)); // "baz"
  const lookupBaz = exports.lookup(symBaz, env);
  assertEquals(getType(lookupBaz), 0x06, 'lookup returns ERROR for undefined symbol');

  // Shadow a binding
  const num99 = exports.make_number(99);
  env = exports.extend(symFoo, num99, env);
  const lookupShadowed = exports.lookup(symFoo, env);
  assertEquals(lookupShadowed, num99, 'lookup returns shadowed value');

  // Test evaluation
  console.log('\n--- Evaluation Engine ---');

  // Reset environment
  env = nilVal;

  // Self-evaluating values
  const evalNil = exports.eval(nilVal, env);
  assertEquals(evalNil, nilVal, 'eval(nil) returns nil');

  const evalNum = exports.eval(num42, env);
  assertEquals(evalNum, num42, 'eval(42) returns 42');

  // Symbol lookup
  env = exports.extend(symFoo, num42, env);
  const evalSym = exports.eval(symFoo, env);
  assertEquals(evalSym, num42, 'eval(foo) looks up value');

  // Undefined symbol
  const evalUndef = exports.eval(symBar, env);
  assertEquals(getType(evalUndef), 0x06, 'eval(undefined) returns error');

  // Built-in functions
  const builtinCons = exports.make_builtin(0);  // cons
  const builtinCar = exports.make_builtin(1);   // car
  const builtinCdr = exports.make_builtin(2);   // cdr
  const builtinAtom = exports.make_builtin(3);  // atom
  const builtinEq = exports.make_builtin(4);    // eq

  assertEquals(getType(builtinCons), 0x07, 'make_builtin creates BUILTIN type');

  // Add built-ins to environment
  const symCons = exports.make_symbol(exports.cons(exports.make_bytes4(0x636F6E73), nilVal)); // "cons"
  const symCar = exports.make_symbol(exports.cons(exports.make_bytes3(0x636172), nilVal));    // "car"

  env = exports.extend(symCons, builtinCons, env);
  env = exports.extend(symCar, builtinCar, env);

  // Evaluate: (cons 1 2)
  const consExpr = exports.cons(symCons,
    exports.cons(num1, exports.cons(exports.make_number(2), nilVal)));
  const consResult = exports.eval(consExpr, env);
  assertEquals(getType(consResult), 0x03, 'eval((cons 1 2)) returns CONS');
  assertEquals(getValue(exports.car(consResult)), 1, 'car of result is 1');
  assertEquals(getValue(exports.cdr(consResult)), 2, 'cdr of result is 2');

  // Evaluate: (car (cons 10 20))
  const num10 = exports.make_number(10);
  const num20 = exports.make_number(20);
  const innerCons = exports.cons(symCons,
    exports.cons(num10, exports.cons(num20, nilVal)));
  const carExpr = exports.cons(symCar, exports.cons(innerCons, nilVal));
  const carResult = exports.eval(carExpr, env);
  assertEquals(getValue(carResult), 10, 'eval((car (cons 10 20))) returns 10');

  // Error propagation
  const errorExpr = exports.cons(symCar, exports.cons(symBaz, nilVal)); // (car baz) where baz is undefined
  const errorResult = exports.eval(errorExpr, env);
  assertEquals(getType(errorResult), 0x06, 'eval propagates errors from arguments');

  // Test special forms
  console.log('\n--- Special Forms ---');

  // Reset environment
  env = nilVal;

  // quote
  const symQuote = exports.make_symbol(exports.cons(exports.make_bytes4(0x71756F74),
    exports.cons(exports.make_bytes1(0x65), nilVal))); // "quote"

  // (quote 42) -> 42 (unevaluated)
  const quoteExpr = exports.cons(symQuote, exports.cons(num42, nilVal));
  const quoteResult = exports.eval(quoteExpr, env);
  assertEquals(quoteResult, num42, 'eval((quote 42)) returns 42 unevaluated');

  // (quote foo) -> foo (symbol, not looked up)
  const quoteSym = exports.cons(symQuote, exports.cons(symFoo, nilVal));
  const quoteSymResult = exports.eval(quoteSym, env);
  assertEquals(quoteSymResult, symFoo, 'eval((quote foo)) returns symbol unevaluated');

  // (quote (1 2 3)) -> (1 2 3) (list, unevaluated)
  const list123 = exports.cons(num1, exports.cons(exports.make_number(2),
    exports.cons(exports.make_number(3), nilVal)));
  const quoteList = exports.cons(symQuote, exports.cons(list123, nilVal));
  const quoteListResult = exports.eval(quoteList, env);
  assertEquals(quoteListResult, list123, 'eval((quote (1 2 3))) returns list unevaluated');

  // if
  const symIf = exports.make_symbol(exports.cons(exports.make_bytes2(0x6966), nilVal)); // "if"

  // (if 1 2 3) -> 2 (condition is truthy)
  const ifTrue = exports.cons(symIf,
    exports.cons(num1,
      exports.cons(exports.make_number(2),
        exports.cons(exports.make_number(3), nilVal))));
  const ifTrueResult = exports.eval(ifTrue, env);
  assertEquals(getValue(ifTrueResult), 2, 'eval((if 1 2 3)) returns then branch');

  // (if nil 2 3) -> 3 (condition is nil)
  const ifFalse = exports.cons(symIf,
    exports.cons(nilVal,
      exports.cons(exports.make_number(2),
        exports.cons(exports.make_number(3), nilVal))));
  const ifFalseResult = exports.eval(ifFalse, env);
  assertEquals(getValue(ifFalseResult), 3, 'eval((if nil 2 3)) returns else branch');

  // label
  const symLabel = exports.make_symbol(exports.cons(exports.make_bytes4(0x6C616265),
    exports.cons(exports.make_bytes1(0x6C), nilVal))); // "label"

  // (label x 42 x) -> 42
  const labelExpr = exports.cons(symLabel,
    exports.cons(symFoo,
      exports.cons(num42,
        exports.cons(symFoo, nilVal))));
  const labelResult = exports.eval(labelExpr, env);
  assertEquals(labelResult, num42, 'eval((label x 42 x)) binds and returns value');

  // (label x 10 (label y 20 (cons x y))) -> (10 . 20)
  const symX = exports.make_symbol(exports.cons(exports.make_bytes1(0x78), nilVal)); // "x"
  const symY = exports.make_symbol(exports.cons(exports.make_bytes1(0x79), nilVal)); // "y"
  env = exports.extend(symCons, builtinCons, nilVal); // Need cons in env

  const innerLabel = exports.cons(symLabel,
    exports.cons(symY,
      exports.cons(num20,
        exports.cons(exports.cons(symCons,
          exports.cons(symX,
            exports.cons(symY, nilVal))), nilVal))));

  const outerLabel = exports.cons(symLabel,
    exports.cons(symX,
      exports.cons(num10,
        exports.cons(innerLabel, nilVal))));

  const nestedLabelResult = exports.eval(outerLabel, env);
  assertEquals(getType(nestedLabelResult), 0x03, 'nested label returns CONS');
  assertEquals(getValue(exports.car(nestedLabelResult)), 10, 'nested label car is 10');
  assertEquals(getValue(exports.cdr(nestedLabelResult)), 20, 'nested label cdr is 20');

  // Test lambda and macro
  console.log('\n--- Lambda ---');

  // Reset environment with built-ins
  env = nilVal;
  env = exports.extend(symCons, builtinCons, env);
  env = exports.extend(symCar, builtinCar, env);

  const symLambda = exports.make_symbol(exports.cons(exports.make_bytes4(0x6C616D62),
    exports.cons(exports.make_bytes2(0x6461), nilVal))); // "lambda"

  // Create lambda: (lambda (x) x)
  const params = exports.cons(symX, nilVal);  // (x)
  const body = symX;  // x
  const lambdaExpr = exports.cons(symLambda,
    exports.cons(params,
      exports.cons(body, nilVal)));

  const identityFn = exports.eval(lambdaExpr, env);
  assertEquals(getType(identityFn), 0x04, 'eval((lambda (x) x)) creates LAMBDA type');

  // Apply identity function: ((lambda (x) x) 42)
  const applyIdentity = exports.cons(lambdaExpr,
    exports.cons(num42, nilVal));
  const identityResult = exports.eval(applyIdentity, env);
  assertEquals(getValue(identityResult), 42, 'eval(((lambda (x) x) 42)) returns 42');

  // Lambda with two params: (lambda (x y) (cons x y))
  const params2 = exports.cons(symX, exports.cons(symY, nilVal));  // (x y)
  const consBody = exports.cons(symCons,
    exports.cons(symX,
      exports.cons(symY, nilVal)));  // (cons x y)
  const pairFnExpr = exports.cons(symLambda,
    exports.cons(params2,
      exports.cons(consBody, nilVal)));

  const applyPair = exports.cons(pairFnExpr,
    exports.cons(num1,
      exports.cons(num2, nilVal)));  // ((lambda (x y) (cons x y)) 1 2)
  const pairResult = exports.eval(applyPair, env);
  assertEquals(getType(pairResult), 0x03, '((lambda (x y) (cons x y)) 1 2) returns CONS');
  assertEquals(getValue(exports.car(pairResult)), 1, 'car is 1');
  assertEquals(getValue(exports.cdr(pairResult)), 2, 'cdr is 2');

  // Lambda with closure: (label x 10 (lambda (y) (cons x y)))
  const closureLambda = exports.cons(symLambda,
    exports.cons(exports.cons(symY, nilVal),
      exports.cons(exports.cons(symCons,
        exports.cons(symX,
          exports.cons(symY, nilVal))), nilVal)));

  const closureExpr = exports.cons(symLabel,
    exports.cons(symX,
      exports.cons(num10,
        exports.cons(closureLambda, nilVal))));  // (label x 10 (lambda (y) (cons x y)))

  const closureFn = exports.eval(closureExpr, env);
  assertEquals(getType(closureFn), 0x04, 'lambda with closure creates LAMBDA');

  // Apply closure: ((label x 10 (lambda (y) (cons x y))) 20)
  const applyClosureExpr = exports.cons(symLabel,
    exports.cons(symX,
      exports.cons(num10,
        exports.cons(exports.cons(closureLambda,
          exports.cons(num20, nilVal)), nilVal))));

  const closureResult = exports.eval(applyClosureExpr, env);
  assertEquals(getType(closureResult), 0x03, 'lambda closure application returns CONS');
  assertEquals(getValue(exports.car(closureResult)), 10, 'closure captured x=10');
  assertEquals(getValue(exports.cdr(closureResult)), 20, 'argument y=20');

  // Test macro
  console.log('\n--- Macro ---');

  const symMacro = exports.make_symbol(exports.cons(exports.make_bytes4(0x6D616372),
    exports.cons(exports.make_bytes1(0x6F), nilVal))); // "macro"

  // Create macro: (macro (x) (cons (quote cons) (cons x (cons x nil))))
  // This builds the list (cons x x) as data
  const macroParams = exports.cons(symX, nilVal);
  const quoteConsExpr = exports.cons(symQuote, exports.cons(symCons, nilVal));  // (quote cons)
  const macroBody = exports.cons(symCons,
    exports.cons(quoteConsExpr,
      exports.cons(exports.cons(symCons,
        exports.cons(symX,
          exports.cons(exports.cons(symCons,
            exports.cons(symX, exports.cons(nilVal, nilVal))), nilVal))), nilVal)));  // (cons (quote cons) (cons x (cons x nil)))

  const dupMacroExpr = exports.cons(symMacro,
    exports.cons(macroParams,
      exports.cons(macroBody, nilVal)));

  const dupMacro = exports.eval(dupMacroExpr, env);
  assertEquals(getType(dupMacro), 0x05, 'eval((macro (x) (cons x x))) creates MACRO type');

  // Store macro in environment
  const symDup = exports.make_symbol(exports.cons(exports.make_bytes3(0x647570), nilVal)); // "dup"
  env = exports.extend(symDup, dupMacro, env);

  // Debug: lookup dup directly
  const lookedUpDup = exports.lookup(symDup, env);
  assertEquals(getType(lookedUpDup), 0x05, 'lookup(dup) returns MACRO type');
  assertEquals(lookedUpDup, dupMacro, 'lookup(dup) returns same macro');

  // Apply macro: (dup 5) -> (cons 5 5) -> (5 . 5)
  const applyMacro = exports.cons(symDup,
    exports.cons(exports.make_number(5), nilVal));

  const macroResult = exports.eval(applyMacro, env);
  assertEquals(getType(macroResult), 0x03, '(dup 5) returns CONS');
  assertEquals(getValue(exports.car(macroResult)), 5, 'macro expansion car is 5');
  assertEquals(getValue(exports.cdr(macroResult)), 5, 'macro expansion cdr is 5');

  // Error cases
  console.log('\n--- Error Cases ---');

  // Too few arguments
  const tooFewArgs = exports.cons(pairFnExpr,
    exports.cons(num1, nilVal));  // ((lambda (x y) ...) 1)
  const tooFewResult = exports.eval(tooFewArgs, env);
  assertEquals(getType(tooFewResult), 0x06, 'too few args returns error');

  // Too many arguments
  const tooManyArgs = exports.cons(pairFnExpr,
    exports.cons(num1,
      exports.cons(num2,
        exports.cons(num3, nilVal))));  // ((lambda (x y) ...) 1 2 3)
  const tooManyResult = exports.eval(tooManyArgs, env);
  assertEquals(getType(tooManyResult), 0x06, 'too many args returns error');

  // Print summary
  console.log('\n' + '='.repeat(50));
  console.log(`Tests passed: ${passCount}`);
  console.log(`Tests failed: ${failCount}`);
  console.log('='.repeat(50));

  process.exit(failCount > 0 ? 1 : 0);
}

runTests().catch(err => {
  console.error('Test runner error:', err);
  process.exit(1);
});
