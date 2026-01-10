package base_model

import (
	"time"
)

type UnitMenu struct {
	Id              string    `json:"id" gorm:"column:id;primaryKey;comment:ID"`
	ParentId        string    `json:"parentId" gorm:"column:parent_id;comment:上级ID"`
	UnitId          string    `json:"unitId" gorm:"->"`
	MenuType        int       `json:"menuType" gorm:"column:menu_type;not null;comment: 菜单类型:0代表菜单、1代表iframe、2代表外链、3代表按钮、4所需额外接口"`
	Title           string    `json:"title" gorm:"column:title;size:255;comment: 菜单标题"`
	Name            string    `json:"name" gorm:"column:name;size:255;comment:	菜单名称"`
	Path            string    `json:"path" gorm:"column:path;size:255;comment: 菜单路径"`
	Component       string    `json:"component" gorm:"column:component;size:255;comment: 组件"`
	Rank            *int      `json:"rank" gorm:"column:rank;comment: 排序"`
	Redirect        string    `json:"redirect" gorm:"column:redirect;size:255;comment: 跳转地址"`
	Icon            string    `json:"icon" gorm:"column:icon;size:255;comment: 图标"`
	ExtraIcon       string    `json:"extraIcon" gorm:"column:extra_icon;size:255;comment: 额外图标"`
	EnterTransition string    `json:"enterTransition" gorm:"column:enter_transition;size:255;comment: 进入过渡"`
	LeaveTransition string    `json:"leaveTransition" gorm:"column:leave_transition;size:255;comment: 离开过渡"`
	ActivePath      string    `json:"activePath" gorm:"column:active_path;size:255;comment: 激活路径"`
	Auths           string    `json:"auths" gorm:"column:auths;comment:	权限"`
	FrameSrc        string    `json:"frameSrc" gorm:"column:frame_src;size:255;comment: iframe 地址"`
	FrameLoading    bool      `json:"frameLoading" gorm:"column:frame_loading;default:true;comment: 是否加载 iframe"`
	KeepAlive       bool      `json:"keepAlive" gorm:"column:keep_alive;default:false;comment: 是否缓存"`
	HiddenTag       bool      `json:"hiddenTag" gorm:"column:hidden_tag;default:false;comment: 是否隐藏Tag"`
	FixedTag        bool      `json:"fixedTag" gorm:"column:fixed_tag;default:false;comment: 是否固定Tag"`
	ShowLink        bool      `json:"showLink" gorm:"column:show_link;default:true;comment: 是否显示Link"`
	ShowParent      bool      `json:"showParent" gorm:"column:show_parent;default:false;comment: 是否显示Parent"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime;comment: 创建时间"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime;comment: 修改时间"`
	Deleted         int       `json:"deleted" gorm:"type:int4;default:0;comment:是否删除"`
	Clone           int       `json:"clone" gorm:"type:int2;default:0;comment:允许克隆：0否1是"`
}

var UNIT_MENU_TYPE_MENU = 0
var UNIT_MENU_TYPE_IFRAME = 1
var UNIT_MENU_TYPE_EXTERNAL = 2
var UNIT_MENU_TYPE_BUTTON = 3
var UNIT_MENU_TYPE_OTHER_API = 4
var UNIT_MENU_TYPE_MAP = map[int]string{
	UNIT_MENU_TYPE_MENU:      "菜单",
	UNIT_MENU_TYPE_IFRAME:    "iframe",
	UNIT_MENU_TYPE_EXTERNAL:  "外链",
	UNIT_MENU_TYPE_BUTTON:    "按钮",
	UNIT_MENU_TYPE_OTHER_API: "所需额外接口",
}
