package cloudevents

import (
	"github.com/duglin/tools/jsonext"
)

// CloudEvent implements the v1.0 version of the CloudEvent micro-spec
type CloudEvent struct {

	// EventID holds the unique ID of this occurrence.
	EventID string `json:"eventId"`

	// Extensions hold additional features not specified by the spec.
	Extensions map[string]interface{} `json:",inline,exts"`
}

// UnmarshalJSON overrides the JSON behavior to add support for
// inline attrubtes. Inline attributes are pending due to
// https://github.com/golang/go/issues/6213
func (e *CloudEvent) UnmarshalJSON(b []byte) error {
	return jsonext.Unmarshal(b, e)
}

// MarshalJSON overrides the JSON behavior to add support for
// inline attrubtes. Inline attributes are pending due to
// https://github.com/golang/go/issues/6213
func (e CloudEvent) MarshalJSON() ([]byte, error) {
	return jsonext.Marshal(e)
}
