package v1alpha1

type Match struct {
	// Any allows specifying resources which will be ORed.
	Any []Any `json:"any,omitempty"`

	// All allows specifying resources which will be ANDed.
	All []Any `json:"all,omitempty"`
}
