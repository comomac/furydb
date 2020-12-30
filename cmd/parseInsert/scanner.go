package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	// If we see a digit then consume as a number.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace(true)
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(ch)
	case ',':
		return COMMA, string(ch)
	case '(':
		return LEFTPAR, string(ch)
	case ')':
		return RIGHTPAR, string(ch)
	case '\'':
		return SINGLEQUO, string(ch)
	case '"':
		return DOUBLEQUO, string(ch)
	case ';':
		return SEMICOL, string(ch)
	}

	return ILLEGAL, string(ch)
}

// Scan returns the next token and literal value.
func (s *Scanner) ScanValue() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	// If we see a digit then consume as a number.
	if isWhitespaceB(ch) {
		s.unread()
		return s.scanWhitespace(false)
	} else if ch != ',' && ch != '\'' && ch != ';' && ch != ' ' && ch != ')' {
		s.unread()
		return s.scanValue()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case ',':
		return COMMA, string(ch)
	case '\'':
		return SINGLEQUO, string(ch)
	case ';':
		return SEMICOL, string(ch)
	case ' ':
		return WS, string(ch)
	case ')':
		return RIGHTPAR, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace(packSpaces bool) (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if packSpaces && !isWhitespace(ch) {
			s.unread()
			break
		} else if !packSpaces && !isWhitespaceB(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// ScanValue consumes the current rune and all contiguous value runes.
func (s *Scanner) ScanValueOLD() (tok Token, lit string) {
	// skip whitespace
	ch := s.read()
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace(true)
	}
	s.unread()

	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	var quoteOn bool
	var isNumber bool

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		// Read the next rune.
		ch := s.read()

		if ch == eof {
			return ILLEGAL, string(ch)

			// } else if !isNumber && !quoteOn && ch == ',' {
		} else if ch == ',' {
			fmt.Println("1111111111111111 COMMA")
			s.unread()
			return COMMA, string(ch)

		} else if !isNumber && !quoteOn && isDigit(ch) || ch == '.' || ch == '-' {
			isNumber = true
			buf.WriteRune(ch)
			continue

		} else if isNumber {
			if isDigit(ch) || ch == '.' || ch == '-' {
				buf.WriteRune(ch)
				continue
			} else {
				return ILLEGAL, fmt.Sprintf("mix number + text. [%+v]", string(ch))
			}

		} else if !quoteOn && ch == '\'' {
			// opening quote
			quoteOn = true
			continue

		} else if quoteOn && ch == '\'' {
			// read ahead
			ch2 := s.read()

			if ch2 == '\'' {
				// escape quote
				buf.WriteRune(ch)
				continue

			} else {
				// closing quote
				quoteOn = false
				break
			}
		} else if quoteOn {
			// any rune
			buf.WriteRune(ch)
			continue

		} else {
			return ILLEGAL, string(ch)
		}
	}

	return VALUE, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanValue() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == ',' {
			s.unread()
			break
		} else if ch == '\'' {
			s.unread()
			break
		} else if ch == ';' {
			s.unread()
			break
		} else if ch == ')' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// Otherwise return as a regular identifier.
	return VALUE, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToUpper(buf.String()) {
	case "INSERT":
		return INSERT, buf.String()
	case "INTO":
		return INTO, buf.String()
	case "VALUES":
		return VALUES, buf.String()
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isWhitespaceB returns true if the rune is a tab, or newline.
func isWhitespaceB(ch rune) bool { return ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
