package colly

import (
	"bytes"
	"net/http"
	"testing"
)

// windows-1251 (CP1251) encoding of "Прометей"
var cp1251Prometei = []byte{0xcf, 0xf0, 0xee, 0xec, 0xe5, 0xf2, 0xe5, 0xe9}

// TestResponseFixCharsetMetaDeclaration reproduces
// https://github.com/gocolly/colly/issues/777 : when the Content-Type header
// has no charset, fixCharset relies purely on chardet's statistical detection
// and ignores the charset declared in <meta http-equiv="Content-Type">.
func TestResponseFixCharsetMetaDeclaration(t *testing.T) {
	body := []byte(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=windows-1251">
<title>` + string(cp1251Prometei) + `</title>
</head>
<body>` + string(cp1251Prometei) + `</body>
</html>`)

	r := &Response{
		Body: body,
		Headers: &http.Header{
			"Content-Type": {"text/html"},
		},
	}

	if err := r.fixCharset(true, ""); err != nil {
		t.Fatalf("fixCharset returned error: %v", err)
	}

	// After decoding the windows-1251 bytes to UTF-8, the body must contain
	// the proper UTF-8 encoding of "Прометей".
	utf8Prometei := []byte("Прометей")
	if !bytes.Contains(r.Body, utf8Prometei) {
		t.Fatalf("body does not contain UTF-8 \"Прометей\" after fixCharset; got %q", r.Body)
	}
	// The raw windows-1251 byte sequence must be gone.
	if bytes.Contains(r.Body, cp1251Prometei) {
		t.Fatalf("body still contains raw windows-1251 bytes after fixCharset: %q", r.Body)
	}
}
