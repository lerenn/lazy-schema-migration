package main

import (
	"fmt"
	"strings"

	csm "github.com/lerenn/continuous-schema-migration"
)

type EntityV1 struct {
	FullName string
}

type EntityV2 struct {
	FirstNames string
	LastName   string
}

type EntityV3 struct {
	FirstNames string
	LastName   string
	Age        int
}

var (
	migrations = []csm.MigrationJSON{
		// From Person V1 to V2
		func(data []byte) ([]byte, error) {
			// Use WrapperJSON to get V1 and V2 forms
			return csm.WrapperJSON(data, func(v1 EntityV1) (EntityV2, error) {
				// Split the fullname between last name and first names
				names := strings.Split(v1.FullName, " ")
				firstNames, lastName := "", ""
				if len(names) >= 2 {
					firstNames = strings.Join(names[:len(names)-1], " ")
					lastName = names[len(names)-1]
				}

				// Set it in the new version
				return EntityV2{
					FirstNames: firstNames,
					LastName:   lastName,
				}, nil
			})
		},
		// From Person V2 to V3
		func(data []byte) ([]byte, error) {
			// Use WrapperJSON to get V2 and V3 forms
			return csm.WrapperJSON(data, func(v2 EntityV2) (EntityV3, error) {
				// Set the new version with no age, as we have no mean to know it
				return EntityV3{
					FirstNames: v2.FirstNames,
					LastName:   v2.LastName,
				}, nil
			})
		},
	}
)

func CreateNewObject() {
	// Create a new migrator
	mig := csm.NewMigratorJSON[EntityV3](migrations)

	// Create a new object
	newEntry := EntityV3{
		FirstNames: "John Robert",
		LastName:   "Reddington",
	}
	fmt.Printf("Creating a new object: %+v\n", newEntry)

	// Export this new object
	data, err := mig.Export(newEntry)
	if err != nil {
		panic(err)
	}
	fmt.Printf(" + Created data: %s\n\n", string(data)) // Ready to be inserted in DB
}

func UseOldObject() {
	// Create a new migrator
	mig := csm.NewMigratorJSON[EntityV3](migrations)

	// Importing an old object, do some modifications and save it as last version
	data := []byte(`{"FullName":"John Robert Reddington","__schema_version":1}`)
	fmt.Printf("Reading an old data from DB: %s\n", string(data))

	migratedEntity, err := mig.Import(data)
	if err != nil {
		panic(migratedEntity)
	}
	fmt.Printf(" + Migrated data: %+v\n", migratedEntity)

	migratedEntity.Age = 45
	fmt.Printf(" + Modified data with age: %+v\n", migratedEntity)

	data, err = mig.Export(migratedEntity)
	if err != nil {
		panic(err)
	}
	fmt.Printf(" + Updated data: %s\n", data) // Ready to be updated in DB
}

func main() {
	CreateNewObject()
	UseOldObject()
}
