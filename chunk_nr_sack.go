package sctp

import (
	"encoding/binary"
	"errors"
	"fmt"
)

/*

chunkNRSack represents an SCTP Chunk of type NR-SACK

This chunk is sent to a peer
endpoint to
(1) acknowledge DATA chunks received in-order,
(2) acknowledge DATA chunks received out-of-order, and
(3) identify DATA chunks received more than once (i.e., duplicate.)
(4) inform the peer endpoint of non-renegable out-of-order DATA chunks.

0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|   Type = 0x10 |  Chunk Flags  |         Chunk Length          |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      Cumulative TSN Ack                       |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Advertised Receiver Window Credit (a_rwnd)           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|Number of R Gap Ack Blocks = N |Number of NR Gap Ack Blocks = M|
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| Number of Duplicate TSNs = X  |           Reserved            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| R Gap Ack Block #1 Start      |   R Gap Ack Block #1 End      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
/                                                               /
\                              ...                              \
/                                                               /
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  R Gap Ack Block #N Start     |  R Gap Ack Block #N End       |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  NR Gap Ack Block #1 Start    |   NR Gap Ack Block #1 End     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
/                                                               /
\                              ...                              \
/                                                               /
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|   NR Gap Ack Block #M Start   |  NR Gap Ack Block #M End      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                       Duplicate TSN 1                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
/                                                               /
\                              ...                              \
/                                                               /
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                       Duplicate TSN X                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

*/

type rgapAckBlock struct {
	start uint16
	end   uint16
}

type nrgapAckBlock struct {
	start uint16
	end   uint16
}

var (
	errChunkTypeNotNRSack           = errors.New("ChunkType is not of type NR-SACK")
	errNRSackSizeNotLargeEnoughInfo = errors.New("NR-SACK Chunk size is not large enough to contain header")
	errNRSackSizeNotMatchPredicted  = errors.New("NR-SACK Chunk size does not match predicted amount from header values")
)

// String makes rgapAckBlock printable
func (g rgapAckBlock) String() string {
	return fmt.Sprintf("%d - %d", g.start, g.end)
}

// String makes rgapAckBlock printable
func (g nrgapAckBlock) String() string {
	return fmt.Sprintf("%d - %d", g.start, g.end)
}

type chunkNRSack struct {
	chunkHeader
	cumulativeTSNAck               uint32
	advertisedReceiverWindowCredit uint32
	rgapAckBlocks                  []rgapAckBlock
	nrgapAckBlocks                 []nrgapAckBlock
	duplicateTSN                   []uint32
}

const (
	NRSackHeaderSize   = 12
	NRSackReservedSize = 4
)

func (s *chunkNRSack) unmarshal(raw []byte) error {
	if err := s.chunkHeader.unmarshal(raw); err != nil {
		return err
	}

	if s.typ != ctNRSack {
		return fmt.Errorf("%w: actually is %s", errChunkTypeNotNRSack, s.typ.String())
	}

	if len(s.raw) < NRSackHeaderSize {
		return fmt.Errorf("%w: %v remaining, needs %v bytes", errNRSackSizeNotLargeEnoughInfo,
			len(s.raw), NRSackHeaderSize)
	}

	s.cumulativeTSNAck = binary.BigEndian.Uint32(s.raw[0:])
	s.advertisedReceiverWindowCredit = binary.BigEndian.Uint32(s.raw[4:])
	s.rgapAckBlocks = make([]rgapAckBlock, binary.BigEndian.Uint16(s.raw[8:]))
	s.nrgapAckBlocks = make([]nrgapAckBlock, binary.BigEndian.Uint16(s.raw[10:]))
	s.duplicateTSN = make([]uint32, binary.BigEndian.Uint16(s.raw[12:]))

	if len(s.raw) != NRSackHeaderSize+NRSackReservedSize+(4*len(s.rgapAckBlocks)+4*len(s.nrgapAckBlocks)+(4*len(s.duplicateTSN))) {
		fmt.Println(len(s.raw))
		fmt.Println(NRSackHeaderSize + (4*len(s.rgapAckBlocks) + 4*len(s.nrgapAckBlocks) + (4 * len(s.duplicateTSN))))
		return errNRSackSizeNotMatchPredicted
	}

	offset := NRSackHeaderSize
	for i := range s.rgapAckBlocks {
		s.rgapAckBlocks[i].start = binary.BigEndian.Uint16(s.raw[offset:])
		s.rgapAckBlocks[i].end = binary.BigEndian.Uint16(s.raw[offset+2:])
		offset += 4
	}

	for i := range s.nrgapAckBlocks {
		s.nrgapAckBlocks[i].start = binary.BigEndian.Uint16(s.raw[offset:])
		s.nrgapAckBlocks[i].end = binary.BigEndian.Uint16(s.raw[offset+2:])
		offset += 4
	}

	for i := range s.duplicateTSN {
		s.duplicateTSN[i] = binary.BigEndian.Uint32(s.raw[offset:])
		offset += 4
	}

	return nil
}

func (s *chunkNRSack) marshal() ([]byte, error) {
	nrsackRaw := make([]byte, NRSackHeaderSize+NRSackReservedSize+(4*len(s.rgapAckBlocks)+4*len(s.nrgapAckBlocks)+(4*len(s.duplicateTSN))))
	binary.BigEndian.PutUint32(nrsackRaw[0:], s.cumulativeTSNAck)
	binary.BigEndian.PutUint32(nrsackRaw[4:], s.advertisedReceiverWindowCredit)
	binary.BigEndian.PutUint16(nrsackRaw[8:], uint16(len(s.rgapAckBlocks)))
	binary.BigEndian.PutUint16(nrsackRaw[10:], uint16(len(s.nrgapAckBlocks)))
	binary.BigEndian.PutUint16(nrsackRaw[12:], uint16(len(s.duplicateTSN)))
	offset := NRSackHeaderSize
	for _, g := range s.rgapAckBlocks {
		binary.BigEndian.PutUint16(nrsackRaw[offset:], g.start)
		binary.BigEndian.PutUint16(nrsackRaw[offset+2:], g.end)
		offset += 4
	}

	for _, g := range s.nrgapAckBlocks {
		binary.BigEndian.PutUint16(nrsackRaw[offset:], g.start)
		binary.BigEndian.PutUint16(nrsackRaw[offset+2:], g.end)
		offset += 4
	}

	for _, t := range s.duplicateTSN {
		binary.BigEndian.PutUint32(nrsackRaw[offset:], t)
		offset += 4
	}

	s.chunkHeader.typ = ctNRSack
	s.chunkHeader.raw = nrsackRaw
	return s.chunkHeader.marshal()
}

func (s *chunkNRSack) check() (abort bool, err error) {
	return false, nil
}

// String makes chunkNRSack printable
func (s *chunkNRSack) String() string {
	res := fmt.Sprintf("SACK cumTsnAck=%d arwnd=%d dupTsn=%d",
		s.cumulativeTSNAck,
		s.advertisedReceiverWindowCredit,
		s.duplicateTSN)

	for _, gap := range s.rgapAckBlocks {
		res = fmt.Sprintf("%s\n r-gap ack: %s", res, gap)
	}

	return res
}
