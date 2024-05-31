// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nodekit

import (
	"bytes"
	"errors"
	"fmt"
	"math"

	"github.com/mr-tron/base58/base58"
)

const checksumLen = 4

var (
	ErrBase58Decoding   = errors.New("base58 decoding error")
	ErrMissingChecksum  = errors.New("input string is smaller than the checksum size")
	ErrBadChecksum      = errors.New("invalid input checksum")
	errEncodingOverFlow = errors.New("encoding overflow")
)

// Encode [bytes] to a string using cb58 format.
// [bytes] may be nil, in which case it will be treated the same as an empty
// slice.
func EncodeCB58(bytes []byte) (string, error) {
	bytesLen := len(bytes)
	if bytesLen > math.MaxInt32-checksumLen {
		return "", errEncodingOverFlow
	}
	checked := make([]byte, bytesLen+checksumLen)
	copy(checked, bytes)
	copy(checked[len(bytes):], Checksum(bytes, checksumLen))
	return base58.Encode(checked), nil
}

// Decode [str] to bytes from cb58.
func DecodeCB58(str string) ([]byte, error) {
	decodedBytes, err := base58.Decode(str)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBase58Decoding, err)
	}
	if len(decodedBytes) < checksumLen {
		return nil, ErrMissingChecksum
	}
	// Verify the checksum
	rawBytes := decodedBytes[:len(decodedBytes)-checksumLen]
	checksum := decodedBytes[len(decodedBytes)-checksumLen:]
	if !bytes.Equal(checksum, Checksum(rawBytes, checksumLen)) {
		return nil, ErrBadChecksum
	}
	return rawBytes, nil
}