module github.com/bored-engineer/git-to-parquet

go 1.26.2

require (
	github.com/bored-engineer/git-fastpack v0.0.0-20260524064109-53c50c6052db
	github.com/edsrzf/mmap-go v1.2.0
	github.com/parquet-go/parquet-go v0.30.1
)

require (
	github.com/4kills/go-libdeflate/v2 v2.2.2 // indirect
	github.com/andybalholm/brotli v1.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/parquet-go/bitpack v1.0.0 // indirect
	github.com/parquet-go/jsonlite v1.5.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.26 // indirect
	github.com/twpayne/go-geom v1.6.1 // indirect
	golang.org/x/sys v0.45.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/go-git/go-git/v5 => github.com/bored-engineer/go-git/v5 v5.0.0-20250617195532-6922309b08c0
