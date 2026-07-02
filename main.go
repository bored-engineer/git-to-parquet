package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	fastpack "github.com/bored-engineer/git-fastpack"
	git "github.com/bored-engineer/git-to-parquet/pkg"
	"github.com/edsrzf/mmap-go"
	"github.com/parquet-go/parquet-go"
	flag "github.com/spf13/pflag"
)

func main() {
	zstd := flag.Bool("zstd", false, "compress the parquet files with zstd")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [--zstd] <packfile>...", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	for _, filename := range args {
		if err := run(filename, *zstd); err != nil {
			log.Fatalf("%s:%s", filename, err)
		}
	}
}

func run(filename string, zstd bool) (rerr error) {
	scanner, err := fastpack.NewScanner(10000)
	if err != nil {
		return fmt.Errorf("fastpack.New failed: %w", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("os.Open for %q failed: %w", filename, err)
	}
	defer func() {
		if err := file.Close(); err != nil && rerr == nil {
			rerr = fmt.Errorf("(*os.File).Close for %q failed: %w", filename, err)
		}
	}()

	packfile, err := mmap.Map(file, mmap.RDONLY, 0)
	if err != nil {
		return fmt.Errorf("mmap.Map for %q failed: %w", filename, err)
	}
	defer func() {
		if err := packfile.Unmap(); err != nil && rerr == nil {
			rerr = fmt.Errorf("(mmap.MMap).Unmap for %q failed: %w", filename, err)
		}
	}()
	scanner.Reset(packfile)

	_, objects, err := scanner.Header()
	if err != nil {
		return fmt.Errorf("(*fastpack.Scanner).Header failed: %w", err)
	}
	checksum, err := scanner.Trailer()
	if err != nil {
		return fmt.Errorf("(*fastpack.Scanner).Trailer failed: %w", err)
	}

	writerOptions := func(bufferPattern string) []parquet.WriterOption {
		opts := []parquet.WriterOption{
			parquet.ColumnPageBuffers(
				parquet.NewFileBufferPool("", bufferPattern),
			),
		}
		if zstd {
			opts = append(opts, parquet.Compression(&parquet.Zstd))
		}
		return opts
	}

	blobFile, err := os.CreateTemp("", "blob-*.parquet")
	if err != nil {
		return fmt.Errorf("os.CreateTemp failed: %w", err)
	}
	defer func() {
		if err := blobFile.Close(); err != nil && rerr == nil {
			rerr = fmt.Errorf("(*os.File).Close failed: %w", err)
		}
	}()
	defer os.Remove(blobFile.Name())

	blobWriter := parquet.NewGenericWriter[git.Blob](blobFile, writerOptions("blob-*.buffer")...)
	defer func() {
		if err := blobWriter.Close(); err != nil && rerr == nil {
			rerr = fmt.Errorf("(*parquet.GenericWriter[git.Blob]).Close failed: %w", err)
			return
		}
		if rerr == nil {
			if err := os.Rename(blobFile.Name(), fmt.Sprintf("blob-%x.parquet", checksum)); err != nil {
				rerr = fmt.Errorf("os.Rename for %q failed: %w", blobFile.Name(), err)
			}
		}
	}()

	commitFile, err := os.CreateTemp("", "commit-*.parquet")
	if err != nil {
		return fmt.Errorf("os.CreateTemp failed: %w", err)
	}
	defer func() {
		if err := commitFile.Close(); err != nil && rerr == nil {
			rerr = fmt.Errorf("(*os.File).Close failed: %w", err)
		}
	}()
	defer os.Remove(commitFile.Name())

	commitWriter := parquet.NewGenericWriter[git.Commit](commitFile, writerOptions("commit-*.buffer")...)
	defer func() {
		if err := commitWriter.Close(); err != nil && rerr == nil {
			rerr = fmt.Errorf("(*parquet.GenericWriter[git.Commit]).Close failed: %w", err)
			return
		}
		if rerr == nil {
			if err := os.Rename(commitFile.Name(), fmt.Sprintf("commit-%x.parquet", checksum)); err != nil {
				rerr = fmt.Errorf("os.Rename for %q failed: %w", commitFile.Name(), err)
			}
		}
	}()

	tagFile, err := os.CreateTemp("", "tag-*.parquet")
	if err != nil {
		return fmt.Errorf("os.CreateTemp failed: %w", err)
	}
	defer func() {
		if err := tagFile.Close(); err != nil && rerr == nil {
			rerr = fmt.Errorf("(*os.File).Close failed: %w", err)
		}
	}()
	defer os.Remove(tagFile.Name())

	tagWriter := parquet.NewGenericWriter[git.Tag](tagFile, writerOptions("tag-*.buffer")...)
	defer func() {
		if err := tagWriter.Close(); err != nil && rerr == nil {
			rerr = fmt.Errorf("(*parquet.GenericWriter[git.Tag]).Close failed: %w", err)
			return
		}
		if rerr == nil {
			if err := os.Rename(tagFile.Name(), fmt.Sprintf("tag-%x.parquet", checksum)); err != nil {
				rerr = fmt.Errorf("os.Rename for %q failed: %w", tagFile.Name(), err)
			}
		}
	}()

	treeFile, err := os.CreateTemp("", "tree-*.parquet")
	if err != nil {
		return fmt.Errorf("os.CreateTemp failed: %w", err)
	}
	defer func() {
		if err := treeFile.Close(); err != nil && rerr == nil {
			rerr = fmt.Errorf("(*os.File).Close failed: %w", err)
		}
	}()
	defer os.Remove(treeFile.Name())

	treeWriter := parquet.NewGenericWriter[git.Tree](treeFile, writerOptions("tree-*.buffer")...)
	defer func() {
		if err := treeWriter.Close(); err != nil && rerr == nil {
			rerr = fmt.Errorf("(*parquet.GenericWriter[git.Tree]).Close failed: %w", err)
			return
		}
		if rerr == nil {
			if err := os.Rename(treeFile.Name(), fmt.Sprintf("tree-%x.parquet", checksum)); err != nil {
				rerr = fmt.Errorf("os.Rename for %q failed: %w", treeFile.Name(), err)
			}
		}
	}()

	for range objects {
		ot, buf, err := scanner.Object()
		if err != nil {
			return fmt.Errorf("(*fastpack.Scanner).Object failed: %w", err)
		}
		oid := fastpack.OID(ot, buf)
		switch ot {
		case fastpack.CommitObject:
			commit := git.Commit{OID: oid}
			if err := commit.UnmarshalBinary(buf); err != nil {
				return fmt.Errorf("(*git.Commit).UnmarshalBinary for %q failed: %w", oid, err)
			}
			if _, err := commitWriter.Write([]git.Commit{commit}); err != nil {
				return fmt.Errorf("(*parquet.GenericWriter[git.Commit]).Write for %q failed: %w", oid, err)
			}
		case fastpack.TreeObject:
			tree := git.Tree{OID: oid}
			if err := tree.UnmarshalBinary(buf); err != nil {
				return fmt.Errorf("(*git.Tree).UnmarshalBinary for %q failed: %w", oid, err)
			}
			if _, err := treeWriter.Write([]git.Tree{tree}); err != nil {
				return fmt.Errorf("(*parquet.GenericWriter[git.Tree]).Write for %q failed: %w", oid, err)
			}
		case fastpack.TagObject:
			tag := git.Tag{OID: oid}
			if err := tag.UnmarshalBinary(buf); err != nil {
				return fmt.Errorf("(*git.Tag).UnmarshalBinary for %q failed: %w", oid, err)
			}
			if _, err := tagWriter.Write([]git.Tag{tag}); err != nil {
				return fmt.Errorf("(*parquet.GenericWriter[git.Tag]).Write for %q failed: %w", oid, err)
			}
		case fastpack.BlobObject:
			blob := git.Blob{OID: oid}
			if err := blob.UnmarshalBinary(buf); err != nil {
				return fmt.Errorf("(*git.Blob).UnmarshalBinary for %q failed: %w", oid, err)
			}
			blob.Contents = nil
			if _, err := blobWriter.Write([]git.Blob{blob}); err != nil {
				return fmt.Errorf("(*parquet.GenericWriter[git.Blob]).Write for %q failed: %w", oid, err)
			}
		default:
			return fmt.Errorf("unknown object type: %s", ot)
		}
	}

	return nil
}
