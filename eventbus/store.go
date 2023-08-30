package eventbus

import (
	"time"

	"github.com/0xDeSchool/gap/ddd"
)

type EventMsg struct {
	ddd.EntityBase[int64] `bson:",inline"`
	Topic                 string            `bson:"topic" json:"topic"`
	EventID               string            `bson:"eventID" json:"eventID"`
	Data                  string            `bson:"data" json:"data"`
	PublishTime           *time.Time        `bson:"publishTime" json:"publishTime"`
	ErrorMsg              string            `bson:"errorMsg" json:"errorMsg"`
	Metadata              map[string]string `bson:"metadata" json:"metadata"`
}
