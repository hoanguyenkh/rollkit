package types

import (
	"bytes"
	"encoding"
	"fmt"
	"time"

	"github.com/celestiaorg/go-header"
)

type Hash = header.Hash

// BaseHeader contains the most basic data of a header
type BaseHeader struct {
	// Height represents the block height (aka block number) of a given header
	Height uint64
	// Time contains Unix nanotime of a block
	Time uint64
	// The Chain ID
	ChainID string
}

// Header defines the structure of Rollkit block header.
type Header struct {
	BaseHeader
	// Block and App version
	Version Version

	// prev block info
	LastHeaderHash Hash

	// hashes of block data
	LastCommitHash Hash // commit from aggregator(s) from the last block
	DataHash       Hash // Block.Data root aka Transactions
	ConsensusHash  Hash // consensus params for current block
	AppHash        Hash // state after applying txs from the current block

	// Root hash of all results from the txs from the previous block.
	// This is ABCI specific but smart-contract chains require some way of committing
	// to transaction receipts/results.
	LastResultsHash Hash

	// Note that the address can be derived from the pubkey which can be derived
	// from the signature when using secp256k.
	// We keep this in case users choose another signature format where the
	// pubkey can't be recovered by the signature (e.g. ed25519).
	ProposerAddress []byte // original proposer of the block
}

func (h *Header) New() *Header {
	return new(Header)
}

func (h *Header) IsZero() bool {
	return h == nil
}

func (h *Header) ChainID() string {
	return h.BaseHeader.ChainID
}

func (h *Header) Height() uint64 {
	return h.BaseHeader.Height
}

func (h *Header) LastHeader() Hash {
	return h.LastHeaderHash[:]
}

// Returns unix time with nanosecond precision
func (h *Header) Time() time.Time {
	return time.Unix(0, int64(h.BaseHeader.Time))
}

func (h *Header) Verify(untrstH *Header) error {
	if bytes.Equal(h.ProposerAddress, untrstH.ProposerAddress) {
		return nil
	}
	return fmt.Errorf("incorrect proposer address")
}

func (h *Header) Validate() error {
	return h.ValidateBasic()
}

// ValidateBasic performs basic validation of a header.
func (h *Header) ValidateBasic() error {
	if len(h.ProposerAddress) == 0 {
		return ErrNoProposerAddress
	}

	return nil
}

var _ header.Header[*Header] = &Header{}
var _ encoding.BinaryMarshaler = &Header{}
var _ encoding.BinaryUnmarshaler = &Header{}
