# Verify Data Packet in go-ndn

> NDN new certificate format uses __pure__ data packet tlv with additional _ValidityPeriod_ and some other fields. go-ndn adopts it early, so both this tutorial and the code are based on the new certificate format.

Data verification in go-ndn mux is implemented with middleware:

- `ChecksumVerifier`: verifies only checksum, including `sha256` and `crc32c` (hardware-accelerated in SSE 4.2)
- `Verifier`: fetches and verifies actual signature, including `RSA` and `ECDSA`, with a list of verify rules

In general, you always want to add `ChecksumVerifier` to `Fetcher`, so you can at least verify integrity of data packets. However, `Verifier` is a bit tricky to set up, so we want to show you how to use it.

## Verify rule

`Verifier` accepts a list of verify rules that are compiled into your application. If a rule matches a data packet (rules are applied in order), that packet and its keys will be recursively verified; However, if no rule matches, that packet will not be verified. This allows `Verifier` to be _chained_ in `Fetcher`.

```go
fetcher.Use(mux.Verifier(
  &mux.VerifyRule{
    ...
  },
  &mux.VerifyRule{
    ...
  },
))
```

Every rule has `DataPattern`, which will match against data name. This pattern syntax documentation can be found on [Google RE2](https://github.com/google/re2/wiki/Syntax). Notice that this syntax is __not compatible__ with [NDN Regular Expression](http://named-data.net/doc/ndn-cxx/current/tutorials/utils-ndn-regex.html).

Then if `KeyPattern` exists, `Verifier` will use _captures_ from `DataPattern`, and compiles a new regexp to match against key locator name. The key will be fetched to verify the current data packet. This whole process is done recursively on the key. Finally, we will match some anchor rule, which has `DataSHA256`.

```go
&mux.VerifyRule{
  DataPattern: "^/go-ndn/([a-z]+)/post",
  KeyPattern: "^/go-ndn/$1", // capture ([a-z]+) in $1
}
...
// trust anchor
&mux.VerifyRule{
  DataPattern: "^/go-ndn",
  DataSHA256: "dba7a16fe89ce323117e9dcb8087dbe2eaf2cebaf5b5b30b2c7f3797c3b52550",
}
```

## Improve verification performance

If you don't want to fetch keys from remote every time, you can add more middleware to `Fetcher`.

- `persist.Cacher`: uses on-disk mmap packet store
- `mux.Cacher`: uses in-memory lru packet store

> Verification in mux is designed to __block__ because each interest is already handled in a separate goroutine. After [Go 1.5](https://golang.org/doc/go1.5), your application throughput should be maximized because by default, Go programs run with GOMAXPROCS set to the number of cores available; in prior releases it defaulted to 1.

