package csm

import (
	"errors"
	"fmt"
)

var (
	// ErrGeneric is the generic error that can be used to isolate
	// continuous-schema-migration errors.
	ErrGeneric = errors.New("continuous schema migration error")
	// ErrNoVersion happens when no version is detected.
	ErrNoVersion = fmt.Errorf("%w: no version detected", ErrGeneric)
	// ErrInvalidVersionFormat happens when the version is not an integer.
	ErrInvalidVersionFormat = fmt.Errorf("%w: version is not an integer", ErrGeneric)
	// ErrVersionNotFound happens when the given version is not found in migrations.
	ErrVersionNotFound = fmt.Errorf("%w: version not found in migrations", ErrGeneric)
	// ErrRunningMigration happens when a migration fails.
	ErrRunningMigration = fmt.Errorf("%w: running migration failed", ErrGeneric)
)
