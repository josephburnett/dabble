#!/usr/bin/env node

const fs = require('fs');
const readline = require('readline');
const { execSync } = require('child_process');
const Reader = require('./reader');
const Printer = require('./printer');

// Compile and load WASM
function loadWasm() {
  console.log('Compiling WAT...');
  execSync('wat2wasm core.wat -o core.wasm', { encoding: 'utf8' });

  const wasmBuffer = fs.readFileSync('core.wasm');
  const wasmModule = new WebAssembly.Module(wasmBuffer);
  const instance = new WebAssembly.Instance(wasmModule);
  return instance.exports;
}

// Create initial environment with built-ins
function makeInitialEnv(exports, reader) {
  let env = exports.nil();

  // Add built-in functions
  const builtins = [
    ['cons', 0],
    ['car', 1],
    ['cdr', 2],
    ['atom', 3],
    ['eq', 4]
  ];

  for (const [name, id] of builtins) {
    const sym = reader.makeSymbol(name);
    const fn = exports.make_builtin(id);
    env = exports.extend(sym, fn, env);
  }

  return env;
}

async function main() {
  const exports = loadWasm();
  const reader = new Reader(exports);
  const printer = new Printer(exports);

  // Create initial environment
  let env = makeInitialEnv(exports, reader);

  console.log('Dabble Lisp REPL');
  console.log('Type expressions to evaluate, or Ctrl+C to exit');
  console.log('');

  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    prompt: '> '
  });

  rl.prompt();

  rl.on('line', (line) => {
    try {
      // Skip empty lines
      if (line.trim() === '') {
        rl.prompt();
        return;
      }

      // Read
      const expr = reader.read(line);

      // Eval
      const result = exports.eval(expr, env);

      // Print
      console.log(printer.print(result));

    } catch (err) {
      console.error('Error:', err.message);
    }

    rl.prompt();
  });

  rl.on('close', () => {
    console.log('\nBye!');
    process.exit(0);
  });
}

main().catch(err => {
  console.error('Fatal error:', err);
  process.exit(1);
});
