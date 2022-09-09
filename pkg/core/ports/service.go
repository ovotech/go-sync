package ports

// Service interfaces are used to allow Sync to communicate with third party services.
type Service interface {
	Get() ([]string, error)                                          // Get users in a service.
	Add(...string) (success []string, failure []error, err error)    // Add users to a service.
	Remove(...string) (success []string, failure []error, err error) // Remove users from a service.
}
