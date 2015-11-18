# Publish Pre-generated Data Packets

Publishing per-generated data packets becomes simple after we introduce `mux.RawCacher`.

> __mux.RawCacher(ndn.Cache, cpy)__ returns a middleware for caching data packets. When `cpy` is true, a data packet is copied when it is written or read from cache. When `cpy` is false, no copying is done; it is useful if you don't need to change data packets and want more performance.

To publish data, you need to create one ndn.Cache. In this case, we choose to use `persist.Cache` because the default in-memory content store is not persistent.

```go
// create new cache
c := persist.NewCache("published.db")

// publish data
c.Add(gopherChunk1)
c.Add(gopherChunk2)

m := mux.New()
// use cacher middleware created from this cache
// cpy is set to false because persist.Cache
// marshals and unmarshals data packet to bytes
m.Use(mux.RawCacher(c, false))
...
```

With this cache `c`, you can publish data packets from any go-routine.

> All ndn.Cache implementations in go-ndn are assumed to be thread-safe.
