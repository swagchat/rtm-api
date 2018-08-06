package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic"
	"github.com/swagchat/rtm-api/config"
	"github.com/swagchat/rtm-api/logger"
)

type elasticsearchProvider struct{}

func (ep *elasticsearchProvider) Run() {
	c := config.Config()
	client, err := elastic.NewClient(elastic.SetURL(c.Metrics.Elasticsearch.URL))
	if err != nil {
		logger.Error(err.Error())
	}

	exec(func() {
		l, _ := time.LoadLocation("Etc/GMT")
		nowTime := time.Unix(time.Now().Unix(), 0).In(l)
		m := makeMetrics(nowTime)
		_, err := client.Index().
			Index(fmt.Sprintf("%s-%s", c.Metrics.Elasticsearch.Index, nowTime.Format(c.Metrics.Elasticsearch.IndexTimeFormat))).
			Type(c.Metrics.Elasticsearch.Type).
			BodyJson(m).
			Do(context.Background())
		if err != nil {
			logger.Error(err.Error())
		}
	}, c.Metrics.Interval)
}
