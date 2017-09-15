### next

- [go-ndn TLV extension](tlv-extension.md)
- propagate errors from `ndn.Sender` and `mux.Handler`
- add dockerfiles and [docker images](https://hub.docker.com/u/gondn/) for nfd and bridge
- move "go-ndn/health" to a private project

### 2017-02-22

- replace the core raft implementation with [etcd/raft](https://github.com/coreos/etcd/tree/master/raft) from core os
- use `mux.Verifier` in go-nfd
- remove `log.Fatalln`
- replace `health.Logger` with `health.InfluxDB`

### 2016-08-18

- support implicit digest lookup
- optimize persist content store
- reduce lpm memory usage
- support udp multicast
- remove global ndn.ContentStore

### 2016-02-21

- Release `1.3`
  - [Raft distributed consensus protocol with NDN transport](https://github.com/go-ndn/raft)
  - optimize scheduled pit removal
  - optimize lpm object allocation
  - add `mux.Publisher` to push data to content store with middleware

### 2015-11-28

- Release `1.2`
  - __ndn "send/push" semantics__ (Oli Gavin): This is an experimental protocol based on interest pattern `<nodeName>/ACK/<dataName>`. A producer will use this interest pattern to ensure that a specific node receives a specific data packet (see `mux.Notify` and `mux.Listener`).
  - add go-ndn node [health monitoring tool](https://github.com/go-ndn/health)
  - refactor go-nfd middleware (Oli Gavin)
  - fix tlv cannot unmarshal time.Time (reported by qhsong)
  - fix stale data in go-nfd content store (reported by Mahyuddin Husairi)
- Tutorial
  - [Generate Data Before Interest](publish.md)

### 2015-10-15

- Release `1.1`
  - refactor `packet` and `tlv` package for significantly less memory allocation. ([before](bench/2015-09-13.svg) and [after](bench/2015-09-27.svg))
  - experimental ndn certificate format
  - update verify (rsa, ecdsa, sha256, crc32c and hmac) and encrypt (RSA-OAEP and AES-CTR) middleware
- Tutorial
  - [Verify Data Packet](verify.md)
  - [Encrypt Data Packet](encrypt.md)

### 2015-09-13

- First public stable release `1.0`
