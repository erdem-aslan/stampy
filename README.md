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

Grab the binary from below link (will be available soon) and execute it with 'stampy -help';

    Usage of stampy:
      -buckets int
        	Number of buckets for keys to be evenly distributed, higher numbers will increase concurrency with a memory overhead (default 64)
      -ip string
        	A valid IPv4 address for serving restful interface, ex: 127.0.0.1 (default "0.0.0.0")
      -port int
        	An unoccupied port for serving restful interface (default 4000)

Executing without providing any parameter is sufficient for stampy, it'll just run with its default configuration.

## Usage

Stampy exposes only one type of interface which is HTTP/RESTful access.

Methods and their effects on paths:

- GET


        resource: /stampy/v1/info

    {
    	"Name": "Stampy, Elephant in the room",
    	"Version": "0.0.1-alpha",
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
    {"payload":"any kind of stringified payload",
     "creationDate":"2015-10-02T00:46:54.062119404+03:00",
     "lastAccessed":"2015-10-02T00:47:17.335877637+03:00",
     "validUntil":"2015-11-01T18:04:59.723630823+03:00"}

Precision is so overrated on date fields, thou they will be tuned eventually.

- PUT


        resource :  /stampy/v1/cache/{key}
        json body:  {"payload","any kind of stringified payload"} // permanent key.value
                    {"payload","any kind of stringified payload", "validUntil":"2015-10-02T00:46:54.062119404+03:00"} // with ttl


- DELETE


        resource: /stampy/v1/cache/{key}

A special case for deleting;

Since stampy is a lazy-ass elephant, he cannot be sure if he has already eaten the key.

So, 'DELETE' operations are idiom potent, key might be expired, cleanup up but Stampy tries to delete it and silently discard if key is missing.

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request or buy me lunch.

## History

Released 0.0.1-alpha [2015/2/10]

## Credits

Credits to Lovely cartoon character Stampy from Simpsons Season V

## License

Since this tiny project is not a first class Telenity Product, there is no licensing.