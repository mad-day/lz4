/*
Copyright (c) 2019, Simon Schmidt
Copyright (c) 2015, Pierre Curto
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of xxHash nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/


package lz4


// SplitCompressedBlock splits the source buffer into the control bytes (d1) and literal (d2).
//
// The destination buffers can be preallocated by doing:
//   d1 := make([]byte,0,1<<16)
//   d2 := make([]byte,0,1<<16)
func SplitCompressedBlock(src []byte, d1,d2 *[]byte) (error) {
	di := 0
	si, sn := 0, len(src)
	if sn == 0 {
		return nil
	}
	osi := 0

	for {
		// literals and match lengths (token)
		lLen := int(src[si] >> 4)
		mLen := int(src[si] & 0xF)
		if si++; si == sn {
			return  ErrInvalidSource
		}

		// literals
		if lLen > 0 {
			if lLen == 0xF {
				for src[si] == 0xFF {
					lLen += 0xFF
					if si++; si == sn {
						return ErrInvalidSource
					}
				}
				lLen += int(src[si])
				if si++; si == sn {
					return ErrInvalidSource
				}
			}
			if si+lLen > sn {
				return ErrShortBuffer
			}
			*d1 = append(*d1,src[osi:si]...) // Copy Control Data
			osi = si
			*d2 = append(*d2,src[si:si+lLen]...) // Copy Literal Data.
			di += lLen
			
			osi = si+lLen
			if si += lLen; si >= sn {
				return nil
			}
			
		}

		if si += 2; si >= sn {
			return ErrInvalidSource
		}
		offset := int(src[si-2]) | int(src[si-1])<<8
		if offset == 0 {
			return ErrInvalidSource
		}

		// match
		if mLen == 0xF {
			for src[si] == 0xFF {
				mLen += 0xFF
				if si++; si == sn {
					return ErrInvalidSource
				}
			}
			mLen += int(src[si])
			if si++; si == sn {
				return ErrInvalidSource
			}
		}
		// minimum match length is 4
		mLen += 4
	}
}

// MergeCompressedBlock splits the source buffers into the destination buffer.
//
// The destination buffer can be preallocated by doing:
//   d1 := make([]byte,0,1<<16)
func MergeCompressedBlock(src, lit []byte, d1 *[]byte) (error) {
	di := 0
	si, sn := 0, len(src)
	if sn == 0 {
		return nil
	}
	osi,li,ln := 0,0,len(lit)

	for {
		// literals and match lengths (token)
		lLen := int(src[si] >> 4)
		mLen := int(src[si] & 0xF)
		if si++; si == sn {
			if li<ln {
				*d1 = append(*d1,src[osi:si]...) // Copy Control Data
				*d1 = append(*d1,lit[li:]...) // Copy Literal Data.
				return nil
			}
			return ErrInvalidSource
		}

		// literals
		if lLen > 0 {
			if lLen == 0xF {
				for src[si] == 0xFF {
					lLen += 0xFF
					if si++; si == sn {
						return ErrInvalidSource
					}
				}
				lLen += int(src[si])
				if si++; si == sn {
					return ErrInvalidSource
				}
			}
			if li+lLen > ln {
				return ErrShortBuffer
			}
			*d1 = append(*d1,src[osi:si]...) // Copy Control Data
			osi = si
			*d1 = append(*d1,lit[li:li+lLen]...) // Copy Literal Data.
			di += lLen
			
			if li += lLen; li >= ln || si >= sn {
				return nil
			}
		}

		if si += 2; si >= sn {
			return ErrInvalidSource
		}
		offset := int(src[si-2]) | int(src[si-1])<<8
		if offset == 0 {
			return ErrInvalidSource
		}

		// match
		if mLen == 0xF {
			for src[si] == 0xFF {
				mLen += 0xFF
				if si++; si == sn {
					return ErrInvalidSource
				}
			}
			mLen += int(src[si])
			if si++; si == sn {
				return ErrInvalidSource
			}
		}
		// minimum match length is 4
		mLen += 4
	}
}


