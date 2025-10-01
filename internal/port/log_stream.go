package port

type LogStream interface {
	GetStreamUrl() string
	PushLog(data string)
	SetNote(data string)
	Close()
}
