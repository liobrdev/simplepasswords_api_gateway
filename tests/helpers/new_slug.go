package helpers

import (
	"testing"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func NewSlug(t *testing.T) string {
	if slug, err := utils.GenerateSlug(32); err != nil {
		t.Errorf("Failed to generate new slug: %s", err.Error())
		panic(err)
	} else {
		return slug
	}
}
