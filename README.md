yamlkeys
========

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/yamlkeys)](https://goreportcard.com/report/github.com/tliron/yamlkeys)

This Go library allows for decoding arbitrary YAML content, which includes maps with complex
keys, into basic Go data types (strings, ints, floats, bools, maps, and slices).

To quote from the [YAML specification](https://yaml.org/spec/1.2/spec.html#tag/repository/map):

> YAML places no restrictions on the type of keys; in particular, they are not restricted to
> being scalars.

Note that there are two notations for specifying complex keys in YAML. You can use a condensed
notation:

```yaml
{complex1: 0, complex2: 1}: value1
```

Or a multiline notation with the key and value specified separately:

```yaml
? complex1: 0
  complex2: 1
: value1
```

This often-overlooked feature of YAML is required by certain YAML-based formats, notably
[TOSCA](http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.3/TOSCA-Simple-Profile-YAML-v1.3.html#_Schema_Definition).

Importantly, the basic Go map is still used here in order to allow for broadest compatibility
with similar parsers, such as Go's default [JSON parser](https://golang.org/pkg/encoding/json/),
albeit with important caveats detailed below.

An alternative solution could use an entirely different map implementation, such as
[this one](https://github.com/cornelk/hashmap) or
[this one](https://godoc.org/github.com/timtadh/data-structures/hashtable). In weighing the
pros vs. the cons we preferred the basic Go map.

This library is intended to be used as an add-on for [go-yaml](https://github.com/go-yaml/yaml),
which was originally developed by Canonical. In the future we may also support
[Masaaki Goshima's go-yaml](https://github.com/goccy/go-yaml).

The former library can decode complex keys into its custom
[Node](https://godoc.org/gopkg.in/yaml.v3#Node) type, but will fail when decoding them into
a Go map ([see this playground](https://play.golang.org/p/TjlTlHeDIy_C)).

The latter library does not fail when decoding complex keys, but instead it silently converts
them to strings ([see this playground](https://play.golang.org/p/wqjFi5FshAd)). Note that converting
complex keys to strings is a workaround, not a solution. It's impossible to distinguish between keys
that are actual strings and complex keys that have been converted to strings, and also you would have
to convert your keys to strings, too. Our solution here does not have these problems.


Features and Limitations
------------------------

### Map Operations

Supporting complex keys is non-trivial in Go, which requires keys to be trivially comparable.
Primitives and structs of primitives are supported, but maps and slices are not. Our trick here
is to wrap complex keys in a special type and to use a *pointer* to it as the actual key.
Pointers "work" in that they can be used as keys without panicking (they are just integers),
but of course the basic Go map operations — get, put, delete — are unable to take into
consideration the actual key value.

For this reason we here provide replacements for basic Go map operations, which handle
wrapping/unwrapping and actual key comparison, as well as utility functions for working with
complex keys. These operations will work on both complex keys and simple keys, so that if you
stick to our versions then you will ensure the broadest compatibility.

Unfortunately, you *must* use our provided map operations. The basic Go get
(`value = map[complexKey]`) won't work. The basic Go put (`map[complexKey] = value`) would
appear to "work" but would allow for duplicates.

This will require discipline on your end, because there is no way to enforce this requirement
via the compiler. It is the cost of our insistence on using the basic Go map.

### Typed Errors

The go-yaml library does not return typed errors, making it difficult to extract error information,
such as the line and column in which an error occurred. For convenience we provide a `DecodeError`
with this information.

We do this not only for yamlkeys errors, but also convert go-yaml errors by parsing the error
message string.

### Multiple Documents

The go-yaml library's `decoder.Decode` function only decodes the first document it finds in the
stream and then stops. For compatibility, we have kept the same behavior here.

However, for convenience we also provide `DecodeAll` and `DecodeStringAll` functions that attempt
to decode the entire stream.


Usage Examples
--------------

```go
text := `
{complex1: 0, complex2: 1}: value1
{complex1: 0, complex2: 2}: value2
`

data, _ := yamlkeys.DecodeString(text)
map_ := data.(yamlkeys.Map)

// Iteration
for k, v := range map_ {
    fmt.Printf("key = %v, value = %v\n", yamlkeys.KeyData(k), v)
}

key := map[any]any{
    "complex1": 0,
    "complex2": 1,
}

// Get
v, _ := yamlkeys.MapGet(map_, key)
fmt.Printf("original value = %v\n", v)

// Put
yamlkeys.MapPut(map_, key, "value3")
v, _ = yamlkeys.MapGet(map_, key)
fmt.Printf("modified value = %v\n", v)

// Delete
yamlkeys.MapDelete(map_, key)

// Force keys to be strings (e.g. for compatibility with JSON)
for k, v := range map_ {
    fmt.Printf("key = %v, value = %v\n", yamlkeys.KeyString(k), v)
}
```

In [playground](https://play.golang.org/p/QYpGZhLnrMB).
