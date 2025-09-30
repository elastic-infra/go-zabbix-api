package zabbix_test

import (
	"testing"

	zapi "github.com/claranet/go-zabbix-api"
)

func TestRoleGet(t *testing.T) {
	api := testGetAPI(t)

	params := zapi.Params{}
	roles, err := api.RolesGet(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) == 0 {
		t.Skip("No roles found to test")
	}

	// Test specific role retrieval by ID
	params = zapi.Params{
		"roleids": []string{roles[0].RoleID},
	}
	specificRoles, err := api.RolesGet(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(specificRoles) != 1 {
		t.Errorf("Expected 1 role, got %d", len(specificRoles))
	}
	if specificRoles[0].RoleID != roles[0].RoleID {
		t.Errorf("Role ID mismatch: expected %s, got %s",
			roles[0].RoleID, specificRoles[0].RoleID)
	}
}

func TestRoleGetWithFilter(t *testing.T) {
	api := testGetAPI(t)

	// Test filter by type (User type = 1)
	params := zapi.Params{
		"filter": map[string]interface{}{
			"type": zapi.RoleTypeUser,
		},
	}
	roles, err := api.RolesGet(params)
	if err != nil {
		t.Fatal(err)
	}

	// Verify all returned roles are of User type
	for _, role := range roles {
		if role.Type != zapi.RoleTypeUser {
			t.Errorf("Expected role type %d, got %d for role %s",
				zapi.RoleTypeUser, role.Type, role.Name)
		}
	}
}

func TestRoleTypes(t *testing.T) {
	// Test role type constants
	if zapi.RoleTypeUser != 1 {
		t.Errorf("Expected RoleTypeUser to be 1, got %d", zapi.RoleTypeUser)
	}
	if zapi.RoleTypeAdmin != 2 {
		t.Errorf("Expected RoleTypeAdmin to be 2, got %d", zapi.RoleTypeAdmin)
	}
	if zapi.RoleTypeSuperAdmin != 3 {
		t.Errorf("Expected RoleTypeSuperAdmin to be 3, got %d", zapi.RoleTypeSuperAdmin)
	}
}
