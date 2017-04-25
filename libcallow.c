#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <inttypes.h>

typedef enum { NIL, SYMBOL, NUMBER, ERROR, LIST, FUNC } type;

typedef int64_t value_t;

typedef struct {
    type type;
    value_t value;
} atom_t;

typedef struct cell_t {
    atom_t car;
    struct cell_t* cdr;
} cell_t;

typedef atom_t (*func_t)(int len, atom[] args);

typedef struct lambda_t {
    int names_len;
    []atom_t names;
    int forms_len;
    []atom_t forms;
    atom_t environment;
    func_t func;
};

atom_t read(FILE *fp);

atom_t list(FILE *fp) {
    atom_t car = read(fp);
    if (car.type == NIL || car.type == ERROR) {
	return car;
    }
    atom_t cdr = list(fp);
    if (cdr.type == ERROR) {
	return cdr;
    }
    cell_t* cl = malloc(sizeof(cell_t));
    cl->car = car;
    cl->cdr = (cell_t*) cdr.value;
    return (atom_t) { LIST, (value_t) cl };
}

atom_t symbol(FILE *fp) {
    atom_t v = (atom_t) { SYMBOL, 0 };
    int ch;
    int index = 0;
    while((ch = getc(fp)) != EOF) {
	if (index == 8) {
	    return (atom_t) { ERROR, 0 };
	}
	switch (ch) {
	case 'a' ... 'z':
	case '0' ... '9':
	case '-':
	    ((char*) &v.value)[index++] = ch;
	    continue;
	default:
	    ungetc(ch, fp);
	    return v;
	}
    }
    return v;
}

atom_t number(FILE *fp) {
    int ch;
    atom_t v = { NUMBER, 0 };
    value_t sign = 0;
    while ((ch = getc(fp)) != EOF) {
	switch (ch) {
	case '-':
	    if (sign != 0) {
		return (atom_t) { ERROR, 0 };
	    }
	    sign = -1;
	    continue;
	case '0' ... '9':
	    if (sign == 0) {
		sign = 1;
	    }
	    v.value = v.value * 10 + (ch - '0');
	    continue;
	case 'a' ... 'z':
	    return (atom_t) { ERROR, 0 };
	default:
	    ungetc(ch, fp);
	    v.value = v.value * sign;
	    return v;
	}
    }
    v.value = v.value * sign;
    return v;
}

atom_t string(FILE *fp) {
    int ch = getc(fp);
    if (ch == EOF) {
	return (atom_t) { ERROR, 0 };
    }
    if (ch == '"') {
	return (atom_t) { NIL, 0 };
    }
    atom_t car = { SYMBOL, 0 };
    ((char*) &car.value)[0] = ch;
    atom_t cdr = string(fp);
    if (cdr.type == ERROR) {
	return cdr;
    }
    cell_t* cl = malloc(sizeof(cell_t));
    cl->car = car;
    cl->cdr = (cell_t*) cdr.value;
    return (atom_t) { LIST, (value_t) cl };
}

atom_t read(FILE *fp) {
    int ch;
    while ((ch = getc(fp)) != EOF) {
	switch(ch) {
	case ' ':
	    continue;
	case '(':
	    return list(fp);
	case ')':
	    return (atom_t) { NIL, 0 };
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
	    return (atom_t) { ERROR, 0 };
	}
    }
    return (atom_t) { ERROR, 0 };
}

void print_index(FILE *fp, atom_t v, int index) {
    if (index > 0) {
        fprintf(fp, " ");
    }
    switch (v.type) {
    case LIST:
	if (index <= 0) {
	    fprintf(fp, "(");
	}
	print_index(fp, ((cell_t*) v.value)->car, 0);
	cell_t* cdr = ((cell_t*) v.value)->cdr;
	if (cdr != 0) {
	    print_index(fp, (atom_t) { LIST, (value_t) cdr }, index + 1);
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

void print(FILE *fp, atom_t v) {
    print_index(fp, v, 0);
}

atom_t atom(int len, atom_t[] args) {
    if (len != 1) {
	return (atom) { ERROR, 0 };
    }
    if (args[0].type == LIST) {
	return (atom) { NIL, 0 };
    }
    return (atom) { SYMBOL, 't' };
}

atom_t car(int len, atom_t[] args) {
    if (len != 1) {
	return (atom) { ERROR, 0 };
    }
    if (args[0].type != LIST) {
	return (atom_t) { ERROR, 0 };
    }
    return ((cell*) args[0].value)->car;
}

atom_t cdr(int len, atom_t[] args) {
    if (len != 1) {
	return (atom) { ERROR, 0 };
    }
    if (args[0].type != LIST) {
	return (atom) { ERROR, 0 };
    }
    cell_t* c = ((cell_t*) args[0].value)->cdr;
    if (c == 0) {
	return (atom_t) { NIL, 0 };
    }
    return (atom_t) { LIST, (value_t) c };
}

atom_t cond(int len, atom_t[] args) {
    if (len == 0 || args[0].type != LIST) {
	return (atom) { ERROR, 0 };
    }
    pred = ((cell*) args[0].value)->car;
    if (pred != NIL) {
	cell* c = ((cell*) args[0].value)->cdr;
	if (c == 0) {
	    return (atom_t) { ERROR, 0 };
	}
	return c->car;
    }
    return cond(len - 1, args + 1);
}

atom_t cons(int len, atom_t[] args) {
    if (len != 2) {
	return (atom_t) { ERROR, 0 };
    }
    if (args[1].type != NIL && args[1].type != LIST) {
	return (atom_t) { ERROR, 0 };
    }
    cell_t* c = malloc(sizeof(cell_t));
    c->car = args[0];
    c->cdr = 0;
    if (args[1].type == LIST) {
	c->cdr = (cell_t*) args[1].value;
    }
    return (atom_t) { LIST, (value_t) c };
}

atom_t eq(int len, atom_t[] args) {
    if (len != 2) {
	return (atom_t) { ERROR, 0 };
    }
    if (args[0].type != args[1].type || args[0].value != args[1].value) {
	return (atom_t) { NIL, 0 };
    }
    if (args[0].type == LIST) {
	return eq(cdr(args[0]), cdr(args[1]));
    }
    return (atom_t) { SYMBOL, 't' };
}

atom_t quote(int len, atom_t[] args);

atom_t lambda(atom_t a) {

}

atom_t label(atom_t a) {

}
