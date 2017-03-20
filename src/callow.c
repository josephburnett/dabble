#include <stdio.h>
#include <stdint.h>

typedef uint64_t symbol;
typedef struct {
    car symbol;
    cdr symbol;
} cell;

cell* new_cell(symbol car, symbol cdr)
{
    cell* c = malloc(sizeof(cell));
    c->car = car;
    c->cdr = cdr;
    return c;
}

cell* NIL = new_cell(0, 0);
cell* T =   new_cell((symbol) {'T',' ',' ',' ',' ',' ',' ',' '}, NIL);

symbol LIST =   (symbol) {'L','I','S','T',' ',' ',' ',' '};
symbol SYMBOL = (symbol) {'S','Y','M','B','O','L',' ',' '};
symbol STRING = (symbol) {'S','T','R','I','N','G',' ',' '};

int main(int argc, char *argv[])
{
    if (argc != 2) {
	printf("Usage: callow <filename>\n");
	return 1;
    }
    FILE *fp;
    if ((fp = fopen(*argv, "r")) == NULL) {
	printf("callow: can't open %s\n", *argv);
	return 1;
    }
    fclose(fp);
    return 0;
}

cell* read_cell(FILE *fp)
{
    int c;
    while ((c = getc(fp)) != EOF) {
	switch (c) {
	case '(':
	    return new_cell(LIST, read_cell(fp));
	case 'a':
	    return new_cell(STRING, read_cell(fp));
	case ')':
	    return 
	default:
	    continue;
	}
    }
}
