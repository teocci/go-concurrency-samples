// Package cmd
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-22
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/teocci/go-concurrency-samples/internal/cmd/cmdapp"
	"github.com/teocci/go-concurrency-samples/internal/config"
	"github.com/teocci/go-concurrency-samples/internal/core"
	"github.com/teocci/go-concurrency-samples/internal/filemgr"
	"github.com/teocci/go-concurrency-samples/internal/logger"
)

var (
	app = &cobra.Command{
		Use:           cmdapp.Name,
		Short:         cmdapp.Short,
		Long:          cmdapp.Long,
		PreRunE:       validate,
		RunE:          runE,
		SilenceErrors: false,
		SilenceUsage:  true,
	}

	filename string
	dest     string
	extract  bool
	merge    bool
)

// Add supported cli commands/flags
func init() {
	dest = cmdapp.DDefault
	extract = cmdapp.EDefault
	merge = cmdapp.MDefault

	cobra.OnInitialize(initConfig)

	app.Flags().StringVarP(&filename, cmdapp.FName, cmdapp.FShort, filename, cmdapp.FDesc)
	app.Flags().StringVarP(&dest, cmdapp.DName, cmdapp.DShort, dest, cmdapp.DDesc)

	app.Flags().BoolVarP(&extract, cmdapp.EName, cmdapp.EShort, extract, cmdapp.EDesc)
	app.Flags().BoolVarP(&merge, cmdapp.MName, cmdapp.MShort, merge, cmdapp.MDesc)

	_ = app.MarkFlagRequired(cmdapp.FName)

	config.AddFlags(app)
}

// Load config
func initConfig() {
	if err := config.LoadConfigFile(); err != nil {
		log.Fatal(err)
	}

	config.LoadLogConfig()
}

func validate(ccmd *cobra.Command, args []string) error {
	if config.Version {
		fmt.Printf(cmdapp.VersionTemplate, cmdapp.Name, cmdapp.Version, cmdapp.Commit)

		return nil
	}

	if !config.Verbose {
		ccmd.HelpFunc()(ccmd, args)

		return fmt.Errorf("")
	}

	return nil
}

func runE(ccmd *cobra.Command, args []string) error {
	var err error
	config.Log, err = logger.New(config.LogConfig)
	if err != nil {
		return ErrCanNotLoadLogger(err)
	}

	if merge || extract {
		if !filemgr.FileExists(filename) {
			return ErrFileDoesNotExist(filename)
		}
	}

	// make channel for errors
	errs := make(chan error)

	go func() {
		mode := core.EMNormal
		if extract {
			mode = core.EMExtract
		}
		if merge {
			mode = core.EMMerge
		}
		errs <- core.Start(filename, dest, mode)
	}()

	// break if any of them return an error (blocks exit)
	if err := <-errs; err != nil {
		config.Log.Fatal(err)
	}

	return err
}

func Execute() {
	err := app.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
