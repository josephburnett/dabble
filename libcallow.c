#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <inttypes.h>

typedef enum { NIL, SYMBOL, NUMBER, ERROR, LIST, LAMBDA, FUNC } type;

typedef int64_t chunk_t;

typedef struct {
    type type;
    chunk_t value;
} value_t;

typedef struct cell_t {
    value_t car;
    struct cell_t* cdr;
} cell_t;

typedef struct {
    value_t names;
    value_t form;
    value_t env;
} lambda_t;

typedef value_t (*func_t)(value_t, value_t);

typedef struct {
    func_t func;
} func_s;

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

value_t eval(value_t v, value_t env);
value_t eval_list(value_t list, value_t env);

value_t atom(value_t args, value_t env) {
    if (len(args, 0) != 1) {
	return (value_t) { ERROR, 0 };
    }
    args = eval_list(args, env);
    args = ((cell_t*) args.value)->car;
    if (args.type == LIST) {
	return (value_t) { NIL, 0 };
    }
    return (value_t) { SYMBOL, 't' };
}

value_t car(value_t args, value_t env) {
    if (len(args, 0) != 1) {
	return (value_t) { ERROR, 0 };
    }
    args = eval_list(args, env);
    args = ((cell_t*) args.value)->car;
    if (args.type != LIST) {
	return (value_t) { ERROR, 0 };
    }
    return ((cell_t*) args.value)->car;
}

value_t cdr(value_t args, value_t env) {
    if (len(args, 0) != 1) {
	return (value_t) { ERROR, 0 };
    }
    args = eval_list(args, env);
    args = ((cell_t*) args.value)->car;
    if (args.type != LIST) {
	return (value_t) { ERROR, 0 };
    }
    cell_t* c = ((cell_t*) args.value)->cdr;
    if (c == 0) {
	return (value_t) { NIL, 0 };
    }
    return (value_t) { LIST, (chunk_t) c };
}

value_t cond(value_t args, value_t env) {
    if (len(args, 0) < 2) {
	return (value_t) { ERROR, 0 };
    }
    args = eval_list(args, env);
    value_t pred = ((cell_t*) args.value)->car;
    value_t val = ((cell_t*) args.value)->cdr->car;
    if (pred.type == SYMBOL ||
	(pred.type == NUMBER && pred.value != 0)) {
	return val;
    }
    cell_t* c = ((cell_t*) args.value)->cdr->cdr;
    if (c == 0) {
	return (value_t) { ERROR, 0 };
    }
    return cond((value_t) { LIST, (chunk_t) c }, env);
}

value_t cons(value_t args, value_t env) {
    if (len(args, 0) != 2) {
	return (value_t) { ERROR, 0 };
    }
    args = eval_list(args, env);
    value_t car = ((cell_t*) args.value)->car;
    value_t cdr = ((cell_t*) args.value)->cdr->car;
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

value_t eq(value_t args, value_t env) {
    if (len(args, 0) != 2) {
	return (value_t) { ERROR, 0 };
    }
    args = eval_list(args, env);
    value_t a = ((cell_t*) args.value)->car;
    value_t b = ((cell_t*) args.value)->cdr->car;
    return eq_internal(a, b);
}

value_t quote(value_t args, value_t env) {
    if (len(args, 0) != 1) {
	return (value_t) { ERROR, 0 };
    }
    return ((cell_t*) args.value)->car;
}

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

value_t bind(value_t name, value_t value, value_t env) {
    if (name.type != SYMBOL) {
	return (value_t) { ERROR, 0 };
    }
    cell_t* first = malloc(sizeof(cell_t));
    cell_t* second = malloc(sizeof(cell_t));
    first->car = name;
    first->cdr = second;
    second->car = value;
    second->cdr = 0;
    cell_t* new_env = malloc(sizeof(cell_t));
    new_env->car = (value_t) { LIST, (chunk_t) first };
    if (env.type == LIST) {
	new_env->cdr = (cell_t*) env.value;
    } else {
	new_env->cdr = 0;
    }
    return (value_t) { LIST, (chunk_t) new_env };
}

value_t eval_list(value_t list, value_t env) {
    if (list.type == NIL) {
	return list;
    }
    cell_t* cell = malloc(sizeof(cell_t));
    value_t car = eval(((cell_t*) list.value)->car, env);
    cell->car = car;
    cell->cdr = 0;
    cell_t* cdr = ((cell_t*) list.value)->cdr;
    if (cdr != 0) {
	value_t l = eval_list((value_t) { LIST, (chunk_t) cdr }, env);
	cell->cdr = (cell_t*) l.value;
    }
    return (value_t) { LIST, (chunk_t) cell };
}

value_t eval(value_t v, value_t env) {
    switch (v.type) {
    case NIL:
    case NUMBER:
    case ERROR:
    case FUNC:
    case LAMBDA:
	return v;
    case SYMBOL:
	return lookup(v, env);
    case LIST: {
	value_t params = (value_t) { NIL, 0 };
	cell_t* cdr = ((cell_t*) v.value)->cdr;
	if (cdr != 0) {
	    params = (value_t) { LIST, (chunk_t) cdr };
	}
	value_t first = ((cell_t*) v.value)->car;
	first = eval(first, env);
	switch (first.type) {
	case FUNC: {
	    func_s* func = (func_s*) first.value;
	    return (*(func->func))(params, env);
	}
	case LAMBDA: {
	    lambda_t* lamb = (lambda_t*) first.value;
	    if (len(params, 0) != len(lamb->names, 0)) {
		return (value_t) { ERROR, 0 };
	    }
	    value_t lambda_env = lamb->env;
	    cell_t* name = (cell_t*) lamb->names.value;
	    cell_t* param = (cell_t*) params.value;
	    while (name != 0) {
		lambda_env = bind(name->car, param->car, lambda_env);
		name = name->cdr;
		param = param->cdr;
	    }
	    return eval(lamb->form, lambda_env);
	}
	default:
	    return (value_t) { ERROR, 0 };
	}
    }
    }
}

value_t wrap_fn(func_t fn) {
    func_s* func = malloc(sizeof(func_s));
    func->func = fn;
    return (value_t) { FUNC, (chunk_t) func };
}

value_t label(value_t args, value_t env) {
    if (len(args, 0) != 3) {
	return (value_t) { ERROR, 0 };
    }
    args = eval_list(args, env);
    value_t name = ((cell_t*) args.value)->car;
    value_t value = ((cell_t*) args.value)->cdr->car;
    value_t form = ((cell_t*) args.value)->cdr->cdr->car;
    env = bind(name, value, env);
    return eval(form, env);
}

value_t lambda(value_t args, value_t env) {
    if (len(args, 0) != 2) {
	return (value_t) { ERROR, 0 };
    }
    value_t names = ((cell_t*) args.value)->car;
    if (names.type != LIST) {
	return (value_t) { ERROR, 0 };
    }
    value_t form = ((cell_t*) args.value)->cdr->car;
    lambda_t* lamb = malloc(sizeof(lambda_t));
    lamb->names = names;
    lamb->form = form;
    lamb->env = env;
    return (value_t) { LAMBDA, (chunk_t) lamb };
}

value_t callow_core() {
    value_t env = (value_t) { NIL, 0 };
    env = bind(read_string("atom"), wrap_fn(&atom), env);
    env = bind(read_string("car"), wrap_fn(&car), env);
    env = bind(read_string("cdr"), wrap_fn(&cdr), env);
    env = bind(read_string("cond"), wrap_fn(&cond), env);
    env = bind(read_string("cons"), wrap_fn(&cons), env);
    env = bind(read_string("eq"), wrap_fn(&eq), env);
    env = bind(read_string("quote"), wrap_fn(&quote), env);
    env = bind(read_string("label"), wrap_fn(&label), env);
    env = bind(read_string("lambda"), wrap_fn(&lambda), env);
    return env;
}
