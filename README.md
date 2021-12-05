# DUMPDB: A LDMUD OBJ_DUMP -> SQLite3 DB Tool

## About

The [LDMUD][ldmud] MUD driver supports dumping a formatted text file with
information about every object loaded at the time of the dump. This is an
invaluable tool for finding object leaks and understanding your game's memory
usage.

The format of `OBJ_DUMP` can be challenging to work with using shell tools as
some fields are optional and may be omitted. Having the contents parsed and
represented in a small database makes the data easier to explore. Use `dumpdb`
to answer questions like:

* What is the most cloned object in memory?
* What is the largest memory usage by base file?
* What object has the largest tick count?
* What base file has the largest tick count?

[ldmud]: https://github.com/ldmud/ldmud

## Usage

1. Install `dumpdb`.
2. Install the `sqlite3` command line tools for your OS:
```bash
# For example...
apt-get install sqlite3
```
3. In-game, collect an `OBJ_DUMP` by calling:
```c
efun::dump_driver_info(DDI_OBJECTS);
```
4. Out of game, run `dumpdb` on your `OBJ_DUMP` file:
```bash
dumpdb /path/to/my/game/lib/OBJ_DUMP
```

Processing a 5.5MB OBJ_DUMP with ~47,000 lines takes approximately 6s and
creates a 6.5MB SQLite database.

5. Explore the db with `sqlite3 OBJ_DUMP.sqlite`:
```bash
sqlite3 OBJ_DUMP.sqlite
```

```sql
.tables
.schema obj_dump
SELECT COUNT(id) FROM obj_dump;
```

Tested with `OBJ_DUMP`'s collected from [LDMud][ldmud] 3.6.4 and 3.5.4.

## Installation (Linux x86_64)

If you're on 64bit Linux you can download a pre-built executable. Other
platforms should build from source.

1. Download [a pre-built release](https://github.com/dune-mud/dumpdb/releases).
2. Extract the `.tar.gz`.
```bash
tar -xvf dumpdb_*_linux_amd64.tar.gz
```
3. Run the `dumpdb` executable.
```bash
./dumpdb /path/to/my/game/lib/OBJ_DUMP
```

## Installation from Source

* Set up Go 1.16+ by following the [Go install docs][golang].
* You will also need a C compiler to support the CGo sqlite dependency. E.g.
```bash
apt-get install gcc
```
* Clone the source code from git and change to the directory:
```bash
git clone https://github.com/dune-mud/dumpdb.git && cd dumpdb
```
* Build and install
```bash
go install -ldflags "-X main.commit=$(git rev-parse HEAD)" ./...
```
* Run the `dumpdb` executable from your Go bin.
```bash
$GOPATH/bin/dumpdb /path/to/my/game/lib/OBJ_DUMP
```

[golang]: https://go.dev/doc/install

## Database Schema

The `obj_dump` table that is created in the SQLite3 database has the following
schema:

```sql
CREATE TABLE IF NOT EXISTS obj_dump (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  name         TEXT NOT NULL,
  base_file    TEXT NOT NULL,
  size         INTEGER NOT NULL,
  full_size    INTEGER NOT NULL,
  refs         INTEGER NOT NULL DEFAULT 0,
  hb           INTEGER NOT NULL DEFAULT 0,
  environment  TEXT NOT NULL DEFAULT "",
  ticks        INTEGER NOT NULL DEFAULT 0,
  swap_status  TEXT NOT NULL DEFAULT "",
  created      TEXT NOT NULL
);
```

## OBJ_DUMP format

The `OBJ_DUMP` format is described in [`man dump_driver_info`][dump_driver_info]:

> For every object, a line is written into the file with the following
> information in the given order:
>  - object name
>  - size in memory, shared data counted only once
>  - size in memory if data wouldn't be shared
>  - number of references
>  - 'HB' if the object has a heartbeat, nothing if not.
>  - the name of the environment, or '--' if the object has no environment
>  - in parentheses the number of execution ticks spent in this object
>  - the swap status: nothing if not swapped, 'PROG SWAPPED' if only the program
>    is swapped, 'VAR SWAPPED' if only the variables are swapped, 'SWAPPED' if
>    both program and variables are swapped
>  - the time the object was created

[dump_driver_info]: https://github.com/ldmud/ldmud/blob/master/doc/efun/dump_driver_info

## Handy Queries

* Total number of objects.
```sql
SELECT COUNT(id) 
  FROM obj_dump;
```

* Total number of base files (E.g. `name` with the "#xxxxxx" clone reference
  removed).
```sql
SELECT COUNT(DISTINCT base_file)
  FROM obj_dump;
```

* Top 25 largest users of memory by object.
```sql
SELECT name, size 
  FROM obj_dump 
  ORDER BY size DESC 
  LIMIT 25;
```

* Top 25 largest users of memory by base file.
```sql
SELECT base_file, SUM(size) AS usage
  FROM obj_dump
  GROUP BY base_file
  ORDER BY usage DESC
  LIMIT 25;
```

* Top 25 most cloned base files.
```sql
SELECT base_file, COUNT(id) AS clones 
  FROM obj_dump
  GROUP BY base_file
  ORDER BY clones DESC
  LIMIT 25;
```

* Top 25 largest tick counts by object:
```sql
SELECT name, ticks
  FROM obj_dump
  ORDER BY ticks DESC
  LIMIT 25;
```

* Top 25 largest tick counts by base file:
```sql
SELECT base_file, SUM(ticks) AS usage
  FROM obj_dump
  GROUP BY base_file
  ORDER BY usage DESC
  LIMIT 25;
```

* Top 25 objects by reference count:
```sql
SELECT name, refs
  FROM obj_dump
  ORDER BY refs DESC
  LIMIT 25;
```

* Top 25 base files by reference count:
```sql
SELECT base_file, SUM(refs) AS refs
  FROM obj_dump
  GROUP BY base_file
  ORDER BY refs DESC
  LIMIT 25;
```

* Number of clones from a particular directory:
```sql
SELECT COUNT(id) 
  FROM obj_dump 
  WHERE name LIKE 'd/Ix/paradox/spire/%';
```
