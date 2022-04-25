package lib

type FreeTextCommand string

type ShortMessage string // SMS

type Mode string

var (
	Url    Mode = "url"
	Cmd    Mode = "cmd"
	Detail Mode = "xfg"
)

type Trigger int

var (
	NewRequest Trigger = 1
	LostFocus  Trigger = 2
	GainFocus  Trigger = 3

	TabLeft  Trigger = 4
	TabRight Trigger = 5
)
