# Lock manager

A simple package for running functions using key-based distributed locks backed by Redis.

## Installation

```bash
go get github.com/kajjagtenberg/lockmanager
```

## Usage

Use the following example to get started.

```golang
package main

import (
	"context"
	"log"
	"time"

	"github.com/kajjagtenberg/lockmanager"
	"github.com/go-redis/redis/v8"
)

func main() {
	rds := redis.NewClient(&redis.Options{})

	manager, err := lockmanager.NewLockManager(rds)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	if err := manager.Lock(ctx, "test", func() error {
		log.Println("Running under lock")

		time.Sleep(time.Second * 20)

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	if err := manager.Lock(ctx, "test", func() error {
		log.Println("Running under lock")

		time.Sleep(time.Second * 20)

		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
```

## Authors

- [@kajjagtenberg](https://www.github.com/kajjagtenberg)

## License

[MIT](https://choosealicense.com/licenses/mit/)
