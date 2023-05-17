package eventbus

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"time"

	"github.com/0xDeSchool/gap/ddd"
	"github.com/0xDeSchool/gap/log"
	"github.com/0xDeSchool/gap/utils/linq"
	"github.com/0xDeSchool/gap/x"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventMsg struct {
	ddd.EntityBase `bson:",inline"`
	Topic          string            `bson:"topic" json:"topic"`
	EventID        string            `bson:"eventID" json:"eventID"`
	Data           string            `bson:"data" json:"data"`
	PublishTime    *time.Time        `bson:"publishTime" json:"publishTime"`
	ErrorMsg       string            `bson:"errorMsg" json:"errorMsg"`
	Metadata       map[string]string `bson:"metadata" json:"metadata"`
}

type EventStore interface {
	// Save saves an event to the event store
	Save(ctx context.Context, data *EventMsg) error

	GetList(ctx context.Context, p *x.PageParam) (x.PagedResult[EventMsg], error)
	Delete(ctx context.Context, id []primitive.ObjectID) error
}

// ListenErrorEvents 监听错误事件，每天指定整点执行
func ListenErrorEvents(clock int) {
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)                                                          //通过now偏移24小时
		next = time.Date(next.Year(), next.Month(), next.Day(), clock, 0, 0, 0, next.Location()) //获取下一个凌晨的日期
		t := time.NewTimer(next.Sub(now))                                                        //计算当前时间到凌晨的时间间隔，设置一个定时器
		<-t.C
		//以下为定时执行的操作
		if store, ok := app.GetOptional[EventStore](); ok {
			page := x.NewPageParam(1, 100)
			for {
				events, err := (*store).GetList(context.Background(), page)
				if err != nil {
					log.Warnf("get error events failed, error: %s", err.Error())
					break
				}
				for _, event := range events.Data {
					if event.ErrorMsg == testError {
						continue
					}
					PublishJSON(context.Background(), event.Topic, event.Data)
				}
				ids := linq.Map(events.Data, func(event *EventMsg) primitive.ObjectID {
					return event.ID
				})
				err = (*store).Delete(context.Background(), ids)
				if err != nil {
					log.Warnf("delete error events failed, error: %s", err.Error())
				}
				if !events.HasMore {
					break
				} else {
					page.Page++
				}
			}
		}
	}
}
