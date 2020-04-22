package user

import (
	"strings"
)

func matchPermission(pattern, target string) bool {
	p1s := strings.Split(pattern, ":")
	p2s := strings.Split(target, ":")

	if len(p2s) > len(p1s) {
		return false
	}

	for idx, v := range p2s {
		if v == "*" {
			continue
		}

		if v != p1s[idx] {
			return false
		}
	}

	return true
}

func (m *User) HasPermission(permission string) bool {
	if m == nil {
		return false
	}

	for _, p := range m.GetPermissions() {
		if matchPermission(permission, p) {
			return true
		}
	}
	return false
}
