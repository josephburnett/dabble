#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <inttypes.h>

typedef enum { LIST, SYMBOL, NUMBER, FUNC, ERROR } type;

typedef int64_t value;
typedef struct cell {
    type t;
    value car;
    struct cell* cdr;
} cell;

cell* NIL = 0;

cell* new_cell(type t, value car, cell* cdr)
{
    cell* c = malloc(sizeof(cell));
    c->t = t;
    c->car = car;
    c->cdr = cdr;
    return c;
}

int is_error(cell* c)
{
    if (c != NIL && c->t == ERROR) return 1;
    return 0;
}

typedef enum { R_WHITESPACE, R_SYMBOL, R_NUMBER } reading;

cell* string(FILE *fp);

cell* read(FILE *fp)
{
    int ch;
    int index = 0;
    value v = 0;
    cell* c = NIL;
    value sign = 1;
    reading rd = R_WHITESPACE;
    while ((ch = getc(fp)) != EOF) {
        switch (ch) {
        case ' ':
            switch (rd) {
            case R_WHITESPACE:
                continue;
            case R_SYMBOL:
                c = read(fp);
                if (is_error(c)) return c;
                return new_cell(SYMBOL, v, c);
            case R_NUMBER:
                c = read(fp);
                if (is_error(c)) return c;
                return new_cell(NUMBER, v * sign, c);
            }
        case '(': {
            c = read(fp);
            if (is_error(c)) return c;
            if (c == NIL) return new_cell(ERROR, 0, NIL);
            cell* cl = new_cell(LIST, (value) c, NIL);
            c = read(fp);
            if (is_error(c)) return c;
            cl->cdr = c;
            return cl;
        }
        case ')':
            switch (rd) {
            case R_SYMBOL:
                return new_cell(SYMBOL, v, NIL);
            case R_NUMBER:
                return new_cell(NUMBER, v * sign, NIL);
            default:
                return NIL;
            }
        case '-':
            switch (rd) {
            case R_WHITESPACE:
                sign = -1;
                rd = R_NUMBER;
                continue;
            default:
                return new_cell(ERROR, 0, NIL);
            }
        case '0' ... '9':
            switch (rd) {
            case R_WHITESPACE:
            case R_NUMBER:
                rd = R_NUMBER;
                ch = ch - '0';
                v = v * 10 + ch;
                continue;
            default:
                return new_cell(ERROR, 0, NIL);
            }
        case 'a' ... 'z':
            switch (rd) {
            case R_WHITESPACE:
            case R_SYMBOL:
                rd = R_SYMBOL;
                if (index == 7) {
                    return new_cell(ERROR, 0, NIL);
                }
                ((char *) &v)[index++] = ch;
                continue;
            }
        case '"': {
            c = string(fp);
            if (is_error(c)) return c;
            if (c == NIL) return new_cell(ERROR, 0, NIL);
            cell* cl = new_cell(LIST, (value) c, NIL);
            c = read(fp);
            if (is_error(c)) return c;
            cl->cdr = c;
            return cl;
        }
        default:
            return new_cell(ERROR, 0, NIL);
        }
    }
    return (cell*) NIL;
}

cell* string(FILE *fp)
{
    int ch;
    value v = 0;
    cell* c = NIL;
    while ((ch = getc(fp)) != EOF) {
        if (ch == '"') {
            return NIL;
        } else {
            ((char*) &v)[0] = ch;
            c = string(fp);
            if (is_error(c)) return c;
            return new_cell(SYMBOL, v, c);
        }
    }
}

void print_index(FILE *fp, cell* c, int index)
{
    if (index > 0) {
        fprintf(fp, " ");
    }
    switch (c->t) {
    case LIST:
        fprintf(fp, "(");
        print_index(fp, (cell*) c->car, 0);
        fprintf(fp, ")");
        break;
    case SYMBOL: {
        char* sym = (char*) &(c->car);
        int i;
        for (i = 0; i < 8; i++) {
            if (sym[i] != 0) {
                fprintf(fp, "%c", sym[i]);
            }
        }
        break;
    }
    case NUMBER:
        fprintf(fp, "%" PRId64, c->car);
        break;
    case FUNC:
        fprintf(fp, "<func>");
        break;
    case ERROR:
        fprintf(fp, "<error>");
        break;
    }
    if (c->cdr != NIL) {
        print_index(fp, (cell*) c->cdr, index + 1);
    }
}

void print(FILE *fp, cell* c)
{
    print_index(fp, c, 0);
}
