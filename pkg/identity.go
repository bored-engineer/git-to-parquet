package git

import (
	"bytes"
	"strconv"
	"time"
)

// Identity represents a Git identity, containing the author's name, email, and (hopefully) a timestamp.
type Identity struct {
	Name  []byte `parquet:"name,dict"`
	Email []byte `parquet:"email,dict,optional"`
	// Timestamp is the parsed unix timestamp, adjusted based on the provided timezone.
	Timestamp *time.Time `parquet:"timestamp,timestamp,optional"`
	// Timezone is the raw timezone string, e.g. "-0800".
	// It can be useful to retain to determine the original timezone of the commit.
	Timezone []byte `parquet:"timezone,dict,optional"`
}

// UnmarshalText replicates the split_ident_line method from git's ident.c
// https://github.com/git/git/blob/845c48a16a7f7b2c44d8cb137b16a4a1f0140229/ident.c#L275C5-L275C21
// However, it has been modified to be more permissive to avoid loss of corrupt data
func (i *Identity) UnmarshalText(text []byte) error {
	// Ex: "Luke Young <github@bored.engineer> 1418751585 -0800"
	mailStart := bytes.IndexByte(text, '<')
	if mailStart == -1 {
		// Retain the malformed text as the "name"
		i.Name = text
		return nil
	}
	i.Name = bytes.Trim(text[:mailStart], " ")

	// We intentionally deviate from split_ident_line and use the _last_ '>'
	// This allows us to handle cases where the email address contains '>'
	mailEnd := bytes.LastIndexByte(text, '>')
	if mailEnd == -1 {
		// If we can't find a closing '>', the email is malformed, give up
		i.Name = text
		return nil
	}
	i.Email = bytes.Trim(text[mailStart+1:mailEnd], " ")

	// Hopefully, the remaining text is a timestamp and timezone
	tsBytes := bytes.Trim(text[mailEnd+1:], " ")
	if len(tsBytes) == 0 {
		// If there's no timestamp, bail now
		return nil
	}

	// Split the timestamp and timezone (if present)
	tsBytes, timezone, ok := bytes.Cut(tsBytes, []byte{' '})
	if tsInt, err := strconv.ParseInt(string(tsBytes), 10, 64); err == nil {
		ts := time.Unix(tsInt, 0).In(time.UTC)
		i.Timestamp = &ts
	}
	if ok {
		i.Timezone = bytes.TrimLeft(timezone, " ")
		if len(i.Timezone) == 5 && i.Timestamp != nil {
			if tz, err := strconv.Atoi(string(i.Timezone)); err == nil {
				hours, mins := tz/100, tz%100
				ts := i.Timestamp.In(time.FixedZone("", int((hours*60*60)+mins*60)))
				i.Timestamp = &ts
			}
		}
	}

	return nil
}
