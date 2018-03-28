package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic"
	"github.com/swagchat/rtm-api/logging"
	"github.com/swagchat/rtm-api/utils"
	"go.uber.org/zap/zapcore"
)

type ElasticsearchProvider struct{}

func (provider *ElasticsearchProvider) Run() {
	c := utils.Config()
	client, err := elastic.NewClient(elastic.SetURL(c.Metrics.Elasticsearch.URL))
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Kind:     "metrics-error",
			Provider: "elasticsearch",
			Message:  err.Error(),
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
				Kind:     "metrics-error",
				Provider: "elasticsearch",
				Message:  err.Error(),
			})
		}
	}, c.Metrics.Interval)
}
