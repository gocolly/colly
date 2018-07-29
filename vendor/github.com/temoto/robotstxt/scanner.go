package robotstxt

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"os"
	"unicode/utf8"
)

type byteScanner struct {
	ErrorCount int
	Quiet      bool

	buf           []byte
	pos           token.Position
	lastChunk     bool
	ch            rune
	keyTokenFound bool
}

var WhitespaceChars = []rune{' ', '\t', '\v'}

func newByteScanner(srcname string, quiet bool) *byteScanner {
	return &byteScanner{
		Quiet: quiet,
		ch:    -1,
		pos:   token.Position{Filename: srcname},
	}
}

func (s *byteScanner) Feed(input []byte, end bool) error {
	s.buf = input
	s.pos.Offset = 0
	s.pos.Line = 1
	s.pos.Column = 1
	s.lastChunk = end

	// Read first char into look-ahead buffer `s.ch`.
	if err := s.nextChar(); err != nil {
		return err
	}

	// Skip UTF-8 byte order mark
	if s.ch == 65279 {
		s.nextChar()
		s.pos.Column = 1
	}

	return nil
}

func (s *byteScanner) GetPosition() token.Position {
	return s.pos
}

func (s *byteScanner) Scan() (string, error) {
	//println("--- Scan(). Offset / len(s.buf): ", s.pos.Offset, len(s.buf))

	for {
		// Note Offset > len, not >=, so we can Scan last character.
		if s.lastChunk && s.pos.Offset > len(s.buf) {
			return "", io.EOF
		}

		s.skipSpace()

		if s.ch == -1 {
			return "", io.EOF
		}

		// EOL
		if s.isEol() {
			s.keyTokenFound = false
			// skip subsequent newline chars
			for s.ch != -1 && s.isEol() {
				s.nextChar()
			}
			// emit newline as separate token
			return "\n", nil
		}

		// skip comments
		if s.ch == '#' {
			s.keyTokenFound = false
			s.skipUntilEol()
			//            s.state = "start"
			if s.ch == -1 {
				return "", io.EOF
			}
			// emit newline as separate token
			return "\n", nil
		}

		// else we found something
		break
	}

	/*
	   if s.state == "start" {
	       s.state = "key"
	   }
	*/

	var tok bytes.Buffer
	tok.WriteRune(s.ch)
	s.nextChar()
	for s.ch != -1 && !s.isSpace() && !s.isEol() {
		// Do not consider ":" to be a token separator if a first key token
		// has already been found on this line (avoid cutting an absolute URL
		// after the "http:")
		if s.ch == ':' && !s.keyTokenFound {
			//            s.state = "pre-value"
			s.nextChar()
			s.keyTokenFound = true
			break
		}

		tok.WriteRune(s.ch)
		s.nextChar()
	}
	return tok.String(), nil
}

func (s *byteScanner) ScanAll() ([]string, error) {
	var results []string
	for {
		t, err := s.Scan()
		if t != "" {
			results = append(results, t)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return results, err
		}
	}
	return results, nil
}

func (s *byteScanner) error(pos token.Position, msg string) {
	s.ErrorCount++
	if !s.Quiet {
		fmt.Fprintf(os.Stderr, "robotstxt from %s: %s\n", pos.String(), msg)
	}
}

func (s *byteScanner) isEol() bool {
	return s.ch == '\n' || s.ch == '\r'
}

func (s *byteScanner) isSpace() bool {
	for _, r := range WhitespaceChars {
		if s.ch == r {
			return true
		}
	}
	return false
}

func (s *byteScanner) skipSpace() {
	//println("--- string(ch): ", s.ch, ".")
	for s.ch != -1 && s.isSpace() {
		s.nextChar()
	}
}

func (s *byteScanner) skipUntilEol() {
	//println("--- string(ch): ", s.ch, ".")
	for s.ch != -1 && !s.isEol() {
		s.nextChar()
	}
	// skip subsequent newline chars
	for s.ch != -1 && s.isEol() {
		s.nextChar()
	}
}

// Reads next Unicode char.
func (s *byteScanner) nextChar() error {
	//println("--- nextChar(). Offset / len(s.buf): ", s.pos.Offset, len(s.buf))

	if s.pos.Offset >= len(s.buf) {
		s.ch = -1
		return io.EOF
	}
	s.pos.Column++
	if s.ch == '\n' {
		s.pos.Line++
		s.pos.Column = 1
	}
	r, w := rune(s.buf[s.pos.Offset]), 1
	if r >= 0x80 {
		r, w = utf8.DecodeRune(s.buf[s.pos.Offset:])
		if r == utf8.RuneError && w == 1 {
			s.error(s.pos, "illegal UTF-8 encoding")
		}
	}
	s.pos.Column++
	s.pos.Offset += w
	s.ch = r
	return nil
}
