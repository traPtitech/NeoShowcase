package domain

import (
	"github.com/samber/lo"
	"github.com/samber/oops"
)

// ApplicationGroup is a set of applications that share owners and may nest their websites with each other.
// Every application belongs to exactly one group.
type ApplicationGroup struct {
	ID   string
	Name string
	// Implicit is true for a group the system creates for a single application that is not in any group a user created.
	// Implicit groups are hidden from users, and their Name is empty.
	Implicit bool
	OwnerIDs []string
}

// NewApplicationGroup returns a group a user created.
func NewApplicationGroup(name string, ownerIDs []string) *ApplicationGroup {
	return &ApplicationGroup{
		ID:       NewID(),
		Name:     name,
		OwnerIDs: ownerIDs,
	}
}

// NewImplicitApplicationGroup returns a group for a single application that is not in any group a user created.
func NewImplicitApplicationGroup(ownerIDs []string) *ApplicationGroup {
	return &ApplicationGroup{
		ID:       NewID(),
		Implicit: true,
		OwnerIDs: ownerIDs,
	}
}

func (g *ApplicationGroup) Validate() error {
	if !g.Implicit && g.Name == "" {
		return oops.New("name is required")
	}
	if len(g.OwnerIDs) == 0 {
		return oops.New("owner_ids cannot be empty")
	}
	return nil
}

func (g *ApplicationGroup) IsOwner(user *User) bool {
	return user.Admin || lo.Contains(g.OwnerIDs, user.ID)
}
