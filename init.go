package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"external-api-service/globals"
	"external-api-service/types"
)

var l zerolog.Logger

// this init functions sets up the logger which is used for this microservice
func init() {
	// set the time format to unix timestamps to allow easier machine handling
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// allow the logger to create an error stack for the logs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// now use the environment variable `LOG_LEVEL` to determine the logging
	// level for the microservice.
	rawLoggingLevel, isSet := os.LookupEnv("LOG_LEVEL")

	// if the value is not set, use the info level as default.
	var loggingLevel zerolog.Level
	if !isSet {
		loggingLevel = zerolog.InfoLevel
	} else {
		// now try to parse the value of the raw logging level to a logging
		// level for the zerolog package
		var err error
		loggingLevel, err = zerolog.ParseLevel(rawLoggingLevel)
		if err != nil {
			// since an error occurred while parsing the logging level, use info
			loggingLevel = zerolog.InfoLevel
			log.Warn().Msg("unable to parse value from environment. using info")
		}
	}
	// since now a logging level is set, configure the logger
	zerolog.SetGlobalLevel(loggingLevel)
	l = log.With().Str("step", "init").Logger()
}

// this function initializes the environment variables used in this microservice
// and validates that the configured variables are present.
func init() {
	l.Info().Msg("loading environment for microservice")

	// now check if the default location for the environment configuration
	// was changed via the `ENV_CONFIG_LOCATION` variable
	location, locationChanged := os.LookupEnv("ENV_CONFIG_LOCATION")
	if !locationChanged {
		// since the location has not changed, set the default value
		location = "./environment.json"
		l.Debug().Msg("location for environment config not changed")
	}
	l.Info().Str("path", location).Msg("loading environment configuration file")
	// open the file
	file, err := os.Open(location)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to open environment configuration file")
	}

	// now parse the configuration file
	var config types.EnvironmentFile
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to parse environment configuration file")
	}
	l.Info().Msg("successfully loaded environment configuration")

	// since the configuration was successfully loaded, check the required
	// environment variables
	l.Info().Msg("validating configuration against current environment")
	l.Info().Msg("loading required environment variables")
	for _, envKey := range config.Required {
		// overwriting the logger in this loop to make the calls shorter
		l := l.With().Str("envKey", envKey).Logger()
		l.Debug().Msg("checking required environment variable")
		value, envSet := os.LookupEnv(envKey)
		if !envSet {
			fileKey := fmt.Sprintf("%s_FILE", envKey)
			l := l.With().Str("fileKey", fileKey).Logger()
			l.Warn().Msg("not found in environment. checking for docker secret")
			filePath, filePathSet := os.LookupEnv(fileKey)
			if !filePathSet {
				// since neither the environment variable was set, nor a
				// docker secret is mapped to the environment variable, stop
				// the microservice
				l.Fatal().Msg("environment variable neither set nor supplied as docker secret")
			}
			// now check if the filePath is empty
			filePath = strings.TrimSpace(filePath)
			if filePath == "" {
				// since the file path is empty, the environment variable cannot
				// be read. this also stops the microservice
				l.Fatal().Msg("docker secret setup found, but filepath is empty")
			}
			// now start loading the docker secret contents
			fileContentBytes, err := os.ReadFile(filePath)
			if err != nil {
				l.Fatal().Err(err).Msg("failed to read discovered docker secret")
			}
			// now check if the number of read bytes exceeds 0
			if len(fileContentBytes) == 0 {
				l.Fatal().Msg("configured docker secret is empty")
			}
			// now convert the file contents to a string
			fileContents := string(fileContentBytes)
			// now trim the contents
			fileContents = strings.TrimSpace(fileContents)
			// now set the file contents as the value for the environment
			// variable under the envKey
			globals.Environment[envKey] = fileContents
			l.Debug().Msg("loaded environment variable from docker secret")
			continue
		}

		// trim away any excess whitespaces from the value
		value = strings.TrimSpace(value)
		if value == "" {
			l.Fatal().Msg("required environment variable empty")
		}
		// set the value to the environment
		globals.Environment[envKey] = value
		l.Debug().Msg("loaded environment variable from real environment")
	}

	// now read the optional variables
	l.Info().Msg("checking optional environment variables")
	for envKey, defaultValue := range config.Optional {
		l := l.With().Str("envKey", envKey).Str("default", defaultValue).Logger()
		l.Debug().Msg("loading optional environment variable")
		availableValue, isSet := os.LookupEnv(envKey)
		if !isSet {
			l.Debug().Msg("using default value")
			globals.Environment[envKey] = defaultValue
			continue
		}
		// trim the available value
		availableValue = strings.TrimSpace(availableValue)
		l.Info().Str("suppliedValue", availableValue).Msg("using supplied value")
		globals.Environment[envKey] = availableValue
	}
	l.Info().Msg("finished processing the environment configuration")
}
