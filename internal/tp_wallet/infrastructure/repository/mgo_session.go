package repository

import (
	"context"
	"gluttonous/internal/wallet/domain/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ repo.MgoSession = (*MgoSession)(nil)

type MgoSession struct {
	repository
}

func (m *MgoSession) MgoSession(ctx context.Context) (mongo.SessionContext, error) {
	session, err := m.mgo.StartSession()
	if err != nil {
		return nil, err
	}
	sessionContext := mongo.NewSessionContext(ctx, session)
	return sessionContext, nil
}
