package types

// this file contains structs or other data types that are used for the
// configuration of this microservice

// OAuthScope describes a json-compatible scope file used for authorization
// mechanisms in this microservice
type OAuthScope struct {
	// JSONSchema contains the uri pointing to the schema the scope was written
	// with
	JSONSchema string `json:"$schema"`
	// Name contains a short name/title for the scope
	Name string `json:"name"`
	// Description may contain a longer description of the scope (e.g., who this
	// scope should be used for)
	Description string `json:"description"`
	// Identifier contains the identifier for the scope in authorization data
	Identifier string `json:"scopeStringValue"`
}

// EnvironmentFile describes the layout of the environment configuration file
// which is used to tell the service which environment variables to read and
// which ones are required for the service to work properly
type EnvironmentFile struct {
	// Required contains an array of environment variable keys that are
	// required for the service to work
	Required []string `json:"required"`
	// Optional contains a mapping of environment variable keys to their default
	// values. The default values need to be strings
	Optional map[string]string `json:"optional"`
}
