package csm

// Migrator is the interface that every migrator should satisfy to respect
// libraries engagements.
// Note: E is the entity type and D is the Data type.
type Migrator[E, D any] interface {
	// Import imports a object from D (data) to E (entity) type object,
	// following migrations set in the migrator.
	Import(data D) (E, error)

	// Export exports a E (entity) type object into a D (data) object,
	// adding the version number to the resulting data.
	Export(entity E) (D, error)

	// LastVersion returns the last version of the migrations
	LastVersion() int
}
