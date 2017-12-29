local io = require "io"
local string = require "string"

-- Types --

local function _cell(car, cdr)
   local c = {
      type = "cell",
      car = car,
      cdr = cdr,
      len = 1,
   }
   if cdr then
      c.len = cdr.len + 1
   end
   return c
end

local function _lambda (fn, env)
   return {
      type = "lambda",
      fn = fn,
      env = env,
   }
end

local function _error (msg)
   return {
      type = "error",
      msg = msg,
   }
end

local function _is_error (v)
   return type(v) == "table" and v.type == "error"
end

local function _is_list (v)
   return type(v) == "table" and v.type == "cell"
end

local function _is_lambda (v)
   return type(v) == "table" and v.type == "lambda"
end

local function _is_symbol_(v)
   return type(v) == "string"
end

local function _is_number (v)
   return type(v) == "number"
end

local function _type (v)
   if _is_list(v) then
      return "list"
   elseif _is_lambda(v) then
      return "lambda"
   elseif _is_symbol(v) then
      return "symbol"
   elseif _is_number(v) then
      return "number"
   elseif _is_error(v) then
      return "error"
   else
      return "nil"
   end
end

-- Parsing --

local function _strip (s)
   return (string.match(s, "%s*(.*)%s*"))
end

local function _read (str)

   local function _read_list (l)
      local list, v, rest = {}, nil, l
      repeat
	 v, rest = _read(rest)
	 if is_error(v) then return v end
	 list[#list + 1] = v
      until rest == nil
      return list
   end

   local a = _strip(str)
   
   if a == "" then
      return nil
   end
   
   local list = string.match(a, "^%((.*)%)")
   if list then
      return _read_list(list)
   end
   
   local symbol, index = string.match(a, "^(%a%w*)()")
   if symbol then
      return symbol, string.sub(a, index, -1)
   end
   
   local number, index = string.match(a, "^(%d)()")
   number = tonumber(number)
   if number then
      return number, string.sub(a, index, -1)
   end

   return _error("invalid input")
end

-- Printing --

local function _print (v)
   if _is_list(v) then
      local n = #v
      io.write("(")
      for i,v in ipairs(v) do
	 _print(v)
	 if i ~= n then
	    io.write(" ")
	 end
      end
      io.write(")")
   elseif _is_symbol(v) or _is_number(v) then
      io.write(v)
   elseif _is_error(v) then
      io.write(string.format("<error %s>" , v._str))
   else
      io.write("<unknown>")
   end
end

-- Internal Functions --

local function _len (v)
   if not _is_list(v) then
      return 0
   end
   return v.len
end

local function _bind (sym, v, env)
   return _cell(_cell(sym, v), env)
end

local function _lookup (sym, env)
   if not env then
      return _error("Lookup of symbol " ..
		       sym .. " failed.")
   end
   if env.car.car == sym then
      return env.car.cdr
   end
   return _lookup(sym, env.cdr)
end

local function _eval (v, env)
   if _is_number(v) or
      _is_error(v) or
      _is_lambda(v)
   then
      return v
   end
   if _is_symbol(v) then
      return _eval(_lookup(v, env))
   end
   if _is_list(v) then
      return _cell(_eval(v.car), _eval(v.cdr))
   end
   return nil
end

-- Lisp Functions --

local function car (args, env)
   args = _eval(args, env)
   if _len(args) != 1 then
      return _error("car requires 1 argument. " ..
		       _len(args) .. " provided.")
   end
   if not _is_list(args.car) then
      return _error("car requires a list argument. " ..
		      _type(args.car) .. " provided.")
   end
   return args.car.car
end

local function cdr (args, env)
   args = _eval(args, env)
   if _len(args) != 1 then
      return _error("cdr requires 1 argument. " ..
		       _len(args) .. " provided.")
   end
   if not _is_list(args.car) then
      return _error("cdr requires a list argument. " ..
		       _type(args.car) .. " provided.")
   end
   return args.car.cdr
end

local function list (args, env)
   args = _eval(args, env)
   if _len(args) != 1 then
      return _error("list requires 1 argument. " ..
		       _len(args) .. " provided.")
   end
   if _is_list(args.car) then
      return "t"
   else
      return nil
   end
end

local function cons (args, env)
   args = _eval(args, env)
   if _len(args) != 2 then 
      return _error("cons requires 2 arguments. " ..
		       _len(args) .. " provided.")
   end
   if not _is_list(args.cdr.car) then
      return _error("cons requires a second list argument.  " ..
		       _type(args.cdr.car) .. " provided.")
   end
   return _cell(args.car, args.cdr.car)
end

local function label (args, env)
   if _len(args) != 3 then
      return _error("label requires 2 arguments. " ..
		       _len(args) .. " provided.")
   end
   if _type(args.car) != "symbol" then
      return _error("label requires first symbol argument. " ..
		       _type(args.car) .. " provided.")
   end
   local label_env = bind(args.car, args.cdr.car, env)
   return _eval(args.cdr.cdr.car, label_env)
end

-- Library Exports --

return {
   read = _read,
   print = _print,
}
