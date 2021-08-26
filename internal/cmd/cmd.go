// Package cmd
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-22
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/teocci/go-concurrency-samples/internal/config"
	"github.com/teocci/go-concurrency-samples/internal/core"
	"github.com/teocci/go-concurrency-samples/internal/dirfiles"
	"github.com/teocci/go-concurrency-samples/internal/logger"
)

const (
	version = "v1.0"
	commit  = "420"
)

var (
	// shaman provides the shaman cli/server functionality
	app = &cobra.Command{
		Use:               "go-concurrency-samples",
		Short:             "Unzip a file and process the data inside",
		Long:              `This application unzip a file containing logs and csv dirfiles that will be merged and then inserted into a database.`,
		PersistentPreRunE: readConfig,
		PreRunE:           preFlight,
		RunE:              startApp,
		SilenceErrors:     false,
		SilenceUsage:      false,
	}

	filename string
	isSplit  = false
)

// add supported cli commands/flags
func init() {
	app.Flags().StringVarP(&filename, "filename", "f", filename, "Zip file that contains the logs.")
	app.Flags().BoolVarP(&isSplit, "is-split", "s", isSplit, "If the zip file has been split")

	config.AddFlags(app)
}

func readConfig(ccmd *cobra.Command, args []string) error {
	if err := config.LoadConfigFile(); err != nil {
		fmt.Printf("sError: %v\n", err)
		return err
	}

	return nil
}

func preFlight(ccmd *cobra.Command, args []string) error {
	if config.Version {
		fmt.Printf("go-concurrency-samples %s (%s)\n", version, commit)

		return fmt.Errorf("")
	}

	if !config.Verbose {
		ccmd.HelpFunc()(ccmd, args)

		return fmt.Errorf("")
	}

	return nil
}

func startApp(ccmd *cobra.Command, args []string) error {
	var err error

	config.Log, err = logger.New(config.LogLevel, config.Verbose, false, config.File)
	if err != nil {
		return err
	}

	if !dirfiles.Exists(filename) {
		return errors.New(fmt.Sprintf("%s file does not exist", filename))
	}

	// make channel for errors
	errs := make(chan error)

	go func() {
		errs <- core.Start(filename, isSplit)
	}()

	// break if any of them return an error (blocks exit)
	if err := <-errs; err != nil {
		config.Log.Fatal(err.Error())
	}

	return err
}

func Execute() {
	err := app.Execute()
	if err != nil {
		return
	}
}
