package dao

import (
	"context"
	"file-server/elasticSearch"
	"file-server/models/meta"
	"reflect"

	"github.com/olivere/elastic/v7"
)

func PutFileMeta(fileMeta meta.File) error {
	_, err := elasticSearch.Client.Index().
		Index("filemeta").BodyJson(fileMeta).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func GetFileMeta(fileName string, version int) (file meta.File, err error) {
	fileNameQuery := elastic.NewTermQuery("name", fileName)
	versionQuery := elastic.NewTermQuery("version", version)

	boolQuery := elastic.NewBoolQuery()
	boolQuery.Filter(fileNameQuery, versionQuery)

	return searchFileMeta(boolQuery)
}

func IsExistsFileMeta(fileName string) (file meta.File, err error) {
	fileNameQuery := elastic.NewTermQuery("name", fileName)
	return searchFileMeta(fileNameQuery)
}

func searchFileMeta(query elastic.Query) (file meta.File, err error) {
	searchResult, err := elasticSearch.Client.Search().
		Index("filemeta").Query(query).Do(context.Background())
	if err != nil {
		return
	}
	for _, item := range searchResult.Each(reflect.TypeOf(file)) {
		file = item.(meta.File)
	}
	return
}
