#include "libcallow.c"
#include <stdio.h>
#include <string.h>

char* test_cases[][3] = {
    {
        "Single symbol list",
        "(abcd)",
        "(abcd)"
    },
    {
        "Multiple symbol list",
        "(a b c d)",
        "(a b c d)"
    },
    {
        "Nested lists",
        "(a (b) ((c)))",
        "(a (b) ((c)))"
    },
    {
        "Nested lists",
        "((a) (b))",
        "((a) (b))"
    },
    {
        "Whitespace ignored",
        "( a  b     c)",
        "(a b c)"
    },
    {
        "Single number list",
        "(1234)",
        "(1234)"
    },
    {
        "Multiple number list",
        "(1 2 3 4)",
        "(1 2 3 4)"
    },
    {
        "Negative numbers",
        "(-1 -2 -3)",
        "(-1 -2 -3)"
    },
    {
        "Single string",
        "\"abcd\"",
        "(a b c d)"
    },
    {
        "Single string list",
        "(\"abcd\")",
        "((a b c d))"
    },
    {
        "Nested string lists",
        "(\"a\" \"b\")",
        "((a) (b))"
    }
};

int check(char name[], char given[], char expect[])
{
    int ch;
    FILE *stream;

    stream = fmemopen(given, strlen(given), "r");
    cell* c = read(stream);

    char* actual;
    size_t size;
    stream = open_memstream(&actual, &size);
    print(stream, c, 0, 1);
    fclose(stream);

    if (strcmp(actual, expect) != 0) {
        printf("\nFAIL: %s\n", name);
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
    int i;
    for (i = 0; i < sizeof(test_cases) / sizeof(test_cases[0]); i++) {
        fail += check(test_cases[i][0], test_cases[i][1], test_cases[i][2]);
    }

    if (fail == 0) {
        printf("\nALL TESTS PASSED!\n\n");
    } else {
        printf("\n%d FAILED TESTS!\n\n", fail);
    }
}
