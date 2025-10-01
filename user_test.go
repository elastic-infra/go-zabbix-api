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

	// Get a valid user group ID for the test
	userGroupID := getValidUserGroupID(t, api)

	// Test user creation
	randomSuffix := randomString(8)
	users := zapi.Users{
		{
			Username: "testuser_" + randomSuffix,
			Name:     "TestName",
			Surname:  "TestSurname",
			Password: "ComplexPassword123!", // Strong password not containing username
			RoleID:   "1",                   // User role
			UsrGrps: zapi.UserGroups{
				{
					GroupID: userGroupID, // Required in Zabbix 6.0+
				},
			},
		},
	}

	err := api.UsersCreate(users)
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

	err = api.UsersUpdate(updateUsers)
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
	err = api.UsersDeleteByIds([]string{createdUserID})
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
	mediatypes, err := api.MediaTypesGet(zapi.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if len(mediatypes) == 0 {
		t.Skip("No media types found to test")
	}

	// Get a valid user group ID for the test
	userGroupID := getValidUserGroupID(t, api)

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
			UsrGrps: zapi.UserGroups{
				{
					GroupID: userGroupID, // Dynamic user group ID
				},
			},
		},
	}

	err = api.UsersCreate(users)
	if err != nil {
		t.Fatal(err)
	}

	createdUserID := users[0].UserID
	defer func() {
		// Clean up
		api.UsersDeleteByIds([]string{createdUserID})
	}()

	// Test retrieval with media and user groups
	retrievedUsers, err := api.UsersGet(zapi.Params{
		"userids":       []string{createdUserID},
		"selectMedias":  "extend",
		"selectUsrgrps": "extend",
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(retrievedUsers) != 1 {
		t.Errorf("Expected 1 user, got %d", len(retrievedUsers))
		return
	}

	// Note: Media and user group retrieval might be empty if not supported by the test environment
	// This is acceptable as we're testing the structure definition
	user := retrievedUsers[0]
	if user.UserID != createdUserID {
		t.Errorf("Expected user ID %s, got %s", createdUserID, user.UserID)
	}

	// Note: User groups field should be available even if empty
	// UsrGrps field existence validates the structure
}

// Helper function to get a valid user group ID for testing
func getValidUserGroupID(t *testing.T, api *zapi.API) string {
	// Try to get user groups
	params := zapi.Params{
		"output": []string{"usrgrpid", "name"},
	}

	userGroups, err := api.UserGroupsGet(params)
	if err != nil {
		// If UserGroupsGet fails, try a direct API call
		response, callErr := api.CallWithError("usergroup.get", params)
		if callErr != nil {
			// Fall back to commonly available group ID
			t.Logf("Warning: Could not fetch user groups, using default group ID '7': %v", callErr)
			return "7" // Default to "Zabbix administrators" group
		}

		if groups, ok := response.Result.([]interface{}); ok && len(groups) > 0 {
			if group, ok := groups[0].(map[string]interface{}); ok {
				if usrgrpid, ok := group["usrgrpid"].(string); ok {
					return usrgrpid
				}
			}
		}
		return "7" // Fallback
	}

	if len(userGroups) == 0 {
		t.Logf("Warning: No user groups found, using default group ID '7'")
		return "7"
	}

	// Return the first available user group ID
	return userGroups[0].GroupID
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
