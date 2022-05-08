package lib

type FreeTextCommand string

type ShortMessage string // SMS

type Mode string

var (
	Url    Mode = "url"
	Cmd    Mode = "cmd"
	Detail Mode = "xfg"
)

type Event int

var (
	NewRequest Event = 1
	LostFocus  Event = 2
	GainFocus  Event = 3

	TabLeft  Event = 4
	TabRight Event = 5

	UpdateHistory   Event = 6
	UpdateBookmarks Event = 7

	SavedResponse Event = 10

	Reset Event = 11
)
