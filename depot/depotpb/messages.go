package depotpb

import (
	"mall/internal/registry"
	"mall/internal/registry/serdes"
)

const (
	CommandChannel = "mall.depot.commands"

	CreateShoppingListCommand = "depotapi.CreateShoppingListCommand"
	CancelShoppingListCommand = "depotapi.CancelShoppingListCommand"
	InitiateShoppingCommand   = "depotapi.InitiateShoppingCommand"

	CreatedShoppingListReply = "depotapi.CreatedShoppingListReply"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// Commands
	if err := serde.Register(&CreateShoppingList{}); err != nil {
		return err
	}
	if err := serde.Register(&CancelShoppingList{}); err != nil {
		return err
	}
	if err := serde.Register(&InitiateShopping{}); err != nil {
		return err
	}

	// Replies
	if err := serde.Register(&CreatedShoppingList{}); err != nil {
		return err
	}

	return nil
}

// Commands
func (*CreateShoppingList) Key() string { return CreateShoppingListCommand }
func (*CancelShoppingList) Key() string { return CancelShoppingListCommand }
func (*InitiateShopping) Key() string   { return InitiateShoppingCommand }

// Replies
func (*CreatedShoppingList) Key() string { return CreatedShoppingListReply }
