#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>
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

typedef value_t (*func_t)(value_t);

typedef struct {
    int len;
    value_t args;
    func_t func;
} lambda_t;

value_t read(FILE *fp);

value_t read_string(char form[]) {
    FILE* stream;
    stream = fmemopen(form, strlen(form), "r");
    value_t v = read(stream);
    fclose(stream);
    return v;
}

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
    cell_t* c = (cell_t*) v.value;
    if (c == 0) {
	return l;
    }
    return len((value_t) { LIST, (chunk_t) c->cdr }, l+1);
}

value_t atom(value_t v) {
    if (len(v, 0) != 1) {
	return (value_t) { ERROR, 0 };
    }
    v = ((cell_t*) v.value)->car;
    if (v.type == LIST) {
	return (value_t) { NIL, 0 };
    }
    return (value_t) { SYMBOL, 't' };
}

value_t car(value_t v) {
    if (len(v, 0) != 1) {
	return (value_t) { ERROR, 0 };
    }
    v = ((cell_t*) v.value)->car;
    if (v.type != LIST) {
	return (value_t) { ERROR, 0 };
    }
    return ((cell_t*) v.value)->car;
}

value_t cdr(value_t v) {
    if (len(v, 0) != 1) {
	return (value_t) { ERROR, 0 };
    }
    v = ((cell_t*) v.value)->car;
    if (v.type != LIST) {
	return (value_t) { ERROR, 0 };
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
    value_t pred = ((cell_t*) v.value)->car;
    value_t val = ((cell_t*) v.value)->cdr->car;
    if (pred.type == SYMBOL ||
	(pred.type == NUMBER && pred.value != 0)) {
	return val;
    }
    cell_t* c = ((cell_t*) v.value)->cdr->cdr;
    if (c == 0) {
	return (value_t) { ERROR, 0 };
    }
    return cond((value_t) { LIST, (chunk_t) c });
}

value_t cons(value_t v) {
    if (len(v, 0) != 2) {
	return (value_t) { ERROR, 0 };
    }
    value_t car = ((cell_t*) v.value)->car;
    value_t cdr = ((cell_t*) v.value)->cdr->car;
    if (cdr.type != LIST && cdr.type != NIL) {
	return (value_t) { ERROR, 0 };
    }
    cell_t* c = malloc(sizeof(cell_t));
    c->car = car;
    c->cdr = 0;
    if (cdr.type == LIST) {
	c->cdr = (cell_t*) cdr.value;
    }
    return (value_t) { LIST, (chunk_t) c };
}

value_t eq_internal(value_t a, value_t b) {
    if (a.type != b.type || a.value != b.value) {
	return (value_t) { NIL, 0 };
    }
    if (a.type == LIST) {
	return eq_internal((value_t) { LIST, (chunk_t) ((cell_t*) a.value)->car.value },
			   (value_t) { LIST, (chunk_t) ((cell_t*) b.value)->car.value });
    }
    return (value_t) { SYMBOL, 't' };
}

value_t eq(value_t v) {
    if (len(v, 0) != 2) {
	return (value_t) { ERROR, 0 };
    }
    value_t a = ((cell_t*) v.value)->car;
    value_t b = ((cell_t*) v.value)->cdr->car;
    return eq_internal(a, b);
}

value_t quote(value_t v) {
    if (len(v, 0) != 1) {
	return (value_t) { ERROR, 0 };
    }
    return ((cell_t*) v.value)->car;
}

value_t lambda(value_t v);

value_t label(value_t v);

value_t lookup(value_t s, value_t env) {
    if (env.type == NIL) {
	return (value_t) { ERROR, 0 };
    }
    value_t binding = ((cell_t*) env.value)->car;
    value_t first = ((cell_t*) binding.value)->car;
    if (eq_internal(first, s).type != NIL) {
	return ((cell_t*) binding.value)->cdr->car;
    }
    cell_t* cdr = ((cell_t*) env.value)->cdr;
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
	value_t func = ((cell_t*) v.value)->car;
	if (func.type == SYMBOL) {
	    func = lookup(func, env);
	}
	if (func.type != FUNC) {
	    return (value_t) { ERROR, 0 };
	}
	lambda_t* lamb = (lambda_t*) func.value;
	cell_t* cdr = ((cell_t*) v.value)->cdr;
	// TODO: eval params in turn
	value_t params = (value_t) { NIL, 0 };
	if (cdr != 0) {
	    params = (value_t) { LIST, (chunk_t) cdr };
	}
	return (*(lamb->func))(params);
    }
    }
}

value_t bind(char name[], func_t fn, value_t env) {
    value_t name_value = read_string(name);
    if (name_value.type != SYMBOL) {
	return (value_t) { ERROR, 0 };
    }
    cell_t* first = malloc(sizeof(cell_t));
    cell_t* second = malloc(sizeof(cell_t));
    lambda_t* lamb = malloc(sizeof(lambda_t));
    first->car = name_value;
    first->cdr = second;
    second->car = (value_t) { FUNC, (chunk_t) lamb };
    second->cdr = 0;
    lamb->func = fn;
    cell_t* new_env = malloc(sizeof(cell_t));
    new_env->car = (value_t) { LIST, (chunk_t) first };
    if (env.type == LIST) {
	new_env->cdr = (cell_t*) env.value;
    } else {
	new_env->cdr = 0;
    }
    return (value_t) { LIST, (chunk_t) new_env };
}

value_t callow_core() {
    value_t env = (value_t) { NIL, 0 };
    env = bind("atom", &atom, env);
    env = bind("car", &car, env);
    env = bind("cdr", &cdr, env);
    env = bind("cond", &cond, env);
    env = bind("cons", &cons, env);
    env = bind("eq", &eq, env);
    env = bind("quote", &quote, env);
    return env;
}
