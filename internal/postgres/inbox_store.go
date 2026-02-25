package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"mall/internal/am"
	"mall/internal/tm"

	"github.com/lib/pq"
)

type InboxStore struct {
	tableName string
	db        DBTX
}

var _ tm.InboxStore = (*InboxStore)(nil)

func NewInboxStore(tableName string, db DBTX) InboxStore {
	return InboxStore{
		tableName: tableName,
		db:        db,
	}
}

func (s InboxStore) Save(ctx context.Context, msg am.IncomingMessage) error {
	const query = "INSERT INTO %s (id, name, subject, data, metadata, sent_at, received_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	metadata, err := json.Marshal(msg.Metadata())
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, s.table(query), msg.ID(), msg.MessageName(), msg.Subject(), msg.Data(), metadata, msg.SentAt(), msg.ReceivedAt())
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			// unique_violation error
			if pgErr.Code == pq.ErrorCode("23505") {
				return tm.ErrDuplicateMessage(msg.ID())
			}
		}
		return err
	}

	return nil
}

func (s InboxStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}
