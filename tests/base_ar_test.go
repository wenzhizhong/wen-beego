package tests

import (
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models_ar/base_ar"
	"fmt"
	"testing"
)

func TestCheckUserHasUnits(t *testing.T) {
	userId := "019917fa-ee48-7ac5-bdec-3bb11a2d54d9"
	data, err := base_ar.GetUserUnitList(userId, &models.Plat{}, &models.PlatUser{})
	for i := 0; i < len(data); i++ {
		fmt.Println(data[i])
	}

	fmt.Println("====== end ========", err)
}
