package main

import (
	"log"
	"os"

	logger "github.com/yourpov/logrite"

	"spotify-cli/app"
	"spotify-cli/internal/logx"
)

func main() {
	logger.SetConfig(logger.Config{
		ShowIcons:    true,
		UppercaseTag: true,
		UseColors:    true,
	})
	log.SetFlags(0)
	if err := app.Run(os.Args[1:]); err != nil {
		_ = logx.Err("app failed: %v", err)
		os.Exit(1)
	}
}
