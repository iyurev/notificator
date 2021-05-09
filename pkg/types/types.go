package types

type Sender interface {
	Send(event Event) error
}

type Event interface {
	Msg(rt ReceiverType) ([]byte, error)
	Recipient() *RecipientRef
}

type RecipientRef struct {
	Project string
	Users   []string
	Groups  []string
}

func (r *RecipientRef) GetProjectName() string {
	return r.Project
}

func NewReceiverRef(projectName string) *RecipientRef {
	return &RecipientRef{
		Project: projectName,
	}
}

type ReceiverType int32

const (
	ReceiverTypeTelegram   ReceiverType = 0
	ReceiverTypeMattermost ReceiverType = 1
)

var (
	ReceiverTypeName = map[int32]string{
		0: "telegram",
		1: "mattermost",
	}
	ReceiverTypeValue = map[string]int32{
		"telegram":   0,
		"mattermost": 1,
	}
)

func TelegramReceiverType() ReceiverType {
	return ReceiverTypeTelegram
}
