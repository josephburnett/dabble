local io = require "io"
local string = require "string"

local function strip (s)
  return (string.match(s, "%s*(.*)%s*"))
end

local function error (msg)
  return {
    _error = true,
    _msg = msg,
  }
end

local function is_error(v)
  return type(v) == "table" and v._error
end

local function is_list(v)
  return type(v) == "table" and not v._error
end

local function is_symbol(v)
  return type(v) == "string"
end

local function is_number(v)
  return type(v) == "number"
end

local function read(str)

  local function read_list (l)
    local list, v, rest = {}, nil, l
    repeat
      v, rest = read(rest)
      if is_error(v) then return v end
      list[#list + 1] = v
    until rest == nil
    return list
  end

  local a = strip(str)
  
  if a == "" then
    return nil
  end
  
  local list = string.match(a, "^%((.*)%)")
  if list then
    return read_list(list)
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

  return error("invalid input")
end

local function callow_print (v)
  if is_list(v) then
    n = #v
    io.write("(")
    for i,v in ipairs(v) do
      callow_print(v)
      if i ~= n then
        io.write(" ")
      end
    end
    io.write(")")
  elseif is_symbol(v) or is_number(v) then
    io.write(v)
  elseif is_error(v) then
    io.write(string.format("<error %s>" , v._str))
  else
    io.write("<unknown>")
  end
end

return {
  read = read,
  print = callow_print,
}