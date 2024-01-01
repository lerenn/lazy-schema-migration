package csm

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

// WrapperBSON is function that will unwrap a BSON object into a type A, run a
// callback that should transform the type A object into type B, and wrap the
// resulting object of type B into a BSON.
func WrapperBSON[A, B any](doc bson.D, fn func(A) (B, error)) (bson.D, error) {
	var a A
	if err := bsonToStruct(doc, &a); err != nil {
		return nil, err
	}

	b, err := fn(a)
	if err != nil {
		return nil, err
	}

	return structToBSON(&b)
}

// MigrationBSON is the function signature for BSON object migration.
type MigrationBSON func(bson.D) (bson.D, error)

// Check migrator interface is respected.
var _ Migrator[any, bson.D] = (*MigratorBSON[any])(nil)

// MigratorBSON is the structure that will contains the migration and apply
// them on BSON to object T type object, and revert.
type MigratorBSON[T any] struct {
	migrations []MigrationBSON
}

// NewMigratorBSON creates a new BSON migrator with given type and migrations.
func NewMigratorBSON[T any](migrations []MigrationBSON) *MigratorBSON[T] {
	return &MigratorBSON[T]{
		migrations: migrations,
	}
}

// LastVersion returns the last version of the migrations.
func (mj *MigratorBSON[T]) LastVersion() int {
	return len(mj.migrations) + 1
}

// Import imports a BSON object into the T type object, following migrations
// set in the migrator.
func (mj *MigratorBSON[T]) Import(data bson.D) (T, error) {
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
	err = bsonToStruct(data, &output)
	return output, err
}

func (mj *MigratorBSON[T]) getVersion(data bson.D) (int, error) {
	var mappedEntry map[string]any
	if err := bsonToStruct(data, &mappedEntry); err != nil {
		return 0, err
	}

	rawVersion, exists := mappedEntry[VersionFieldKey]
	if !exists {
		return 0, fmt.Errorf("%w: %+v", ErrNoVersion, data)
	}

	version, ok := rawVersion.(int32)
	if !ok {
		return 0, fmt.Errorf("%w: %q", ErrInvalidVersionFormat, reflect.TypeOf(rawVersion))
	}

	// TODO: check version is unit

	return int(version), nil
}

func (mj *MigratorBSON[T]) getMigrationsFromVersion(version int) ([]MigrationBSON, error) {
	if version <= 0 && version > mj.LastVersion() {
		return nil, fmt.Errorf("%w: %d", ErrVersionNotFound, version)
	}

	return mj.migrations[version-1:], nil
}

// Export exports a T type object into a BSON object, adding the version number
// to the BSON object.
func (mj *MigratorBSON[T]) Export(entity T) (bson.D, error) {
	b, err := structToBSON(entity)
	if err != nil {
		return b, err
	}

	return append(b, bson.E{
		Key:   VersionFieldKey,
		Value: mj.LastVersion(),
	}), nil
}

func structToBSON(v any) (doc bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func bsonToStruct(doc bson.D, v any) error {
	bsonBytes, err := bson.Marshal(doc)
	if err != nil {
		return err
	}

	return bson.Unmarshal(bsonBytes, v)
}
