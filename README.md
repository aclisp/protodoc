# protodoc

Generate a markdown document for a Protobuf file.

Usage:

```
go get -u github.com/aclisp/protodoc
protodoc demo.proto > demo.md
```

## Example

See [demo.md](demo.md)

## Design

The Protobuf file generally has four sections:

1. Service definitions, containing RPC methods
2. Enumerations
3. User defined types, called *object*
4. RPC Request and Response objects

When generating a markdown document,
* <1> and <4> are combined together
* <2> and <3> are referenced by <4>

Inspired by [protobuf 为经络，gRPC为骨架](https://mp.weixin.qq.com/s/jMrkrLpPxzJA4GsHFHKs-Q).

## Limitation

Only a subset of [proto3](https://developers.google.com/protocol-buffers/docs/proto3) is supported as for best practice.

* No [map](https://developers.google.com/protocol-buffers/docs/proto3#maps), use the alternative [equivalent syntax](https://developers.google.com/protocol-buffers/docs/proto3#backwards_compatibility).
* Service rpc parameter should not has [nested types](https://developers.google.com/protocol-buffers/docs/proto3#nested).
