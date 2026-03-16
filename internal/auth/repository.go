package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	redisClient *redis.Client
}

func NewRepository(redisClient *redis.Client) *Repository {
	return &Repository{
		redisClient: redisClient,
	}
}

func (r *Repository) CreateSession(ctx context.Context, session Session, sessionTtlSeconds int) error {
	sessionKey := createSessionKey(session.ID)
	sessionData := sessionToModel(session)

	sessionDataJson, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, sessionKey, sessionDataJson, time.Duration(sessionTtlSeconds)*time.Second).Err()
}

func (r *Repository) GetSessionByID(ctx context.Context, sessionId uuid.UUID) (*Session, error) {
	sessionKey := createSessionKey(sessionId)

	sessionDataJson, err := r.redisClient.Get(ctx, sessionKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrorSessionNotFound
	}
	if err != nil {
		return nil, err
	}

	var sessionData sessionModel
	err = json.Unmarshal([]byte(sessionDataJson), &sessionData)
	if err != nil {
		return nil, err
	}

	sessionEntity, err := modelToSession(sessionData)
	if err != nil {
		return nil, err
	}

	return &sessionEntity, nil
}

func (r *Repository) DeleteSession(ctx context.Context, sessionId uuid.UUID) error {
	sessionKey := createSessionKey(sessionId)
	return r.redisClient.Del(ctx, sessionKey).Err()
}

func (r *Repository) RefreshSession(ctx context.Context, sessionId uuid.UUID, sessionTtlSeconds int) error {
	sessionKey := createSessionKey(sessionId)
	return r.redisClient.Expire(ctx, sessionKey, time.Duration(sessionTtlSeconds)*time.Second).Err()
}

func sessionToModel(session Session) sessionModel {
	return sessionModel{
		ID:        session.ID.String(),
		UserID:    session.UserId.String(),
		CreatedAt: session.CreatedAt.Unix(),
		UpdatedAt: session.UpdatedAt.Unix(),
	}
}

func modelToSession(sessionModel sessionModel) (Session, error) {
	userId, err := uuid.Parse(sessionModel.UserID)
	if err != nil {
		return Session{}, err
	}

	sessionId, err := uuid.Parse(sessionModel.ID)
	if err != nil {
		return Session{}, err
	}

	return Session{
		ID:        sessionId,
		UserId:    userId,
		CreatedAt: time.Unix(sessionModel.CreatedAt, 0).UTC(),
		UpdatedAt: time.Unix(sessionModel.UpdatedAt, 0).UTC(),
	}, nil
}

type sessionModel struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func createSessionKey(sessionId uuid.UUID) string {
	return "session:" + sessionId.String()
}
