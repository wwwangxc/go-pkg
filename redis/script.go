package redis

var (
	luaScriptLock = `
if (redis.call('EXISTS', KEYS[1]) == 0)
then
  redis.call('HSET', KEYS[1], 'UUID', ARGV[2])
  redis.call('PEXPIRE', KEYS[1], ARGV[1])
  return redis.call('HINCRBY', KEYS[1], 'COUNT', 1)
end

if (redis.call('HGET', KEYS[1], 'UUID') == ARGV[2])
then
  redis.call('PEXPIRE', KEYS[1], ARGV[1])
  return redis.call('HINCRBY', KEYS[1], 'COUNT', 1)
end

return 0
`

	luaScriptUnlock = `
ret_key_not_exist = 0
ret_invalid_uuid = 1
ret_del_fail = 2
ret_success = 666

if (redis.call('EXISTS', KEYS[1]) == 0)
then
  return ret_key_not_exist
end

if (redis.call('HGET', KEYS[1], 'UUID') ~= ARGV[1])
then
  return ret_invalid_uuid
end

if (redis.call('HGET', KEYS[1], 'COUNT') > 1)
  redis.call('HINCRBY', KEYS[1], 'COUNT', -1)
  return ret_success
end

if (redis.call('DEL', KEYS[1]) == 0)
then
  return ret_del_fail
end

redis.call('PUBLISH', KEYS[1], 1)
return ret_success
`
)
