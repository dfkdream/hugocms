package user

import "testing"

func TestMatchPermission(t *testing.T) {
	for idx, v := range []struct {
		pattern string
		target  string
		result  bool
	}{
		{"", "", true},
		{"hugocms:list", "*", true},
		{"hugocms:list", "hugocms:*", true},
		{"hugocms:list", "", false},
		{"hugocms:list", "hugocms:user", false},
		{"hugocms:list:read", "hugocms:list:*", true},
		{"hugocms:list:read", "hugocms:*:read", true},
		{"hugocms:list:read", "hugocms:*:write", false},
	} {
		if result := matchPermission(v.pattern, v.target); result != v.result {
			t.Errorf("Test %d failed. expected: %v, result:%v", idx, v.result, result)
		}
	}
}
