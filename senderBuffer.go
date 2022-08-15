package sctp

type senderBuffer struct {
	buffer                         [][]byte       // Main sender buffer needed
	splits                         [][][]byte     // splits for buffer for each path
	paths                          []*Association // Individual paths - can access MTU of each path from paths.mtu
	bufferedAmount                 []uint32       // bufferredAmount[i] is the buffer size occupied by chunks on paths[i]
	bufferedAmountOutstandingBytes []uint32       // bufferedAmountOutstandingBytes[i] is the buffer size occupied by outstanding chunks on paths[i]
	bufferSize                     uint32         // total buffer size
}