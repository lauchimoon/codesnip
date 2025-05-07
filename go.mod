module github.com/lauchimoon/codesnip

go 1.24.1

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/lauchimoon/codesnip/lexer v0.0.0-00010101000000-000000000000
)

require golang.org/x/image v0.26.0

replace github.com/lauchimoon/codesnip/lexer => ./lexer
