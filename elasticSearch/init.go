package elasticSearch

import (
	"context"
	"file-server/setting"
	"log"

	"github.com/olivere/elastic/v7"
)

var Client *elastic.Client

func Init(cfg *setting.ElasticSearchConfig) {
	url := "http://" + cfg.Host + ":" + cfg.Port
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		log.Fatalln("init elasticSearch err: ", err)
	}

	Client = client
}

func CreateIndex(index string, mappings string) error {
	exists, err := Client.IndexExists(index).Do(context.Background())
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	createIndex, err := Client.CreateIndex(index).BodyString(mappings).Do(context.Background())
	if err != nil {
		return err
	}
	if !createIndex.Acknowledged {

	}
	return nil
}
