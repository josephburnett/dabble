#include "callowlib.c"
#include <stdio.h>
#include <string.h>

int check(char name[], char given[], char expect[])
{
    int ch;
    FILE *stream;

    stream = fmemopen(given, strlen(given), "r");
    cell* c = read_cell(stream);

    char* actual;
    size_t size;
    stream = open_memstream(&actual, &size);
    print(stream, c);
    fclose(stream);

    if (strcmp(actual, expect) != 0) {
	printf("FAIL: %s\n", name);
	printf("    Given:    %s\n", given);
	printf("    Expected: %s\n", expect);
	printf("    Got:      %s\n", actual);
	return 1;
    }
    return 0;
}

int main(int argc, char *argv[])
{
    int fail = 0;
  
    fail += check("Single symbol list", "(abcd)", "(abcd)");

    if (fail == 0) {
	printf("\nALL TESTS PASSED!\n\n");
    } else {
	printf("\n%d FAILED TESTS!\n\n", fail);
    }
} 




























