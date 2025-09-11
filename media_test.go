package zabbix_test

import (
	"testing"

	zapi "github.com/claranet/go-zabbix-api"
)

func TestMediaTypeGet(t *testing.T) {
	api := testGetAPI(t)

	params := zapi.Params{}
	mediatypes, err := api.MediaTypeGet(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(mediatypes) == 0 {
		t.Skip("No media types found to test")
	}

	// Test specific media type retrieval by ID
	params = zapi.Params{
		"mediatypeids": []string{mediatypes[0].MediaTypeID},
	}
	specificMediatypes, err := api.MediaTypeGet(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(specificMediatypes) != 1 {
		t.Errorf("Expected 1 media type, got %d", len(specificMediatypes))
	}
	if specificMediatypes[0].MediaTypeID != mediatypes[0].MediaTypeID {
		t.Errorf("Media type ID mismatch: expected %s, got %s",
			mediatypes[0].MediaTypeID, specificMediatypes[0].MediaTypeID)
	}
}

func TestMediaTypeGetWithFilter(t *testing.T) {
	api := testGetAPI(t)

	// Test filter by type
	params := zapi.Params{
		"filter": map[string]interface{}{
			"type": zapi.MediaTypeEmail,
		},
	}
	mediatypes, err := api.MediaTypeGet(params)
	if err != nil {
		t.Fatal(err)
	}

	// Verify all returned media types are email type
	for _, mt := range mediatypes {
		if mt.Type != zapi.MediaTypeEmail {
			t.Errorf("Expected media type %d, got %d", zapi.MediaTypeEmail, mt.Type)
		}
	}
}
