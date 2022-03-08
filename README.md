# speedtest-lib

Speed test lib provides mechanism for network bandwidth check.

### Supported provider

* Ookla's https://www.speedtest.net/
* Netflix's https://fast.com/

### How to use

Install library:

```shell
go get github.com/goforbroke1006/speedtest-lib

```

Below sample on small application:

```go
package main

import (
	"context"
	"fmt"
	"time"

	speedtest_lib "github.com/goforbroke1006/speedtest-lib"
	"github.com/goforbroke1006/speedtest-lib/pkg/content"
)

func main() {
	start := time.Now()

	var (
		download float64
		upload   float64
		err      error
	)

	w := speedtest_lib.New()
	download, upload, err = w.DoRequest(context.Background(), speedtest_lib.ProviderKindNetflix)
	if err != nil {
		panic(err)
	}
	fmt.Println("Netflix",
		"D", content.DataLen(download*content.Bit).MegaBites(),
		"U", content.DataLen(upload*content.Bit).MegaBites(),
		"Spend:", time.Since(start).Seconds())
}

```

### How to run benchmark tests

Run in terminal:

```shell
make setup
make benchmark
```

### How to run basic example

Run in terminal:

```shell
make setup
make dep gen
go run examples/basic/main.go

```
