package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	fastpack "github.com/bored-engineer/git-fastpack"
	git "github.com/bored-engineer/git-to-parquet/pkg"
	"github.com/parquet-go/parquet-go"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <packfile>", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	scanner, err := fastpack.NewScanner(10000)
	if err != nil {
		log.Fatalf("fastpack.New failed: %v", err)
	}

	packfile, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("os.ReadFile for %q failed: %v", os.Args[1], err)
	}
	scanner.Reset(packfile)

	_, objects, err := scanner.Header()
	if err != nil {
		log.Fatalf("(*fastpack.Scanner).Header failed: %v", err)
	}
	checksum, err := scanner.Trailer()
	if err != nil {
		log.Fatalf("(*fastpack.Scanner).Trailer failed: %v", err)
	}

	blobFile, err := os.CreateTemp("", "blob-*.parquet")
	if err != nil {
		log.Fatalf("os.CreateTemp failed: %v", err)
	}
	defer func() {
		if err := blobFile.Close(); err != nil {
			log.Fatalf("(*os.File).Close failed: %v", err)
		}
	}()
	defer os.Remove(blobFile.Name())

	blobWriter := parquet.NewGenericWriter[git.Blob](blobFile,
		parquet.ColumnPageBuffers(
			parquet.NewFileBufferPool("", "blob-*.buffer"),
		),
	)
	defer func() {
		if err := blobWriter.Close(); err != nil {
			log.Fatalf("(*parquet.GenericWriter[git.Blob]).Close failed: %v", err)
		}
		if err := os.Rename(blobFile.Name(), fmt.Sprintf("blob-%x.parquet", checksum)); err != nil {
			log.Fatalf("os.Rename for %q failed: %v", blobFile.Name(), err)
		}
	}()

	commitFile, err := os.CreateTemp("", "commit-*.parquet")
	if err != nil {
		log.Fatalf("os.CreateTemp failed: %v", err)
	}
	defer func() {
		if err := commitFile.Close(); err != nil {
			log.Fatalf("(*os.File).Close failed: %v", err)
		}
	}()
	defer os.Remove(commitFile.Name())

	commitWriter := parquet.NewGenericWriter[git.Commit](commitFile,
		parquet.ColumnPageBuffers(
			parquet.NewFileBufferPool("", "commit-*.buffer"),
		),
	)
	defer func() {
		if err := commitWriter.Close(); err != nil {
			log.Fatalf("(*parquet.GenericWriter[git.Commit]).Close failed: %v", err)
		}
		if err := os.Rename(commitFile.Name(), fmt.Sprintf("commit-%x.parquet", checksum)); err != nil {
			log.Fatalf("os.Rename for %q failed: %v", commitFile.Name(), err)
		}
	}()

	tagFile, err := os.CreateTemp("", "tag-*.parquet")
	if err != nil {
		log.Fatalf("os.CreateTemp failed: %v", err)
	}
	defer func() {
		if err := tagFile.Close(); err != nil {
			log.Fatalf("(*os.File).Close failed: %v", err)
		}
	}()
	defer os.Remove(tagFile.Name())

	tagWriter := parquet.NewGenericWriter[git.Tag](tagFile,
		parquet.ColumnPageBuffers(
			parquet.NewFileBufferPool("", "tag-*.buffer"),
		),
	)
	defer func() {
		if err := tagWriter.Close(); err != nil {
			log.Fatalf("(*parquet.GenericWriter[git.Tag]).Close failed: %v", err)
		}
		if err := os.Rename(tagFile.Name(), fmt.Sprintf("tag-%x.parquet", checksum)); err != nil {
			log.Fatalf("os.Rename for %q failed: %v", tagFile.Name(), err)
		}
	}()

	treeFile, err := os.CreateTemp("", "tree-*.parquet")
	if err != nil {
		log.Fatalf("os.CreateTemp failed: %v", err)
	}
	defer func() {
		if err := treeFile.Close(); err != nil {
			log.Fatalf("(*os.File).Close failed: %v", err)
		}
	}()
	defer os.Remove(treeFile.Name())

	treeWriter := parquet.NewGenericWriter[git.Tree](treeFile,
		parquet.ColumnPageBuffers(
			parquet.NewFileBufferPool("", "tree-*.buffer"),
		),
	)
	defer func() {
		if err := treeWriter.Close(); err != nil {
			log.Fatalf("(*parquet.GenericWriter[git.Tree]).Close failed: %v", err)
		}
		if err := os.Rename(treeFile.Name(), fmt.Sprintf("tree-%x.parquet", checksum)); err != nil {
			log.Fatalf("os.Rename for %q failed: %v", treeFile.Name(), err)
		}
	}()

	for range objects {
		ot, buf, err := scanner.Object()
		if err != nil {
			log.Fatalf("(*fastpack.Scanner).Object failed: %v", err)
		}
		oid := fastpack.OID(ot, buf)
		switch ot {
		case fastpack.CommitObject:
			commit := git.Commit{OID: oid}
			if err := commit.UnmarshalBinary(buf); err != nil {
				log.Fatalf("(*git.Commit).UnmarshalBinary for %q failed: %v", oid, err)
			}
			if _, err := commitWriter.Write([]git.Commit{commit}); err != nil {
				log.Fatalf("(*parquet.GenericWriter[git.Commit]).Write for %q failed: %v", oid, err)
			}
		case fastpack.TreeObject:
			tree := git.Tree{OID: oid}
			if err := tree.UnmarshalBinary(buf); err != nil {
				log.Fatalf("(*git.Tree).UnmarshalBinary for %q failed: %v", oid, err)
			}
			if _, err := treeWriter.Write([]git.Tree{tree}); err != nil {
				log.Fatalf("(*parquet.GenericWriter[git.Tree]).Write for %q failed: %v", oid, err)
			}
		case fastpack.TagObject:
			tag := git.Tag{OID: oid}
			if err := tag.UnmarshalBinary(buf); err != nil {
				log.Fatalf("(*git.Tag).UnmarshalBinary for %q failed: %v", oid, err)
			}
			if _, err := tagWriter.Write([]git.Tag{tag}); err != nil {
				log.Fatalf("(*parquet.GenericWriter[git.Tag]).Write for %q failed: %v", oid, err)
			}
		case fastpack.BlobObject:
			blob := git.Blob{OID: oid}
			if err := blob.UnmarshalBinary(buf); err != nil {
				log.Fatalf("(*git.Blob).UnmarshalBinary for %q failed: %v", oid, err)
			}
			blob.Contents = nil
			if _, err := blobWriter.Write([]git.Blob{blob}); err != nil {
				log.Fatalf("(*parquet.GenericWriter[git.Blob]).Write for %q failed: %v", oid, err)
			}
		default:
			log.Fatalf("unknown object type: %s", ot)
		}
	}

}
