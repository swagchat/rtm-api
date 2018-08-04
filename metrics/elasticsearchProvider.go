package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic"
	"github.com/swagchat/rtm-api/config"
	"github.com/swagchat/rtm-api/logging"
	"go.uber.org/zap/zapcore"
)

type elasticsearchProvider struct{}

func (ep *elasticsearchProvider) Run() {
	c := config.Config()
	client, err := elastic.NewClient(elastic.SetURL(c.Metrics.Elasticsearch.URL))
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Kind:     "metrics",
			Provider: "elasticsearch",
			Message:  fmt.Sprintf("%s endpoint[%s]", err.Error(), c.Metrics.Elasticsearch.URL),
		})
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
			logging.Log(zapcore.ErrorLevel, &logging.AppLog{
				Kind:     "metrics",
				Provider: "elasticsearch",
				Message:  fmt.Sprintf("%s endpoint[%s] index[%s] indexTimeFormat[%s]", err.Error(), c.Metrics.Elasticsearch.URL, c.Metrics.Elasticsearch.Index, c.Metrics.Elasticsearch.IndexTimeFormat),
			})
		}
	}, c.Metrics.Interval)
}
