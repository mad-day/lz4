# lz4
LZ4 compression and decompression in pure Go.

This library is derived from [github.com/pierrec/lz4](https://github.com/pierrec/lz4/tree/revert-8-master), namely the `revert-8-master` branch.
Unlike pierrec's version, this one exposes Low level routines only, that are meant to be used directly.

The goal of this fork is, to retain the Low level interface exposed by the `revert-8-master` version of `github.com/pierrec/lz4` in order to
allow the implementation of efficient delta-compression techniques.

### Additional code.

The package contains two further functions, namely `SplitCompressedBlock` and `MergeCompressedBlock`,
which split a compressed block into literal and non-literal data. The code is derived from the
function `UncompressBlock`.

