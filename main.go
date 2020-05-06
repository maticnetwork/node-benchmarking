// package main

// import (
// 	// "github.com/maticnetwork/monitoring-tools/scripts"
// 	"github.com/maticnetwork/monitoring-tools/benchmarking"
// )

// func main() {
// 	// scripts.SignerCount()
// 	// scripts.Deposits()
// 	// scripts.RapidFire()
// 	benchmarking.RapidFire5()
// }

package main

import (
	// "fmt"
	"log"
	"os"
	"time"

	"github.com/maticnetwork/monitoring-tools/benchmarking"
	"github.com/urfave/cli/v2"
)

func main() {
  app := &cli.App{}
  app.UseShortOptionHandling = true
  app.Commands = []*cli.Command{
    {
      Name:  "txcount",
      Usage: "Subscribe to blocks and print tx count",
      Action: func(c *cli.Context) error {
				benchmarking.SubscribeBlocks()
        return nil
      },
		},
		{
      Name:  "fire",
      Usage: "fire txs",
      Flags: []cli.Flag{
        &cli.IntFlag{
					Name: "txs",
					Required: true,
          Usage:   "number of txs to fire",
          Value: 500,
				},
				&cli.IntFlag{
					Name: "clients",
					Required: true,
          Usage:   "Number of nodes to connect to",
          Value: 1,
        },
        &cli.Int64Flag{
					Name: "seed",
          Usage:   "seed to generate a random private key",
          Value: time.Now().Unix(),
        },
        &cli.IntFlag{
					Name: "delay",
          Usage:   "seed to generate a random private key",
          Value: 0,
        },
			},
      Action: func(c *cli.Context) error {
				benchmarking.RapidFire(
          c.Int("txs"),
          c.Int("clients"),
          c.Int64("seed") + time.Now().Unix(),
          c.Int("delay"),
        )
        return nil
      },
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
