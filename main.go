package main

import (
	"flag"
	"fmt"
	"os"
)

// errPrint prints a formatted message to stderr and exits non-zero.
func errPrint(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+msg, args...)
	os.Exit(1)
}

// infoPrint prints a formatted message to stderr but does not exit.
func infoPrint(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "INFO: "+msg, args...)
}

func main() {
	dbFile := flag.String("db", "OBJ_DUMP.sqlite", "SQLite Database file path to populate")
	forceFlag := flag.Bool("force", false, "Add OBJ_DUMP data to an existing database")
	flag.Parse()

	if _, err := os.Stat(*dbFile); err == nil && !(*forceFlag) {
		errPrint(
			"db %q exists and -force was not specified. Remove the database or add -force\n",
			*dbFile)
	}

	// Open the database.
	db, err := openDB(*dbFile)
	if err != nil {
		errPrint("opening db: %v\n", err)
	}

	infoPrint("Using database file %q\n", *dbFile)

	defer func() { _ = db.Close() }()

	// Open the OBJ_DUMP file.
	rest := flag.Args()
	if len(rest) != 1 {
		errPrint("usage %s [flags] <OBJ_DUMP file path>\n", os.Args[0])
	}

	objDumpFileName := rest[0]
	infoPrint("Reading OBJ_DUMP from %q\n", objDumpFileName)

	objDumpFile, err := os.Open(objDumpFileName)
	if err != nil {
		errPrint("%v\n", err)
	}

	// Backfill the DB by parsing the OBJ_DUMP file.
	count, err := backfill(objDumpFile, db)
	if err != nil {
		errPrint("backfilling db: %v\n", err)
	}

	infoPrint("Loaded %d objects from %q into DB %q\n", count, objDumpFileName, *dbFile)
}
