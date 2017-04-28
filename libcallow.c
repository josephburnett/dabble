#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <inttypes.h>

typedef enum { NIL, SYMBOL, NUMBER, ERROR, LIST, FUNC } type;

typedef int64_t chunk_t;

typedef struct {
    type type;
    chunk_t value;
} value_t;

typedef struct cell_t {
    value_t car;
    struct cell_t* cdr;
} cell_t;

typedef value_t (*func_t)(int len, atom[] args);

typedef struct lambda_t {
    int len;
    []value_t args;
    func_t func;
};

value_t read(FILE *fp);

value_t list(FILE *fp) {
    value_t car = read(fp);
    if (car.type == NIL || car.type == ERROR) {
	return car;
    }
    value_t cdr = list(fp);
    if (cdr.type == ERROR) {
	return cdr;
    }
    cell_t* cl = malloc(sizeof(cell_t));
    cl->car = car;
    cl->cdr = (cell_t*) cdr.value;
    return (value_t) { LIST, (chunk_t) cl };
}

value_t symbol(FILE *fp) {
    value_t v = (value_t) { SYMBOL, 0 };
    int ch;
    int index = 0;
    while((ch = getc(fp)) != EOF) {
	if (index == 8) {
	    return (value_t) { ERROR, 0 };
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

value_t number(FILE *fp) {
    int ch;
    value_t v = { NUMBER, 0 };
    chunk_t sign = 0;
    while ((ch = getc(fp)) != EOF) {
	switch (ch) {
	case '-':
	    if (sign != 0) {
		return (value_t) { ERROR, 0 };
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
	    return (value_t) { ERROR, 0 };
	default:
	    ungetc(ch, fp);
	    v.value = v.value * sign;
	    return v;
	}
    }
    v.value = v.value * sign;
    return v;
}

value_t string(FILE *fp) {
    int ch = getc(fp);
    if (ch == EOF) {
	return (value_t) { ERROR, 0 };
    }
    if (ch == '"') {
	return (value_t) { NIL, 0 };
    }
    value_t car = { SYMBOL, 0 };
    ((char*) &car.value)[0] = ch;
    value_t cdr = string(fp);
    if (cdr.type == ERROR) {
	return cdr;
    }
    cell_t* cl = malloc(sizeof(cell_t));
    cl->car = car;
    cl->cdr = (cell_t*) cdr.value;
    return (value_t) { LIST, (chunk_t) cl };
}

value_t read(FILE *fp) {
    int ch;
    while ((ch = getc(fp)) != EOF) {
	switch(ch) {
	case ' ':
	    continue;
	case '(':
	    return list(fp);
	case ')':
	    return (value_t) { NIL, 0 };
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
	    return (value_t) { ERROR, 0 };
	}
    }
    return (value_t) { ERROR, 0 };
}

void print_index(FILE *fp, value_t v, int index) {
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
	    print_index(fp, (value_t) { LIST, (chunk_t) cdr }, index + 1);
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

void print(FILE *fp, value_t v) {
    print_index(fp, v, 0);
}

int len(value_t v, int l) {
    if (v.type != LIST) {
	return l;
    }
    cell* c = ((cell*) v.value)->cdr;
    if (c == 0) {
	return l;
    }
    return len((value_t) { LIST, (chunk_t) c }, l+1);
}

value_t atom(value_t v) {
    if (len(v, 0) != 1) {
	return (atom) { ERROR, 0 };
    }
    v = ((cell*) v.value)->car;
    if (v.type == LIST) {
	return (atom) { NIL, 0 };
    }
    return (atom) { SYMBOL, 't' };
}

value_t car(value_t v) {
    if (len(v, 0) != 1) {
	return (atom) { ERROR, 0 };
    }
    v = ((cell*) v.value)->car;
    if (args[0].type != LIST) {
	return (value_t) { ERROR, 0 };
    }
    return ((cell*) args[0].value)->car;
}

value_t cdr(value_t v) {
    if (len(v, 0) != 1) {
	return (atom) { ERROR, 0 };
    }
    v = ((cell*) v.value)->car;
    if (v.type != LIST) {
	return (atom) { ERROR, 0 };
    }
    cell_t* c = ((cell_t*) v.value)->cdr;
    if (c == 0) {
	return (value_t) { NIL, 0 };
    }
    return (value_t) { LIST, (chunk_t) c };
}

value_t cond(value_t v) {
    if (len(v, 0) < 2) {
	return (value_t) { ERROR, 0 };
    }
    value_t pred = ((cell*) v.value)->car;
    value_t val = ((cell*) v.value)->cdr->car;
    if (truthy(pred)) {
	return val;
    }
    cell* c = ((cell*) v.value)->cdr->cdr;
    if (c == 0) {
	return (value_t) { ERROR, 0 };
    }
    return cond((value_t) { LIST, (chunk_t) c });
}

value_t cons(value_t v) {
    if (len(v, 0) != 2) {
	return (value_t) { ERROR, 0 };
    }
    value_t car = ((cell*) v.value)->car;
    value_t cdr = ((cell*) v.value)->cdr->car;
    if (cdr.type != LIST && cdr.type != NIL) {
	return (value_t) { ERROR, 0 };
    }
    cell_t* c = malloc(sizeof(cell_t));
    c->car = car;
    c->cdr = 0;
    if (cdr.type == LIST) {
	c->cdr = (cell_t*) cdr;
    }
    return (value_t) { LIST, (chunk_t) c };
}

value_t eq_internal(value_t a, value_t b) {
    if (a.type != b.type || a.value != b.value) {
	return (value_t) { NIL, 0 };
    }
    if (a.type == LIST) {
	return eq_internal((value_t) { LIST, ((cell*) a.value)->car },
			   (value_t) { LIST, ((cell*) b.value)->car });
    }
    return (value_t) { SYMBOL, 't' };
}

value_t eq(value_t v)) {
    if (len(v, 0) != 2) {
	return (value_t) { ERROR, 0 };
    }
    value_t a = ((cell*) v.value)->car;
    value_t b = ((cell*) v.value)->cdr->car;
    return eq_internal(a, b);
}

value_t quote(value_t v) {
    return v;
}

value_t lambda(int len, value_t[] args);

value_t label(value_t a);

value_t lookup(value_t s, value_t env) {
    if (env.type == NIL) {
	return (value_t) { ERROR, 0 };
    }
    value_t binding = ((cell*) env.value)->car;
    value_t first = ((cell*) binding.value)->car;
    if (eq_internal(first, s).type != NIL) {
	return ((cell*) binding.value)->cdr->car;
    }
    cell* cdr = ((cell*) env.value)->cdr;
    if (cdr == 0) {
	return (value_t) { ERROR, 0 };
    }
    return lookup(s, (value_t) { LIST, (chunk_t) cdr});
}

value_t eval(value_t v, value_t env) {
    switch (v.type) {
    case NIL:
    case NUMBER:
    case ERROR:
    case FUNC:
	return v;
    case SYMBOL:
	return lookup(v, env);
    case LIST: {
	value_t func = ((cell*) v.value)->car;
	if (func.type == SYMBOL) {
	    func = lookup(func, env);
	}
	if (func.type != FUNC) {
	    return (value_t) { ERROR, 0 };
	}
	cell* cdr = ((cell*) v.value)->cdr;
	// TODO: eval params in turn
	value_t params = (value_t) { NIL, 0 };
	if (cdr != 0) {
	    params = cdr->car;
	}
	return ((func_t) func)(params);
    }
    }
}
