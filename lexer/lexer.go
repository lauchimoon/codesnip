package lexer

import (
    "strings"
    "slices"
    "unicode"
)

const (
    TOKEN_END = iota
    TOKEN_PREPROC
    TOKEN_SYMBOL
    TOKEN_OPEN_PAREN
    TOKEN_CLOSE_PAREN
    TOKEN_OPEN_CURLY
    TOKEN_CLOSE_CURLY
    TOKEN_SEMICOLON
    TOKEN_KEYWORD
    TOKEN_COMMENT
    TOKEN_CHAR
    TOKEN_STRING
    TOKEN_SPACE
    TOKEN_NEWLINE
    TOKEN_NUMBER
    TOKEN_PUNCT
    TOKEN_OPEN_BRACKET
    TOKEN_CLOSE_BRACKET
    TOKEN_OPERATOR
    TOKEN_INVALID
)

type TokenKind int

type Token struct {
    Kind    TokenKind
    Text    string
}

type Lexer struct {
    Content    string
    ContentLen int
    Cursor     int
}

var (
    keywords = []string{
        "auto", "break", "case", "char",
        "const", "continue", "default", "do",
        "double", "else", "enum", "extern",
        "float", "for", "goto", "if",
        "int", "long", "register", "return",
        "short", "signed", "sizeof", "static",
        "struct", "switch", "typedef", "union",
        "unsigned", "void", "volatile", "while",
    }
)

func LexerNew(source string) *Lexer {
    return &Lexer{
        Content: source,
        ContentLen: len(source),
        Cursor: 0,
    }
}

func (l *Lexer) Chop() rune {
    if l.Cursor >= l.ContentLen {
        return ' '
    }

    c := l.Content[l.Cursor]
    l.Cursor++
    return rune(c)
}

func (l *Lexer) PeekChar() rune {
    if l.Cursor >= l.ContentLen {
        return ' '
    }

    return rune(l.Content[l.Cursor + 1])
}

func (l *Lexer) Lex() []Token {
    tokens := []Token{}

    c := l.Chop()

    for l.Cursor < l.ContentLen {
        if isSymbol(c) {
            var buffer []rune
            for isSymbol(c) {
                buffer = append(buffer, c)
                c = l.Chop()
            }

            stringBuffer := string(buffer)
            if isKeyword(stringBuffer) {
                tokens = append(tokens, NewToken(TOKEN_KEYWORD, stringBuffer))
            } else if isNumber(stringBuffer) {
                tokens = append(tokens, NewToken(TOKEN_NUMBER, stringBuffer))
            } else {
                tokens = append(tokens, NewToken(TOKEN_SYMBOL, stringBuffer))
            }
        }

        if c == '/' && l.Content[l.Cursor] == '/' {
            // Get rid of slashes
            c = l.Chop()
            c = l.Chop()
            tokens = append(tokens, NewToken(TOKEN_COMMENT, "//"))
            var buffer []rune
            for c != '\n' {
                buffer = append(buffer, c)
                c = l.Chop()
            }

            tokens = append(tokens, NewToken(TOKEN_COMMENT, string(buffer)))
        }

        if c == '#' {
            tokens = append(tokens, NewToken(TOKEN_PREPROC, "#"))
            c = l.Chop()
            var buffer []rune
            for c != '\n' {
                buffer = append(buffer, c)
                c = l.Chop()
            }

            tokens = append(tokens, NewToken(TOKEN_PREPROC, string(buffer)))
        }

        if c == '"' {
            tokens = append(tokens, NewToken(TOKEN_STRING, "\""))
            c = l.Chop()
            var buffer []rune
            for c != '"' && l.Cursor < l.ContentLen {
                buffer = append(buffer, c)
                c = l.Chop()
            }

            tokens = append(tokens, NewToken(TOKEN_STRING, string(buffer)))
            if c == '"' {
                tokens = append(tokens, NewToken(TOKEN_STRING, "\""))
            }
        }

        if c == '\'' {
            tokens = append(tokens, NewToken(TOKEN_CHAR, "'"))
            c = l.Chop()
            var buffer []rune
            for c != '\'' && l.Cursor < l.ContentLen {
                buffer = append(buffer, c)
                c = l.Chop()
            }

            tokens = append(tokens, NewToken(TOKEN_CHAR, string(buffer)))
            if c == '\'' {
                tokens = append(tokens, NewToken(TOKEN_CHAR, "'"))
            }
        }

        if c == ' ' {
            tokens = append(tokens, NewToken(TOKEN_SPACE, " "))
        } else if c == '\n' {
            tokens = append(tokens, NewToken(TOKEN_NEWLINE, "\n"))
        } else if c == '(' {
            tokens = append(tokens, NewToken(TOKEN_OPEN_PAREN, "("))
        } else if c == ')' {
            tokens = append(tokens, NewToken(TOKEN_CLOSE_PAREN, ")"))
        } else if c == '{' {
            tokens = append(tokens, NewToken(TOKEN_OPEN_CURLY, "{"))
        } else if c == '}' {
            tokens = append(tokens, NewToken(TOKEN_CLOSE_CURLY, "}"))
        } else if c == ';' {
            tokens = append(tokens, NewToken(TOKEN_SEMICOLON, ";"))
        } else if isPunct(c) {
            tokens = append(tokens, NewToken(TOKEN_PUNCT, string(c)))
        } else if c == '[' {
            tokens = append(tokens, NewToken(TOKEN_OPEN_BRACKET, "["))
        } else if c == ']' {
            tokens = append(tokens, NewToken(TOKEN_OPEN_BRACKET, "]"))
        } else if isOperator(c) {
            tokens = append(tokens, NewToken(TOKEN_OPERATOR, string(c)))
        }
        c = l.Chop()
    }

    return tokens
}

func NewToken(kind TokenKind, text string) Token {
    return Token{
        Kind: kind,
        Text: text,
    }
}

func isSymbol(c rune) bool {
    return unicode.IsDigit(c) || unicode.IsLetter(c) || c == '_' || c == '*';
}

func isKeyword(s string) bool {
    return slices.Contains(keywords, s) || (strings.Contains(s, "*") && slices.Contains(keywords, strings.Trim(s, "*")))
}

func isNumber(s string) bool {
    for _, c := range s {
        if !(c >= '0' && c <= '9') {
            return false
        }
    }

    return true
}

func isPunct(c rune) bool {
    return c == ',' || c == '.' || c == '?' || c == ':' || c == '!'
}

func isOperator(c rune) bool {
    return c == '+' || c == '-' || c == '*' || c == '/' || c == '^' ||
        c == '<' || c == '>' || c == '=' || c == '&' || c == '|'
}
