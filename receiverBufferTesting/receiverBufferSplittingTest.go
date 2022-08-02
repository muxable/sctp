package main

type receiverBuffer struct {
	buffer [][]byte   // Main receiver buffer needed
	splits [][][]byte // splits for buffer for each path
	//bufferedAmount                 []uint32   // bufferredAmount[i] is the buffer size occupied by chunks on paths[i]
	//bufferedAmountOutstandingBytes []uint32   // bufferedAmountOutstandingBytes[i] is the buffer size occupied by outstanding chunks on paths[i]
	//bufferSize                     uint32        // total buffer size
}

func (rb *receiverBuffer) getTotalBufferSize() int {

	//Assuming size of uint as 8 bytes
	size_of_int := 8
	s := 0
	for _, packet := range rb.buffer {
		s += len(packet) * size_of_int
	}
	return s
}

func getBufferSize(rb []byte) int {

	//Assuming size of uint as 8 bytes
	size_of_int := 8
	return size_of_int * int(len(rb))
}

func (receiverBuff *receiverBuffer) bufferSplitting() {

	tempBuff := receiverBuff.buffer
	splitSize := uint(int(receiverBuff.getTotalBufferSize()) / 3) // Hard coded no. of paths as 3
	println("SPLIT SIZE: ", splitSize)
	i := 0

	// To avoid packet fragmentation
	var receiverSplits [][][]byte

	for i < len(tempBuff) {
		j := i
		currentLength := uint(len(tempBuff[j]))
		var currentBuffer [][]byte
		for (j < len(tempBuff)) && currentLength <= splitSize {
			currentBuffer = append(currentBuffer, tempBuff[j])
			currentLength += uint(getBufferSize((tempBuff[j])))
			j += 1
		}
		i = j
		receiverSplits = append(receiverSplits, currentBuffer)

	}

	receiverBuff.splits = receiverSplits

}

func main() {
	rb := receiverBuffer{
		buffer: [][]byte{{1, 2, 3, 4}, {1, 2}, {4, 5, 7, 8}, {1}, {2}, {4, 6}},
	}

	println("Total Buffer Size: ", rb.getTotalBufferSize())
	rb.bufferSplitting()

	for _, buf := range rb.splits {

		for _, b := range buf {
			for _, c := range b {
				print(c, " ")
			}
			print(", ")
		}
		println()
	}
	// print("\n", rb.splits)
}
