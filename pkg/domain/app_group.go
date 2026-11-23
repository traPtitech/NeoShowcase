package domain

import (
	"unicode/utf8"

	"github.com/samber/lo"
	"github.com/samber/oops"

	"github.com/traPtitech/neoshowcase/pkg/util/optional"
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

// maxApplicationGroupNameLength is the length of the name column in characters.
const maxApplicationGroupNameLength = 100

func (g *ApplicationGroup) Validate() error {
	if !g.Implicit && g.Name == "" {
		return oops.New("name is required")
	}
	if utf8.RuneCountInString(g.Name) > maxApplicationGroupNameLength {
		return oops.Errorf("name must be at most %d characters", maxApplicationGroupNameLength)
	}
	if len(g.OwnerIDs) == 0 {
		return oops.New("owner_ids cannot be empty")
	}
	return nil
}

// VisibleID returns the ID users know the group by, or none for an implicit group, which users do not see.
func (g *ApplicationGroup) VisibleID() optional.Of[string] {
	if g.Implicit {
		return optional.None[string]()
	}
	return optional.From(g.ID)
}

func (g *ApplicationGroup) IsOwner(user *User) bool {
	return user.Admin || lo.Contains(g.OwnerIDs, user.ID)
}
