package protocol

var emptyMultiBulkBytes = []byte("*0\r\n")

type EmptyMultiBulkReply struct{}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

var nullBulkBytes = []byte("$-1\r\n")

type NullBulkReply struct{}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

func (r *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}
