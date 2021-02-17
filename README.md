# MongoDB Ops Manager Go Client

[![Go Reference](https://pkg.go.dev/badge/go.mongodb.org/ops-manager.svg)](https://pkg.go.dev/go.mongodb.org/ops-manager)
![CI](https://github.com/mongodb/go-client-mongodb-ops-manager/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/go.mongodb.org/ops-manager)](https://goreportcard.com/report/go.mongodb.org/ops-manager)

A go client for [Ops Manager](https://docs.opsmanager.mongodb.com/master/reference/api/) 
and [Cloud Manager](https://docs.cloudmanager.mongodb.com/reference/api/) API.

Note that ops-manager only supports the two most recent major versions of Go.

## Usage

```go
import "go.mongodb.org/ops-manager/opsmngr"
```

Construct a new Ops Manager client, then use the various services on the client to
access different parts of the Ops Manager API. For example:

```go
client := opsmngr.NewClient(nil)
```

The services of a client divide the API into logical chunks and correspond to
the structure of the [Ops Manager API documentation](https://docs.opsmanager.mongodb.com/v4.4/reference/api/).

**NOTE:** Using the [context](https://godoc.org/context) package, one can easily
pass cancellation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then `context.Background()`
can be used as a starting point.

### Authentication

The ops-manager library does not directly handle authentication. Instead, when
creating a new client, pass an http.Client that can handle Digest Access authentication for
you. The easiest way to do this is using the [digest](https://github.com/mongodb-forks/digest)
library, but you can always use any other library that provides an `http.Client`.
If you have a private and public API token pair, you can use it with the digest library using:
```go
import (
    "context"
    "log"

    "github.com/mongodb-forks/digest"
    "go.mongodb.org/ops-manager/opsmngr"
)

func main() {
    t := digest.NewTransport("your public key", "your private key")
    tc, err := t.Client()
    if err != nil {
        log.Fatalf(err.Error())
    }

    client := opsmngr.NewClient(tc)
    orgs, _, err := client.Organizations.List(context.Background(), nil)
}
```

Note that when using an authenticated Client, all calls made by the client will
include the specified tokens. Therefore, authenticated clients should
almost never be shared between different users.

To get your API Keys please refer to our docs for,
[Ops Manager](https://docs.opsmanager.mongodb.com/current/tutorial/configure-public-api-access/),
or [Cloud Manager](https://docs.cloudmanager.mongodb.com/tutorial/manage-programmatic-api-keys/).

## Roadmap

This library is being initially developed for [mongocli](https://github.com/mongodb/mongocli),
so API methods will likely be implemented in the order that they are
needed by that application.

## Contributing

See our [CONTRIBUTING.md](CONTRIBUTING.md) Guide.

## License

MongoDB Ops Manager Go Client is released under the Apache 2.0 license. See [LICENSE](LICENSE)
