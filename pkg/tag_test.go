package git

import (
	"bytes"
	"testing"
)

var tagBytes = []byte(`object 78d41cde5cbcadc851b6a9426d78bc88fc633fc2
type tag
tag foobarz
tagger Tag Ger <tagger@example.com> 1465981137 +0000
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

func TestTag(t *testing.T) {
	var tag Tag
	if err := tag.UnmarshalBinary(tagBytes); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	/*if expected := "0c4ac5dc0660929de0be541e5bf19f4dc10c7a01"; tag.OID.String() != expected {
	  t.Errorf("expected OID %s, got %s", expected, tag.OID)
	 }*/
	if expected := "78d41cde5cbcadc851b6a9426d78bc88fc633fc2"; tag.TargetOID.String() != expected {
		t.Errorf("expected TargetOID %s, got %s", expected, tag.TargetOID)
	}
	if !bytes.Equal(tag.Name, []byte("foobarz")) {
		t.Errorf("expected Name 'foobarz', got '%s'", tag.Name)
	}
	if !bytes.Equal(tag.Tagger.Name, []byte("Tag Ger")) || !bytes.Equal(tag.Tagger.Email, []byte("tagger@example.com")) {
		t.Errorf("expected tagger 'Tag Ger <tagger@example.com>', got '%s <%s>'", tag.Tagger.Name, tag.Tagger.Email)
	}
	if len(tag.Headers) != 1 {
		t.Fatalf("expected 1 header, got %d", len(tag.Headers))
	}
	if expected := []byte("gpgsig"); !bytes.Equal(tag.Headers[0].Name, expected) {
		t.Errorf("expected header %q, got %q", expected, tag.Headers[0].Name)
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
-----END PGP SIGNATURE-----`); !bytes.Equal(tag.Headers[0].Value, expected) {
		t.Errorf("expected header value %q, got %q", expected, tag.Headers[0].Value)
	}
	if expected := []byte("message body"); !bytes.Equal(tag.Message, expected) {
		t.Errorf("expected message %q, got %q", expected, tag.Message)
	}
}
