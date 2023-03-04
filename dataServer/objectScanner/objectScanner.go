package objectscanner

import (
	"file-server/dataServer/locate"
)

func start() {
	files := locate.Get()

	for _, hash := range files {
		verify(hash)
	}
}

func verify(hash string) {
	// size, err := filemeta.GetFileSize(hash)
	// if err != nil {
	// 	log.Println("getFileSize err: ", err)
	// 	return
	// }
	//GetStream(hash, size)
}
