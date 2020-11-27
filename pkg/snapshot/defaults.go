package snapshot

// DefaultTransform when snapshot transform option not set
var DefaultTransform Transform

// DefaultMarshal when snapshot marshal option not set
var DefaultMarshal Marshal

// ResetDefaults to initial values.
func ResetDefaults() {
	DefaultTransform = TransformJSON
	DefaultMarshal = MarshalTextOrJSON
}

func init() {
	ResetDefaults()
}
