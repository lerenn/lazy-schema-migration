package csm

type EntryV1 struct {
	FirstName, LastName string
}

type EntryV2 struct {
	FullName string
}

type EntryV3 struct {
	FullName string
	Age      uint
}
