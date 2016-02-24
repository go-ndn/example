# go-ndn TLV Extension

These changes are not part of the original TLV specification. In general, they are added to assist application development. While many can be introduced within `ndn.Data.Content` (a layer above ndn), they are so common that it is better to add in ndn layer to allow different implementations to process in the same way.

## EncryptionType

This provides encryption meta data for `ndn.Data.Content`. EncryptionKeyLocator specifies content key that can be used to decrypt `ndn.Data.Content`.

Currently only AES with CTR mode is implemented.

```go
const (
	EncryptionTypeNone       uint64 = 0
	EncryptionTypeAESWithCTR        = 1
)
```


## CompressionType

This provides compression meta data for `ndn.Data.Content`.

Currently only GZIP, a widely-supported compression method, is implemented.

```go
const (
	CompressionTypeNone uint64 = 0
	CompressionTypeGZIP        = 1
)
```

## ndn.Data.SignatureInfo.ValidityPeriod

This specifies a period when the signature of a data packet is valid.

```go
type ValidityPeriod struct {
	NotBefore string `tlv:"254"`
	NotAfter  string `tlv:"255"`
}

const (
	ISO8601 = "20060102T150405"
)
```

## Extra SignatureType

`CRC32C` is supported as a faster hash function than `SHA256`.

`HMAC` symmetric signature can be helpful if a node cannot generate asymmetric signatures for various reasons.

```go
const (
	SignatureTypeDigestCRC32C   uint64 = 2
	SignatureTypeSHA256WithHMAC        = 4
)
```

## CacheHint

This specifies how producer wants each node to cache data. There are many kinds of data packets that have little value to be cached; for those kinds, producer is encouraged to use `NoCache` to let intermediate nodes reserve cache space for other useful packets.

```go
const (
	CacheHintNone    uint64 = 0
	CacheHintNoCache        = 1
)
```

## `{Max,Min}Components`

The range of both MaxSuffixComponents and MinSuffixComponents  is `[0, +Inf)`. The range takes the whole value space of unsigned integer. However, in many systems, 0 (zero value) is often a convenient default, which is a result of calling `memset`. Semantically, 0 should be an invalid range value, indicating that these selectors are disabled.

After they are replaced them with MaxComponents and MinComponents, the range becomes `[1, +Inf)`, and 0 is available as a default value.
