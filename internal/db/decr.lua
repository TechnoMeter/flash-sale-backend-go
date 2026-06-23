local stock = redis.call('DECR', KEYS[1])
if stock < 0 then
    redis.call('INCR', KEYS[1])  -- rollback to 0
    return -2
end
return stock