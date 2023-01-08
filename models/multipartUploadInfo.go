package models

type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadId   string
	ChunkSize  int
	ChunkCount int
}
