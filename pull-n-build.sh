#!/bin/bash
# Git pull
git pull

# Build the main
go build main.go

# Rename main as a proclogs
mv -v main proclogs
#cp -v proctel /home/rtt/apps/proctel
