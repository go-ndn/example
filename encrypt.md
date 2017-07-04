# Encrypt Data Packet

In addition to basic middleware usage that applies to all routes, you can futher use route-level middleware in `mux`.

Here is the basic route usage shown before:

```go
	// serve hello message
	m.HandleFunc("/hello", func(w ndn.Sender, i *ndn.Interest) error {
		return w.SendData(&ndn.Data{
			Name:    ndn.NewName("/hello"),
			Content: []byte(time.Now().UTC().String()),
		})
	})
```

To specify more middleware for this route, simply put more middleware in the parameters. After `mux` applies _application-level_ middleware, it will then apply `Middleware1` and `Middleware2`, before the route handler is called.

```go
m.HandleFunc("/hello", func(w ndn.Sender, i *ndn.Interest) error {
  ...
}, Middleware1, Middleware2)
```

## Encryption and decryption middleware

`Encryptor` and `Decryptor` use RSA-OAEP to distribute shared AES-128 key, which will be used to encrypt and decrypt in CTR mode.

To use `Encryptor` and `Decryptor`, you need to prepare RSA key for the parameters.

> `*ndn.RSAKey` implements `ndn.Key`. You might need to type assert: `key.(*ndn.RSAKey)`.

For example, Cathy wants to share secrets with only Alice and Bob. In this case, Cathy will use `Encryptor`, and Alice/Bob will use `Decryptor`.

> In this example, `Encryptor` does not register "/cathy/encrypt".

```go
// Cathy shares only with Alice and Bob
m.Use(mux.Encryptor("/cathy/encrypt", AlicePublicKey, BobPublicKey))

// Alice receives secrets from Cathy
m.Use(mux.Decryptor(AlicePrivateKey))
```

## Use `Encryptor` at route level

Now we use `Encryptor` at application level, but there is another case that a producer generates secrets for different groups of people. This can be done easily with _route-level_ middleware.

For example, a producer wants to share secret1 with (A, B) and secret2 with (C).

> In this example, `Encryptor` does not register "/producer/encrypt".

```go
m.HandleFunc("/secret1", func(w ndn.Sender, i *ndn.Interest) error {
  ...
}, mux.Encryptor("/producer/encrypt", A, B))

m.HandleFunc("/secret2", func(w ndn.Sender, i *ndn.Interest) error {
  ...
}, mux.Encryptor("/producer/encrypt", C))
```
