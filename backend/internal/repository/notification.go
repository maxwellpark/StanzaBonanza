package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
)

type NotificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{pool: pool}
}

func (r *NotificationRepository) Create(ctx context.Context, notif *domain.Notification) error {
	if notif.ID == uuid.Nil {
		notif.ID = uuid.New()
	}
	notif.CreatedAt = time.Now().UTC()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO notifications (id, recipient_id, actor_id, type, poem_id, read, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		notif.ID, notif.RecipientID, notif.ActorID, notif.Type, notif.PoemID, notif.Read, notif.CreatedAt,
	)
	return err
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Notification, int, error) {
	page.Normalize()

	var totalCount int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE recipient_id = $1`, userID,
	).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT n.id, n.recipient_id, n.actor_id, n.type, n.poem_id, n.read, n.created_at,
		        u.id, u.display_name, u.avatar_url
		 FROM notifications n
		 LEFT JOIN users u ON u.id = n.actor_id
		 WHERE n.recipient_id = $1
		 ORDER BY n.created_at DESC LIMIT $2 OFFSET $3`,
		userID, page.PageSize, page.Offset(),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notifications := make([]domain.Notification, 0)
	for rows.Next() {
		var n domain.Notification
		var actorID *uuid.UUID
		var actorDisplayName *string
		var actorAvatarURL *string
		if err := rows.Scan(
			&n.ID, &n.RecipientID, &n.ActorID, &n.Type, &n.PoemID, &n.Read, &n.CreatedAt,
			&actorID, &actorDisplayName, &actorAvatarURL,
		); err != nil {
			return nil, 0, err
		}
		if actorID != nil {
			n.Actor = &domain.User{
				ID:          *actorID,
				DisplayName: derefStr(actorDisplayName),
				AvatarURL:   derefStr(actorAvatarURL),
			}
		}
		notifications = append(notifications, n)
	}
	return notifications, totalCount, rows.Err()
}

func (r *NotificationRepository) MarkRead(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE notifications SET read = true WHERE recipient_id = $1 AND id = ANY($2)`,
		userID, ids,
	)
	return err
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
