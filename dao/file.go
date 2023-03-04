package dao

import (
	"context"
	"file-server/elasticSearch"
	"file-server/models"
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

func GetFileMetaByName_Version(fileName string, version int) (file meta.File, err error) {
	fileNameQuery := elastic.NewTermQuery("name", fileName)
	versionQuery := elastic.NewTermQuery("version", version)

	boolQuery := elastic.NewBoolQuery()
	boolQuery.Filter(fileNameQuery, versionQuery)

	return searchFileMeta(boolQuery)
}

func GetFileMetaByHash(fileHash string) (file meta.File, err error) {
	fileHashQuery := elastic.NewTermQuery("hash", fileHash)
	return searchFileMeta(fileHashQuery)
}

func DelMetaData(fileName string, version int) (int64, error) {
	fileNameQuery := elastic.NewTermQuery("name", fileName)
	versionQuery := elastic.NewTermQuery("version", version)

	boolQuery := elastic.NewBoolQuery()
	boolQuery.Filter(fileNameQuery, versionQuery)

	resp, err := elasticSearch.Client.DeleteByQuery("filemeta").Query(boolQuery).Refresh("true").Do(context.Background())
	if err != nil {
		return 0, err
	}
	return resp.Deleted, nil
}

func IsExistsFileName(fileName string) (file meta.File, err error) {
	fileNameQuery := elastic.NewTermQuery("name", fileName)
	return searchFileMeta(fileNameQuery)
}

func IsExistFileHash(fileHash string) (ok bool, err error) {
	fileHashQuery := elastic.NewTermQuery("hash", fileHash)
	meta, err := searchFileMeta(fileHashQuery)
	if err != nil {
		return false, err
	}
	if meta.Hash == "" {
		return false, nil
	}
	return true, nil
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

type aggregateResult struct {
	Aggregations struct {
		Group_by_name struct {
			Buckets []models.Bucket
		}
	}
}

func SearchVersionStatus(min_doc_count int) ([]models.Bucket, error) {
	//aggs := elastic.NewTermsAggregation().Field("name")
	//elasticSearch.Client.Search().Index("filemeta").Aggregation()
	return nil, nil
}
