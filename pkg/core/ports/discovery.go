package ports

// Discovery interfaces are intended to be used when there's multiple methods of converting an email to a service ID,
// and back.
type Discovery interface {
	GetUsernameFromEmail(...string) ([]string, error)
	GetEmailFromUsername(...string) ([]string, error)
}
