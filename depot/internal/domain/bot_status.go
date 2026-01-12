package domain

type BotStatus string

const (
	BotUnknown  BotStatus = ""
	BotIsIdle   BotStatus = "idle"
	BotIsActive BotStatus = "active"
)

func (s BotStatus) String() string {
	switch s {
	case BotIsIdle, BotIsActive:
		return string(s)
	default:
		return ""
	}
}
