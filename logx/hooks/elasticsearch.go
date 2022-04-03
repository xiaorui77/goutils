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
	logger *logx.LogX
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
		App:       entry.Logger.Name,
		Instance:  entry.Logger.Instance,
		Level:     entry.Level.String(),
		Message:   entry.Message,
		Fields:    entry.Fields,
		Timestamp: entry.Time.Format(timeutils.RFC3339Milli),
	}
	return nil
}

func (hook *esHook) Levels() []logx.Level {
	return hookLevels
}

func (hook *esHook) SetLogger(logger *logx.LogX) {
	hook.logger = logger
}

type LogDoc struct {
	App       string                 `json:"app"`
	Instance  string                 `json:"instance"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp string                 `json:"timestamp"`
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
	if _, err := hook.client.Index().Index("logx_" + hook.logger.Name).BodyJson(l).Do(hook.ctx); err != nil {
		hook.errCount++
	}
	hook.totalCount++
}
