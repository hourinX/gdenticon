// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"hash"
)

const DefaultDigest = "md5" // digest used when a caller names none

// MD5 returns the MD5 digest of id.
func MD5(id string) []byte { return sum(md5.New(), id) }

// SHA1 returns the SHA-1 digest of id.
func SHA1(id string) []byte { return sum(sha1.New(), id) }

// SHA224 returns the SHA-224 digest of id.
func SHA224(id string) []byte { return sum(sha256.New224(), id) }

// SHA256 returns the SHA-256 digest of id.
func SHA256(id string) []byte { return sum(sha256.New(), id) }

// SHA384 returns the SHA-384 digest of id.
func SHA384(id string) []byte { return sum(sha512.New384(), id) }

// SHA512 returns the SHA-512 digest of id.
func SHA512(id string) []byte { return sum(sha512.New(), id) }

// SHA512_224 returns the SHA-512/224 digest of id.
func SHA512_224(id string) []byte { return sum(sha512.New512_224(), id) }

// SHA512_256 returns the SHA-512/256 digest of id.
func SHA512_256(id string) []byte { return sum(sha512.New512_256(), id) }

// SHA3_224 returns the SHA3-224 digest of id.
func SHA3_224(id string) []byte { return sum(sha3.New224(), id) }

// SHA3_256 returns the SHA3-256 digest of id.
func SHA3_256(id string) []byte { return sum(sha3.New256(), id) }

// SHA3_384 returns the SHA3-384 digest of id.
func SHA3_384(id string) []byte { return sum(sha3.New384(), id) }

// SHA3_512 returns the SHA3-512 digest of id.
func SHA3_512(id string) []byte { return sum(sha3.New512(), id) }

var digests = map[string]func(string) []byte{
	DefaultDigest: MD5,
	"sha1":        SHA1,
	"sha224":      SHA224,
	"sha256":      SHA256,
	"sha384":      SHA384,
	"sha512":      SHA512,
	"sha512-224":  SHA512_224,
	"sha512-256":  SHA512_256,
	"sha3-224":    SHA3_224,
	"sha3-256":    SHA3_256,
	"sha3-384":    SHA3_384,
	"sha3-512":    SHA3_512,
}

// Digest returns the digest function registered under name.
func Digest(name string) (func(string) []byte, bool) {
	fn, ok := digests[name]
	return fn, ok
}

// sum writes id into h and returns the digest it holds afterwards.
func sum(h hash.Hash, id string) []byte {
	h.Write([]byte(id))
	return h.Sum(nil)
}
