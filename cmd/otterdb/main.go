package main

import (
	"os"

	"github.com/bbengfort/otterdb/pkg"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "otterdb"
	app.Usage = "serve and manage an otterdb replica"
	app.Version = pkg.Version()
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		{
			Name:     "serve",
			Usage:    "serve the TRISA Envoy node server configured from the environment",
			Action:   serve,
			Category: "server",
		},
	}

	app.Run(os.Args)
}

//===========================================================================
// Server Commands
//===========================================================================

func serve(c *cli.Context) (err error) {
	// var conf config.Config
	// if conf, err = config.New(); err != nil {
	// 	return cli.Exit(err, 1)
	// }

	// var trisa *node.Node
	// if trisa, err = node.New(conf); err != nil {
	// 	return cli.Exit(err, 1)
	// }

	// if err = trisa.Serve(); err != nil {
	// 	return cli.Exit(err, 1)
	// }
	return nil
}
