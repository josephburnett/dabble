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

if fail == 0 then
   print("ALL TEST PASSED!")
end
