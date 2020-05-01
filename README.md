# MongoDB Ops Manager Go Client
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/go.mongodb.org/ops-manager/opsmngr)
[![Build Status](https://travis-ci.org/mongodb/go-client-mongodb-ops-manager.svg?branch=master)](https://travis-ci.org/mongodb/go-client-mongodb-ops-manager)

An HTTP client for [Ops Manager](https://docs.opsmanager.mongodb.com/master/reference/api/) 
and [Cloud Manager](https://docs.cloudmanager.mongodb.com/reference/api/) Public API endpoints.

You can view the Official API docs at: 
- https://docs.opsmanager.mongodb.com/master/reference/api/
- https://docs.cloudmanager.mongodb.com/master/reference/api/

## Installation

To get the latest version run this command:

```bash
go get go.mongodb.org/ops-manager
```

## Usage

```go
import "go.mongodb.org/ops-manager/opsmngr"
```

## Authentication 

The Ops Manager API uses [HTTP Digest Authentication](https://docs.opsmanager.mongodb.com/master/core/api/#authentication). 
Provide your PUBLIC_KEY as the username and PRIVATE_KEY as the password as part of the HTTP request. 
See [how to set up public API access](https://docs.opsmanager.mongodb.com/master/tutorial/configure-public-api-access/) for more information.

We use the following library to get HTTP Digest Auth:

- https://github.com/Sectorbob/mlab-ns2/gae/ns/digest

## Example Usage

```go 
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Sectorbob/mlab-ns2/gae/ns/digest"
	"go.mongodb.org/ops-manager/opsmngr"
)

func newClient(publicKey, privateKey string) (*mongodbatlas.Client, error) {

	// Setup a transport to handle digest
	transport := digest.NewTransport(publicKey, privateKey)

	// Initialize the client
	client, err := transport.Client()
	if err != nil {
		return nil, err
	}

	//Initialize the MongoDB API Client.
	return opsmngr.NewClient(client), nil
}

func main() {
	publicKey := os.Getenv("MONGODB_PUBLIC_KEY")
	privateKey := os.Getenv("MONGODB_PRIVATE_KEY")
	projectID := os.Getenv("MONGODB_PROJECT_ID")

	if publicKey == "" || privateKey == "" || projectID == "" {
		log.Fatalln("MONGODB_PUBLIC_KEY, MONGODB_PRIVATE_KEY and MONGODB_PROJECT_ID must be set to run this example")
	}

	client, err := newClient(publicKey, privateKey)
	if err != nil {
		log.Fatalf(err.Error())
	}

	atmStatus, _, err := client.AutomationStatus.Get(context.Background(), projectID)

	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%+v \n", atmStatus)

}
```

## Roadmap

This library is being initially developed for [mongocli](https://github.com/mongodb/mongocli),
so API methods will likely be implemented in the order that they are
needed by that application.

## Contributing

See our [CONTRIBUTING.md](CONTRIBUTING.md) Guide.

## License

MongoDB Ops Manager Go Client is released under the Apache 2.0 license. See [LICENSE](LICENSE)
