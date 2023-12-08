package model

type Message struct {
	Id       int64  `gorm:"primaryKey;column:id;" json:"id"`
	Uid      int64  `gorm:"column:uid;type:int(11);default:(-)" json:"uid"`
	ToUid    int64  `gorm:"column:to_uid;type:int(11);default:(-)" json:"to_uid"`
	Message  string `gorm:"column:message;type:varchar(100);default:(-)" json:"message"`
	Status   uint8  `gorm:"column:status;type:tinyint(1);default:(1)" json:"status"`
	Datetime int64  `gorm:"column:datetime;autoCreateTime" json:"datetime"`
}

// TableName 自定义表名
func (*Message) TableName() string {
	return "zhyu_message"
}
