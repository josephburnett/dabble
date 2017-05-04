#include "libcallow.c"
#include <stdio.h>
#include <string.h>

char* parse_test_cases[][3] = {
    {
	"Single symbol",
	"abcd",
	"abcd"
    },
    {
	"Symbol with numbers",
	"a1b2c3d4",
	"a1b2c3d4"
    },
    {
	"Symbol with dash",
	"a-b-c-d",
	"a-b-c-d"
    },
    {
	"Multiple symbols. One read.",
	"abcd efgh",
	"abcd"
    },
    {
	"Symbol whitespace ignored",
	"   abcd  ",
	"abcd"
    },
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
        "Single number",
        "1234",
        "1234"
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
	"Invalid number",
	"1-234",
	"<error>"
    },
    {
	"Invalid number",
	"1a2b3c4d",
	"<error>"
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
    },
    {
        "Error on symbol too long",
        "abcdefghijklmn",
        "<error>"
    },
    {
        "Error on mismatched parens",
        "(abc]",
        "<error>"
    },
    {
        "Error on open parens",
        "(abc",
        "<error>"
    }
};

char* lookup_test_cases[][4] = {
    {
	"Lookup in empty list",
	"()",
	"a",
	"<error>"
    },
    {
	"Lookup symbol literal",
	"((a b))",
	"a",
	"b"
    },
    {
	"Lookup shadowing symbol literal",
	"((a c) (a b))",
	"a",
	"c"
    },
    {
	"Lookup number",
	"((a 1))",
	"a",
	"1"
    }
};

int check_result(char name[], char actual[], char expect[]) {
    if (strcmp(actual, expect) != 0) {
        printf("\nFAIL: %s\n", name);
        printf("    Expected: %s\n", expect);
        printf("    Got:      %s\n", actual);
        return 1;
    }
    return 0;
}

int check_parse(char name[], char given[], char expect[])
{
    FILE* stream;
    stream = fmemopen(given, strlen(given), "r");
    value_t v = read(stream);
    fclose(stream);

    char* actual;
    size_t size;
    stream = open_memstream(&actual, &size);
    print(stream, v);
    fclose(stream);

    return check_result(name, actual, expect);
}

int check_lookup(char name[], char env[], char symbol[], char expect[]) {
    FILE* stream;
    stream = fmemopen(env, strlen(env), "r");
    value_t env_value = read(stream);
    fclose(stream);

    stream = fmemopen(symbol, strlen(symbol), "r");
    value_t symbol_value = read(stream);
    fclose(stream);

    value_t actual_value = lookup(symbol_value, env_value);
    char* actual;
    size_t size;
    stream = open_memstream(&actual, &size);
    print(stream, actual_value);
    fclose(stream);

    return check_result(name, actual, expect);
}

int main(int argc, char *argv[])
{
    int fail = 0;
    int i;
    for (i = 0; i < sizeof(parse_test_cases) / sizeof(parse_test_cases[0]); i++) {
	char** args = parse_test_cases[i];
        fail += check_parse(args[0], args[1], args[2]);
    }
    for (i = 0; i < sizeof(lookup_test_cases) / sizeof(lookup_test_cases[0]); i++) {
	char** args = lookup_test_cases[i];
	fail += check_lookup(args[0], args[1], args[2], args[3]);
    }

    if (fail == 0) {
        printf("\nALL TESTS PASSED!\n\n");
    } else {
        printf("\n%d FAILED TESTS!\n\n", fail);
    }
}
