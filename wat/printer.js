// Printer: Convert WAT values to S-expression text

class Printer {
  constructor(exports) {
    this.exports = exports;
  }

  // Extract type tag
  getType(val) {
    return Number(val >> 56n) & 0xFF;
  }

  // Extract value/pointer
  getValue(val) {
    return Number(val & 0xFFFFFFFFn);
  }

  // Check if nil
  isNil(val) {
    return val === 0n;
  }

  // Print a value
  print(val) {
    const type = this.getType(val);

    // NIL
    if (type === 0x00) {
      return 'nil';
    }

    // NUMBER
    if (type === 0x01) {
      const value = this.getValue(val);
      // Handle as signed 32-bit integer
      return value > 0x7FFFFFFF ? (value - 0x100000000).toString() : value.toString();
    }

    // SYMBOL
    if (type === 0x02) {
      return this.printSymbol(val);
    }

    // CONS
    if (type === 0x03) {
      return this.printList(val);
    }

    // LAMBDA
    if (type === 0x04) {
      return '#<lambda>';
    }

    // MACRO
    if (type === 0x05) {
      return '#<macro>';
    }

    // ERROR
    if (type === 0x06) {
      return '#<error: ' + this.printBinary(this.exports.error_message(val)) + '>';
    }

    // BUILTIN
    if (type === 0x07) {
      const id = this.getValue(val);
      const names = ['cons', 'car', 'cdr', 'atom', 'eq'];
      return '#<builtin:' + (names[id] || id) + '>';
    }

    // BYTES*
    if (type >= 0x08 && type <= 0x0B) {
      return '#<bytes>';
    }

    return '#<unknown:' + type.toString(16) + '>';
  }

  // Print a symbol
  printSymbol(sym) {
    // Convert SYMBOL to CONS to access binary chain
    const chain = (sym & 0x00000000FFFFFFFFn) | 0x0300000000000000n;
    return this.printBinary(chain);
  }

  // Print binary data as string
  printBinary(chain) {
    const bytes = [];
    let current = chain;

    while (!this.isNil(current)) {
      const car = this.exports.car(current);
      const carType = this.getType(car);
      const carValue = this.getValue(car);

      // Extract bytes based on type
      const count = carType - 7; // BYTES1=8 -> 1, BYTES2=9 -> 2, etc.

      for (let i = 0; i < count; i++) {
        const byte = (carValue >> (i * 8)) & 0xFF;
        bytes.push(byte);
      }

      current = this.exports.cdr(current);
    }

    // Convert bytes to string
    return String.fromCharCode(...bytes);
  }

  // Print a list (handles both proper lists and dotted pairs)
  printList(list) {
    // Check if it's a proper list
    if (this.isProperList(list)) {
      const elements = [];
      let current = list;

      while (!this.isNil(current)) {
        elements.push(this.print(this.exports.car(current)));
        current = this.exports.cdr(current);
      }

      return '(' + elements.join(' ') + ')';
    } else {
      // Dotted pair
      const car = this.exports.car(list);
      const cdr = this.exports.cdr(list);
      return '(' + this.print(car) + ' . ' + this.print(cdr) + ')';
    }
  }

  // Check if a value is a proper list (ends in nil)
  isProperList(list) {
    let current = list;

    while (!this.isNil(current)) {
      if (this.getType(current) !== 0x03) {
        // Not a cons cell - improper
        return false;
      }
      current = this.exports.cdr(current);
    }

    return true;
  }
}

module.exports = Printer;
