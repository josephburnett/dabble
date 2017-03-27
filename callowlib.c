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

cell* read_cell(FILE*fp);

cell* read_symbol(FILE *fp)
{
    cell* l = new_cell(SYMBOL, 0, NIL);
    int i = 0;
    int c;
    while((c = getc(fp)) != EOF) {
	if (c >= 97 && c <= 122) {
	    if (i == 7) {
		printf("Error parsing. Symbol is too long: %c\n", c);
		exit(1);
	    }
	    ((char*) &(l->car))[i] = c;
	    i++;
	    continue;
	}
	switch (c) {
	case ' ':
	    l->cdr = (value) read_cell(fp);
	    return l;
	case ')':
	    ungetc(c, fp);
	    return l;
	default:
	    printf("Error parsing. Invalid symbol character: %c\n", c);
	    exit(1);
	}
    }
}

cell* read_number(FILE *fp)
{
    value num = 0;
    int c;
    while((c = getc(fp)) != EOF) {
	if (c == ' ' || c == ')') {
	    return new_cell(NUMBER, num, NIL);
	}
	if (c < 0 || c > 9) {
	    printf("Error parsing. Invalid number: %c\n", c);
	    exit(1);
	}
	num = num * 10 + ((int64_t) c - 48);
    }
}

cell* read_cell(FILE *fp)
{
    int c;
    while ((c = getc(fp)) != EOF) {
	if (c >= 97 && c <= 122) {
	    ungetc(c, fp);
	    return read_symbol(fp);
	}
	switch (c) {
	case '(':
	    return new_cell(LIST, (value) read_cell(fp), NIL);
	case ')':
	    return (cell*) NIL;
	case ' ':
	    continue;
	default:
	    printf("Error parsing. Invalid char: %c\n", c);
	    exit(1);
	}
    }
}

void print(FILE *fp, cell* c)
{
    if (c == (cell*) NIL) {
      fprintf(fp, ")");
      return;
    }
    switch (c->t) {
    case LIST:
	fprintf(fp, "(");
	print(fp, (cell*) c->car);
	if (c->cdr == NIL) {
	    fprintf(fp, ")");
	    break;
	}
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
	print(fp, (cell*) c->cdr);
    }
}

