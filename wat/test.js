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
