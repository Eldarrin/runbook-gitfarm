package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("This is a message!")
	log.Error().Msg("This is an error message")
}
