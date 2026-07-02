# git-to-parquet [![Go Reference](https://pkg.go.dev/badge/github.com/bored-engineer/git-to-parquet.svg)](https://pkg.go.dev/github.com/bored-engineer/git-to-parquet)
Export every object in a git repository as parquet files for data analysis

`git-to-parquet` reads a single git packfile and writes its objects out as four parquet files (`blob-<checksum>.parquet`, `commit-<checksum>.parquet`, `tag-<checksum>.parquet`, `tree-<checksum>.parquet`), one row per object, ready to be queried with tools like [DuckDB](https://duckdb.org/).

## Install
```console
go install github.com/bored-engineer/git-to-parquet@latest
```

## Usage
```console
git-to-parquet <packfile>...
```

## Example: clone a repository and query it with DuckDB

Clone just the `master` branch of a repository. A fresh clone always contains at least one packfile under `.git/objects/pack`:
```console
$ git clone --bare --single-branch --branch master https://github.com/octocat/Hello-World.git
$ ls Hello-World.git/objects/pack/*.pack
Hello-World.git/objects/pack/pack-4ffbafdd692b576efb2d6c2abfcb10e36e41e418.pack
```

Run `git-to-parquet` against the packfile, producing one parquet file per object type:
```console
$ git-to-parquet Hello-World.git/objects/pack/*.pack
$ ls *.parquet
blob-4ffbafdd692b576efb2d6c2abfcb10e36e41e418.parquet
commit-4ffbafdd692b576efb2d6c2abfcb10e36e41e418.parquet
tag-4ffbafdd692b576efb2d6c2abfcb10e36e41e418.parquet
tree-4ffbafdd692b576efb2d6c2abfcb10e36e41e418.parquet
```

Query the commit history with [DuckDB](https://duckdb.org/). Object IDs are stored as raw 20-byte SHA-1 hashes, so wrap them in `hex()` to get familiar git OID strings:
```console
$ duckdb -c "
  SELECT
    hex(oid) AS oid,
    author.name AS author,
    author.timestamp AS authored_at,
    message
  FROM 'commit-*.parquet'
  ORDER BY author.timestamp DESC;
"
┌──────────────────────────────────────────┬─────────────────────────┬──────────────────────────┬───────────────────────────────────────────────────────────────────────────────┐
│                   oid                    │         author          │       authored_at        │                                    message                                    │
│                 varchar                  │          blob           │ timestamp with time zone │                                     blob                                      │
├──────────────────────────────────────────┼─────────────────────────┼──────────────────────────┼───────────────────────────────────────────────────────────────────────────────┤
│ 7FD1A60B01F91B314F59955A4E4D4E80D8EDF11D │ The Octocat             │ 2012-03-06 17:06:50-06   │ Merge pull request #6 from Spaceghost/patch-1\x0A\x0ANew line at end of file. │
│ 762941318EE16E59DABBACB1B4049EEC22F0D303 │ Johnneylee Jack Rollins │ 2011-09-13 23:42:41-05   │ New line at end of file. --Signed off by Spaceghost                           │
│ 553C2077F0EDC3D5DC5D17262F6AA498E69D6F8E │ cameronmcefee           │ 2011-01-26 13:06:08-06   │ first commit                                                                  │
└──────────────────────────────────────────┴─────────────────────────┴──────────────────────────┴───────────────────────────────────────────────────────────────────────────────┘
```

Because every object type shares the same `oid` (and, for commits, `tree_oid`/`parents`) encoding, DuckDB can join across the commit, tree, and blob parquet files to walk history without ever needing a `git` checkout.
