package zabbix_test

import (
	"testing"

	zapi "github.com/claranet/go-zabbix-api"
)

func TestUsersGet(t *testing.T) {
	api := testGetAPI(t)

	params := zapi.Params{
		"filter": map[string]interface{}{
			"alias":    "Admin", // Under 5.4
			"username": "Admin", // 5.4 or higher
		},
	}
	users, err := api.UsersGet(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
	
	user := users[0]
	
	// Test basic fields
	if user.Username != "Admin" {
		t.Errorf("Expected username 'Admin', got '%s'", user.Username)
	}
	if user.UserID == "" {
		t.Error("UserID should not be empty")
	}
	
	// Test extended fields
	if user.RoleID == "" {
		t.Error("RoleID should not be empty")
	}
	if user.Lang == "" {
		t.Error("Lang should not be empty")
	}
	if user.Theme == "" {
		t.Error("Theme should not be empty")
	}
	if user.Timezone == "" {
		t.Error("Timezone should not be empty")
	}
}

func TestUsersGetAll(t *testing.T) {
	api := testGetAPI(t)

	params := zapi.Params{}
	users, err := api.UsersGet(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) == 0 {
		t.Error("Expected at least one user")
	}
	
	// Test that all users have required fields
	for _, user := range users {
		if user.UserID == "" {
			t.Errorf("User %s has empty UserID", user.Username)
		}
		if user.Username == "" {
			t.Errorf("User %s has empty Username", user.UserID)
		}
		if user.RoleID == "" {
			t.Errorf("User %s has empty RoleID", user.Username)
		}
	}
}

func TestUsersGetWithSpecificFields(t *testing.T) {
	api := testGetAPI(t)

	params := zapi.Params{
		"output": []string{"userid", "username", "roleid", "name", "surname"},
		"filter": map[string]interface{}{
			"username": "Admin",
		},
	}
	users, err := api.UsersGet(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
	
	user := users[0]
	if user.Username != "Admin" {
		t.Errorf("Expected username 'Admin', got '%s'", user.Username)
	}
	if user.RoleID == "" {
		t.Error("RoleID should not be empty")
	}
}

func TestMediaStatus(t *testing.T) {
	// Test MediaStatus constants
	if zapi.MediaStatusActive != 0 {
		t.Errorf("Expected MediaStatusActive to be 0, got %d", zapi.MediaStatusActive)
	}
	if zapi.MediaStatusDisabled != 1 {
		t.Errorf("Expected MediaStatusDisabled to be 1, got %d", zapi.MediaStatusDisabled)
	}
}

func TestUserTypes(t *testing.T) {
	// Test UserType constants (legacy compatibility)
	if zapi.ZabbixUser != 0 {
		t.Errorf("Expected ZabbixUser to be 0, got %d", zapi.ZabbixUser)
	}
	if zapi.ZabbixAdmin != 1 {
		t.Errorf("Expected ZabbixAdmin to be 1, got %d", zapi.ZabbixAdmin)
	}
	if zapi.ZabbixSuperAdmin != 2 {
		t.Errorf("Expected ZabbixSuperAdmin to be 2, got %d", zapi.ZabbixSuperAdmin)
	}
}
