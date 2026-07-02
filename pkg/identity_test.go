package git

import (
	"bytes"
	"testing"
	"time"
)

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestIdentity(t *testing.T) {
	// Test cases for the Identity struct
	tests := map[string]struct {
		input    string
		expected Identity
	}{
		"well-formed": {
			input: "Luke Young <github@bored.engineer> 1418751585 -0800",
			expected: Identity{
				Name:      []byte("Luke Young"),
				Email:     []byte("github@bored.engineer"),
				Timestamp: ptrTime(time.Date(2014, 12, 16, 9, 39, 45, 0, time.FixedZone("", -8*60*60))),
				Timezone:  []byte("-0800"),
			},
		},
		"no email": {
			input: "Luke Young <> 1418751585 -0800",
			expected: Identity{
				Name:      []byte("Luke Young"),
				Email:     nil,
				Timestamp: ptrTime(time.Date(2014, 12, 16, 9, 39, 45, 0, time.FixedZone("", -8*60*60))),
				Timezone:  []byte("-0800"),
			},
		},
		"no timestamp": {
			input: "Luke Young <github@bored.engineer>",
			expected: Identity{
				Name:  []byte("Luke Young"),
				Email: []byte("github@bored.engineer"),
			},
		},
		"malformed email": {
			input: "Luke Young <foo<>bar@bored.engineer> 1418751585 -0800",
			expected: Identity{
				Name:      []byte("Luke Young"),
				Email:     []byte("foo<>bar@bored.engineer"),
				Timestamp: ptrTime(time.Date(2014, 12, 16, 9, 39, 45, 0, time.FixedZone("", -8*60*60))),
				Timezone:  []byte("-0800"),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var identity Identity
			if err := identity.UnmarshalText([]byte(tc.input)); err != nil {
				t.Fatalf("UnmarshalText failed: %v", err)
			}
			if !bytes.Equal(identity.Name, tc.expected.Name) {
				t.Errorf("expected Name %q, got %q", tc.expected.Name, identity.Name)
			}
			if !bytes.Equal(identity.Email, tc.expected.Email) {
				t.Errorf("expected Email %q, got %q", tc.expected.Email, identity.Email)
			}
			if tc.expected.Timestamp == nil {
				if identity.Timestamp != nil {
					t.Errorf("expected nil Timestamp, got %v", identity.Timestamp)
				}
			} else {
				if identity.Timestamp == nil {
					t.Errorf("expected Timestamp %v, got nil", tc.expected.Timestamp)
				} else if !identity.Timestamp.Equal(*tc.expected.Timestamp) {
					t.Errorf("expected Timestamp %v, got %v", tc.expected.Timestamp, identity.Timestamp)
				}
			}
			if !bytes.Equal(identity.Timezone, tc.expected.Timezone) {
				t.Errorf("expected Timezone %q, got %q", tc.expected.Timezone, identity.Timezone)
			}
		})
	}
}
