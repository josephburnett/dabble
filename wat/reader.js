// Reader: Parse S-expression text into WAT values

class Reader {
  constructor(exports) {
    this.exports = exports;
  }

  // Tokenize input string
  tokenize(input) {
    // Simple regex-based tokenizer
    const tokens = [];
    const regex = /\s*([()]|'|[^\s()]+)/g;
    let match;
    while ((match = regex.exec(input)) !== null) {
      if (match[1]) tokens.push(match[1]);
    }
    return tokens;
  }

  // Parse a single expression from token stream
  parseExpr(tokens) {
    if (tokens.length === 0) {
      throw new Error('Unexpected end of input');
    }

    const token = tokens.shift();

    // Quote shorthand
    if (token === "'") {
      const quoted = this.parseExpr(tokens);
      // Build (quote <expr>)
      const quoteSymbol = this.makeSymbol('quote');
      return this.exports.cons(quoteSymbol, this.exports.cons(quoted, this.exports.nil()));
    }

    // List
    if (token === '(') {
      const list = [];
      while (tokens.length > 0 && tokens[0] !== ')') {
        list.push(this.parseExpr(tokens));
      }
      if (tokens.length === 0) {
        throw new Error('Missing closing paren');
      }
      tokens.shift(); // consume ')'

      // Build list from right to left
      let result = this.exports.nil();
      for (let i = list.length - 1; i >= 0; i--) {
        result = this.exports.cons(list[i], result);
      }
      return result;
    }

    if (token === ')') {
      throw new Error('Unexpected closing paren');
    }

    // Atom - number or symbol
    return this.parseAtom(token);
  }

  parseAtom(token) {
    // Try parsing as number
    const num = parseInt(token, 10);
    if (!isNaN(num) && token === num.toString()) {
      return this.exports.make_number(num);
    }

    // Negative number
    if (token.startsWith('-')) {
      const num = parseInt(token, 10);
      if (!isNaN(num)) {
        return this.exports.make_number(num);
      }
    }

    // Symbol
    return this.makeSymbol(token);
  }

  // Build a symbol from a string
  makeSymbol(str) {
    const bytes = [];
    for (let i = 0; i < str.length; i++) {
      bytes.push(str.charCodeAt(i));
    }

    // Build binary chain from bytes
    let chain = this.exports.nil();

    // Process bytes in chunks of 4, 3, 2, 1
    let i = 0;
    while (i < bytes.length) {
      const remaining = bytes.length - i;

      if (remaining >= 4) {
        // Pack 4 bytes (little-endian)
        const val = bytes[i] | (bytes[i+1] << 8) | (bytes[i+2] << 16) | (bytes[i+3] << 24);
        chain = this.exports.cons(this.exports.make_bytes4(val), chain);
        i += 4;
      } else if (remaining === 3) {
        const val = bytes[i] | (bytes[i+1] << 8) | (bytes[i+2] << 16);
        chain = this.exports.cons(this.exports.make_bytes3(val), chain);
        i += 3;
      } else if (remaining === 2) {
        const val = bytes[i] | (bytes[i+1] << 8);
        chain = this.exports.cons(this.exports.make_bytes2(val), chain);
        i += 2;
      } else {
        chain = this.exports.cons(this.exports.make_bytes1(bytes[i]), chain);
        i += 1;
      }
    }

    // Reverse the chain (we built it backwards)
    let reversed = this.exports.nil();
    while (!this.isNil(chain)) {
      reversed = this.exports.cons(this.exports.car(chain), reversed);
      chain = this.exports.cdr(chain);
    }

    return this.exports.make_symbol(reversed);
  }

  // Helper to check if value is nil
  isNil(val) {
    return val === 0n;
  }

  // Main read function
  read(input) {
    const tokens = this.tokenize(input);
    if (tokens.length === 0) {
      return this.exports.nil();
    }
    return this.parseExpr(tokens);
  }
}

module.exports = Reader;
