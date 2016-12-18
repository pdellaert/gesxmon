package main

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/pdellaert/gesxmon/vsphere"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gesxmon"
	app.Usage = "Monitors ESXi events and takes appropriate actions in the Nuage VSD."
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "listen",
			Aliases: []string{"l"},
			Usage:   "Listens to events on ESXi",
			Action:  listenAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "vsphere-url, u",
					Value:  "https://username:password@hostname/sdk",
					Usage:  "The vSphere `URL` for govmomi to connect to a vSphere host.",
					EnvVar: "GESXMON_VSPHERE_URL",
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "/etc/gesxmon/gesxmon.cfg",
			Usage:  "A config `FILE` used to get other values.",
			EnvVar: "GESXMON_CONFIG",
		},
		cli.BoolFlag{
			Name:   "debug, D",
			Usage:  "Enable debug mode for debug output.",
			EnvVar: "GESXMON_DEBUG",
		},
		cli.BoolFlag{
			Name:   "verbose, V",
			Usage:  "Enable verbose mode for more verbose output.",
			EnvVar: "GESXMON_VERBOSE",
		},
	}

	app.Run(os.Args)
}

func listenAction(c *cli.Context) error {
	// Setting up logging
	logger := logrus.New()

	// Setting log level
	if c.GlobalBool("debug") {
		logger.Level = logrus.DebugLevel
	} else if c.GlobalBool("verbose") {
		logger.Level = logrus.InfoLevel
	} else {
		logger.Level = logrus.WarnLevel
	}

	// Starting logging
	logger.WithFields(logrus.Fields{"Command": os.Args}).Debug("Command executed")
	logger.WithFields(logrus.Fields{"vsphere-url": c.String("vsphere-url")}).Debug("Listen command arguments")
	logger.Info("gesxmon listener starting")

	// Creating context
	ctx, cancel := context.WithCancel(context.Background())

	// Setting up SIGTERM handling
	logger.Debug("Setting up SIGTERM handling")
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		logger.Debug("Received SIGTERM, ending process")
		logger.Debug("Canceling context")
		cancel()
		logger.Info("Exiting")
		os.Exit(0)
	}()

	if !c.IsSet("vsphere-url") && !c.GlobalIsSet("config") {
		logger.Fatal("Missing options, please specify a vsphere-url or global config file, either through CLI arguments or ENV variables")
		return cli.NewExitError("Missing options, please specify a vsphere-url or global config file, either through CLI arguments or ENV variables", 1)
	}

	url, err := url.Parse(c.String("vsphere-url"))
	if err != nil {
		logger.Error("Unable to parse the vsphere-url")
		logger.Debug(err)
		return cli.NewExitError("Unable to parse the vsphere-url", 10)
	}

	esxiEventListener := vsphere.NewESXiEventListener(url, true, logger)

	err = esxiEventListener.Run(ctx)
	if err != nil {
		logger.Error("Error while running the event listener")
		logger.Debug(err)
		return cli.NewExitError("Error while running the event listener", 20)
	}

	return nil
}
