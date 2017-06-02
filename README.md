# harapr-jc/hashgen

## Summary

A package (hashgen) that supplies a crypto hash server and utilties for reuse.

## Dependencies

Linux supported. Depends on:

* Go standard library
* OS commands `uuidgen` and `grep`

## Usage

See project `jchash` for usage.

# Development

## Testing

The `hashgen` package has tests.
```
$ cd hashgen
$ go test -v
```

### Code Coverage

Sample coverage check:

```
$ cd hashgen
$ go test -coverprofile=cover.out
...
$ go tool cover -func=cover.out
github.com/harapr-jc/hashgen/dao.go:24:		New			100.0%
github.com/harapr-jc/hashgen/dao.go:29:		Append			71.4%
github.com/harapr-jc/hashgen/dao.go:65:		Get			90.9%
github.com/harapr-jc/hashgen/hasher.go:15:	getSalt			80.0%
github.com/harapr-jc/hashgen/hasher.go:26:	getCryptoHash		100.0%
github.com/harapr-jc/hashgen/lru.go:39:		NewCache		100.0%
github.com/harapr-jc/hashgen/lru.go:53:		Add			77.8%
github.com/harapr-jc/hashgen/lru.go:72:		Get			87.5%
github.com/harapr-jc/hashgen/server.go:51:	HandleGetHashRequest	100.0%
github.com/harapr-jc/hashgen/server.go:73:	HandleHashRequest	93.1%
github.com/harapr-jc/hashgen/server.go:133:	HandleStatsRequest	100.0%
github.com/harapr-jc/hashgen/server.go:138:	StartServer		90.3%
github.com/harapr-jc/hashgen/stats.go:38:	Accumulate		72.7%
github.com/harapr-jc/hashgen/stats.go:59:	GetJson			85.7%
github.com/harapr-jc/hashgen/uuid.go:11:	getUuid			75.0%
total:						(statements)		86.8%
```

### Benchmarking

Some of the tests support benchmarks.

```
$ cd hashgen
$ go test -bench=".*"
```

## Implementation Choices

### Job Id
The use of uuid instead of monotonically increasing id for a job identifier allows server
restart without having to persist the next id. Furthermore, no synchronization is required
on next id.

A drawback is that, without a database back end, ordering of uuid is less efficient for
persistent storage. A consequence is that records that are evicted from the cache will
take longer to find.

### LRU cache
The purpose for using the cache is to put an upper bound on memory usage for the server.

The design intent is that n most recently used jobs are held in cache, all are passed thru
to persistent storage. On a cache miss, job data could be fetched from persistent storage.
The default backing file is called 'backup.json'. One would prefer to use a performant data
store instead. The Go standard library has sql API, but no database drivers are included.

## Outstanding Development Items

* Raise code coverage
* Add safety for when backing file size approaches partition size
