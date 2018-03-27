package metrics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/olivere/elastic"
	"github.com/swagchat/rtm-api/utils"
)

type ElasticsearchProvider struct{}

func (provider *ElasticsearchProvider) Run() {
	c := utils.GetConfig()
	client, err := elastic.NewClient(elastic.SetURL(c.Metrics.Elasticsearch.URL))
	if err != nil {
		log.Println(err)
	}

	exec(func() {
		l, _ := time.LoadLocation("Etc/GMT")
		nowTime := time.Unix(time.Now().Unix(), 0).In(l)
		m := makeMetrics(nowTime)
		_, err := client.Index().
			Index(fmt.Sprintf("%s-%s", utils.AppName, nowTime.Format("2006.01.02"))).
			Type("metrics").
			BodyJson(m).
			Do(context.Background())
		if err != nil {
			log.Println(err)
		}
	}, c.Metrics.Interval)
}
