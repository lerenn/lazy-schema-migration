package csm

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMigratorBSONSuite(t *testing.T) {
	suite.Run(t, new(MigratorBSONSuite))
}

type MigratorBSONSuite struct {
	suite.Suite
}

func bsonEntryV1ToV2(data bson.D) (bson.D, error) {
	return WrapperBSON(data, func(v1 EntryV1) (EntryV2, error) {
		return EntryV2{
			FullName: v1.FirstName + " " + v1.LastName,
		}, nil
	})
}

func bsonEntryV2ToV3(data bson.D) (bson.D, error) {
	return WrapperBSON(data, func(v2 EntryV2) (EntryV3, error) {
		return EntryV3{
			FullName: v2.FullName,
		}, nil
	})
}

var (
	testBSONMigrations = []MigrationBSON{
		bsonEntryV1ToV2,
		bsonEntryV2ToV3,
	}
)

func (suite *MigratorBSONSuite) TestFromBSON() {
	input := struct {
		FirstName, LastName string
		Version             int `bson:"__schema_version"`
	}{
		FirstName: "first",
		LastName:  "last",
		Version:   1,
	}
	inputByte, err := structToBSON(input)
	suite.Require().NoError(err)

	expectedOutput := EntryV3{
		FullName: "first last",
	}

	mig := NewMigratorBSON[EntryV3](testBSONMigrations)
	entry, err := mig.Import(inputByte)
	suite.Require().NoError(err)
	suite.Require().IsType(expectedOutput, entry)
	suite.Require().Equal(expectedOutput, entry)
}

func (suite *MigratorBSONSuite) TestToBSON() {
	input := EntryV3{
		FullName: "first last",
		Age:      34,
	}

	expectedOutput := bson.D{
		{Key: "fullname", Value: "first last"},
		{Key: "age", Value: int64(34)},
		{Key: VersionFieldKey, Value: int(3)},
	}

	mig := NewMigratorBSON[EntryV3](testBSONMigrations)
	output, err := mig.Export(input)
	suite.Require().NoError(err)
	suite.Require().Equal(expectedOutput, output)
}
