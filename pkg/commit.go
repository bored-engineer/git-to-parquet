package git

import (
	"bytes"
	"fmt"
)

type Header struct {
	Name  []byte `parquet:"name,dict"`
	Value []byte `parquet:"value"`
}

type Commit struct {
	OID       OID      `parquet:"oid"`
	TreeOID   OID      `parquet:"tree_oid"`
	Author    Identity `parquet:"author"`
	Committer Identity `parquet:"committer"`
	Parents   []OID    `parquet:"parents"`
	Encoding  []byte   `parquet:"encoding,dict"`
	Headers   []Header `parquet:"headers"`
	Message   []byte   `parquet:"message"`
	Trailers  []Header `parquet:"trailers"`
}

func (c *Commit) UnmarshalBinary(data []byte) error {
	/*c.OID = sha1.Sum(data)
	prefixBytes, data, ok := bytes.Cut(data, []byte{0x00})
	if !ok {
		return fmt.Errorf("invalid commit data: missing \\0 separator")
	}
	sizeBytes, ok := bytes.CutPrefix(prefixBytes, []byte("commit "))
	if !ok {
		return errors.New("invalid prefix: expected 'commit '")
	}
	size, err := strconv.ParseInt(string(sizeBytes), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid commit size: %w", err)
	} else if size != int64(len(data)) {
		return fmt.Errorf("commit size mismatch: expected %d bytes, got %d", size, len(data))
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
		case "tree":
			if err := c.TreeOID.UnmarshalText(value); err != nil {
				return fmt.Errorf("invalid tree OID: %w", err)
			}
		case "author":
			if err := c.Author.UnmarshalText(value); err != nil {
				return fmt.Errorf("invalid author: %w", err)
			}
		case "committer":
			if err := c.Committer.UnmarshalText(value); err != nil {
				return fmt.Errorf("invalid committer: %w", err)
			}
		case "encoding":
			c.Encoding = value
		case "parent":
			var parent OID
			if err := parent.UnmarshalText(value); err != nil {
				return fmt.Errorf("invalid parent OID: %w", err)
			}
			c.Parents = append(c.Parents, parent)
		default:
			if len(header) == 0 {
				// Because we're stripping the ' ' we have to allocate a new slice to hold the modified value.
				last := &c.Headers[len(c.Headers)-1]
				clone := make([]byte, len(last.Value)+1+len(value))
				copy(clone, last.Value)
				clone[len(last.Value)] = '\n'
				copy(clone[len(last.Value)+1:], value)
				last.Value = clone
			} else {
				c.Headers = append(c.Headers, Header{
					Name:  header,
					Value: value,
				})
			}
		}
	}
	c.Message, c.Trailers = parseTrailers(bytes.TrimSpace(data))
	return nil
}

// parseTrailers splits a commit message into the body and any trailing git
// trailers. Trailers are the last paragraph of lines matching "key: value"
// separated from the body by a blank line.
func parseTrailers(msg []byte) (_ []byte, trailers []Header) {
	for {
		idx := bytes.LastIndexByte(msg, '\n')
		if idx == -1 {
			break
		}
		line := msg[idx+1:]
		key, value, ok := bytes.Cut(line, []byte{':'})
		if !ok || bytes.ContainsAny(key, " \t") || len(key) == 0 {
			break
		}
		msg = msg[:idx]
		trailers = append(trailers, Header{
			Name:  key,
			Value: bytes.TrimLeft(value, " "),
		})
	}
	return bytes.TrimRight(msg, "\n"), trailers
	/*lines := bytes.Split(bytes.TrimRight(msg, "\n"), []byte{'\n'})
	// Find the start of the last paragraph (after the last blank line).
	start := len(lines)
	for i := len(lines) - 1; i >= 0; i-- {
		if len(bytes.TrimSpace(lines[i])) == 0 {
			break
		}
		start = i
	}
	if start >= len(lines) {
		return msg, nil
	}
	// Every line in the candidate block must match "token: value".
	for _, line := range lines[start:] {
		key, value, ok := bytes.Cut(line, []byte{':'})
		if !ok || bytes.ContainsAny(key, " \t") || len(key) == 0 {
			return msg, nil
		}
		trailers = append(trailers, Header{
			Name:  key,
			Value: bytes.TrimLeft(value, " "),
		})
	}
	// Trim the footer paragraph (and the blank line before it) from the body.
	body = bytes.TrimRight(bytes.Join(lines[:start], []byte{'\n'}), "\n")
	return body, trailers*/
}
