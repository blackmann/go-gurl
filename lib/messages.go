package lib

type FreeText string

type Trigger int

var (
	NewRequest Trigger = 1
	LostFocus  Trigger = 2
	GainFocus  Trigger = 3

	TabLeft  Trigger = 4
	TabRight Trigger = 5
)
