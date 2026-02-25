package domain

type StoreCreated struct {
	Name     string
	Location string
}

type StoreParticipationToggled struct {
	Participation bool
}

type StoreRebranded struct {
	Name string
}
