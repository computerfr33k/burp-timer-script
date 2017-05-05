# [Burp](https://github.com/grke/burp) Timer Script (Go)

[![Build Status](https://travis-ci.org/computerfr33k/burp-timer-script-go.svg?branch=master)](https://travis-ci.org/computerfr33k/burp-timer-script-go)

# About

This is a port of the default [timer script](https://github.com/grke/burp/blob/master/configs/server/timer_script) written in go for cross platform portability without being dependent on system commands such as gnu date or bash, which is not installed on FreeBSD by default.

# Getting Started

This program is a complete drop-in replacement for the default bash script included with burp installation.
In order to use it, just download and copy it to your existing burp scripts folder (default: `/usr/share/burp/scripts/`), then change the `timer_script` line in your server's config file to point to the program.

For example `timer_script = /usr/share/burp/scripts/timer_script-linux-amd64`
