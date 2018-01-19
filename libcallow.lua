
local string = require "string"

-- Environment --

local callow_root = os.getenv("CALLOW_ROOT")
if not callow_root then
   print("Please set CALLOW_ROOT to point to the repo.")
   os.exit(1)
end

local callow_path = os.getenv("CALLOW_PATH")
if not callow_path then
   callow_path = callow_root
end

-- Types --

local function _nil ()
   return {
   	type = "nil",
   }
end

local function is_nil (v)
   return type(v) == "table" and v.type == "nil"
end

local function _error (msg)
   return {
      type = "error",
      msg = msg,
   }
end

local function is_list (v)
   return type(v) == "table" and v.type == "list"
end

local function _list (car, cdr)
   if not cdr then cdr = _nil() end
   if not is_nil(cdr) and not is_list(cdr) then
      return _error("non-list cdr")
   end
   local c = {
      type = "list",
      car = car,
      cdr = cdr,
      len = 1,
   }
   if not is_nil(cdr) then
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

local function _lambda (names, body, env, loop)
   return {
      type = "lambda",
      names = names,
      body = body,
      env = env,
      loop = loop,
   }
end

local function _macro (names, body, env)
   return {
      type = "macro",
      names = names,
      body = body,
      env = env,
   }
end

local function is_fn (v)
   return type(v) == "table" and v.type == "fn"
end

local function is_lambda (v)
   return type(v) == "table" and v.type == "lambda"
end

local function is_macro (v)
   return type(v) == "table" and v.type == "macro"
end

local function is_error (v)
   return type(v) == "table" and v.type == "error"
end

local function is_symbol (v)
   return type(v) == "table" and v.type == "symbol"
end

local function is_number (v)
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
      if is_error(v) then return v end
      if not v then return end
      if rest then
         local cdr, rest = _read_list(rest)
         return _list(v, cdr), rest
      else
         return _list(v)
      end
   end

   local function _read_string (s)
      if string.len(s) == 0 then
	 return _nil()
      else
	 local sym = _symbol(string.sub(s, 1, 1))
	 local rest = string.sub(s, 2, -1)
	 return _list(sym, _read_string(rest))
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

   local str, index = string.match(a, "^%\"([^\"]*)%\"()")
   if str then
      return _read_string(str), string.sub(a, index + 1, -1)
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

local function _read_all (str)
   local v, extra = _read(str)
   if extra then
      extra = _strip(extra)
      if extra ~= "" then
         return _error("extra input found: " .. extra)
      end
   end
   return v
end

-- Printing --

local function _write (v)
   local out = ""
   if is_list(v) then
      out = out .. "("
      repeat
         out = out .. _write(v.car)
         if v.len > 1 then
            out = out .. " "
         end
         v = v.cdr
      until is_nil(v)
      out = out .. ")"
   elseif is_symbol(v) then
      out = out .. v.sym
   elseif is_number(v) then
      out = out .. v.num
   elseif is_nil(v) then
      out = out .. "()"
   elseif is_error(v) then
      out = out .. string.format("<error: %s>" , v.msg)
   elseif type(v) == "table" and v.type then
      out = out .. "<" .. v.type .. ">"
   else
      out = out .. "<unknown: " .. type(v) .. ">"
   end
   return out
end

-- Internal Functions --

local function list_len (v)
   if not is_list(v) then
      return 0
   end
   return v.len
end

local function equals (a, b)
   if _type(a) ~= _type(b) then return false end
   if is_symbol(a) and a.sym == b.sym then return true end
   if is_number(a) and a.num == b.num then return true end
   if is_nil(a) then return true end
   if is_list(a) and
      equals(a.car, b.car) and
      equals(a.cdr, b.cdr)
   then
      return true
   end
   return false
end

local function _lookup (sym, env)
   if is_nil(env) then
      return _error("Lookup of symbol " ..
		       sym.sym .. " failed.")
   end
   if equals(env.car.car, sym) then
      return env.car.cdr.car
   end
   return _lookup(sym, env.cdr)
end

local _eval, _eval_lambda, _expand_macro

local function _bind (sym, val, env)
   local v = _eval(val, env)
   return _list(_list(sym, _list(v)), env)
end

_eval = function (v, env, loop)
   if is_number(v) or
      is_error(v) or
      is_fn(v) or
      is_lambda(v) or
      is_macro(v)
   then
      return v
   end
   if is_symbol(v) then
      local value= _lookup(v, env)
      if is_symbol(value) then return value end
      return _eval(value, env, loop)
   end
   if is_list(v) then
      local first = _eval(v.car, env, loop)
      if is_fn(first) then
         return first.fn(v.cdr, env, loop)
      end
      if is_lambda(first) then
         local args = _eval(v.cdr, env, loop)
         if is_error(args) then return args end
         return _eval_lambda(first, args, loop)
      end
      if is_macro(first) then
         local body = _expand_macro(first, v.cdr, env, loop)
         if is_error(body) then return body end
         return _eval(body, env, loop)
      end
      return _list(first, _eval(v.cdr, env, loop))
   end
   return v
end

_eval_lambda = function (l, args, loop)
   if list_len(args) ~= list_len(l.names) then
      return _error("<lambda> requires " .. list_len(l.names) ..
                    " args. " .. list_len(args) .. " provided.")
   end
   local n, a, env = l.names, args, l.env
   while not is_nil(n) do
      env = _bind(n.car, a.car, env)
      n = n.cdr
      a = a.cdr
   end
   if l.loop then
      return _eval(l.body, env, l.loop)
   else
      return _eval(l.body, env, l)
   end
end

function _sub (arg, sym, val)
   if is_list(arg) then
      return _list(_sub(arg.car, sym, val), 
                   _sub(arg.cdr, sym, val))
   end
   if is_symbol(arg) and equals(arg, sym) then
      return val
   end
   return arg
end

function _thunk (body, env, loop)
   return _list(_lambda(_nil(), body, env, loop))
end

_expand_macro = function (m, args, env, loop)
   if list_len(args) < list_len(m.names) - 1 then
      return _error("<macro> requires at least " ..
                    list_len(m.names) - 1 .. " args. " ..
                    list_len(args) .. " provided.")
   end
   local n, a, b, i = m.names, args, m.body, 1
   while i <= list_len(m.names) - 1 do
      arg = _thunk(a.car, env, loop)
      b = _sub(b, n.car, arg)
      if is_error(b) then return b end
      n = n.cdr
      a = a.cdr
      i = i + 1
   end
   b = _sub(b, n.car, _thunk(a, env), loop)
   return _thunk(b, m.env, m)
end

local function list_to_string (l)
   if not is_list(l) then
      return nil, "list_to_string requires callow list. " ..
	 _type(l) .. " provided."
   end
   local str = ""
   while not is_nil(l) do
      local sym = l.car
      if not is_symbol(sym) then
	 return nil, "list_to_string requres all symbols. " ..
	    _type(sym) .. " found."
      end
      str = str .. _write(sym)
      l = l.cdr
   end
   return str, nil
end

local function _import_file (filename, env)
   local file, err = io.open(callow_root .. "/src/" ..
		             filename .. ".clw", "r")
   if not file then
      return _error("could not import " .. filename ..
		    ": " .. err)
   end
   local data = _read_all(file:read("a"))
   if is_error(data) then return data end
   if not is_list(data) then
      return _error("import expects top level list in " ..
		       filename .. ". provided  " .. _type(data))
   end
   local import_env = env
   repeat
      local pair = data.car
      if not is_list(pair) then
	 return _error("import expects list pairs. " ..
			  _type(pair) .. " provided.")
      end
      if list_len(pair) ~= 2 then
	 return _error("import expects list pairs. list " ..
			  "of length " .. list_len(pair) ..
			  "provided.")
      end
      local sym, value = pair.car, pair.cdr.car
      if not is_symbol(sym) then
	 return _error("import expects first symbol in pair. " ..
			  _type(sym) .. " provided.")
      end
      import_env = _bind(sym, value, import_env)
      data = data.cdr
   until is_nil(data)
   return import_env
end

-- Lisp Functions --

local function car (args, env, loop)
   args = _eval(args, env, loop)
   if list_len(args) ~= 1 then
      return _error("car requires 1 argument. " ..
		       list_len(args) .. " provided.")
   end
   if not is_list(args.car) then
      return _error("car requires a list argument. " ..
		      _type(args.car) .. " provided.")
   end
   return args.car.car
end

local function cdr (args, env, loop)
   args = _eval(args, env, loop)
   if list_len(args) ~= 1 then
      return _error("cdr requires 1 argument. " ..
		       list_len(args) .. " provided.")
   end
   if not is_list(args.car) then
      return _error("cdr requires a list argument. " ..
		       _type(args.car) .. " provided.")
   end
   return args.car.cdr
end

local function list (args, env, loop)
   args = _eval(args, env, loop)
   if list_len(args) ~= 1 then
      return _error("list requires 1 argument. " ..
		       list_len(args) .. " provided.")
   end
   if is_list(args.car) then
      return _symbol("t")
   else
      return _nil()
   end
end

local function cons (args, env, loop)
   args = _eval(args, env, loop)
   if list_len(args) ~= 2 then 
      return _error("cons requires 2 arguments. " ..
		       list_len(args) .. " provided.")
   end
   if is_nil(args.cdr.car) then
      return _list(args.car)
   end
   if not is_list(args.cdr.car) then
      return _error("cons requires a second list argument.  " ..
		       _type(args.cdr.car) .. " provided.")
   end
   return _list(args.car, args.cdr.car)
end

local function eq (args, env)
   args = _eval(args, env, loop)
   if list_len(args) ~= 2 then
      return _error("eq requires 2 arguments. " ..
		       list_len(args) .. " provided.")
   end
   if equals(args.car, args.cdr.car) then
      return _symbol("t")
   else
      return _nil()
   end
end

local function cond (args, env, loop)
   if list_len(args) == 0 then
      return _error("no matching condition")
   end
   if list_len(args) % 2 == 1 then
      return _error("cond requires an even number of arguments. " ..
		       list_len(args) .. " provided.")
   end
   local test = args.car
   local value = args.cdr.car
   if not equals(_nil(), _eval(test, env, loop)) then
      return _eval(value, env, loop)
   else
      return cond(args.cdr.cdr, env, loop)
   end
end

local function quote (args, env, loop)
   if list_len(args) ~= 1 then
      return _error("quote requires 1 arguments. " ..
		       list_len(args) .. " provided.")
   end
   return args.car
end

local function label (args, env, loop)
   if list_len(args) ~= 3 then
      return _error("label requires 3 arguments. " ..
		       list_len(args) .. " provided.")
   end
   if _type(args.car) ~= "symbol" then
      return _error("label requires first symbol argument. " ..
		       _type(args.car) .. " provided.")
   end
   local label_env = _bind(args.car, args.cdr.car, env)
   return _eval(args.cdr.cdr.car, label_env, loop)
end

local function lambda (args, env, loop)
   if list_len(args) ~= 2 then
      return _error("lambda requires 2 arguments. " ..
		       list_len(args) .. " provided.")
   end
   local names = args.car
   if not is_list(names) and
      not is_nil(names)
   then
      return _error("lambda requires first list or nil argument. " ..
		       _type(names) .. " provided.")
   end
   while not is_nil(names) do
      if not is_symbol(names.car) then
         return _error("lambda names must be symbols ." ..
                       _type(names.car) .. " provided.")
      end
      names = names.cdr
   end
   return _lambda(args.car, args.cdr.car, env)
end

local function macro (args, env, loop)
   if list_len(args) ~= 2 then
      return _error("macro requires 2 arguments. " ..
		       list_len(args) .. " provided.")
   end
   local names = args.car
   if not is_list(names) then
      return _error("macro requires first list argument. " ..
		       _type(names) .. " provided.")
   end
   while not is_nil(names) do
      if not is_symbol(names.car) then
         return _error("macro names must be symbols ." ..
                       _type(names.car) .. " provided.")
      end
      names = names.cdr
   end
   return _macro(args.car, args.cdr.car, env, loop)
end

local function recur (args, env, loop)
   if is_nil(loop) then return loop end
   return _eval(_list(loop, args), env, loop)
end

local function try (args, env, loop)
   if list_len(args) ~= 1 then
      return _error("try requires one argument. " ..
		       list_len(args) .. " provided.")
   end
   local result = _eval(args.car, env, loop)
   if is_error(result) then
      return _list(_nil(), _list(result))
   else
      return _list(_symbol("t"), _list(result))
   end
end

local function throw (args, env, loop)
   if list_len(args) ~= 1 then
      return _error("throw requires one argumnent. " ..
		       list_len(args) .. " provided.")
   end
   if not is_list(args.car) then
      return _error("throw requires a list argument. " ..
		       _type(args.car) .. " provided.")
   end
   local str, err = list_to_string(args.car)
   if err then
      return _error(err)
   end
   return _error(str)
end

local function import (args, env, loop)
   if list_len(args) ~= 2 then
      return _error("import requires 2 arguments. " ..
		       list_len(args) .. " provided.")
   end
   local file_list = args.car
   local body = args.cdr.car
   if not is_list(file_list) then
      return _error("import requires a first list argument. " ..
		       _type(file_list) .. " provided.")
   end
   local import_env = env
   repeat
      local filename, err = list_to_string(file_list.car)
      if err then
	 return _error(err)
      end
      import_env = _import_file(filename, env)
      if is_error(import_env) then return import_env end
      file_list = file_list.cdr
   until is_nil(file_list)
   return _eval(body, import_env, loop)
end

-- Library Exports --

local function _eval_std (v)
   local env = _nil()
   env = _bind(_symbol("car"), _fn(car), env)
   env = _bind(_symbol("cdr"), _fn(cdr), env)
   env = _bind(_symbol("list"), _fn(list), env)
   env = _bind(_symbol("cons"), _fn(cons), env)
   env = _bind(_symbol("eq"), _fn(eq), env)
   env = _bind(_symbol("cond"), _fn(cond), env)
   env = _bind(_symbol("quote"), _fn(quote), env)
   env = _bind(_symbol("label"), _fn(label), env)
   env = _bind(_symbol("lambda"), _fn(lambda), env)
   env = _bind(_symbol("macro"), _fn(macro), env)
   env = _bind(_symbol("recur"), _fn(recur), env)
   env = _bind(_symbol("try"), _fn(try), env)
   env = _bind(_symbol("throw"), _fn(throw), env)
   env = _bind(_symbol("import"), _fn(import), env)
   return _eval(v, env, _nil())
end

return {
   read = _read_all,
   write = _write,
   eval = _eval_std,
   is_nil = is_nil,
   is_symbol = is_symbol,
   is_number = is_number,
   is_list = is_list,
   is_error = is_error,
   is_fn = is_fn,
   is_lambda = is_lambda,
   is_macro = is_macro,
   list_len = list_len,
   list_to_string = list_to_string,
   equals = equals,
}
