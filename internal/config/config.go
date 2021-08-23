// Package config
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-23
package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/teocci/go-concurrency-samples/internal/logger"
	"path/filepath"
)

var (
	LogLevel = "info" // Log level to output [fatal|error|info|debug|trace]
	Verbose  = true   // Run in server mode
	File     = ""     // Configuration file to load
	TempDir  = "./tmp/output-folder" // Temporal directory
	Version  = false // Print version info and exit

	Log *logger.Logger // Central logger for the app
)

// AddFlags adds the available cli flags
func AddFlags(cmd *cobra.Command) {
	// core
	cmd.Flags().StringVarP(&LogLevel, "log-level", "l", LogLevel, "Log level to output [fatal|error|info|debug|trace]")
	cmd.Flags().BoolVarP(&Verbose, "verbose", "v", Verbose, "Run in debug mode")
	cmd.PersistentFlags().StringVarP(&File, "config-file", "c", File, "Configuration file to load")
	cmd.PersistentFlags().StringVarP(&TempDir, "temp-dir", "t", TempDir, "Temporal directory where the app will work")

	cmd.Flags().BoolVarP(&Version, "version", "V", Version, "Print version info and exit")
}

// LoadConfigFile reads the specified config file
func LoadConfigFile() error {
	if File == "" {
		return nil
	}

	// Set defaults to whatever might be there already
	viper.SetDefault("log-level", LogLevel)
	viper.SetDefault("verbose", Verbose)
	viper.SetDefault("temp-dir", TempDir)

	filename := filepath.Base(File)
	viper.SetConfigName(filename[:len(filename)-len(filepath.Ext(filename))])
	viper.AddConfigPath(filepath.Dir(File))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config file - %v", err)
	}

	// Set values. Config file will override commandline
	LogLevel = viper.GetString("log-level")
	Verbose = viper.GetBool("verbose")
	TempDir = viper.GetString("temp-dir")

	return nil
}
