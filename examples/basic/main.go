package main

import (
	"context"
	"fmt"
	"time"

	speedtest_lib "github.com/goforbroke1006/speedtest-lib"
	"github.com/goforbroke1006/speedtest-lib/pkg/content_len"
)

func main() {
	start := time.Now()

	var (
		download float64
		upload   float64
		err      error
	)

	w := speedtest_lib.New()

	download, upload, err = w.DoRequest(context.Background(), speedtest_lib.ProviderKindOokla)
	if err != nil {
		panic(err)
	}
	fmt.Println("Ookla  ",
		"D", content_len.DataLen(download).MegaBites(),
		"U", content_len.DataLen(upload).MegaBites(),
		"Spend:", time.Since(start).Seconds())

	download, upload, err = w.DoRequest(context.Background(), speedtest_lib.ProviderKindNetflix)
	if err != nil {
		panic(err)
	}
	fmt.Println("Netflix",
		"D", content_len.DataLen(download).MegaBites(),
		"U", content_len.DataLen(upload).MegaBites(),
		"Spend:", time.Since(start).Seconds())
}
