#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>

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

cell* read(FILE *fp);

cell* list(FILE *fp) {
    cell *cl = new_cell(LIST, 0, NIL);
    cl->car = (value) read(fp);
    cl->cdr = (value) read(fp);
    return cl;
}

cell* symbol(FILE *fp) {
    cell* cl = new_cell(SYMBOL, 0, NIL);
    int ch;
    int i = 0;
    while ((ch = getc(fp)) != EOF) {
	switch (ch) {
	case ' ':
	    cl->cdr = (value) read(fp);
	    return cl;
	case ')':
	    cl->cdr = NIL;
	    return cl;
	default:
	    if (ch >= 97 && ch <= 122) {
		if (i == 7) {
		    printf("Error parsing. Symbol is too long: %c\n", ch);
		    exit(1);
		}
		((char*) &(cl->car))[i] = ch;
		i++;
		continue;
	    } else {
		printf("Error parsing. Invalid symbol character: %c\n", ch);
		exit(1);
	    }
	}
    }
}

cell* read(FILE *fp) {
    int ch;
    while ((ch = getc(fp)) != EOF) {
	switch (ch) {
	case '(':
	    return list(fp);
	case ')':
	    return (cell*) NIL;
	case ' ':
	    continue;
	default:
	    if (ch >= 97 && ch <= 122) {
		ungetc(ch, fp);
		return symbol(fp);
	    } else {
		printf("Error parsing. Invalid char: %c\n", ch);
		exit(1);
	    }
	}
    }
    return (cell*) NIL;
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
	fprintf(fp, "%ldn", c->car);
	break;
    case FUNC:
	fprintf(fp, "<func>");
	break;
    }
    if (c->cdr != NIL) {
	print(fp, (cell*) c->cdr, index + 1, depth);
    }
}
