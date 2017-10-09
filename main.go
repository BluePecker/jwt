package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/BluePecker/jwt/pkg/term"
)

func main()  {
	stdIn, stdOut, stdErr := term.StdStreams()
	logrus.SetOutput(stdErr)
}