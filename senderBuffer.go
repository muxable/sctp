package sctp

type senderBuffer struct {
	buffer                         []*chunkPayloadData   // Main sender buffer needed
	splits                         [][]*chunkPayloadData // splits for buffer for each path
	paths                          []*Association        // Individual paths - can access MTU of each path from paths.mtu
	bufferedAmount                 []uint32              // bufferredAmount[i] is the buffer size occupied by chunks on paths[i]
	bufferedAmountOutstandingBytes []uint32              // bufferedAmountOutstandingBytes[i] is the buffer size occupied by outstanding chunks on paths[i]
	bufferSize                     uint32                // total buffer size
}

func newSenderBuffer(paths []*Association) *senderBuffer {
	senderBuff := senderBuffer{buffer: []*chunkPayloadData{}, paths: paths, splits: [][]*chunkPayloadData{}}
	senderBuff.bufferSize = uint32(len(senderBuff.buffer))
	a := make([][]*chunkPayloadData, len(paths))
	senderBuff.bufferedAmount = make([]uint32, len(paths))
	senderBuff.bufferedAmountOutstandingBytes = make([]uint32, len(paths))
	senderBuff.splits = a
	return &senderBuff
}

func getPathIdx(q *senderBuffer, path *Association) int {
	for i := range q.paths {
		if q.paths[i].name == path.name {
			return i
		}
	}
	return -1
}

func (q *senderBuffer) push(c *chunkPayloadData) {
	for pathIdx := range q.paths {
		if sendCond_bufferedBytes_pathIdx(q, pathIdx) {
			q.splits[pathIdx] = append(q.splits[pathIdx], c)
			return
		}
	}
}

func (q *senderBuffer) pop(pathIdx int) *chunkPayloadData {
	if len(q.splits[pathIdx]) == 0 {
		return nil
	}
	data := q.splits[pathIdx][0]
	q.splits[pathIdx] = q.splits[pathIdx][1:]
	return data
}

func (q *senderBuffer) get(pathIdx int, pos int) *chunkPayloadData {
	if len(q.splits[pathIdx]) == 0 || pos < 0 || pos >= len(q.splits[pathIdx]) {
		return nil
	}
	return q.splits[pathIdx][pos]
}

func (q *senderBuffer) size(pathIdx int) int {
	return len(q.splits[pathIdx])
}