package lib

type Status int32

var (
	IDLE       Status = 0
	PROCESSING Status = 1
	ERROR      Status = 3
)

func (s Status) GetValue() int {
	return int(s)
}
