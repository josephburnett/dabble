#include "libcallow.c"
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
    value_t form = read(fp);
    value_t env = callow_core();
    value_t result = eval(form, env);
    print(stdout, result);
    printf("\n");
    fclose(fp);
    return 0;
}

