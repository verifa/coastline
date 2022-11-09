package store

import (
	"github.com/verifa/coastline/ent"
	entUser "github.com/verifa/coastline/ent/user"
	"github.com/verifa/coastline/server/oapi"
)

func (s *Store) getEntUser(user *oapi.User) (*ent.User, error) {
	return s.client.User.Query().
		Where(entUser.And(entUser.Sub(user.Sub), entUser.Iss(user.Iss))).
		First(s.ctx)
}

func dbUserToAPI(dbUser *ent.User) *oapi.User {
	dbGroups := make([]string, len(dbUser.Edges.Groups))
	for i, entGroup := range dbUser.Edges.Groups {
		dbGroups[i] = entGroup.Name
	}

	return &oapi.User{
		Sub:     dbUser.Sub,
		Iss:     dbUser.Iss,
		Name:    dbUser.Name,
		Email:   dbUser.Email,
		Picture: dbUser.Picture,
		Groups:  dbGroups,
	}
}
