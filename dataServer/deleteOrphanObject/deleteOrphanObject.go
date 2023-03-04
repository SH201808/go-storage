package deleteorphanobject

import (
	"file-server/dao"
	"file-server/dataServer/locate"
	"file-server/setting"
	"log"
	"os"
)

var (
	garbageFileLoc = "../garbageFile/" + setting.Conf.Port + "/"
)

func init() {
	locate.Mkdir(garbageFileLoc)
}

func Background(fileStoreLoc string) {
	files := locate.Get()

	for _, fileHash := range files {
		hashInMetaData, err := dao.IsExistFileHash(fileHash)
		if err != nil {
			log.Println("find fileMeta by hash err: ", err)
			return
		}
		if !hashInMetaData {
			del(fileHash)
		}
	}
}

func del(hash string) {
	log.Println("deleteorphanobject del: " + hash)
	locate.Delete(hash)
	os.Rename(locate.FileLoc, garbageFileLoc)
}
