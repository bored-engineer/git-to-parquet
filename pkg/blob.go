package git

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

// Blob represents a Git 'blob' object, which optionally contains the raw file contents.
type Blob struct {
	OID      OID               `parquet:"oid"`
	Size     int64             `parquet:"size"`
	MD5      [md5.Size]byte    `parquet:"md5"`
	SHA1     [sha1.Size]byte   `parquet:"sha1"`
	SHA256   [sha256.Size]byte `parquet:"sha256"`
	SHA512   [sha512.Size]byte `parquet:"sha512"`
	Contents []byte            `parquet:"bytes,optional"`
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Blob.
func (b *Blob) UnmarshalBinary(data []byte) error {
	/*b.OID = sha1.Sum(data)
	prefixBytes, data, ok := bytes.Cut(data, []byte{0x00})
	if !ok {
		return fmt.Errorf("invalid blob data: missing \\0 separator")
	}
	sizeBytes, ok := bytes.CutPrefix(prefixBytes, []byte("blob "))
	if !ok {
		return errors.New("invalid prefix: expected 'blob '")
	}
	size, err := strconv.ParseInt(string(sizeBytes), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid blob size: %w", err)
	} else if size != int64(len(data)) {
		return fmt.Errorf("blob size mismatch: expected %d bytes, got %d", size, len(data))
	}*/
	b.Size = int64(len(data))
	b.Contents = data
	b.MD5 = md5.Sum(data)
	b.SHA1 = sha1.Sum(data)
	b.SHA256 = sha256.Sum256(data)
	b.SHA512 = sha512.Sum512(data)
	return nil
}
