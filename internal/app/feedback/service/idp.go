package service

// IdentityProvider specifies an API for generating unique identifiers.
type NanoIdentityProvider interface {
	// ID generates the unique identifier.
	ID() (string, error)
}
