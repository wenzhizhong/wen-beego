package auth

import "WenBeego/apps/common/models/base_model"

// 获取所有模型常量
func GetAllModelConstant() map[string]interface{} {
	return map[string]interface{}{
		"unit_status_map":       base_model.UNIT_STATUS_MAP,
		"unit_role_status_map":  base_model.UNIT_ROLE_STATUS_MAP,
		"unit_user_profile_map": base_model.UNIT_USER_PROFILE_MAP,
		"unit_user_default_map": base_model.UNIT_USER_DEFAULT_MAP,
		"unit_gender_map":       base_model.UNIT_GENDER_MAP,
		"unit_card_type_map":    base_model.UNIT_CARD_TYPE_MAP,
		"unit_user_source_map":  base_model.UNIT_USER_SOURCE_MAP,
	}
}
