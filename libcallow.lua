local string = require "string"

-- Types --

local function _nil ()
   return {
   	type = "nil",
   }
end

local function _is_nil (v)
   return type(v) == "table" and v.type == "nil"
end

local function _error (msg)
   return {
      type = "error",
      msg = msg,
   }
end

local function _is_list (v)
   return type(v) == "table" and v.type == "list"
end

local function _list (car, cdr)
   if not cdr then cdr = _nil() end
   if not _is_nil(cdr) and not _is_list(cdr) then
      return _error("non-list cdr")
   end
   local c = {
      type = "list",
      car = car,
      cdr = cdr,
      len = 1,
   }
   if not _is_nil(cdr) then
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
      num = num,
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

local function _is_symbol (v)
   return type(v) == "table" and v.type == "symbol"
end

local function _is_number (v)
   return type(v) == "table" and v.type == "number"
end

local function _type (v)
   if type(v) == "table" then
      return v.type
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
      if not v then return end
      if rest then
         local cdr, rest = _read_list(rest)
         return _list(v, cdr), rest
      else
         return _list(v)
      end
   end

   local a = _strip(str)
   
   if a == "" then
      return nil
   end
   
   if string.find(a, "^%(%)") then
      return _nil(), string.sub(a, 3, -1)
   end
   
   local list, index = string.match(a, "^(%b())()")
   if list then
      list = string.match(list, "^%((.*)%)")
      return _read_list(list), string.sub(a, index, -1)
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
   local out = ""
   if _is_list(v) then
      out = out .. "("
      repeat
         out = out .. _write(v.car)
         if v.len > 1 then
            out = out .. " "
         end
         v = v.cdr
      until _is_nil(v)
      out = out .. ")"
   elseif _is_symbol(v) then
      out = out .. v.sym
   elseif _is_number(v) then
      out = out .. v.num
   elseif _is_nil(v) then
      out = out .. "()"
   elseif _is_error(v) then
      out = out .. string.format("<error %s>" , v.msg)
   elseif type(v) == "table" and v.type then
      out = out .. "<" .. v.type .. ">"
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

local function _equals (a, b)
   if _type(a) ~= _type(b) then return false end
   if _is_symbol(a) and a.sym == b.sym then return true end
   if _is_number(a) and a.num == b.sym then return true end
   if _is_nil(a) then return true end
   if _is_list(a) and
      _equals(a.car, b.car) and
      _equals(a.cdr, b.cdr)
   then
      return true
   end
   return false
end

local function _bind (sym, v, env)
   return _list(_list(sym, _list(v)), env)
end

local function _lookup (sym, env)
   if _is_nil(env) then
      return _error("Lookup of symbol " ..
		       sym.sym .. " failed.")
   end
   if _equals(env.car.car, sym) then
      return env.car.cdr.car
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
      return _eval(_lookup(v, env), env)
   end
   if _is_list(v) then
      local eval_list = _list(_eval(v.car, env), _eval(v.cdr, env))
      if _is_fn(eval_list.car) then
         local fn = eval_list.car.fn
         local args = eval_list.cdr
         return fn(args, env)
      end
      -- TODO call lambda
      -- TODO expand macro
   end
   return v
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

local function _read_all (str)
   local v, extra = _read(str)
   if extra then
      extra = _strip(extra)
      if extra ~= "" then
         return _error("Extra input found: " .. extra)
      end
   end
   return v
end

local function _eval_std (str)
   local v = _read_all(str)
   local env = _nil()
   env = _bind(_symbol("car"), _fn(car), env)
   env = _bind(_symbol("cdr"), _fn(cdr), env)
   return _eval(v, env)
end

return {
   read = _read_all,
   write = _write,
   eval = _eval_std,
}
