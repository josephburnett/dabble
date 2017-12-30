local string = require "string"

-- Types --

local function _list (car, cdr)
   local c = {
      type = "list",
      car = car,
      cdr = cdr,
      len = 1,
   }
   if cdr then
      c.len = cdr.len + 1
   end
   return c
end

local function _symbol (sym)
   return {
      type = "symbol",
      sym = sym,
   }
end

local function _number (num)
   return {
      type = "number",
      num = num
   }
end

local function _fn (fn)
   return {
      type = "fn",
      fn  = fn,
   }
end

local function _lambda (args, body, env)
   return {
      type = "lambda",
      args = args,
      body = body,
      env = env,
   }
end

local function _macro (args, body, env)
   return {
      type = "macro",
      args = args,
      body = body,
      env = env,
   }
end

local function _error (msg)
   return {
      type = "error",
      msg = msg,
   }
end

local function _is_fn (v)
   return type(v) == "table" and v.type == "fn"
end

local function _is_lambda (v)
   return type(v) == "table" and v.type == "lambda"
end

local function _is_macro (v)
   return type(v) == "table" and v.type == "macro"
end

local function _is_error (v)
   return type(v) == "table" and v.type == "error"
end

local function _is_list (v)
   return type(v) == "table" and v.type == "list"
end

local function _is_symbol (v)
   return type(v) == "table" and v.type == "symbol"
end

local function _is_number (v)
   return type(v) == "table" and v.type == "number"
end

local function _type (v)
   if type(v) == "table" then
      return v.type
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
      local v, rest = _read(l)
      if _is_error(v) then return v end
      if not v then return v end
      local cdr = nil
      if rest then
	 cdr = _read_list(rest)
      end
      return _list(v, cdr)
   end

   local a = _strip(str)
   
   if a == "" then
      return nil
   end
   
   local list = string.match(a, "^%((.*)%)")
   if list then
      return _read_list(list)
   end
   
   local sym, index = string.match(a, "^(%a%w*)()")
   if sym then
      return _symbol(sym), string.sub(a, index, -1)
   end
   
   local num, index = string.match(a, "^(%d)()")
   num = tonumber(num)
   if num then
      return _number(num), string.sub(a, index, -1)
   end

   return _error("invalid input")
end

-- Printing --

local function _write (v)
   out = ""
   if _is_list(v) then
      out = out .. "("
      repeat
	 _write(v.car)
	 if v.len > 1 then
	    out = out .. " "
	 end
	 v = v.cdr
      until not v
      out = out .. ")"
   elseif _is_symbol(v) then
      out = out .. v.sym
   elseif _is_number(v) then
      out = out .. v.num
   elseif _is_error(v) then
      out = out .. string.format("<error %s>" , v.msg)
   elseif type(v) == "table" and v.type then
      out = out .. "<" .. v.type .. ">"
   elseif v == nil then
      out = out .. "<nil>"
   else
      out = out .. "<unknown>"
   end
   return out
end

-- Internal Functions --

local function _len (v)
   if not _is_list(v) then
      return 0
   end
   return v.len
end

local function _bind (sym, v, env)
   return _list(_list(sym, v), env)
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
      _is_fn(v) or
      _is_lambda(v) or
      _is_macro(v)
   then
      return v
   end
   if _is_symbol(v) then
      return _eval(_lookup(v, env))
   end
   if _is_list(v) then
      local eval_list = _list(_eval(v.car), _eval(v.cdr))
      if _is_fn(eval_list.car) then
	 local fn = eval_list.car.fn
	 local args = eval_list.cdr
	 return fn(args, env)
      end
      -- TODO call lambda
      -- TODO expand macro
   end
   return nil
end

-- Lisp Functions --

local function car (args, env)
   args = _eval(args, env)
   if _len(args) ~= 1 then
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
   if _len(args) ~= 1 then
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
   if _len(args) ~= 1 then
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
   if _len(args) ~= 2 then 
      return _error("cons requires 2 arguments. " ..
		       _len(args) .. " provided.")
   end
   if not _is_list(args.cdr.car) then
      return _error("cons requires a second list argument.  " ..
		       _type(args.cdr.car) .. " provided.")
   end
   return _list(args.car, args.cdr.car)
end

local function label (args, env)
   if _len(args) ~= 3 then
      return _error("label requires 2 arguments. " ..
		       _len(args) .. " provided.")
   end
   if _type(args.car) ~= "symbol" then
      return _error("label requires first symbol argument. " ..
		       _type(args.car) .. " provided.")
   end
   local label_env = bind(args.car, args.cdr.car, env)
   return _eval(args.cdr.cdr.car, label_env)
end

local function lambda (args, env)
   if _len(args) ~= 2 then
      return _error("lambda requires 2 arguments. " ..
		       _len(args) .. " provided.")
   end
   if not _is_list(args.car) then
      return _error("lambda requires first list argument. " ..
		       _type(args.car) .. " provided.")
   end
   -- TODO verify all names are symbols
   return _lambda(args.car, args.cdr.car, env)
end

local function macro (args, env)
   if _len(args) ~= 2 then
      return _error("macro requires 2 arguments. " ..
		       _len(args) .. " provided.")
   end
   if not _is_list(args.car) then
      return _error("macro requires first list argument. " ..
		       _type(args.car) .. " provided.")
   end
end

-- Library Exports --

return {
   read = _read,
   write = _write,
}
