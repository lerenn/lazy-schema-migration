package csm

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestMigratorSuite(t *testing.T) {
	suite.Run(t, new(MigratorSuite))
}

type MigratorSuite struct {
	suite.Suite
}

type EntryV1 struct {
	FirstName, LastName string
}

func entryV1ToV2(data []byte) ([]byte, error) {
	return WrapperJSON(data, func(v1 EntryV1) (EntryV2, error) {
		return EntryV2{
			FullName: v1.FirstName + " " + v1.LastName,
		}, nil
	})
}

type EntryV2 struct {
	FullName string
}

func entryV2ToV3(data []byte) ([]byte, error) {
	return WrapperJSON(data, func(v2 EntryV2) (EntryV3, error) {
		return EntryV3{
			FullName: v2.FullName,
		}, nil
	})
}

type EntryV3 struct {
	FullName string
	Age      uint
}

var (
	testMigrations = []MigrationJSON{
		entryV1ToV2,
		entryV2ToV3,
	}
)

func (suite *MigratorSuite) TestFromJSON() {
	input := struct {
		FirstName, LastName string
		Version             int `json:"__schema_version"`
	}{
		FirstName: "first",
		LastName:  "last",
		Version:   1,
	}
	inputByte, err := json.Marshal(input)
	suite.Require().NoError(err)

	expectedOutput := EntryV3{
		FullName: "first last",
	}

	mig := NewMigratorJSON[EntryV3](testMigrations)
	entry, err := mig.Import(inputByte)
	suite.Require().NoError(err)
	suite.Require().IsType(expectedOutput, entry)
	suite.Require().Equal(expectedOutput, entry)
}

func (suite *MigratorSuite) TestToJSON() {
	input := EntryV3{
		FullName: "first last",
		Age:      34,
	}
	expectedOutput := []byte(
		fmt.Sprintf("{%q:%q,%q:%d,%q:%d}", "FullName", "first last", "Age", 34, VersionFieldKey, 3),
	)

	mig := NewMigratorJSON[EntryV3](testMigrations)
	output, err := mig.Export(input)
	suite.Require().NoError(err)
	suite.Require().Equal(string(expectedOutput), string(output))
}
