# Stampy
Started as a pet project, ended up high performance, low footprint, non-heap native key-value store.

 - Embedded RESTFul interface with fast JSON processing
 - Non-intrusive TTL support for keys
 - Configurable Cache Partitioning
 - All in a nutshell, a single binary without any dependencies
 - Builds/Runs on OSX, Linux, Windows
 - Written with GO

#### Whats with the naming?
Stampy, Elephant In The Room !

(hint: Simpsons)

Also, Elephants never forget.

## Installation

`go get github.com/Gladmir/stampy`

## Usage

stampy -help;

    Usage of stampy:
      -buckets int
            Number of buckets for keys to be evenly distributed, higher numbers will increase concurrency with additional memory overhead (default 64)
      -configFile string
            All options are also configurable via config file in YAML format.
      -ip string
            A valid IPv4 address for serving restful interface, ex: 127.0.0.1 (default "0.0.0.0")
      -port int
            An unoccupied port for serving restful interface (default 4000)
Executing without providing any parameter is sufficient for stampy, it'll just run with its default configuration. You can also provide a configuration file written in YAML format.


## Interfaces

Stampy exposes only one (for now, at least) type of interface which is HTTP/RESTful access.

Methods and their effects on paths:

- GET


        resource: /stampy/v1/info

            {
                "Name": "Stampy, Elephant in the room",
                "Version": "0.0.1",
                "Os": "darwin-amd64",
                "CpuCores": 8,
                "MemoryUsage": "214.55 kb",
                "StampyBucketCount": 64,
                "Started": "2015-10-02T00:33:27.264431875+03:00"
            }

        resource: /stampy/v1/cache

            {
                "keyPuts": 0,
                "keyDeletes": 0,
                "keyHits": 0,
                "absentKeyHits": 0,
                "expiredKeys": 0,
                "expiredKeyHits": 0
            }

        resource: /stampy/v1/cache/{key}

            {
                "value":"123123",
                "creationDate":"2015-10-04T18:35:44.954754704+03:00",
                "lastAccessed":"2015-10-04T18:37:12.277701533+03:00",
                "expiryTime":"2015-10-04T18:41:44.956127762+03:00"
            }

Precision is so overrated on date fields, thou they will be tuned eventually.

- PUT


        resource :  /stampy/v1/cache/{key}
        json body:  {"value":"123123"} // permanent key.value
                    {"value":"123123", "timeToLive":360} // with ttl in seconds


- DELETE


        resource: /stampy/v1/cache/{key}

A special case for deleting;

Since stampy is a lazy-ass elephant, he cannot be sure if he has already deleted the key.

So, 'DELETE' operations are idiom potent, key might be expired, cleaned up but Stampy tries to delete it and silently discard if key is missing.

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request or buy me lunch.

## History

  - Released 0.0.1          [2015/4/10] :

  Refactoring, changes on REST interface payloads and YAML Configuration support

  - Released 0.0.1-alpha    [2015/2/10] :

  Initial working prototype

## Credits

Credits to Lovely cartoon character Stampy from Simpsons Season V

## License

MIT
