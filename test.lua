local callow = require "libcallow"

local function check_read (name, test, expect)
   local t = callow.read(test)
   local actual = callow.write(t)
   if actual ~= expect then
      if not actual then
	 actual = "<nil>"
      end
      print("FAIL " .. name)
      print("Expected " .. expect .. " but got " .. actual)
      return 1
   end
   return 0
end

local function test_read ()
   local fail = 0
   fail = fail + check_read("read symbol", "joe", "joe")
   return fail
end

local fail_total = 0
fail_total = fail_total + test_read()
if fail_total == 0 then
   print("ALL TEST PASSED!")
end

-- callow.print(callow.read("(joe 1 2 3 (4  5)   )"))
