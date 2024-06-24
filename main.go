/*
Copyright © 2024 Mark Stella <mark.stella@gammadata.io>
*/
package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	cmd "github.com/mpstella/dcloud/cmd/cli"
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Format the time
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	// Create a buffer to write the formatted log entry
	var b bytes.Buffer

	// Write the formatted log entry
	b.WriteString(fmt.Sprintf("%s %s: [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), "collab", entry.Message))

	return b.Bytes(), nil
}

func main() {
	cmd.Execute()
}
