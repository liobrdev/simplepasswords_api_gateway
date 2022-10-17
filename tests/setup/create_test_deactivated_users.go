package setup

import (
	"testing"
	"time"

	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/helpers"
)

func createTestDeactivatedUsers(
	t *testing.T,
	dbs *databases.Databases,
) (deactivatedUsers []models.DeactivatedUser) {
	deactivatedUsers = []models.DeactivatedUser{
		{
			Slug:         helpers.NewSlug(t),
			Name:         "Doris Jones",
			EmailAddress: "doris.jones@email.co",
			CreatedAt:    time.Date(2022, time.April, 15, 12, 30, 0, 0, time.Local),
			UpdatedAt:    time.Date(2022, time.April, 15, 13, 30, 0, 0, time.Local),
		},
		{
			Slug:         helpers.NewSlug(t),
			Name:         "Dean Jones",
			EmailAddress: "dean.jones@email.co",
			CreatedAt:    time.Date(2022, time.April, 15, 12, 30, 0, 0, time.Local),
			UpdatedAt:    time.Date(2022, time.April, 15, 13, 30, 0, 0, time.Local),
		},
	}

	if result := dbs.ApiGateway.Create(&deactivatedUsers); result.Error != nil {
		t.Fatalf("Create test deactivated users failed: %s", result.Error.Error())
	}

	return
}
