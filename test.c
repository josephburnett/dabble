#include "libcallow.c"
#include <stdio.h>
#include <string.h>

int
check_result (char name[], char actual[], char expect[])
{
  if (actual == 0 || expect == 0)
    {
      printf ("\nERROR: %s\n", name);
      printf ("    Test parameters missing.\n");
      printf ("    Missing a comma?\n");
      return 1;
    }
  if (strcmp (actual, expect) != 0)
    {
      printf ("\nFAIL: %s\n", name);
      printf ("    Expected: %s\n", expect);
      printf ("    Got:      %s\n", actual);
      return 1;
    }
  return 0;
}

int
check_parse (char name[], char given[], char expect[])
{
  value_t v = read_string (given);
  char *actual;
  size_t size;
  FILE *stream = open_memstream (&actual, &size);
  print (stream, v);
  fclose (stream);
  return check_result (name, actual, expect);
}

char *parse_test_cases[][3] = {
  {
   "Single symbol",
   "abcd",
   "abcd"},
  {
   "Symbol with numbers",
   "a1b2c3d4",
   "a1b2c3d4"},
  {
   "Symbol with dash",
   "a-b-c-d",
   "a-b-c-d"},
  {
   "Multiple symbols. One read.",
   "abcd efgh",
   "abcd"},
  {
   "Symbol whitespace ignored",
   "   abcd  ",
   "abcd"},
  {
   "Single symbol list",
   "(abcd)",
   "(abcd)"},
  {
   "Multiple symbol list",
   "(a b c d)",
   "(a b c d)"},
  {
   "Nested lists",
   "(a (b) ((c)))",
   "(a (b) ((c)))"},
  {
   "Nested lists",
   "((a) (b))",
   "((a) (b))"},
  {
   "Whitespace ignored",
   "( a  b     c)",
   "(a b c)"},
  {
   "Single number",
   "1234",
   "1234"},
  {
   "Multiple number list",
   "(1 2 3 4)",
   "(1 2 3 4)"},
  {
   "Negative numbers",
   "(-1 -2 -3)",
   "(-1 -2 -3)"},
  {
   "Invalid number",
   "1-234",
   "<error: Invalid character '-' in number.>"},
  {
   "Invalid number",
   "1a2b3c4d",
   "<error: Invalid letter in number.>"},
  {
   "Single string",
   "\"abcd\"",
   "(a b c d)"},
  {
   "Single string list",
   "(\"abcd\")",
   "((a b c d))"},
  {
   "Nested string lists",
   "(\"a\" \"b\")",
   "((a) (b))"},
  {
   "Error on symbol too long",
   "abcdefghijklmn",
   "<error: Symbol too long.>"},
  {
   "Error on invalid delimiter",
   "(abc]",
   "<error: Invalid character.>"},
  {
   "Error on open parens",
   "(abc",
   "<error: Unexpected EOF.>"},
  {
   "Nil in a list",
   "(1 () 2)",
   "(1 () 2)"}
};

int
check_lookup (char name[], char env[], char symbol[], char expect[])
{
  value_t env_value = read_string (env);
  value_t symbol_value = read_string (symbol);
  value_t actual_value = lookup (symbol_value, env_value);
  char *actual;
  size_t size;
  FILE *stream = open_memstream (&actual, &size);
  print (stream, actual_value);
  fclose (stream);
  return check_result (name, actual, expect);
}

char *lookup_test_cases[][4] = {
  {
   "Lookup in empty list",
   "()",
   "a",
   "<error: Nil environment.>"},
  {
   "Lookup symbol literal",
   "((a b))",
   "a",
   "b"},
  {
   "Lookup shadowing symbol literal",
   "((a c) (a b))",
   "a",
   "c"},
  {
   "Lookup number",
   "((a 1))",
   "a",
   "1"},
  {
   "Lookup deep number",
   "((a 1) (b 2) (c 3) (d 4))",
   "d",
   "4"}
};

int
check_eval (char name[], value_t env, char form[], char expect[])
{
  value_t form_value = read_string (form);
  value_t actual_value = eval (form_value, env);
  char *actual;
  size_t size;
  FILE *stream = open_memstream (&actual, &size);
  print (stream, actual_value);
  fclose (stream);
  return check_result (name, actual, expect);
}

char *eval_test_cases[][4] = {
  {
   "Eval number literal",
   "()",
   "1",
   "1"},
  {
   "Eval nil literal",
   "()",
   "()",
   "()"},
  {
   "Eval lookup number literal",
   "((a 1))",
   "a",
   "1"}
};

char *core_test_cases[][4] = {
  {
   "Atom number",
   "(atom 1)",
   "t"},
  {
   "Atom list",
   "(atom (1))",
   "()"},
  {
   "Atom empty list",
   "(atom ())",
   "t"},
  {
   "Car single element list",
   "(car (1))",
   "1"},
  {
   "Car of a non-list is error",
   "(car 1)",
   "<error: Non-list argument to car.>"},
  {
   "Car of nil is error",
   "(car ())",
   "<error: Non-list argument to car.>"},
  {
   "Cdr single element list",
   "(cdr (1))",
   "()"},
  {
   "Cdr of two element list",
   "(cdr (1 2))",
   "(2)"},
  {
   "Cdr of a non-list is error",
   "(cdr 1)",
   "<error: Non-list argument to crd.>"},
  {
   "Cond simple case",
   "(cond (eq 1 1) 2)",
   "2"},
  {
   "Cond match at beginning",
   "(cond (eq 1 1) 2 (eq 3 4) 5 (eq 6 6) 7)",
   "2"},
  {
   "Cond match in middle",
   "(cond (eq 1 2) 3 (eq 4 4) 5 (eq 6 6) 7)",
   "5"},
  {
   "Cond match at end",
   "(cond (eq 1 2) 3 (eq 4 5) 6 (eq 7 7) 8)",
   "8"},
  {
   "Cons with empty list",
   "(cons 1 ())",
   "(1)",
   },
  {
   "Cons with single element list",
   "(cons 1 (quote (2)))",
   "(1 2)"},
  {
   "Eq with two numbers",
   "(eq 1 1)",
   "t"},
  {
   "Eq with nil",
   "(eq () ())",
   "t"},
  {
   "List eq",
   "(eq (1) (1))",
   "t"},
  {
   "List deep eq",
   "(eq (1 (2 3)) (1 (2 3)))",
   "t"},
  {
   "List deep not eq",
   "(eq (1 (2 3)) (1 (2 4)))",
   "()"},
  {
   "Quote symbol literal",
   "(quote a)",
   "a"},
  {
   "Label number literal",
   "(label a 1 a)",
   "1"},
  {
   "Lambda with cond",
   "((lambda (a) (cond (eq a 1) 2 (eq a 3) 4)) 3)",
   "4"},
  {
   "Lambda with no args",
   "((lambda () 1))",
   "1"},
  {
   "Lambda called from label",
   "(label a (lambda () 1) (a))",
   "1"},
  {
   "Identity macro",
   "(label y 1 (label m (macro (x) x) (m y)))",
   "1"},
  {
   "Wrapping macro",
   "(label y 1 (label m (macro (x) (label y 2 x)) (m y)))",
   "2"},
  {
   "List identity",
   "(1 2)",
   "(1 2)"}
};

int
main (int argc, char *argv[])
{
  value_t core_env = callow_core ();

  int fail = 0;
  int i;
  for (i = 0; i < sizeof (parse_test_cases) / sizeof (parse_test_cases[0]);
       i++)
    {
      char **args = parse_test_cases[i];
      fail += check_parse (args[0], args[1], args[2]);
    }
  for (i = 0;
       i < sizeof (lookup_test_cases) / sizeof (lookup_test_cases[0]); i++)
    {
      char **args = lookup_test_cases[i];
      fail += check_lookup (args[0], args[1], args[2], args[3]);
    }
  for (i = 0; i < sizeof (eval_test_cases) / sizeof (eval_test_cases[0]); i++)
    {
      char **args = eval_test_cases[i];
      value_t test_env = read_string (args[1]);
      fail += check_eval (args[0], test_env, args[2], args[3]);
    }
  for (i = 0; i < sizeof (core_test_cases) / sizeof (core_test_cases[0]); i++)
    {
      char **args = core_test_cases[i];
      fail += check_eval (args[0], core_env, args[1], args[2]);
    }

  if (fail == 0)
    {
      printf ("\nALL TESTS PASSED!\n\n");
    }
  else
    {
      printf ("\n%d FAILED TESTS!\n\n", fail);
    }
}
