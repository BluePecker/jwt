package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/BluePecker/jwt/pkg/term"
	"github.com/BluePecker/jwt/cmd"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	_, _, stdErr := term.StdStreams()
	logrus.SetOutput(stdErr)

	if err := cmd.RootCmd.Execute(); err != nil {
		logrus.Error(err)
	}
}
