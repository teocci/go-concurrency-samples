// Package core
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-22
package core

import (
	"context"
	"fmt"
	"os"

	"github.com/teocci/go-concurrency-samples/src/logger"
)

// Core is an instance of rtsp-simple-server.
type Core struct {
	ctx         context.Context
	ctxCancel   func()
	confPath    string
	conf        *conf.Conf
	confFound   bool
	stats       *stats
	logger      *logger.Logger
	metrics     *metrics
	pprof       *pprof
	pathManager *pathManager
	rtspServer  *rtspServer
	rtspsServer *rtspServer
	rtmpServer  *rtmpServer
	hlsServer   *hlsServer
	api         *api
	confWatcher *confwatcher.ConfWatcher

	// in
	apiConfigSet chan *conf.Conf

	// out
	done chan struct{}
}

// New allocates a core.
func New(args []string) (*Core, bool) {
	k := kingpin.New("rtsp-simple-server", "rtsp-simple-server "+version+"\n\nRTSP server.")

	argVersion := k.Flag("version", "print version").Bool()
	argConfPath := k.Arg("confpath", "path to a config file. The default is rtsp-simple-server.yml.").Default("rtsp-simple-server.yml").String()

	kingpin.MustParse(k.Parse(args))

	if *argVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// on Linux, try to raise the number of file descriptors that can be opened
	// to allow the maximum possible number of clients
	// do not check for errors
	rlimit.Raise()

	ctx, ctxCancel := context.WithCancel(context.Background())

	p := &Core{
		ctx:          ctx,
		ctxCancel:    ctxCancel,
		confPath:     *argConfPath,
		apiConfigSet: make(chan *conf.Conf),
		done:         make(chan struct{}),
	}

	var err error
	p.conf, p.confFound, err = conf.Load(p.confPath)
	if err != nil {
		fmt.Printf("ERR: %s\n", err)
		return nil, false
	}

	err = p.createResources(true)
	if err != nil {
		p.Log(logger.Info, "ERR: %s", err)
		p.closeResources(nil)
		return nil, false
	}

	if p.confFound {
		p.confWatcher, err = confwatcher.New(p.confPath)
		if err != nil {
			p.Log(logger.Info, "ERR: %s", err)
			p.closeResources(nil)
			return nil, false
		}
	}

	go p.run()

	return p, true
}