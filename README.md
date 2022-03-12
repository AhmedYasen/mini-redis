This is a challenge for Go solutions to the
["Build Your Own Redis" Challenge](https://codecrafters.io/challenges/redis).

- I implemented the Redis Protocol to serve the following commands
-  PING, ECHO, SET with expire time in ms, and GET


# Commands to run
PING
ECHO hey
SET k v
GET k
SET k2 v2 100
GET k2 (call this before the expire will return the value other will return nil)
