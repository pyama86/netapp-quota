package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_urlFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./netapp-quota -url", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_userFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./netapp-quota -user", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_passwordFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./netapp-quota -password", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_prefixFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./netapp-quota -prefix", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_svmFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./netapp-quota -svm", " ")

	status := cli.Run(args)
	_ = status
}
