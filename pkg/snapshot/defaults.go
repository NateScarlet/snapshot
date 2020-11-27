package snapshot

// DefaultTransform when snapshot transform option not set
var DefaultTransform Transform

// DefaultMarshal when snapshot marshal option not set
var DefaultMarshal Marshal

// DefaultAssertEqual when snapshot assertEqual option not set
var DefaultAssertEqual AssertEqual

// ResetDefaults to initial values.
func ResetDefaults() {
	DefaultTransform = TransformJSON
	DefaultMarshal = MarshalTextOrJSON
	DefaultAssertEqual = AssertEqualBytes
}

func init() {
	ResetDefaults()
}
