package main

import (
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	if os.Getenv("PRODUCTION") != "" {
		log.Info().Msg("I'm in production")
	}
	log.Info().Msg("This is a message!")
	log.Error().Msg("This is an error message")
}
