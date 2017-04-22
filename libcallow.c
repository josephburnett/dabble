#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <inttypes.h>

typedef enum { NIL, LIST, SYMBOL, NUMBER, FUNC, ERROR } type;

typedef int64_t value;

typedef struct {
    type type;
    value value;
} atom;

typedef struct cell {
    atom car;
    struct cell* cdr;
} cell;

atom read(FILE *fp);

atom list(FILE *fp) {
    atom car = read(fp);
    if (car.type == NIL || car.type == ERROR) {
	return car;
    }
    atom cdr = list(fp);
    if (cdr.type == ERROR) {
	return cdr;
    }
    cell* cl = malloc(sizeof(cell));
    cl->car = car;
    cl->cdr = (cell*) cdr.value;
    return (atom) { LIST, (value) cl };
}

atom symbol(FILE *fp) {
    atom v = (atom) { SYMBOL, 0 };
    int ch;
    int index = 0;
    while((ch = getc(fp)) != EOF) {
	if (index == 8) {
	    return (atom) { ERROR, 0 };
	}
	switch (ch) {
	case 'a' ... 'z':
	    ((char*) &v.value)[index++] = ch;
	    continue;
	default:
	    ungetc(ch, fp);
	    return v;
	}
    }
    return v;
}

atom number(FILE *fp) {
    int ch;
    atom v = { NUMBER, 0 };
    value sign = 0;
    while ((ch = getc(fp)) != EOF) {
	switch (ch) {
	case '-':
	    if (sign != 0) {
		return (atom) { ERROR, 0 };
	    }
	    sign = -1;
	    continue;
	case '0' ... '9':
	    if (sign == 0) {
		sign = 1;
	    }
	    v.value = v.value * 10 + (ch - '0');
	    continue;
	default:
	    ungetc(ch, fp);
	    v.value = v.value * sign;
	    return v;
	}
    }
}

atom string(FILE *fp) {
    int ch = getc(fp);
    if (ch == EOF) {
	return (atom) { ERROR, 0 };
    }
    if (ch == '"') {
	return (atom) { NIL, 0 };
    }
    atom car = { SYMBOL, 0 };
    ((char*) &car.value)[0] = ch;
    atom cdr = string(fp);
    if (cdr.type == ERROR) {
	return cdr;
    }
    cell* cl = malloc(sizeof(cell));
    cl->car = car;
    cl->cdr = (cell*) cdr.value;
    return (atom) { LIST, (value) cl };
}

atom read(FILE *fp) {
    int ch;
    while ((ch = getc(fp)) != EOF) {
	switch(ch) {
	case ' ':
	    continue;
	case '(':
	    return list(fp);
	case ')':
	    return (atom) { NIL, 0 };
	case '-':
	case '0' ... '9':
	    ungetc(ch, fp);
	    return number(fp);
	case 'a' ... 'z':
	    ungetc(ch, fp);
	    return symbol(fp);
	case '"':
	    return string(fp);
	default:
	    return (atom) { ERROR, 0 };
	}
    }
    return (atom) { ERROR, 0 };
}

void print_index(FILE *fp, atom v, int index)
{
    if (index > 0) {
        fprintf(fp, " ");
    }
    switch (v.type) {
    case LIST:
	if (index <= 0) {
	    fprintf(fp, "(");
	}
	print_index(fp, ((cell*) v.value)->car, 0);
	cell* cdr = ((cell*) v.value)->cdr;
	if (cdr != 0) {
	    print_index(fp, (atom) { LIST, (value) cdr }, index + 1);
	}
	if (index <= 0) {
	    fprintf(fp, ")");
	}
        break;
    case SYMBOL: {
        char* sym = (char*) &(v.value);
        int i;
        for (i = 0; i < 8; i++) {
            if (sym[i] != 0) {
                fprintf(fp, "%c", sym[i]);
            }
        }
        break;
    }
    case NUMBER:
        fprintf(fp, "%" PRId64, v.value);
        break;
    case FUNC:
        fprintf(fp, "<func>");
        break;
    case ERROR:
        fprintf(fp, "<error>");
        break;
    case NIL:
	fprintf(fp, "()");
	break;
    }
}

void print(FILE *fp, atom v)
{
    print_index(fp, v, 0);
}
