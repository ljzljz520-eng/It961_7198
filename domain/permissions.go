package domain

import "strings"

type Role string

const (
	RoleCoach    Role = "coach"
	RoleReviewer Role = "reviewer"
	RoleManager  Role = "manager"
)

type AccessPolicy struct {
	Campus       string
	Roles        []Role
	AllowArchive bool
	AllowPublish bool
}

func (p AccessPolicy) Allows(role Role, campus string, action string) bool {
	if p.Campus != "" && !strings.EqualFold(p.Campus, campus) {
		return false
	}
	roleAllowed := false
	for _, candidate := range p.Roles {
		if candidate == role {
			roleAllowed = true
			break
		}
	}
	if !roleAllowed {
		return false
	}
	switch action {
	case "archive":
		return p.AllowArchive
	case "publish":
		return p.AllowPublish
	case "read", "search":
		return true
	default:
		return false
	}
}

func (p AccessPolicy) Normalize() AccessPolicy {
	p.Campus = strings.TrimSpace(strings.ToLower(p.Campus))
	seen := map[Role]bool{}
	roles := make([]Role, 0, len(p.Roles))
	for _, role := range p.Roles {
		if role == "" || seen[role] {
			continue
		}
		seen[role] = true
		roles = append(roles, role)
	}
	p.Roles = roles
	return p
}

func DefaultCoachPolicy(campus string) AccessPolicy {
	return AccessPolicy{Campus: campus, Roles: []Role{RoleCoach}, AllowArchive: false, AllowPublish: false}.Normalize()
}
