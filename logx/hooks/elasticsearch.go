package hooks

import (
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/xiaorui77/goutils/logx"
	"github.com/xiaorui77/goutils/timeutils"
)

var hookLevels = []logx.Level{
	logx.InfoLevel,
	logx.WarnLevel,
	logx.ErrorLevel,
	logx.FatalLevel,
	logx.PanicLevel,
}

type esHook struct {
	client *elastic.Client
	ctx    context.Context

	buff       chan *LogDoc
	errCount   int
	totalCount int
}

func NewEsHook(host string) *esHook {
	client, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetSniff(false),
	)
	if err != nil {
		logx.Fatalf("failed to create Elastic V7 Client: %v", err)
		return nil
	}
	h := &esHook{
		client: client,
		ctx:    context.Background(),
		buff:   make(chan *LogDoc, 1000),
	}
	go h.run()
	return h
}

func (hook *esHook) Fire(entry *logx.Entry) error {
	hook.buff <- &LogDoc{
		Level:   entry.Level.String(),
		Time:    entry.Time.Format(timeutils.RFC3339Milli),
		Message: entry.Message,
	}
	return nil
}

func (hook *esHook) Levels() []logx.Level {
	return hookLevels
}

type LogDoc struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

func (hook *esHook) run() {
	for {
		select {
		case <-hook.ctx.Done():
			return
		case l := <-hook.buff:
			hook.send(l)
		}
	}
}

func (hook *esHook) send(l *LogDoc) {
	if _, err := hook.client.Index().Index("logx_" + l.Level).BodyJson(l).Do(hook.ctx); err != nil {
		hook.errCount++
	}
	hook.totalCount++
}
