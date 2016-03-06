# Generate Data Before Interest

`mux` is a reactive framework; it assumes that data packet is produced on interest. However sometimes we want to pre-generate data packets. There are two ways to do this with `mux`:

1. Directly add data packets to content store
2. Use `mux.Publisher`

## Idea 1: directly add data packet to content store

The idea is simple: after we add data packets to content store, and use the proper cacher middleware, `mux` should be able to answer interests directly. In this tutorial, we are going to use `mux.RawCacher`.

> __mux.RawCacher(ndn.Cache, cpy)__ returns a middleware for caching data packets. When `cpy` is true, a data packet is copied when it is written or read from cache. When `cpy` is false, no copying is done; it is useful if you don't need to change data packets and want more performance.

> All ndn.Cache implementations in go-ndn are assumed to be thread-safe.

To publish data, you need to create one ndn.Cache. In this case, we choose to use `persist.Cache` because the default in-memory content store is not persistent.

```go
// create new cache
// With this cache `c`, you can publish data packets from any go-routine.
c := persist.NewCache("published.db")

// publish data
c.Add(chunk1)
c.Add(chunk2)

m := mux.New()
// use cacher middleware created from this cache
// cpy is set to false because persist.Cache
// marshals and unmarshals data packet to bytes
m.Use(mux.RawCacher(c, false))
...
```

## Idea 2: use `mux.Publisher`

This is an extension to idea 1 by using existing middleware to post-process data.

```go
// create a publisher with cache
publisher := mux.NewPublisher(c)

// compress
publisher.Use(mux.Gzipper)
// after compress, segment
publisher.Use(mux.Segmentor(10))

// this blob will be compressed and then segmented
publisher.Publish(blob)
```
