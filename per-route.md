# Per-route Middleware and Encryption

In addition to basic middleware usage that applies to all routes, you can futher use route-specific middleware in `mux`.

Here is the basic route usage shown before:

```go
	// serve hello message
	m.HandleFunc("/hello", func(w ndn.Sender, i *ndn.Interest) {
		w.SendData(&ndn.Data{
			Name:    ndn.NewName("/hello"),
			Content: []byte(time.Now().UTC().String()),
		})
	})
```

To specify more middleware for this route, simply put more middleware in the parameters. After `mux` applies _global_ middleware, it will then apply `Middleware1`, `Middleware2` and `Middleware3`, before the route handler is called.

```go
m.HandleFunc("/hello", func(w ndn.Sender, i *ndn.Interest) {
  ...
}, Middleware1, Middleware2, Middleware3)
```

## `Encryptor` and `Decryptor`

`Encryptor` and `Decryptor` replace `AESEncryptor` and `AESDecryptor` in the new release of `mux`. They use RSA-OAEP to distribute shared AES-128 key, which will be used to encrypt and decrypt in CTR mode.

To use `Encryptor` and `Decryptor`, you need to prepare RSA key for the parameters.

> `*ndn.RSAKey` implements `ndn.Key`. You might need to type assert: `key.(*ndn.RSAKey)`.

For example, Cathy wants to share secrets with only Alice and Bob. In this case, Cathy will use `Encryptor`, and Alice/Bob will use `Decryptor`.

```go
// Cathy shares only with Alice and Bob
Encryptor(AlicePublicKey, BobPublicKey)

// Alice receives secrets from Cathy
Decryptor(AlicePrivateKey)
```

## Use `Encryptor` per-route

Now we can use `Encryptor` globally, but you can also imagine a situation where a producer generates secrets for different groups of people. This can be done easily with _per-route_ middleware.

For example, a producer wants to share secret1 with (A, B, C) and secret2 with (D, E, F).

```go
m.HandleFunc("/secret1", func(w ndn.Sender, i *ndn.Interest) {
  ...
}, Encryptor(A, B, C))

m.HandleFunc("/secret2", func(w ndn.Sender, i *ndn.Interest) {
  ...
}, Encryptor(D, E, F))
```
