package main

import (
	"fmt"
	"io"

	"github.com/alecthomas/participle/v2"
)

func backfill(objDump io.Reader, db *database) (int, error) {
	// Build a parser with our AST and lexer
	parser := participle.MustBuild(
		&ObjectDump{Objects: nil},
		participle.Lexer(objDumpLexer))

	ast := &ObjectDump{Objects: nil}

	infoPrint("Parsing OBJ_DUMP into AST\n")

	// Build the AST by parsing the contents of the OBJ_DUMP file.
	err := parser.Parse("", objDump, ast)
	if err != nil {
		return 0, fmt.Errorf("parsing: %w", err)
	}

	infoPrint("Inserting object info into DB\n")

	// Insert the objects from the AST to the DB
	err = db.Insert(ast.Objects)
	if err != nil {
		return 0, fmt.Errorf("inserting: %w", err)
	}

	infoPrint("Finished inserts.\n")

	return len(ast.Objects), nil
}
