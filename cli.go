package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/buckhx/gofence/geofence"
	"github.com/codegangsta/cli"
)

var version string

func client(args []string) {
	app := cli.NewApp()
	app.Name = "fence"
	app.Usage = "Fence geojson point features"
	app.ArgsUsage = "path/to/geojson/dir"
	app.Version = version // set with go tool link -X
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "fence",
			Value: "rtree",
			Usage: "Type of fence to use " + strings.Join(geofence.FenceLabels, "|"),
		},
		cli.IntFlag{
			Name:  "zoom, z",
			Value: 18,
			Usage: "Some fences require a zoom level",
		},
		cli.StringFlag{
			Name:  "port, p",
			Value: "8080",
			Usage: "Port to bind to",
		},
	}
	app.Action = func(c *cli.Context) {
		args := c.Args()
		if len(args) < 1 || args[0] == "" {
			die(c, "fences_path required")
		}
		z := c.Int("zoom")
		if z < 0 || z > 23 {
			die(c, "required 0 <= -z <= 23")
		}
		path := args[0]
		label := c.String("fence")
		fences, err := geofence.LoadFenceIndex(path, label, z)
		if err != nil {
			die(c, err.Error())
		}
		port := fmt.Sprintf(":%s", c.String("port"))
		err = geofence.ListenAndServe(port, fences)
		die(c, err.Error())
	}
	app.Run(args)
}

func main() {
	client(os.Args)
}

func die(c *cli.Context, msg string) {
	cli.ShowAppHelp(c)
	fmt.Println(msg)
	os.Exit(1)
}
