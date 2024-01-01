package csm

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// WrapperJSON is function that will unwrap a JSON object into a type A, run a
// callback that should transform the type A object into type B, and wrap the
// resulting object of type B into a JSON.
func WrapperJSON[A, B any](data []byte, fn func(A) (B, error)) ([]byte, error) {
	var a A
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, err
	}

	b, err := fn(a)
	if err != nil {
		return nil, err
	}

	return json.Marshal(b)
}

// MigrationJSON is the function signature for JSON object migration.
type MigrationJSON func([]byte) ([]byte, error)

// Check migrator interface is respected.
var _ Migrator[any, []byte] = (*MigratorJSON[any])(nil)

// MigratorJSON is the structure that will contains the migration and apply
// them on JSON to object T type object, and revert.
type MigratorJSON[T any] struct {
	migrations []MigrationJSON
}

// NewMigratorJSON creates a new JSON migrator with given type and migrations.
func NewMigratorJSON[T any](migrations []MigrationJSON) *MigratorJSON[T] {
	return &MigratorJSON[T]{
		migrations: migrations,
	}
}

// LastVersion returns the last version of the migrations.
func (mj *MigratorJSON[T]) LastVersion() int {
	return len(mj.migrations) + 1
}

// Import imports a JSON object into the T type object, following migrations
// set in the migrator.
func (mj *MigratorJSON[T]) Import(data []byte) (T, error) {
	version, err := mj.getVersion(data)
	if err != nil {
		return *new(T), err
	}

	migrations, err := mj.getMigrationsFromVersion(version)
	if err != nil {
		return *new(T), err
	}

	for i, m := range migrations {
		data, err = m(data)
		if err != nil {
			return *new(T), fmt.Errorf("%w: %d", ErrRunningMigration, i)
		}
	}

	var output T
	err = json.Unmarshal(data, &output)
	return output, err
}

func (mj *MigratorJSON[T]) getVersion(data []byte) (int, error) {
	var mappedEntry map[string]any
	if err := json.Unmarshal(data, &mappedEntry); err != nil {
		return 0, err
	}

	rawVersion, exists := mappedEntry[VersionFieldKey]
	if !exists {
		return 0, fmt.Errorf("%w: %s", ErrNoVersion, string(data))
	}

	version, ok := rawVersion.(float64)
	if !ok {
		return 0, fmt.Errorf("%w: %q", ErrInvalidVersionFormat, reflect.TypeOf(rawVersion))
	}

	// TODO: check version is unit

	return int(version), nil
}

func (mj *MigratorJSON[T]) getMigrationsFromVersion(version int) ([]MigrationJSON, error) {
	if version <= 0 && version > mj.LastVersion() {
		return nil, fmt.Errorf("%w: %d", ErrVersionNotFound, version)
	}

	return mj.migrations[version-1:], nil
}

// Export exports a T type object into a JSON object, adding the version number
// to the JSON object.
func (mj *MigratorJSON[T]) Export(entity T) ([]byte, error) {
	b, err := json.Marshal(entity)
	if err != nil {
		return b, err
	}

	str := string(b)
	str = fmt.Sprintf("{%s,%q:%d}", b[1:len(str)-1], VersionFieldKey, mj.LastVersion())

	return []byte(str), nil
}
