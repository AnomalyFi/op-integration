package nodekit

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"
)

// TODO rewrite this with new base cases
func removeWhitespace(s string) string {
	// Split the string on whitespace then concatenate the segments
	return strings.Join(strings.Fields(s), "")
}

var ReferenceNmtRoot NmtRoot = NmtRoot{
	Root: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

var ReferenceL1BLockInfo L1BlockInfo = L1BlockInfo{
	Number:    123,
	Timestamp: *NewU256().SetUint64(0x456),
	Hash:      common.Hash{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef},
}

var ReferenceHeader Header = Header{
	Height:           42,
	Timestamp:        789,
	L1Head:           124,
	TransactionsRoot: ReferenceNmtRoot,
}

func TestNodeKitTypesNmtRootJson(t *testing.T) {
	data := []byte(removeWhitespace(`{
		"root": [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]
	}`))

	// Check encoding.
	encoded, err := json.Marshal(ReferenceNmtRoot)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	require.Equal(t, encoded, data)

	// Check decoding
	var decoded NmtRoot
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	require.Equal(t, decoded, ReferenceNmtRoot)

	CheckJsonRequiredFields[NmtRoot](t, data, "root")
}

func TestNodeKitTypesL1BLockInfoJson(t *testing.T) {
	data := []byte(removeWhitespace(`{
		"number": 123,
		"timestamp": "0x456",
		"hash": "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	}`))

	// Check encoding.
	encoded, err := json.Marshal(ReferenceL1BLockInfo)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	require.Equal(t, encoded, data)

	// Check decoding
	var decoded L1BlockInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	require.Equal(t, decoded, ReferenceL1BLockInfo)

	CheckJsonRequiredFields[L1BlockInfo](t, data, "number", "timestamp", "hash")
}

func TestNodeKitTypesHeaderJson(t *testing.T) {
	data := []byte(removeWhitespace(`{
		"height": 42,
		"timestamp": 789,
		"l1_head": 124,
		"transactions_root": {
			"root": [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]
		}
	}
	`))

	// Check encoding.
	encoded, err := json.Marshal(ReferenceHeader)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	require.Equal(t, encoded, data)

	// Check decoding
	var decoded Header
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	require.Equal(t, decoded, ReferenceHeader)

	CheckJsonRequiredFields[Header](t, data, "height", "timestamp", "l1_head", "transactions_root")
}

func TestNodeKitTransactionJson(t *testing.T) {
	data := []byte(removeWhitespace(`{
		"vm": "0",
		"payload": [1,2,3,4,5]
	}`))
	tx := Transaction{
		ChainId: "0",
		Data:    []byte{1, 2, 3, 4, 5},
	}

	// Check encoding.
	encoded, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	require.Equal(t, encoded, data)

	// Check decoding
	var decoded Transaction
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	require.Equal(t, decoded, tx)

	CheckJsonRequiredFields[Transaction](t, data, "vm", "payload")
}

// func TestNodeKitTypesNmtRootCommit(t *testing.T) {
// 	require.Equal(t, ReferenceNmtRoot.Commit(), Commitment{251, 80, 232, 195, 91, 2, 138, 18, 240, 231, 31, 172, 54, 204, 90, 42, 215, 42, 72, 187, 15, 28, 128, 67, 149, 117, 26, 114, 232, 57, 190, 10})
// }

// func TestNodeKitTypesL1BlockInfoCommit(t *testing.T) {
// 	require.Equal(t, ReferenceL1BLockInfo.Commit(), Commitment{224, 122, 115, 150, 226, 202, 216, 139, 51, 221, 23, 79, 54, 243, 107, 12, 12, 144, 113, 99, 133, 217, 207, 73, 120, 182, 115, 84, 210, 230, 126, 148})
// }

func TestNodeKitTypesHeaderCommit(t *testing.T) {
	require.Equal(t, ReferenceHeader.Commit(), Commitment{69, 70, 204, 173, 194, 81, 14, 28, 209, 104, 204, 53, 32, 43, 75, 233, 35, 99, 95, 128, 155, 14, 46, 17, 217, 191, 252, 217, 29, 252, 131, 61})
}

func TestNodeKitCommitmentNewConversion(t *testing.T) {

	comm_expected := Commitment{193, 98, 70, 80, 45, 4, 82, 113, 146, 158, 194, 61, 72, 64, 34, 217, 173, 46, 78, 63, 115, 159, 115, 122, 219, 58, 120, 223, 9, 52, 140, 166}

	root := [32]byte{88, 48, 50, 99, 140, 172, 117, 35, 116, 212, 26, 123, 187, 199, 189, 130, 55, 219, 55, 144, 86, 21, 30, 68, 214, 253, 157, 141, 160, 54, 5, 190}

	h := &Header{
		Height:    2539,
		Timestamp: 1703696824,
		L1Head:    252,
		TransactionsRoot: NmtRoot{
			Root: root[:],
		},
	}

	comm := h.Commit()

	require.Equal(t, comm, comm_expected)

}

func TestNodeKitCommitmentNewNewConversion(t *testing.T) {

	comm_expected := Commitment{246, 72, 71, 162, 203, 235, 120, 113, 123, 165, 56, 167, 19, 161, 196, 4, 180, 153, 56, 201, 83, 59, 235, 187, 93, 21, 26, 126, 35, 145, 94, 0}

	root := [32]byte{143, 94, 54, 190, 59, 183, 179, 98, 143, 189, 103, 102, 125, 111, 136, 122, 104, 120, 85, 160, 58, 205, 77, 206, 184, 91, 77, 23, 80, 68, 26, 35}

	h := &Header{
		Height:            840,
		TimestampOriginal: 1703778845527,
		L1Head:            109,
		TransactionsRoot: NmtRoot{
			Root: root[:],
		},
	}

	comm := h.Commit()

	require.Equal(t, comm, comm_expected)

}

func TestNodeKitCommitmentNewNewNewConversion(t *testing.T) {

	comm_expected := Commitment{246, 72, 71, 162, 203, 235, 120, 113, 123, 165, 56, 167, 19, 161, 196, 4, 180, 153, 56, 201, 83, 59, 235, 187, 93, 21, 26, 126, 35, 145, 94, 0}

	newRoot, _ := DecodeCB58("2GQrkp2R78AGQRPNue6oFBD4L2rBbWs7GgbFTgJ6CTrdAKJmKH")

	root := [32]byte{143, 94, 54, 190, 59, 183, 179, 98, 143, 189, 103, 102, 125, 111, 136, 122, 104, 120, 85, 160, 58, 205, 77, 206, 184, 91, 77, 23, 80, 68, 26, 35}

	require.Equal(t, newRoot, root)

	h := &Header{
		Height:            840,
		TimestampOriginal: 1703778845527,
		L1Head:            109,
		TransactionsRoot: NmtRoot{
			Root: newRoot[:],
		},
	}

	comm := h.Commit()

	require.Equal(t, comm, comm_expected)

}

// {2GQrkp2R78AGQRPNue6oFBD4L2rBbWs7GgbFTgJ6CTrdAKJmKH 1703778845527 109 840}

//2023/12/28 09:54:51 expected header header &{840 1703778845527 109 {[143 94 54 190 59 183 179 98 143 189 103 102 125 111 136 122 104 120 85 160 58 205 77 206 184 91 77 23 80 68 26 35]}}
//2023/12/28 09:54:51 expected comm comm {false [6707296349103152640 13013495035599776699 8909589728062456836 17746513096184920177]}
//2023/12/28 09:54:51 expected comm in bytes comm [246 72 71 162 203 235 120 113 123 165 56 167 19 161 196 4 180 153 56 201 83 59 235 187 93 21 26 126 35 145 94 0]

func TestNodeKitCommitmentFromU256TrailingZero(t *testing.T) {
	comm := Commitment{209, 146, 197, 195, 145, 148, 17, 211, 52, 72, 28, 120, 88, 182, 204, 206, 77, 36, 56, 35, 3, 143, 77, 186, 69, 233, 104, 30, 90, 105, 48, 0}
	roundTrip, err := CommitmentFromUint256(comm.Uint256())
	require.Nil(t, err)
	require.Equal(t, comm, roundTrip)
}

func CheckJsonRequiredFields[T any](t *testing.T, data []byte, fields ...string) {
	// Parse the JSON object into a map so we can selectively delete fields.
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	for _, field := range fields {
		data, err := json.Marshal(withoutKey(obj, field))
		require.Nil(t, err, "failed to marshal JSON")

		var dec T
		err = json.Unmarshal(data, &dec)
		require.NotNil(t, err, "unmarshalling without required field %s should fail", field)
	}
}

func withoutKey[K comparable, V any](m map[K]V, key K) map[K]V {
	copied := make(map[K]V)
	for k, v := range m {
		if k != key {
			copied[k] = v
		}
	}
	return copied
}
