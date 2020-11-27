package snapshot

import "os"

// DefaultTransform option
var DefaultTransform Transform

// DefaultMarshal option
var DefaultMarshal Marshal

// DefaultAssertEqual option
var DefaultAssertEqual AssertEqual

// DefaultUpdate option
var DefaultUpdate bool

// ResetDefaults to initial values.
func ResetDefaults() {
	DefaultTransform = TransformJSON
	DefaultMarshal = MarshalTextOrJSON
	DefaultAssertEqual = AssertEqualBytes
	DefaultUpdate = os.Getenv("SNAPSHOT_UPDATE") == "true"
}

func init() {
	ResetDefaults()
}
