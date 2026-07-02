package git

import (
	"bytes"
	"testing"
)

var commitBytes = []byte(`tree eebfed94e75e7760540d1485c740902590a00332
parent 04b871796dc0420f8e7561a895b52484b701d51a
author A U Thor <author@example.com> 1465981137 +0000
committer C O Mitter <committer@example.com> 1465981137 +0000
gpgsig -----BEGIN PGP SIGNATURE-----
 Version: GnuPG v1
 $
 iQEcBAABAgAGBQJXYRjRAAoJEGEJLoW3InGJ3IwIAIY4SA6GxY3BjL60YyvsJPh/
 HRCJwH+w7wt3Yc/9/bW2F+gF72kdHOOs2jfv+OZhq0q4OAN6fvVSczISY/82LpS7
 DVdMQj2/YcHDT4xrDNBnXnviDO9G7am/9OE77kEbXrp7QPxvhjkicHNwy2rEflAA
 zn075rtEERDHr8nRYiDh8eVrefSO7D+bdQ7gv+7GsYMsd2auJWi1dHOSfTr9HIF4
 HJhWXT9d2f8W+diRYXGh4X0wYiGg6na/soXc+vdtDYBzIxanRqjg8jCAeo1eOTk1
 EdTwhcTZlI0x5pvJ3H0+4hA2jtldVtmPM4OTB0cTrEWBad7XV6YgiyuII73Ve3I=
 =jKHM
 -----END PGP SIGNATURE-----

message body
`)

func TestCommit(t *testing.T) {
	var commit Commit
	if err := commit.UnmarshalBinary(commitBytes); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	/*if expected := "0c4ac5dc0660929de0be541e5bf19f4dc10c7a01"; commit.OID.String() != expected {
		t.Errorf("expected OID %s, got %s", expected, commit.OID)
	}*/
	if expected := "eebfed94e75e7760540d1485c740902590a00332"; commit.TreeOID.String() != expected {
		t.Errorf("expected TreeOID %s, got %s", expected, commit.TreeOID)
	}
	if len(commit.Parents) != 1 {
		t.Fatalf("expected 1 parent, got %d", len(commit.Parents))
	}
	if expected := "04b871796dc0420f8e7561a895b52484b701d51a"; commit.Parents[0].String() != expected {
		t.Errorf("expected parent OID %s, got %s", expected, commit.Parents[0])
	}
	if !bytes.Equal(commit.Author.Name, []byte("A U Thor")) || !bytes.Equal(commit.Author.Email, []byte("author@example.com")) {
		t.Errorf("expected author 'A U Thor <author@example.com>', got '%s <%s>'", commit.Author.Name, commit.Author.Email)
	}
	if !bytes.Equal(commit.Committer.Name, []byte("C O Mitter")) || !bytes.Equal(commit.Committer.Email, []byte("committer@example.com")) {
		t.Errorf("expected committer 'C O Mitter', got '%s <%s>'", commit.Committer.Name, commit.Committer.Email)
	}
	if len(commit.Headers) != 1 {
		t.Fatalf("expected 1 header, got %d", len(commit.Headers))
	}
	if expected := []byte("gpgsig"); !bytes.Equal(commit.Headers[0].Name, expected) {
		t.Errorf("expected header %q, got %q", expected, commit.Headers[0].Name)
	}
	if expected := []byte(`-----BEGIN PGP SIGNATURE-----
Version: GnuPG v1
$
iQEcBAABAgAGBQJXYRjRAAoJEGEJLoW3InGJ3IwIAIY4SA6GxY3BjL60YyvsJPh/
HRCJwH+w7wt3Yc/9/bW2F+gF72kdHOOs2jfv+OZhq0q4OAN6fvVSczISY/82LpS7
DVdMQj2/YcHDT4xrDNBnXnviDO9G7am/9OE77kEbXrp7QPxvhjkicHNwy2rEflAA
zn075rtEERDHr8nRYiDh8eVrefSO7D+bdQ7gv+7GsYMsd2auJWi1dHOSfTr9HIF4
HJhWXT9d2f8W+diRYXGh4X0wYiGg6na/soXc+vdtDYBzIxanRqjg8jCAeo1eOTk1
EdTwhcTZlI0x5pvJ3H0+4hA2jtldVtmPM4OTB0cTrEWBad7XV6YgiyuII73Ve3I=
=jKHM
-----END PGP SIGNATURE-----`); !bytes.Equal(commit.Headers[0].Value, expected) {
		t.Errorf("expected header value %q, got %q", expected, commit.Headers[0].Value)
	}
	if expected := []byte("message body"); !bytes.Equal(commit.Message, expected) {
		t.Errorf("expected message %q, got %q", expected, commit.Message)
	}
	if len(commit.Trailers) != 0 {
		t.Errorf("expected 0 footers, got %d", len(commit.Trailers))
	}
}

var commitWithTrailersBytes = []byte(`tree eebfed94e75e7760540d1485c740902590a00332
author A U Thor <author@example.com> 1465981137 +0000
committer C O Mitter <committer@example.com> 1465981137 +0000

Fix the thing

This commit fixes the thing by doing the other thing.

Signed-off-by: A U Thor <author@example.com>
Co-authored-by: C O Mitter <committer@example.com>
Fixes: #42
`)

func TestCommitWithTrailers(t *testing.T) {
	var commit Commit
	if err := commit.UnmarshalBinary(commitWithTrailersBytes); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if expected := []byte("Fix the thing\n\nThis commit fixes the thing by doing the other thing."); !bytes.Equal(commit.Message, expected) {
		t.Errorf("expected message %q, got %q", expected, commit.Message)
	}
	cases := []Header{
		{Name: []byte("Multi-line"), Value: []byte("This is a multi-line trailer\nthat continues on the next line\nand the next line")},
		{Name: []byte("Fixes"), Value: []byte("#42")},
		{Name: []byte("Co-authored-by"), Value: []byte("C O Mitter <committer@example.com>")},
		{Name: []byte("Signed-off-by"), Value: []byte("A U Thor <author@example.com>")},
	}
	if len(commit.Trailers) != len(cases) {
		t.Fatalf("expected %d footers, got %d: %+v", len(cases), len(commit.Trailers), commit.Trailers)
	}
	for i, expected := range cases {
		if !bytes.Equal(commit.Trailers[i].Name, expected.Name) || !bytes.Equal(commit.Trailers[i].Value, expected.Value) {
			t.Errorf("footer[%d]: expected %+v, got %+v", i, expected, commit.Trailers[i])
		}
	}
}
