package syntax

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type tokenKind int

const (
	tokenEOF tokenKind = 256 + iota
	tokenIdentifier
	tokenString
	tokenInteger
	tokenNewline
)

type token struct {
	kind  tokenKind
	text  string
	value *Value
	start int
	end   int
}

type lexer struct {
	filename string
	source   string
	offset   int
	comments []Comment
	lines    []int
}

func newLexer(filename, source string) *lexer {
	l := &lexer{filename: filename, source: source, lines: []int{0}}
	for i := 0; i < len(source); i++ {
		if source[i] == '\r' {
			if i+1 < len(source) && source[i+1] == '\n' {
				i++
			}
			l.lines = append(l.lines, i+1)
		} else if source[i] == '\n' {
			l.lines = append(l.lines, i+1)
		}
	}
	return l
}

func (l *lexer) position(offset int) Position {
	// Binary search keeps source locations inexpensive even for large documents.
	lo, hi := 0, len(l.lines)
	for lo+1 < hi {
		mid := (lo + hi) / 2
		if l.lines[mid] <= offset {
			lo = mid
		} else {
			hi = mid
		}
	}
	return Position{Offset: offset, Line: lo + 1, Column: utf8.RuneCountInString(l.source[l.lines[lo]:offset]) + 1}
}

func (l *lexer) span(start, end int) Range {
	return Range{Filename: l.filename, Start: l.position(start), End: l.position(end)}
}

func (l *lexer) fail(start, end int, message string) error {
	return &Diagnostic{Range: l.span(start, end), Message: message}
}

func (l *lexer) next() (token, error) {
	for l.offset < len(l.source) {
		start := l.offset
		c := l.source[start]
		if c == ' ' || c == '\t' {
			l.offset++
			continue
		}
		if strings.HasPrefix(l.source[start:], "//") {
			for l.offset < len(l.source) && l.source[l.offset] != '\n' && l.source[l.offset] != '\r' {
				l.offset++
			}
			l.comments = append(l.comments, Comment{Text: l.source[start:l.offset], Range: l.span(start, l.offset)})
			continue
		}
		l.offset++
		if c == '\n' || c == '\r' {
			if c == '\r' && l.offset < len(l.source) && l.source[l.offset] == '\n' {
				l.offset++
			}
			return token{kind: tokenNewline, start: start, end: l.offset}, nil
		}
		if strings.ContainsRune("{}[]=,;", rune(c)) {
			return token{kind: tokenKind(c), text: string(c), start: start, end: l.offset}, nil
		}
		if identifierStart(c) {
			for l.offset < len(l.source) && identifierContinue(l.source[l.offset]) {
				l.offset++
			}
			return token{kind: tokenIdentifier, text: l.source[start:l.offset], start: start, end: l.offset}, nil
		}
		if c >= '0' && c <= '9' || c == '-' && l.offset < len(l.source) && l.source[l.offset] >= '0' && l.source[l.offset] <= '9' {
			for l.offset < len(l.source) && l.source[l.offset] >= '0' && l.source[l.offset] <= '9' {
				l.offset++
			}
			value, err := strconv.ParseInt(l.source[start:l.offset], 10, 64)
			if err != nil {
				return token{}, l.fail(start, l.offset, "integer is outside the supported 64-bit range")
			}
			return token{kind: tokenInteger, value: &Value{Kind: Integer, Int: value, Range: l.span(start, l.offset)}, start: start, end: l.offset}, nil
		}
		if c == '"' || c == '\'' {
			l.offset = start
			return l.stringToken()
		}
		r, width := utf8.DecodeRuneInString(l.source[start:])
		l.offset = start + width
		return token{}, l.fail(start, l.offset, fmt.Sprintf("unexpected character %q", r))
	}
	return token{kind: tokenEOF, start: len(l.source), end: len(l.source)}, nil
}

func identifierStart(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func identifierContinue(c byte) bool {
	return identifierStart(c) || c >= '0' && c <= '9' || c == '-'
}

func validIdentifier(s string) bool {
	if s == "" || !identifierStart(s[0]) {
		return false
	}
	for i := 1; i < len(s); i++ {
		if !identifierContinue(s[i]) {
			return false
		}
	}
	return true
}

func (l *lexer) stringToken() (token, error) {
	start := l.offset
	quote := l.source[start]
	triple := strings.HasPrefix(l.source[start:], strings.Repeat(string(quote), 3))
	literal := quote == '\''
	if literal && !triple {
		return token{}, l.fail(start, start+1, "literal strings use triple single quotes")
	}
	delimiter := string(quote)
	if triple {
		delimiter = strings.Repeat(delimiter, 3)
	}
	l.offset += len(delimiter)
	if triple {
		if l.offset >= len(l.source) || l.source[l.offset] != '\n' && l.source[l.offset] != '\r' {
			return token{}, l.fail(start, l.offset, "triple-quoted string must begin with a newline")
		}
		if l.source[l.offset] == '\r' && l.offset+1 < len(l.source) && l.source[l.offset+1] == '\n' {
			l.offset++
		}
		l.offset++
	}
	bodyStart := l.offset
	for l.offset < len(l.source) {
		if strings.HasPrefix(l.source[l.offset:], delimiter) {
			body := l.source[bodyStart:l.offset]
			if triple {
				lineStart := l.offset
				for lineStart > bodyStart && l.source[lineStart-1] != '\n' && l.source[lineStart-1] != '\r' {
					lineStart--
				}
				indent := l.source[lineStart:l.offset]
				if strings.Trim(indent, " \t") != "" {
					return token{}, l.fail(lineStart, l.offset, "closing triple quotes must be on their own line")
				}
				body = strings.ReplaceAll(l.source[bodyStart:lineStart], "\r\n", "\n")
				body = strings.ReplaceAll(body, "\r", "\n")
				body = strings.TrimSuffix(body, "\n")
				lines := strings.Split(body, "\n")
				for i, line := range lines {
					if strings.Trim(line, " \t") == "" {
						lines[i] = ""
					} else if !strings.HasPrefix(line, indent) {
						return token{}, l.fail(start, l.offset, "multiline content has less indentation than the closing delimiter")
					} else {
						lines[i] = strings.TrimPrefix(line, indent)
					}
				}
				body = strings.Join(lines, "\n")
			}
			l.offset += len(delimiter)
			parts, err := decodeString(body, literal)
			if err != nil {
				return token{}, l.fail(start, l.offset, err.Error())
			}
			return token{kind: tokenString, value: &Value{Kind: String, Parts: parts, Literal: literal, Range: l.span(start, l.offset)}, start: start, end: l.offset}, nil
		}
		c := l.source[l.offset]
		if c < 0x20 && c != '\t' && (c != '\r' && c != '\n' || !triple) {
			return token{}, l.fail(l.offset, l.offset+1, "control character or newline in a quoted string")
		}
		if c == '\\' && !literal {
			l.offset++
			if l.offset >= len(l.source) {
				break
			}
		}
		l.offset++
	}
	return token{}, l.fail(start, len(l.source), "unterminated string")
}

func decodeString(body string, literal bool) ([]StringPart, error) {
	if literal {
		return []StringPart{{Text: body}}, nil
	}
	var parts []StringPart
	var text strings.Builder
	flush := func() {
		if text.Len() > 0 {
			parts = append(parts, StringPart{Text: text.String()})
			text.Reset()
		}
	}
	for i := 0; i < len(body); {
		if body[i] == '\\' {
			i++
			if i == len(body) {
				return nil, fmt.Errorf("incomplete string escape")
			}
			switch body[i] {
			case '\\', '"':
				text.WriteByte(body[i])
			case 'n':
				text.WriteByte('\n')
			case 'r':
				text.WriteByte('\r')
			case 't':
				text.WriteByte('\t')
			case '$':
				if i+1 >= len(body) || body[i+1] != '{' {
					return nil, fmt.Errorf("only \\${ escapes an interpolation opener")
				}
				text.WriteString("${")
				i++
			default:
				return nil, fmt.Errorf("unsupported string escape \\%c", body[i])
			}
			i++
			continue
		}
		if strings.HasPrefix(body[i:], "${") {
			end := strings.IndexByte(body[i+2:], '}')
			if end < 0 {
				return nil, fmt.Errorf("unterminated interpolation; expected ${name}")
			}
			name := body[i+2 : i+2+end]
			if !validIdentifier(name) {
				return nil, fmt.Errorf("interpolation must name a binding, not an expression")
			}
			flush()
			parts = append(parts, StringPart{Reference: name})
			i += end + 3
			continue
		}
		text.WriteByte(body[i])
		i++
	}
	flush()
	return parts, nil
}
