package lz4

// #cgo CFLAGS: -O3
// #include "src/lz4hc.h"
// #include "src/lz4hc.c"
import "C"

import (
	"fmt"
)

// CompressHC compresses in and puts the content in out. len(out)
// should have enough space for the compressed data (use CompressBound
// to calculate). Returns the number of bytes in the out slice. Determines
// the compression level automatically.
func CompressHC(in, out []byte) (int, error) {
	// 0 automatically sets the compression level.
	return CompressHCLevel(in, out, 0)
}

// CompressHCLevel compresses in at the given compression level and puts the
// content in out. len(out) should have enough space for the compressed data
// (use CompressBound to calculate). Returns the number of bytes in the out
// slice. To automatically choose the compression level, use 0. Otherwise, use
// any value in the inclusive range 1 (worst) through 16 (best). Most
// applications will prefer CompressHC.
func CompressHCLevel(in, out []byte, level int) (outSize int, err error) {
	// LZ4HC does not handle empty buffers. Pass through to Compress.
	if len(in) == 0 || len(out) == 0 {
		return Compress(in, out)
	}

	outSize = int(C.LZ4_compressHC2_limitedOutput(p(in), p(out), clen(in), clen(out), C.int(level)))
	if outSize == 0 {
		err = fmt.Errorf("insufficient space for compression")
	}
	return
}
