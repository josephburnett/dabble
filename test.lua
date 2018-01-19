local callow = require "libcallow"

fail = 0

local function check_read (name, test, expect)
   local t = callow.read(test)
   local actual = callow.write(t)
   if actual ~= expect then
      if not actual then
         actual = "<nil>"
      end
      print("FAIL " .. name)
      print("Expected " .. expect .. " but got " .. actual)
      fail = fail + 1
   end
end

check_read("read symbol", "joe", "joe")
check_read("read number", "1", "1")
check_read("read nil", "()", "()")
check_read("read list with one symbol", "(a)", "(a)")
check_read("read list with multiple symbols",
	   "(a b)", "(a b)")
check_read("read list with one number", "(1)", "(1)")
check_read("read list with multiple numbers",
	   "(1 2)", "(1 2)")
check_read("read list with one nil", "(())", "(())")
check_read("read list with multiple nils",
	   "(() ())", "(() ())")
check_read("read list with all types",
       "(a 1 (b 2) ())", "(a 1 (b 2) ())")
check_read("read nested lists", "(((a)))", "(((a)))")
check_read("read string", "\"abc\"", "(a b c)")
check_read("read two strings", "(\"ab\" \"cd\")",
	   "((a b) (c d))")
check_read("read empty string as nil", "\"\"", "()")

local function check_read_error (name, test)
   local t = callow.read(test)
   local actual = callow.write(t)
   local expect = "<error"
   if string.sub(actual, 1, string.len(expect)) ~= expect then
      if not actual then
         actual = "<nil>"
      end
      print("FAIL " .. name)
      print("Expected error but got " .. actual)
      fail = fail + 1
   end
end

check_read_error("unbalanced parens left", "(()")
check_read_error("unbalanced parens right", "())")
check_read_error("symbol beginning with number",
                 "123bad")
check_read_error("number with two decimals", "1.2.3")
check_read_error("invalid characters", "[]")

local function check_eval (name, test, expect)
   local t = callow.eval(test)
   local actual = callow.write(t)
   if actual ~= expect then
      if not actual then
         actual = "<nil>"
      end
      print("FAIL " .. name)
      print("Expected " .. expect .. " but got " .. actual)
      fail = fail + 1
   end
end

check_eval("car of single element list",
           "(car (1))", "1")
check_eval("car of multiple element list",
           "(car (1 2))", "1")
check_eval("car evals args",
           "(label a 1 (car (a 2)))", "1")

check_eval("cdr of single element list",
           "(cdr (1))", "()")
check_eval("cdr of multiple element list",
           " (cdr (1 2))", "(2)")
check_eval("cdr evals args",
           "(label a 2 (cdr (1 a)))", "(2)")

check_eval("list of nil", "(list ())", "()")
check_eval("list of number", "(list 1)", "()")
check_eval("list of symbol", "(list t)", "()")
check_eval("list of single element list",
           "(list (1))", "t")
check_eval("list evals args",
           "(label a 1 (list (a)))", "t")

check_eval("cons number with nil",
           "(cons 1 ())", "(1)")
check_eval("cons number with one element list",
           "(cons 1 (2))", "(1 2)")
check_eval("cons list with list",
           "(cons (1) (2))", "((1) 2)")
check_eval("cons nil with list",
           "(cons () (1))", "(() 1)")
check_eval("cons nil with nil",
           "(cons () ())", "(())")
check_eval("cons evals args",
           "(label a 1 (cons a ()))", "(1)")

check_eval("eq nils", "(eq () ())", "t")
check_eval("eq numbers", "(eq 1 1)", "t")
check_eval("eq lists", "(eq (1) (1))", "t")
check_eval("eq nested lists",
           "(eq (1 (2 3) 4) (1 (2 3) 4))", "t")
check_eval("eq not number and nil",
           "(eq 1 ())", "()")
check_eval("eq not number and list",
           "(eq 1 (1))", "()")
check_eval("eq evals args",
           "(label a 1 (eq a 1))", "t")

check_eval("cond one condition",
           "(cond (quote t) 1)", "1")
check_eval("cond evals match",
           "(label a 1 (cond (quote t) a))", "1")
check_eval("cond eval only match",
           "(cond () a (quote t) 1)", "1")
check_eval("cond any number truthy",
           "(cond 0 1)", "1")
check_eval("cond any symbol truthy",
           "(cond (quote a) 1)", "1")
check_eval("cond any list truthy",
           "(cond (0) 1)", "1")

check_eval("quote number", "(quote 1)", "1")
check_eval("quote nil", "(quote ())", "()")
check_eval("quote list", "(quote (1 2))", "(1 2)")
check_eval("quote symbol", "(quote a)", "a")
check_eval("quote not eval args",
           "(label a 1 (quote a))", "a")

check_eval("label bind number",
           "(label a 1 a)", "1")
check_eval("label bind nil",
           "(label a () a)", "()")
check_eval("label bind list",
           "(label a (1) a)", "(1)")
check_eval("label shadowing binding",
           "(label a 1 (label a 2 a))", "2")
check_eval("label nested scope",
           "(label a 1 (label b 2 a))", "1")
check_eval("label resolving symbol",
           "(label a 1 (label b a b))", "1")
check_eval("label evals bound value",
           "(label a 1 (label b a b))", "1")

check_eval("lambda with no args",
           "(label f (lambda () 1) (f))", "1")
check_eval("lambda with one arg",
           "(label f (lambda (a) a) (f 1))", "1")
check_eval("lambda inline execution",
           "((lambda () 1))", "1")
check_eval("lambda captures environment",
           "(label a 1 (label f (lambda () a) (label a 2 (f))))",
           "1")
check_eval("lambda evals args",
           "(label a 1 ((lambda (a) a) a))", "1")

check_eval("macro indentity",
           "(label y 1 (label m (macro (x xs) x) (m y)))",
           "1")
check_eval("macro cannot capture variables",
           "(label y 1 (label m (macro (x xs) (label y 2 x)) (m y)))",
           "1")
check_eval("macro with cond",
           "(label m (macro (x xs) (cond (eq x (quote y)) 1 (eq x (quote z)) 2)) (m (quote z)))",
           "2")
check_eval("macro captures env from definition",
           "(label m (label a 1 (macro (x xs) (cons a x))) (m (2)))",
           "(1 2)")
check_eval("macro captures env from args",
           "(label a 1 (label m (macro (x xs) (label a 2 x)) (m a)))",
           "1")

check_eval("recur calls lambda",
	   "((lambda (x)" ..
	   "   (cond" ..
	   "     (list (cdr x)) (recur (cdr x))" ..
	   "     (quote t) (car x))) (1 2 3))",
	   "3")
check_eval("recur calls macro",
	   "((macro (x xs)" ..
	   "   (cond" ..
	   "     (list (cdr x)) (recur (cdr x))" ..
	   "     (quote t) (car x))) (1 2 3))",
	   "3")
check_eval("recur returns nil outside lambda",
	   "(recur)",
	   "()")
check_eval("recur in lambda with macro",
	   "(label wrap (macro (xs) xs)" ..
	   "  ((lambda (a)" ..
	   "     (cond" ..
	   "       (list a) (cons (wrap (car a))" ..
           "                       (recur (cdr a)))" ..
	   "       (quote t) a))" ..
	   "   (1 2 3)))",
	   "((1) (2) (3))")
check_eval("recur in lambda with recursive macro",
	   "(label wrap" ..
           "  (macro (x xs)" ..
	   "    (cond" ..
	   "      (list x) (cons (cons (car x) ())" ..
	   "                      (recur (cdr x)))" ..
	   "      (quote t) x))" ..
	   "  ((lambda (a)" ..
	   "     (cond" ..
	   "       (list a) (cons (wrap (car a))" ..
           "                       (recur (cdr a)))" ..
	   "       (quote t) a))" ..
	   "     ((1 2 3) (4 5 6) (7 8 9))))",
           "(((1) (2) (3)) ((4) (5) (6)) ((7) (8) (9)))")

local function check_eval_error (name, test)
   local t = callow.eval(test)
   local actual = callow.write(t)
   local expect = "<error"
   if string.sub(actual, 1, string.len(expect)) ~= expect then
      if not actual then
         actual = "<nil>"
      end
      print("FAIL " .. name)
      print("Expected error but got " .. actual)
      fail = fail + 1
   end
end

check_eval_error("car of nil", "(car ())")
check_eval_error("car of number", "(car 1)")
check_eval_error("car with no args", "(car)")
check_eval_error("car with multiple args",
                 "(car (1) (1))")

check_eval_error("cdr of nil", "(cdr ())")
check_eval_error("cdr of number", "(cdr 1)")
check_eval_error("cdr with no args", "(cdr)")
check_eval_error("cdr with multiple args",
                 "(cdr (1) (1))")

check_eval_error("list with no args", "(list)")
check_eval_error("list with multiple args",
                 "(list 1 2)")

check_eval_error("cons with non-list",
                 "(cons 1 2)")
check_eval_error("cons with no args", "(cons)")
check_eval_error("cons with one arg", "(cons 1)")
check_eval_error("cons with multiple args",
                 "(cons () () ())")

check_eval_error("eq with no args", "(eq)")
check_eval_error("eq with one arg", "(eq 1)")
check_eval_error("eq with multiple args",
                 "(eq 1 1 1)")

check_eval_error("cond no matching",
                 "(cond (() 1))")
check_eval_error("cond no args", "(cond)")
check_eval_error("cond short clause", "(cond (1))")
check_eval_error("cond long clause",
                 "(cond (1 1 1))")

check_eval_error("quote with no args", "(quote)")
check_eval_error("quote with multiple args",
                 "(quote 1 2)")

check_eval_error("label with no args", "(label)")
check_eval_error("label with one arg", "(label a)")
check_eval_error("label with two args",
                 "(label a 1)")
check_eval_error("label with multiple args",
                 "(label a 1 a a)")

check_eval_error("lambda with no args", "(lambda)")
check_eval_error("lambda with too many args",
                 "(lambda () 1 1)")
check_eval_error("lambda with non-list first arg",
                 "(lambda 1 1)")
check_eval_error("lambda with non-symbol name",
                 "(lambda (1) 1)")

check_eval_error("macro with no args", "(macro)")
check_eval_error("macro with no names",
                 "(macro () 1)")
check_eval_error("macro with no body",
                 "(macro (x xs))")
check_eval_error("macro with non-symbol args",
                 "(macro (1) 2)")

local function test (name)
   local expect = "t"
   local t = callow.eval("(import (\"tst/" .. name .. "\") (test))")
   local actual = callow.write(t)
   if actual ~= expect then
      if not actual then
         actual = "<nil>"
      end
      print("FAIL test " .. name)
      print("Expected " .. expect .. " but got " .. actual)
      fail = fail + 1
   end
end

test("import_simple_binding")
test("core/and")

if fail == 0 then
   print("ALL TEST PASSED!")
end

