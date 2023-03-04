package filemeta

import (
	"file-server/dao"
	"file-server/models/meta"

	"github.com/olivere/elastic/v7"
)

func Put(file meta.File) error {
	newFileMeta, err := dao.IsExistsFileName(file.Name)
	if err != nil {
		if !elastic.IsNotFound(err) {
			return err
		}
		file.Version = 1
	} else {
		file.Version = newFileMeta.Version + 1
	}
	return dao.PutFileMeta(file)
}

func Get(fileName string, fileVersion int) (meta.File, error) {
	return dao.GetFileMetaByName_Version(fileName, fileVersion)
}

func GetFileSize(fileHash string) (string, error) {
	meta, err := dao.GetFileMetaByHash(fileHash)
	if err != nil {
		return "", err
	}
	return meta.Size, nil
}
