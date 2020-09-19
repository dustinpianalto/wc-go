package wc

func GetLineCount(chunk []byte) int64 {
	var count int64
	for _, b := range chunk {
		if b == '\n' {
			count++
		}
	}
	return count
}

func ConcurrentChunkCounter(chunks <-chan []byte, counts chan<- int64) {
	var totalCount int64
	for chunk := range chunks {
		totalCount += GetLineCount(chunk)
	}
	counts <- totalCount
}
