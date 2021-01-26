--2021-01-23 13:21:06
--NGINX openResty, 聊天socket链接分发,除了RabbitMQ之外的第二套替代方案

local server = require "resty.websocket.server"
local cjson = require("cjson") --json操作
local restyRedis = require("resty.redis") --Redis操作

--封装Redis
local redis = restyRedis.new()
local MyRedis = { db_index = 0, use_pool = false }
--redis connect
function MyRedis:connect()
    --设置超时(毫秒)
    redis:set_timeout(2000)
    --建立连接
    local ok, err = redis:connect("127.0.0.1", 6379)
    if not ok then
        self.errMsg("connect", err)
        return false
    end
    -- 如果有密码就用这个local res, err = redis:auth("password")
    --连接池状态无法使用库选择
    if not self.use_pool then
        redis:select(self.db_index)
    end
    return true
end

--redis get
function MyRedis:get(key)
    local resp, err = redis:get(key)
    if not resp then
        self.close()
        self.errMsg("get", err)
        return "0"
    end
    return resp
end

--redis set
function MyRedis:set(key, value)
    local ok, err = redis:set(key, value)
    if not ok then
        self.close()
        self.errMsg("set", err)
        return false
    end
    return true
end


--处理websocket
local wb, err = server:new {
    timeout = 5000, -- in milliseconds
    max_payload_len = 65535,
}
if not wb then
    ngx.log(ngx.ERR, "failed to new websocket: ", err)
    return ngx.exit(444)
end

local data, typ, err = wb:recv_frame()

if not data then
    if not string.find(err, "timeout", 1, true) then
        ngx.log(ngx.ERR, "failed to receive a frame: ", err)
        return ngx.exit(444)
    end
end

if typ == "close" then
    -- for typ "close", err contains the status code
    local code = err

    -- send a close frame back:

    local bytes, err = wb:send_close(1000, "enough, enough!")
    if not bytes then
        ngx.log(ngx.ERR, "failed to send the close frame: ", err)
        return
    end
    ngx.log(ngx.INFO, "closing with status code ", code, " and message ", data)
    return
end

if typ == "ping" then
    -- send a pong frame back:

    local bytes, err = wb:send_pong(data)
    if not bytes then
        ngx.log(ngx.ERR, "failed to send frame: ", err)
        return
    end
elseif typ == "pong" then
    -- just discard the incoming pong frame

else
    ngx.log(ngx.INFO, "received a frame of type ", typ, " and payload ", data)
end

wb:set_timeout(1000) -- change the network timeout to 1 second

bytes, err = wb:send_text("Hello world")
if not bytes then
    ngx.log(ngx.ERR, "failed to send a text frame: ", err)
    return ngx.exit(444)
end

bytes, err = wb:send_binary("blah blah blah...")
if not bytes then
    ngx.log(ngx.ERR, "failed to send a binary frame: ", err)
    return ngx.exit(444)
end

local bytes, err = wb:send_close(1000, "enough, enough!")
if not bytes then
    ngx.log(ngx.ERR, "failed to send the close frame: ", err)
    return
end

