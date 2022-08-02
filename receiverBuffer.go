package sctp

type receiverBuffer struct {
	buffer                         [][]byte      // Main receiver buffer needed
	splits                         [][][]byte    // splits for buffer for each path
	paths                          []Association // Individual paths - can access MTU of each path from paths.mtu
	bufferedAmount                 []uint32      // bufferredAmount[i] is the buffer size occupied by chunks on paths[i]
	bufferedAmountOutstandingBytes []uint32      // bufferedAmountOutstandingBytes[i] is the buffer size occupied by outstanding chunks on paths[i]
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
