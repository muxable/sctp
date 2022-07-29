package sctp

func bufferSplitting(senderBuff senderBuffer) {

	tempBuff := senderBuff.buffer
	splitSize := uint(int(senderBuff.bufferSize) / len(senderBuff.paths))
	var senderSplits [][]byte
	for (len(tempBuff) > 0) && (uint(len(tempBuff)) > splitSize) {
		senderSplits = append(senderSplits, senderBuff.buffer[0:splitSize])
		tempBuff = tempBuff[splitSize:]
	}
	if len(tempBuff) > 0 {
		senderSplits = append(senderSplits, tempBuff)
	}
	senderBuff.splits = senderSplits
}

func sendCond_bufferedBytes(senderBuff senderBuffer, path *Association) bool {

	var pathIdx int
	for i := range senderBuff.paths {
		if senderBuff.paths[i].name == path.name {
			pathIdx = i
			break
		}
	}
	return senderBuff.bufferedAmount[pathIdx]+senderBuff.paths[pathIdx].mtu <= uint32(int(senderBuff.bufferSize)/len(senderBuff.paths))
}

func sendCond_receiveBuffer_bufferedBytes(senderBuff senderBuffer, path *Association) bool {

	var pathIdx int
	for i := range senderBuff.paths {
		if senderBuff.paths[i].name == path.name {
			pathIdx = i
			break
		}
	}
	return senderBuff.bufferedAmount[pathIdx] <= uint32((int(senderBuff.paths[pathIdx].rwnd)+len(senderBuff.bufferedAmountOutstandingBytes))/len(senderBuff.paths))
}

func sendCond_outstandingBytes(senderBuff senderBuffer, path *Association) bool {

	var pathIdx int
	for i := range senderBuff.paths {
		if senderBuff.paths[i].name == path.name {
			pathIdx = i
			break
		}
	}
	return senderBuff.bufferedAmountOutstandingBytes[pathIdx]+senderBuff.paths[pathIdx].mtu <= uint32(int(senderBuff.bufferSize)/len(senderBuff.paths))
}

func sendCond_receiveBuffer_outstandingBytes(senderBuff senderBuffer, path *Association) bool {

	var pathIdx int
	for i := range senderBuff.paths {
		if senderBuff.paths[i].name == path.name {
			pathIdx = i
			break
		}
	}
	return senderBuff.bufferedAmountOutstandingBytes[pathIdx] <= uint32((int(senderBuff.paths[pathIdx].rwnd)+len(senderBuff.bufferedAmountOutstandingBytes))/len(senderBuff.paths))
}