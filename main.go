package main

import (
	"log"
	"os"
	"time"

	"github.com/maticnetwork/monitoring-tools/benchmarking"
	"github.com/maticnetwork/monitoring-tools/scripts"
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
					Name:     "txs",
					Required: true,
					Usage:    "number of txs to fire",
					Value:    500,
				},
				&cli.IntFlag{
					Name:     "clients",
					Required: true,
					Usage:    "Number of nodes to connect to",
					Value:    1,
				},
				&cli.Int64Flag{
					Name:  "seed",
					Usage: "seed to generate a random private key",
					Value: 0,
				},
				&cli.IntFlag{
					Name:  "delay",
					Usage: "seed to generate a random private key",
					Value: 0,
				},
			},
			Action: func(c *cli.Context) error {
				benchmarking.RapidFire(
					c.Int("txs"),
					c.Int("clients"),
					c.Int64("seed")+time.Now().Unix(),
					c.Int("delay"),
				)
				return nil
			},
		},
		{
			Name:  "roothash",
			Usage: "Merkle root of the block headers",
			Flags: []cli.Flag{
				&cli.Uint64Flag{
					Name:     "start",
					Required: true,
					Usage:    "start block",
				},
				&cli.Uint64Flag{
					Name:     "end",
					Required: true,
					Usage:    "end block",
				},
			},
			Action: func(c *cli.Context) error {
				scripts.RootHash(c.Uint64("start"), c.Uint64("end"))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
