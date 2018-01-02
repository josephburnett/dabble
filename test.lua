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

if fail == 0 then
   print("ALL TEST PASSED!")
end
