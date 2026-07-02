package git

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
)

type TreeEntry struct {
	Name string `parquet:"name,dict"`
	Mode uint32 `parquet:"mode"`
	OID  OID    `parquet:"oid"`
}

type Tree struct {
	OID     OID         `parquet:"oid"`
	Entries []TreeEntry `parquet:"entries,list"`
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Tree.
func (t *Tree) UnmarshalBinary(data []byte) error {
	/*prefixBytes, data, ok := bytes.Cut(data, []byte{0x00})
	if !ok {
		return fmt.Errorf("invalid tree data: missing \\0 separator")
	}
	sizeBytes, ok := bytes.CutPrefix(prefixBytes, []byte("tree "))
	if !ok {
		return errors.New("invalid prefix: expected 'tree '")
	}
	size, err := strconv.ParseInt(string(sizeBytes), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid tree size: %w", err)
	} else if size != int64(len(data)) {
		return fmt.Errorf("tree size mismatch: expected %d bytes, got %d", size, len(data))
	}*/
	for {
		var entry TreeEntry
		// Parse the filemode (ex: "100644")
		var modeBytes []byte
		var ok bool
		modeBytes, data, ok = bytes.Cut(data, []byte(" "))
		if !ok {
			if len(modeBytes) == 0 {
				return nil // End of tree entries
			}
			return io.ErrUnexpectedEOF
		}
		if mode, err := strconv.ParseUint(string(modeBytes), 8, 32); err != nil {
			return fmt.Errorf("invalid file mode: %w", err)
		} else {
			entry.Mode = uint32(mode)
		}
		// Parse the name (ex: "README.md")
		var nameBytes []byte
		nameBytes, data, ok = bytes.Cut(data, []byte{0x00})
		if !ok {
			return io.ErrUnexpectedEOF
		}
		entry.Name = string(nameBytes)
		// Parse the OID
		if len(data) < sha1.Size {
			return io.ErrUnexpectedEOF
		}
		copy(entry.OID[:], data[:sha1.Size])
		data = data[sha1.Size:]
		// Add the entry to the tree
		t.Entries = append(t.Entries, entry)
	}
}
