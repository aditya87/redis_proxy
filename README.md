### Redis Proxy

An HTTP Proxy for Redis, written in Go. The endpoints available are:

* GET `http://<proxy_address>?key=k` : Maps to the Redis GET command. Returns the value stored at the key `k`.
* POST `http://<proxy_address>` : Maps to the Redis SET command. Takes in a JSON body of the form `{"k": "v"}`, which SETs the value `v` for key `k`.

##### Running

You will need to have Docker installed. Run the following: <br>
```
$ make run redis_host=<REDIS_HOST> redis_port=<REDIS_PORT> port=<PORT> redis_pass=<REDIS_PASSWORD> capacity=<MAX_NUM_OF_CACHE_KEYS> expiry=<CACHE_EXPIRY_TIME_IN_SEC>
```

Where:
* REDIS_HOST=the host address of the backing Redis instance (default: localhost)
* REDIS_PORT=the TCP port on which the Redis instance accepts connections (default: 7000)
* PORT=the port on which the Proxy will be accepting connections (default: 3000)
* REDIS_PASSWORD=the password to connect to the Redis instance (no default)
* MAX_NUM_OF_CACHE_KEYS=the maximum number of keys that can be stored at a time in the local cache of the Proxy. (default: 20)
* EXPIRY=time in seconds for which a key can be stored in the local cache. (default: 30s)

##### Key architectural features

*Cache*:<br>
A GET request, directed at the proxy, returns the value of the specified key from the proxy’s local cache if the local cache contains a value for that key. If the local cache does not contain a value for the specified key, it fetches the value from the backing Redis instance, using the Redis GET command, and stores it in the local cache, associated with the specified key.

The cache capacity is set by the `capacity` parameter above. Keys expire after `expiry` seconds, after which they are removed from the cache. A subsequent GET on a removed key will hit the backing Redis instance.

The cache uses a Least Recently Used (LRU) eviction policy when it has hit its capacity and needs to evict a key in order to make room for a new key. Note that this does not supersede expiration time. In other words, expiration time is relative to when a key first enters the cache, and is not reset when a key is accessed. In hindsight I could have decided to reset the expiration timer, but I decided against it in the interest of time. It's a small change in the cache source code in order to implement that.

Tracking the keys in order of how recently they were accessed is done by the help of a simple linked list. Keys are moved to the back (tail) of the list when they are accessed, so that the head of the list is the least recently used key in the cache. New keys are inserted at the back of the list. When the cache hits its capacity, the list is popped so that the the LRU key gets evicted.

Removing keys that are expired is done with the help of a concurrent threaded subroutine that runs as a GoRoutine, executing in a loop once every `expiry/10` seconds. Keys that were inserted into the cache at least `expiry` seconds ago are removed.

I used the linked list available in the `container/list` package in Go. The `container/list` package provides a circularly linked list, so that inserting at the back/removing from the front are both constant time operations.

Therefore, the algorithmic complexity of inserting a key into the cache is O(1) (since we simply insert at the tail of the list), and the complexity of looking up a key is O(n) relative to the current size of the cache (since we have to iterate through the list to find the key).

*Sequential concurrent processing*:<br>
Multiple clients are able to concurrently connect to the proxy without adversely impacting the functional behavior of the proxy. When multiple clients make concurrent requests to the proxy, it is acceptable for them to be processed sequentially i.e. a request from the second only starts processing after the first request has completed and a response has been returned to the first client.

I enforced concurrency by the use of mutexes which control access to critical sections of the code. Mutexes are used both in the proxy server code and in the cache code. The proxy server mutex ensures that request processing is serialized, i.e. only one request can be processed at a time. The cache mutex enforces that the expiration time thread does not collide with getting/setting keys in the cache. This ensures isolated access to the cache itself within the application.

##### Testing

Unit tests: The code is unit tested with the help of the Ginkgo test runner (https://github.com/onsi/ginkgo). In order to run unit tests, simply run `make unit_test`

Integration tests: Integration tests are written as a set of functions in `integration/integration.go`. In order to run them, simply run `make test`
