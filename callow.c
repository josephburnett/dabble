#include "callowlib.c"
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>

int main(int argc, char *argv[])
{
    if (argc != 2) {
	printf("Usage: callow <filename>\n");
	return 1;
    }
    FILE *fp;
    if ((fp = fopen(argv[1], "r")) == NULL) {
	printf("callow: can't open %s\n", *argv);
	return 1;
    }
    cell* c = read_cell(fp);
    print(stdout, c);
    printf("\n");
    fclose(fp);
    return 0;
}

