package base_model

import (
	"time"
)

type UnitMenu struct {
	Id              string    `gorm:"column:id;primaryKey" json:"id"`
	ParentId        string    `gorm:"column:parent_id" json:"parentId"`
	UnitId          string    `gorm:"column:unit_id" json:"unitId"`
	MenuType        int       `gorm:"column:menu_type;not null" json:"menuType"`
	Title           string    `gorm:"column:title;size:255" json:"title"`
	Name            string    `gorm:"column:name;size:255" json:"name"`
	Path            string    `gorm:"column:path;size:255" json:"path"`
	Component       string    `gorm:"column:component;size:255" json:"component"`
	Rank            *int      `gorm:"column:rank" json:"rank"`
	Redirect        string    `gorm:"column:redirect;size:255" json:"redirect"`
	Icon            string    `gorm:"column:icon;size:255" json:"icon"`
	ExtraIcon       string    `gorm:"column:extra_icon;size:255" json:"extraIcon"`
	EnterTransition string    `gorm:"column:enter_transition;size:255" json:"enterTransition"`
	LeaveTransition string    `gorm:"column:leave_transition;size:255" json:"leaveTransition"`
	ActivePath      string    `gorm:"column:active_path;size:255" json:"activePath"`
	Auths           string    `gorm:"column:auths" json:"auths"`
	FrameSrc        string    `gorm:"column:frame_src;size:255" json:"frameSrc"`
	FrameLoading    bool      `gorm:"column:frame_loading;default:true" json:"frameLoading"`
	KeepAlive       bool      `gorm:"column:keep_alive;default:false" json:"keepAlive"`
	HiddenTag       bool      `gorm:"column:hidden_tag;default:false" json:"hiddenTag"`
	FixedTag        bool      `gorm:"column:fixed_tag;default:false" json:"fixedTag"`
	ShowLink        bool      `gorm:"column:show_link;default:true" json:"showLink"`
	ShowParent      bool      `gorm:"column:show_parent;default:false" json:"showParent"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Deleted         int       `json:"deleted" gorm:"type:int4;default:0;comment:是否删除"`
	Clone           int       `json:"clone" gorm:"type:int2;default:0;comment:允许克隆：0否1是"`
}
