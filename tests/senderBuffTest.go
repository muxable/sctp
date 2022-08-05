package main

import (
	"fmt"
)

type senderBuffer struct {
	buffer                         [][]byte       // Main sender buffer needed
	splits                         [][][]byte     // splits for buffer for each path
	paths                          int// Individual paths - can access MTU of each path from paths.mtu
	//bufferedAmount                 []uint32       // bufferredAmount[i] is the buffer size occupied by chunks on paths[i]
	//bufferedAmountOutstandingBytes []uint32       // bufferedAmountOutstandingBytes[i] is the buffer size occupied by outstanding chunks on paths[i]
	bufferSize                     uint32         // total buffer size
}

func bufferSplitting(senderBuff senderBuffer) senderBuffer{

	tempBuff := senderBuff.buffer
	splitSize := uint(int(senderBuff.bufferSize) / senderBuff.paths)
	var senderSplits [][][]byte
	for (len(tempBuff) > 0) && (uint(len(tempBuff)) >= splitSize) {
		senderSplits = append(senderSplits, senderBuff.buffer[0:splitSize])
		tempBuff = tempBuff[splitSize:]
	}
	if len(tempBuff) > 0{
		senderSplits = append(senderSplits, tempBuff)
	}
	/*
	// As an alternative, we can add if anything extra to the first buffer split and increase its size
	if len(tempBuff) > 0{
		senderSplits[0] = append(senderSplits[0], tempBuff...)
	}*/
	senderBuff.splits = senderSplits
	return senderBuff
}

func main() {
	buff := make([][]byte, 20)
	for i := range buff {
		buff[i] = make([]byte, 10)
	}
	
	paths := 5
	sendBuff := senderBuffer{buffer: buff, paths: paths,bufferSize: uint32(len(buff))}
	sendBuff = bufferSplitting(sendBuff)
	for i := range sendBuff.splits{
		fmt.Print(sendBuff.splits[i],"\n")
	}
}