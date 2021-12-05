package main

import "github.com/alecthomas/participle/v2/lexer"

// objDumpLexer describes the lexical elements in the OBJ_DUMP file.
// nolint: govet
var objDumpLexer = lexer.MustSimple([]lexer.Rule{
	// Filenames start with a letter, a slash, or on Dune, an @ for virtual stuff.
	{`Filename`, `[a-zA-Z/@][^\s]*`, nil},
	// Datestamps in the OBJ_DUMP are YYYY:MM:DD-HH:MM:SS
	{`Date`, `\d{4}\.\d{2}\.\d{2}-\d{2}\:\d{2}\:\d{2}`, nil},
	// Integers are just repeated digits. No floats or negative values here.
	{"Int", `\d+`, nil},
	// Punctuation is parens and the "--" used to indicate no environment.
	{`Punct`, `[()]|--`, nil},
	// Whitespace is ... whitespace.
	{"whitespace", `\s+`, nil},
})

// ObjectDump instances are zero or more Objects.
// nolint: govet
type ObjectDump struct {
	Objects []*Object `@@*`
}

// Object instances describe traits of a dumped object.
// See `man dump_driver_info` and the description of DDI_OBJECTS for more.
// nolint: govet
type Object struct {
	// Name is the object's name.
	Name string ` @Filename `
	// Basefile is Name with any trailing "#xxxxx" object reference stripped.
	// It has no grammar because it isn't present in OBJ_DUMP, we post-process it
	// from the 'Name' that is.
	Basefile string
	// Size in memory, shared data counted only once.
	Size int ` @Int `
	// Size in memory if data wouldn't be shared.
	FullSize int ` "(" @Int ")" `
	// References is the count of references to the object.
	References int ` "ref" @Int`
	// HB is true if the object has a heartbeat.
	HB bool ` @"HB"? `
	// Environment contains the name of the objects environment (if it had one).
	Environment string ` ("--" | @Filename) `
	// Ticks is the number of execution ticks spent in this object.
	Ticks int ` "(" @Int ")"`
	// SwapStatus is a string describing the swap status (if any).
	SwapStatus string ` @("PROG" "SWAPPED"|"VAR" "SWAPPED"|"SWAPPED")? `
	// The time the object was created.
	Created string ` @Date `
}
