package git

import (
	"bytes"
	"fmt"
)

type Tag struct {
	OID       OID      `parquet:"oid"`
	Name      []byte   `parquet:"name"`
	Type      Type     `parquet:"type,dict,enum"`
	Tagger    Identity `parquet:"tagger"`
	TargetOID OID      `parquet:"target_oid"`
	Headers   []Header `parquet:"headers"`
	Message   []byte   `parquet:"message"`
}

func (t *Tag) UnmarshalBinary(data []byte) error {
	/*t.OID = sha1.Sum(data)
	prefixBytes, data, ok := bytes.Cut(data, []byte{0x00})
	if !ok {
		return fmt.Errorf("invalid tag data: missing \\0 separator")
	}
	sizeBytes, ok := bytes.CutPrefix(prefixBytes, []byte("tag "))
	if !ok {
		return errors.New("invalid prefix: expected 'tag '")
	}
	size, err := strconv.ParseInt(string(sizeBytes), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid tag size: %w", err)
	} else if size != int64(len(data)) {
		return fmt.Errorf("tag size mismatch: expected %d bytes, got %d", size, len(data))
	}*/
	for {
		var line []byte
		var ok bool
		line, data, ok = bytes.Cut(data, []byte{'\n'})
		if !ok || len(line) == 0 {
			break
		}
		header, value, _ := bytes.Cut(line, []byte{' '})
		switch string(header) {
		case "object":
			if err := t.TargetOID.UnmarshalText(value); err != nil {
				return fmt.Errorf("invalid tree OID: %w", err)
			}
		case "tag":
			t.Name = value
		case "type":
			if err := t.Type.UnmarshalText(value); err != nil {
				return fmt.Errorf("invalid tag type: %w", err)
			}
		case "tagger":
			if err := t.Tagger.UnmarshalText(value); err != nil {
				return fmt.Errorf("invalid tagger: %w", err)
			}
		default:
			if len(header) == 0 {
				// Because we're stripping the ' ' we have to allocate a new slice to hold the modified value.
				last := &t.Headers[len(t.Headers)-1]
				clone := make([]byte, len(last.Value)+1+len(value))
				copy(clone, last.Value)
				clone[len(last.Value)] = '\n'
				copy(clone[len(last.Value)+1:], value)
				last.Value = clone
			} else {
				t.Headers = append(t.Headers, Header{
					Name:  header,
					Value: value,
				})
			}
		}
	}
	t.Message = bytes.TrimSpace(data)
	return nil
}
