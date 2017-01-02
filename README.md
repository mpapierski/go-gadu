go-gadu
=======

[![Build Status](https://travis-ci.org/mpapierski/go-gadu.svg?branch=master)](https://travis-ci.org/mpapierski/go-gadu)

`libgadu` bindings for **Go**.

# Usage

```sh
go get github.com/mpapierski/go-gadu
```

Example code:

```go
package main

import "github.com/mpapierski/go-gadu"
import "fmt"

func main() {
    fmt.Printf("Version: %s\n", gadu.Version())
}
```

# Authors

* Michał Papierski <michal@papierski.net>
