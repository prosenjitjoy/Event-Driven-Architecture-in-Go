package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"mall/internal/am"
	"mall/internal/tm"

	"github.com/lib/pq"
)

type OutboxStore struct {
	tableName string
	db        DBTX
}

type outboxMessage struct {
	id      string
	name    string
	subject string
	data    []byte
}

var _ tm.OutboxStore = (*OutboxStore)(nil)
var _ am.RawMessage = (*outboxMessage)(nil)

func NewOutboxStore(tableName string, db DBTX) OutboxStore {
	return OutboxStore{
		tableName: tableName,
		db:        db,
	}
}

func (s OutboxStore) Save(ctx context.Context, msg am.RawMessage) error {
	const query = "INSERT INTO %s (id, name, subject, data) VALUES ($1, $2, $3, $4)"

	_, err := s.db.ExecContext(ctx, s.table(query), msg.ID(), msg.MessageName(), msg.Subject(), msg.Data())
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

func (s OutboxStore) FindUnpublished(ctx context.Context, limit int) ([]am.RawMessage, error) {
	const query = "SELECT id, name, subject, data FROM %s WHERE published_at IS NULL LIMIT %d"

	rows, err := s.db.QueryContext(ctx, s.table(query, limit))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("closing event rows: %w", err)
		}
	}(rows)

	var msgs []am.RawMessage

	for rows.Next() {
		msg := outboxMessage{}
		err = rows.Scan(&msg.id, &msg.name, &msg.subject, &msg.data)
		if err != nil {
			return msgs, err
		}

		msgs = append(msgs, msg)
	}

	return msgs, rows.Err()
}

func (s OutboxStore) MarkPublished(ctx context.Context, ids ...string) error {
	const query = "UPDATE %s SET published_at = CURRENT_TIMESTAMP WHERE id = ANY($1)"

	_, err := s.db.ExecContext(ctx, s.table(query), pq.Array(ids))
	if err != nil {
		return err
	}

	return nil
}

func (s OutboxStore) table(query string, args ...any) string {
	params := []any{s.tableName}
	params = append(params, args...)

	return fmt.Sprintf(query, params...)
}

func (m outboxMessage) ID() string          { return m.id }
func (m outboxMessage) Subject() string     { return m.subject }
func (m outboxMessage) MessageName() string { return m.name }
func (m outboxMessage) Data() []byte        { return m.data }
