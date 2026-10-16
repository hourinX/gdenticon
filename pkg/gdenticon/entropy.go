// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"crypto/sha256"
	"encoding/binary"
)

// entropy reads bits from a digest, most significant bit first.
type entropy struct {
	seed    []byte
	block   []byte
	bit     int
	counter uint64
}

func newEntropy(seed []byte) *entropy {
	return &entropy{seed: seed, block: seed}
}

// Bool returns the next bit of the stream.
func (e *entropy) Bool() bool {
	if e.bit == len(e.block)*8 {
		e.extend()
	}
	b := e.block[e.bit/8]
	shift := 7 - uint(e.bit%8)
	e.bit++
	return b>>shift&1 == 1
}

// extend replaces the block with the next counter-mode one, so the stream never
// runs dry however many bits a grid asks for.
func (e *entropy) extend() {
	e.counter++
	var c [8]byte
	binary.BigEndian.PutUint64(c[:], e.counter)

	h := sha256.New()
	h.Write(e.seed)
	h.Write(c[:])
	e.block = h.Sum(nil)
	e.bit = 0
}
