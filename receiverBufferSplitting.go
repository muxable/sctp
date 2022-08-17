package sctp

func bufferSplitting(receiverBuff receiverBuffer) {

	tempBuff := receiverBuff.buffer
	splitSize := uint(int(receiverBuff.getTotalBufferSize()) / len(receiverBuff.paths))
	i := 0

	// To avoid packet fragmentation
	var receiverSplits [][][]byte

	for i < len(tempBuff) {
		j := i
		currentLength := uint(len(tempBuff[j]))
		var currentBuffer [][]byte
		for (j < len(tempBuff)) && currentLength <= splitSize {
			j += 1
			currentBuffer = append(currentBuffer, tempBuff[j])
		}
		receiverSplits = append(receiverSplits, currentBuffer)
		i = j
	}

	receiverBuff.splits = receiverSplits

}

func sendCond_bufferedBytes(receiverBuff receiverBuffer, path *Association) bool {

	var pathIdx int
	for i := range receiverBuff.paths {
		if receiverBuff.paths[i].name == path.name {
			pathIdx = i
			break
		}
	}
	return receiverBuff.bufferedAmount[pathIdx]+receiverBuff.paths[pathIdx].mtu <= uint32(int(receiverBuff.getTotalBufferSize())/len(receiverBuff.paths))
}

func sendCond_receiveBuffer_bufferedBytes(receiverBuff receiverBuffer, path *Association) bool {

	var pathIdx int
	for i := range receiverBuff.paths {
		if receiverBuff.paths[i].name == path.name {
			pathIdx = i
			break
		}
	}
	return receiverBuff.bufferedAmount[pathIdx] <= uint32((int(receiverBuff.paths[pathIdx].rwnd)+len(receiverBuff.bufferedAmountOutstandingBytes))/len(receiverBuff.paths))
}

func sendCond_outstandingBytes(receiverBuff receiverBuffer, path *Association) bool {

	var pathIdx int
	for i := range receiverBuff.paths {
		if receiverBuff.paths[i].name == path.name {
			pathIdx = i
			break
		}
	}
	return receiverBuff.bufferedAmountOutstandingBytes[pathIdx]+receiverBuff.paths[pathIdx].mtu <= uint32(int(receiverBuff.getTotalBufferSize())/len(receiverBuff.paths))
}

func sendCond_receiveBuffer_outstandingBytes(receiverBuff receiverBuffer, path *Association) bool {

	var pathIdx int
	for i := range receiverBuff.paths {
		if receiverBuff.paths[i].name == path.name {
			pathIdx = i
			break
		}
	}
	return receiverBuff.bufferedAmountOutstandingBytes[pathIdx] <= uint32((int(receiverBuff.paths[pathIdx].rwnd)+len(receiverBuff.bufferedAmountOutstandingBytes))/len(receiverBuff.paths))
}
