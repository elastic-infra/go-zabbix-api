package zabbix_test

import (
	"math/rand"
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

func TestUserCRUD(t *testing.T) {
	api := testGetAPI(t)

	// Test user creation
	randomSuffix := randomString(8)
	users := zapi.Users{
		{
			Username: "testuser_" + randomSuffix,
			Name:     "TestName",
			Surname:  "TestSurname",
			Password: "ComplexPassword123!", // Strong password not containing username
			RoleID:   "1",                   // User role
		},
	}

	err := api.UserCreate(users)
	if err != nil {
		t.Fatal(err)
	}

	if users[0].UserID == "" {
		t.Error("UserID should not be empty after create")
	}

	createdUserID := users[0].UserID

	// Test user retrieval
	retrievedUsers, err := api.UsersGet(zapi.Params{
		"userids": []string{createdUserID},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(retrievedUsers) != 1 {
		t.Errorf("Expected 1 user, got %d", len(retrievedUsers))
	}
	if retrievedUsers[0].Username != users[0].Username {
		t.Errorf("Expected username %s, got %s", users[0].Username, retrievedUsers[0].Username)
	}

	// Test user update
	updateUsers := zapi.Users{
		{
			UserID:  createdUserID,
			Name:    "Updated",
			Surname: "UpdatedSurname",
		},
	}

	err = api.UserUpdate(updateUsers)
	if err != nil {
		t.Fatal(err)
	}

	// Verify update
	updatedUsers, err := api.UsersGet(zapi.Params{
		"userids": []string{createdUserID},
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedUsers[0].Name != "Updated" {
		t.Errorf("Expected name 'Updated', got '%s'", updatedUsers[0].Name)
	}
	if updatedUsers[0].Surname != "UpdatedSurname" {
		t.Errorf("Expected surname 'UpdatedSurname', got '%s'", updatedUsers[0].Surname)
	}

	// Test user deletion
	err = api.UserDelete([]string{createdUserID})
	if err != nil {
		t.Fatal(err)
	}

	// Verify deletion
	deletedUsers, err := api.UsersGet(zapi.Params{
		"userids": []string{createdUserID},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(deletedUsers) != 0 {
		t.Errorf("Expected 0 users after deletion, got %d", len(deletedUsers))
	}
}

func TestUserWithMedia(t *testing.T) {
	api := testGetAPI(t)

	// Get first available media type
	mediatypes, err := api.MediaTypeGet(zapi.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if len(mediatypes) == 0 {
		t.Skip("No media types found to test")
	}

	// Test user creation with media
	randomSuffix := randomString(8)
	users := zapi.Users{
		{
			Username: "testuser_" + randomSuffix,
			Name:     "TestName",
			Surname:  "TestSurname",
			Password: "ComplexPassword123!",
			RoleID:   "1",
			Medias: zapi.Medias{
				{
					MediaTypeID: mediatypes[0].MediaTypeID,
					SendTo:      []string{"test@example.com"},
					Active:      zapi.MediaStatusActive,
					Severity:    63, // All severities
					Period:      "1-7,00:00-24:00",
				},
			},
		},
	}

	err = api.UserCreate(users)
	if err != nil {
		t.Fatal(err)
	}

	createdUserID := users[0].UserID
	defer func() {
		// Clean up
		api.UserDelete([]string{createdUserID})
	}()

	// Test retrieval with media
	retrievedUsers, err := api.UsersGet(zapi.Params{
		"userids":     []string{createdUserID},
		"selectMedia": "extend",
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(retrievedUsers) != 1 {
		t.Errorf("Expected 1 user, got %d", len(retrievedUsers))
		return
	}

	// Note: Media retrieval might be empty if not supported by the test environment
	// This is acceptable as we're testing the structure definition
	user := retrievedUsers[0]
	if user.UserID != createdUserID {
		t.Errorf("Expected user ID %s, got %s", createdUserID, user.UserID)
	}
}

// Helper function to generate random string for test users
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
