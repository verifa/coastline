package store

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/verifa/coastline/ent"
	"github.com/verifa/coastline/ent/group"
	entSession "github.com/verifa/coastline/ent/session"
	entUser "github.com/verifa/coastline/ent/user"
	"github.com/verifa/coastline/server/oapi"
)

func (s *Store) NewSession(user *oapi.User) (uuid.UUID, error) {

	dbUser, err := s.createUpdateUser(user)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("getting user from database: %w", err)
	}

	dbSession, err := s.client.Session.Create().
		SetUser(dbUser).
		SetDuration(int64(s.config.SessionDuration)).
		Save(s.ctx)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("creating new session: %w", err)
	}

	return dbSession.ID, nil
}

func (s *Store) GetSession(uuid uuid.UUID) (*oapi.User, error) {
	// Get the session
	dbSession, err := s.client.Session.Query().Where(entSession.ID(uuid)).WithUser().First(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, fmt.Errorf("checking is session exists: %w", err)
		}
		// Else session does not exist
		return nil, fmt.Errorf("session does not exist")
	}
	// Validate the session (i.e. has it expired).
	// Calculate the time when it expires and check if current time is after the
	// expiration time. If so, it has expired.
	expiresAt := dbSession.CreateTime.Add(time.Duration(dbSession.Duration))
	if time.Now().After(expiresAt) {
		// Delete the session and return error
		if err := s.client.Session.DeleteOne(dbSession).Exec(s.ctx); err != nil {
			return nil, fmt.Errorf("deleting expired session: %w", err)
		}
		return nil, fmt.Errorf("session has expired")
	}
	dbUser, err := s.client.User.Query().
		Where(entUser.ID(dbSession.Edges.User.ID)).
		WithGroups().
		First(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, fmt.Errorf("checking is session exists: %w", err)
		}
		// Else session does not exist
		return nil, fmt.Errorf("session does not exist")
	}

	return dbUserToAPI(dbUser), nil
}

func (s *Store) EndSession(uuid uuid.UUID) error {
	err := s.client.Session.DeleteOneID(uuid).Exec(s.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// It's fine, the session never existed so it's gone now
			return nil
		}
		return fmt.Errorf("deleting session: %w", err)
	}
	return nil
}

func (s *Store) createUpdateUser(user *oapi.User) (*ent.User, error) {
	var (
		dbUser *ent.User
		err    error
	)
	// Get database groups from user
	dbGroups, err := s.createReadGroups(user.Groups)
	if err != nil {
		return nil, fmt.Errorf("creating user groups: %w", err)
	}
	dbUser, err = s.client.User.Query().
		Where(entUser.And(entUser.Sub(user.Sub), entUser.Iss(user.Iss))).
		First(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, fmt.Errorf("checking if user exists: %w", err)
		}
		// Else we need to create the user
		dbUser, err = s.client.User.Create().
			SetSub(user.Sub).
			SetIss(user.Iss).
			SetName(user.Name).
			SetNillableEmail(user.Email).
			SetNillablePicture(user.Picture).
			AddGroups(dbGroups...).
			Save(s.ctx)
		if err != nil {
			return nil, fmt.Errorf("creating new user: %w", err)
		}

		return dbUser, nil
	}

	// Get external groups the user is currently a member of
	extGroups, err := s.client.User.Query().Where(entUser.ID(dbUser.ID)).QueryGroups().Where(group.IsExternal(true)).All(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("getting user external groups: %w", err)
	}

	// If user already exists, update info for that user.
	dbUser, err = s.client.User.UpdateOne(dbUser).
		SetName(user.Name).
		SetNillableEmail(user.Email).
		SetNillablePicture(user.Picture).
		RemoveGroups(extGroups...).
		AddGroups(dbGroups...).
		Save(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	return dbUser, nil
}

func (s *Store) createReadGroups(groups []string) ([]*ent.Group, error) {
	var dbGroups []*ent.Group

	bulkGroups := make([]*ent.GroupCreate, 0)
	for _, name := range groups {
		entGroup, err := s.client.Group.Query().Where(group.Name(name)).First(s.ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return nil, fmt.Errorf("checking if group exists: %s: %w", name, err)
			}
			// Else let's add the group for creation
			bulkGroups = append(bulkGroups, s.client.Group.Create().SetName(name).SetIsExternal(true))
			continue
		}
		// If the group already exists, just add it to the list of groups to return
		dbGroups = append(dbGroups, entGroup)
	}

	newGroups, err := s.client.Group.CreateBulk(bulkGroups...).Save(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("creating groups: %w", err)
	}

	return append(dbGroups, newGroups...), nil
}
