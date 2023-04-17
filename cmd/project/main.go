package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/vstdy/otus-highload/cmd/project/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Application crashed")
		os.Exit(1)
	}
}
