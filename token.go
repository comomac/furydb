package furydb

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // main
	VALUE // sql value; 'bla', 1, 1.23

	// Misc characters
	ASTERISK  // *
	COMMA     // ,
	LEFTPAR   // (
	RIGHTPAR  // )
	SINGLEQUO // '
	DOUBLEQUO // "
	SEMICOL   // ;

	// Keywords
	CREATE
	TABLE
	INSERT
	INTO
	VALUES

	// Column Types
	BOOL
	INT
	FLOAT
	STRING
	TIME
	BYTES
	UUID

	// Constraints
	NOTNULL
	PRIMARYKEY
	FOREIGNKEY
)
