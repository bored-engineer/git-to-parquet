package git

import "fmt"

type Type string

const (
	TypeTag    Type = "tag"
	TypeCommit Type = "commit"
	TypeTree   Type = "tree"
	TypeBlob   Type = "blob"
)

func (t Type) String() string {
	return string(t)
}

func (t *Type) UnmarshalText(text []byte) error {
	switch string(text) {
	case "tag", "commit", "tree", "blob":
		*t = Type(string(text))
	default:
		return fmt.Errorf("invalid type: expected 'tag', 'commit', 'tree', or 'blob', got '%s'", string(text))
	}
	return nil
}
