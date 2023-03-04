package deleteoldmetadata

import (
	"file-server/dao"
	"log"
)

const MIN_VERSION_COUNT = 5

func main() {
	buckets, err := dao.SearchVersionStatus(MIN_VERSION_COUNT + 1)
	if err != nil {
		log.Println(err)
		return
	}
	for i := range buckets {
		bucket := buckets[i]
		for v := 0; v < bucket.Doc_count-MIN_VERSION_COUNT; v++ {
			dao.DelMetaData(bucket.Key, v+int(bucket.Min_version.Value))
		}
	}
}
