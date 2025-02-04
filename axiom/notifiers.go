package axiom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=Type -linecomment -output=notifiers_string.go

// Type represents the type of a notifier.
type Type uint8

// All available notifier types.
const (
	emptyType Type = iota //

	Pagerduty // pagerduty
	Slack     // slack
	Email     // email
	Webhook   // webhook
)

func typeFromString(s string) (t Type, err error) {
	switch s {
	case emptyType.String():
		t = emptyType
	case Pagerduty.String():
		t = Pagerduty
	case Slack.String():
		t = Slack
	case Email.String():
		t = Email
	case Webhook.String():
		t = Webhook
	default:
		err = fmt.Errorf("unknown type %q", s)
	}

	return t, err
}

// MarshalJSON implements `json.Marshaler`. It is in place to marshal the
// Type to its string representation because that's what the server
// expects.
func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON implements `json.Unmarshaler`. It is in place to unmarshal the
// Type from the string representation the server returns.
func (t *Type) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}

	*t, err = typeFromString(s)

	return err
}

// A Notifier alerts users by using the configured service to reach out to them.
type Notifier struct {
	// ID is the unique ID of the notifier.
	ID string `json:"id"`
	// Name is the display name of the notifier.
	Name string `json:"name"`
	// Type of a notifier.
	Type Type `json:"type"`
	// Properties of the notifier.
	Properties interface{} `json:"properties"`
	// DisabledUntil is the time until the notifier is being executed again.
	DisabledUntil time.Time `json:"disabledUntil"`
	// CreatedAt is the time the notifer was created.
	CreatedAt time.Time `json:"metaCreated"`
	// ModifiedAt is the time the notifer was last modified.
	ModifiedAt time.Time `json:"metaModified"`
	// Version of the notifier.
	Version int64 `json:"metaVersion"`
}

// NotifiersService handles communication with the notifier related operations of
// the Axiom API.
//
// Axiom API Reference: /api/v1/notifiers
type NotifiersService service

// List all available notifiers.
func (s *NotifiersService) List(ctx context.Context) ([]*Notifier, error) {
	var res []*Notifier
	if err := s.client.call(ctx, http.MethodGet, s.basePath, nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get a notifier by id.
func (s *NotifiersService) Get(ctx context.Context, id string) (*Notifier, error) {
	path := s.basePath + "/" + id

	var res Notifier
	if err := s.client.call(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Create a notifier with the given properties.
func (s *NotifiersService) Create(ctx context.Context, req Notifier) (*Notifier, error) {
	var res Notifier
	if err := s.client.call(ctx, http.MethodPost, s.basePath, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Update the notifier identified by the given id with the given properties.
func (s *NotifiersService) Update(ctx context.Context, id string, req Notifier) (*Notifier, error) {
	path := s.basePath + "/" + id

	var res Notifier
	if err := s.client.call(ctx, http.MethodPut, path, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Delete the notifier identified by the given id.
func (s *NotifiersService) Delete(ctx context.Context, id string) error {
	return s.client.call(ctx, http.MethodDelete, s.basePath+"/"+id, nil, nil)
}
