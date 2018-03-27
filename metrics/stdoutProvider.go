package metrics

import (
	"time"

	"github.com/swagchat/rtm-api/utils"
	"go.uber.org/zap"
)

type StdoutProvider struct{}

func (provider *StdoutProvider) Run() {
	c := utils.GetConfig()
	exec(func() {
		l, _ := time.LoadLocation("Etc/GMT")
		nowTime := time.Unix(time.Now().Unix(), 0).In(l)
		m := makeMetrics(nowTime)
		sb := utils.NewStringBuilder()
		mStr := sb.PrintStruct("config", m)
		utils.AppLogger.Info("",
			zap.String("metrics", mStr),
		)
	}, c.Metrics.Interval)
}
