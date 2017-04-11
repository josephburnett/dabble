#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <inttypes.h>

typedef enum { LIST, SYMBOL, NUMBER, FUNC } type;

typedef int64_t value;
typedef struct {
    type t;
    value car;
    value cdr;
} cell;

cell* new_cell(type t, value car, value cdr)
{
    cell* c = malloc(sizeof(cell));
    c->t = t;
    c->car = car;
    c->cdr = cdr;
    return c;
}

value NIL = 0;


typedef enum { R_WHITESPACE, R_SYMBOL, R_NUMBER } reading;

cell* string(FILE *fp);

cell* read(FILE *fp) {
    int ch;
    int index = 0;
    value v = 0;
    value sign = 1;
    reading rd = R_WHITESPACE;
    while ((ch = getc(fp)) != EOF) {
        switch (ch) {
        case ' ':
            switch (rd) {
            case R_WHITESPACE:
                continue;
            case R_SYMBOL:
                return new_cell(SYMBOL, v, (value) read(fp));
            case R_NUMBER:
                return new_cell(NUMBER, v * sign, (value) read(fp));
            }
        case '(': {
            cell* cl = new_cell(LIST, (value) read(fp), NIL);
            cl->cdr = (value) read(fp);
            return cl;
        }
        case ')':
            switch (rd) {
            case R_SYMBOL:
                return new_cell(SYMBOL, v, NIL);
            case R_NUMBER:
                return new_cell(NUMBER, v * sign, NIL);
            default:
                return (cell*) NIL;
            }
        case '-':
            switch (rd) {
            case R_WHITESPACE:
                sign = -1;
                rd = R_NUMBER;
                continue;
            default:
                printf("Error parsing. Invalid '-' (only allowed before numbers).\n");
                exit(1);
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
                printf("Error parsing. Invalid character: %c\n", ch);
                exit(1);
            }
        case 'a' ... 'z':
            switch (rd) {
            case R_WHITESPACE:
            case R_SYMBOL:
                rd = R_SYMBOL;
                if (index == 7) {
                    printf("Error parsing. Symbol is too long.\n");
                    exit(1);
                }
                ((char *) &v)[index++] = ch;
                continue;
            }
        case '"': {
            cell* cl = new_cell(LIST, (value) string(fp), NIL);
            cl->cdr = (value) read(fp);
            return cl;
        }
        default:
            exit(1);
        }
    }
    return (cell*) NIL;
}

cell* string(FILE *fp) {
    int ch;
    value v = 0;
    while ((ch = getc(fp)) != EOF) {
        if (ch == '"') {
            return (cell*) NIL;
        } else {
            ((char*) &v)[0] = ch;
            return new_cell(SYMBOL, v, (value) string(fp));
        }
    }
}

void print(FILE *fp, cell* c, int index, int depth)
{
    if (index > 0) {
        fprintf(fp, " ");
    }
    switch (c->t) {
    case LIST:
        fprintf(fp, "(");
        print(fp, (cell*) c->car, 0, depth + 1);
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
    }
    if (c->cdr != NIL) {
        print(fp, (cell*) c->cdr, index + 1, depth);
    }
}
