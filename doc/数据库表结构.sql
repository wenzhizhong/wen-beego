-- public.config definition

-- Drop table

-- DROP TABLE public.config;

CREATE TABLE public.config (
	id bpchar(36) NOT NULL, -- ID
	"name" varchar(255) NOT NULL, -- 配置名称
	value text NULL, -- 配置值
	description text NULL, -- 备注
	value_type varchar(20) NULL DEFAULT 'string'::character varying, -- 数据类型
	category varchar(50) NULL, -- 分类
	is_readonly bool NULL DEFAULT false, -- 是否可修改
	"version" int4 NULL DEFAULT 1, -- 版本
	created_at timestamptz NULL DEFAULT now(), -- 创建时间
	updated_at timestamptz NULL DEFAULT now(), -- 更新时间
	created_by varchar(50) NULL, -- 创建者
	updated_by varchar(50) NULL, -- 更新人
	deleted bool NULL DEFAULT false, -- 是否删除
	CONSTRAINT config_name_key UNIQUE (name),
	CONSTRAINT config_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_config_category ON public.config USING btree (category);
CREATE INDEX idx_config_deleted ON public.config USING btree (deleted);
COMMENT ON TABLE public.config IS '配置参数表';

-- Column comments

COMMENT ON COLUMN public.config.id IS 'ID';
COMMENT ON COLUMN public.config."name" IS '配置名称';
COMMENT ON COLUMN public.config.value IS '配置值';
COMMENT ON COLUMN public.config.description IS '备注';
COMMENT ON COLUMN public.config.value_type IS '数据类型';
COMMENT ON COLUMN public.config.category IS '分类';
COMMENT ON COLUMN public.config.is_readonly IS '是否可修改';
COMMENT ON COLUMN public.config."version" IS '版本';
COMMENT ON COLUMN public.config.created_at IS '创建时间';
COMMENT ON COLUMN public.config.updated_at IS '更新时间';
COMMENT ON COLUMN public.config.created_by IS '创建者';
COMMENT ON COLUMN public.config.updated_by IS '更新人';
COMMENT ON COLUMN public.config.deleted IS '是否删除';


-- public.file definition

-- Drop table

-- DROP TABLE public.file;

CREATE TABLE public.file (
	id bpchar(36) NOT NULL, -- ID
	batch_id varchar(32) NULL,
	real_name varchar(255) NULL, -- 文件真实的名称
	"name" varchar(255) NULL, -- 文件名
	suffix varchar(255) NULL, -- 后缀
	"path" varchar(512) NULL, -- 路径
	"type" varchar(255) NULL, -- 类型
	"size" varchar(100) NULL, -- 大小
	file_md5 varchar(255) NULL, -- md5校验码
	unit_id varchar(100) NULL, -- 组织单位id
	unit_type varchar(100) NULL, -- 组织单位类别
	create_by varchar(255) NULL, -- 创建者
	update_by varchar(255) NULL, -- 更新者
	createtime timestamp NULL, -- 创建时间
	updatetime timestamp NULL, -- 更新时间
	CONSTRAINT file_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public.file IS '文件上传信息';

-- Column comments

COMMENT ON COLUMN public.file.id IS 'ID';
COMMENT ON COLUMN public.file.real_name IS '文件真实的名称';
COMMENT ON COLUMN public.file."name" IS '文件名';
COMMENT ON COLUMN public.file.suffix IS '后缀';
COMMENT ON COLUMN public.file."path" IS '路径';
COMMENT ON COLUMN public.file."type" IS '类型';
COMMENT ON COLUMN public.file."size" IS '大小';
COMMENT ON COLUMN public.file.file_md5 IS 'md5校验码';
COMMENT ON COLUMN public.file.unit_id IS '组织单位id';
COMMENT ON COLUMN public.file.unit_type IS '组织单位类别';
COMMENT ON COLUMN public.file.create_by IS '创建者';
COMMENT ON COLUMN public.file.update_by IS '更新者';
COMMENT ON COLUMN public.file.createtime IS '创建时间';
COMMENT ON COLUMN public.file.updatetime IS '更新时间';


-- public.file_slice definition

-- Drop table

-- DROP TABLE public.file_slice;

CREATE TABLE public.file_slice (
	id bpchar(36) NOT NULL, -- ID
	file_name varchar(255) NULL, -- 文件名
	current_size int8 NULL, -- 当前分片大小
	slice_size int8 NULL, -- 分片大小
	total_size int8 NULL, -- 文件总大小
	slice_index int8 NULL, -- 当前分片，从1开始
	slice_total int8 NULL, -- 总分片数
	file_md5 varchar(45) NULL, -- 文件标识
	file_path varchar(255) NULL, -- 路径
	create_by varchar(255) NULL, -- 创建人
	createtime timestamp NULL, -- 创建时间
	updatetime timestamp NULL, -- 更新时间
	CONSTRAINT file_slice_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public.file_slice IS '文件上传分片信息';

-- Column comments

COMMENT ON COLUMN public.file_slice.id IS 'ID';
COMMENT ON COLUMN public.file_slice.file_name IS '文件名';
COMMENT ON COLUMN public.file_slice.current_size IS '当前分片大小';
COMMENT ON COLUMN public.file_slice.slice_size IS '分片大小';
COMMENT ON COLUMN public.file_slice.total_size IS '文件总大小';
COMMENT ON COLUMN public.file_slice.slice_index IS '当前分片，从1开始';
COMMENT ON COLUMN public.file_slice.slice_total IS '总分片数';
COMMENT ON COLUMN public.file_slice.file_md5 IS '文件标识';
COMMENT ON COLUMN public.file_slice.file_path IS '路径';
COMMENT ON COLUMN public.file_slice.create_by IS '创建人';
COMMENT ON COLUMN public.file_slice.createtime IS '创建时间';
COMMENT ON COLUMN public.file_slice.updatetime IS '更新时间';


-- public.mchnt definition

-- Drop table

-- DROP TABLE public.mchnt;

CREATE TABLE public.mchnt (
	id bpchar(36) NOT NULL, -- ID
	pid bpchar(36) NOT NULL DEFAULT ''::bpchar, -- 上级组织id
	logo varchar(512) NULL DEFAULT ''::character varying, -- logo
	"name" varchar(100) NOT NULL, -- 单位名称
	code varchar(100) NOT NULL, -- 组织机构代码
	corporation varchar(100) NOT NULL, -- 法人
	license varchar(512) NOT NULL DEFAULT ''::character varying, -- 营业执照
	address varchar(255) NULL DEFAULT ''::character varying, -- 地址
	status int4 NOT NULL DEFAULT 0, -- 0未审核，1审核通过，2审核不通过，3禁用
	plat_status int4 NOT NULL DEFAULT 0, -- 平台审核状态(同上)
	deleted int4 NULL DEFAULT 0, -- 是否删除：0否1是
	created_at int8 NULL, -- 创建时间
	updated_at int8 NULL, -- 更新时间
	deleted_at int8 NULL, -- 删除时间
	created_by bpchar(36) NOT NULL DEFAULT ''::bpchar, -- 创建人
	updated_by bpchar(36) NULL, -- 更新人
	deleted_by bpchar(36) NULL, -- 删除人
	sort int4 NOT NULL DEFAULT 0, -- 排序
	CONSTRAINT org_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public.mchnt IS '组织机构表-Merchant';

-- Column comments

COMMENT ON COLUMN public.mchnt.id IS 'ID';
COMMENT ON COLUMN public.mchnt.pid IS '上级组织id';
COMMENT ON COLUMN public.mchnt.logo IS 'logo';
COMMENT ON COLUMN public.mchnt."name" IS '单位名称';
COMMENT ON COLUMN public.mchnt.code IS '组织机构代码';
COMMENT ON COLUMN public.mchnt.corporation IS '法人';
COMMENT ON COLUMN public.mchnt.license IS '营业执照';
COMMENT ON COLUMN public.mchnt.address IS '地址';
COMMENT ON COLUMN public.mchnt.status IS '0未审核，1审核通过，2审核不通过，3禁用';
COMMENT ON COLUMN public.mchnt.plat_status IS '平台审核状态(同上)';
COMMENT ON COLUMN public.mchnt.deleted IS '是否删除：0否1是';
COMMENT ON COLUMN public.mchnt.created_at IS '创建时间';
COMMENT ON COLUMN public.mchnt.updated_at IS '更新时间';
COMMENT ON COLUMN public.mchnt.deleted_at IS '删除时间';
COMMENT ON COLUMN public.mchnt.created_by IS '创建人';
COMMENT ON COLUMN public.mchnt.updated_by IS '更新人';
COMMENT ON COLUMN public.mchnt.deleted_by IS '删除人';
COMMENT ON COLUMN public.mchnt.sort IS '排序';


-- public.mchnt_api_statistics definition

-- Drop table

-- DROP TABLE public.mchnt_api_statistics;

CREATE TABLE public.mchnt_api_statistics (
	id bpchar(36) NOT NULL, -- ID
	perms_id bpchar(36) NULL DEFAULT ''::bpchar, -- menu_perms.id
	uri varchar NOT NULL DEFAULT ''::character varying, -- URI
	pv int8 NOT NULL DEFAULT 0, -- 当日PV
	uv int8 NOT NULL DEFAULT 0, -- 单日UV
	"date" int8 NOT NULL, -- 日期
	unit_id bpchar(36) NOT NULL, -- 组织单位id
	modulename varchar NOT NULL DEFAULT ''::character varying, -- 模块
	CONSTRAINT mchnt_api_statistics_pk PRIMARY KEY (id)
);
COMMENT ON TABLE public.mchnt_api_statistics IS 'api统计';

-- Column comments

COMMENT ON COLUMN public.mchnt_api_statistics.id IS 'ID';
COMMENT ON COLUMN public.mchnt_api_statistics.perms_id IS 'menu_perms.id';
COMMENT ON COLUMN public.mchnt_api_statistics.uri IS 'URI';
COMMENT ON COLUMN public.mchnt_api_statistics.pv IS '当日PV';
COMMENT ON COLUMN public.mchnt_api_statistics.uv IS '单日UV';
COMMENT ON COLUMN public.mchnt_api_statistics."date" IS '日期';
COMMENT ON COLUMN public.mchnt_api_statistics.unit_id IS '组织单位id';
COMMENT ON COLUMN public.mchnt_api_statistics.modulename IS '模块';


-- public.mchnt_customer definition

-- Drop table

-- DROP TABLE public.mchnt_customer;

CREATE TABLE public.mchnt_customer (
	id bpchar(36) NOT NULL, -- ID
	mchnt_id bpchar(36) NOT NULL, -- 组织id
	customer_id bpchar(36) NOT NULL, -- 用户id
	status int4 NOT NULL DEFAULT 1, -- 当前组织内状态：0禁用,1正常
	deleted int2 NOT NULL DEFAULT 0, -- 是否已删除
	CONSTRAINT org_customer_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX org_customer_org_id_customer_id_idx ON public.mchnt_customer USING btree (mchnt_id, customer_id);
COMMENT ON TABLE public.mchnt_customer IS '组织用户关系表-Merchant';

-- Column comments

COMMENT ON COLUMN public.mchnt_customer.id IS 'ID';
COMMENT ON COLUMN public.mchnt_customer.mchnt_id IS '组织id';
COMMENT ON COLUMN public.mchnt_customer.customer_id IS '用户id';
COMMENT ON COLUMN public.mchnt_customer.status IS '当前组织内状态：0禁用,1正常';
COMMENT ON COLUMN public.mchnt_customer.deleted IS '是否已删除';


-- public.mchnt_dept definition

-- Drop table

-- DROP TABLE public.mchnt_dept;

CREATE TABLE public.mchnt_dept (
	id bpchar(36) NOT NULL, -- ID
	pid bpchar(36) NULL, -- 上级部门id
	unit_id bpchar(36) NOT NULL, -- 组织单位id
	"name" varchar(100) NOT NULL, -- 部门名称
	principal_id bpchar(36) NULL, -- 负责人id
	sort int4 NOT NULL DEFAULT 0, -- 排序
	status int4 NOT NULL DEFAULT 0, -- 状态：0禁用1启用
	deleted int4 NOT NULL DEFAULT 0, -- 是否删除：0否1是
	updated_at int8 NOT NULL, -- 更新时间
	deleted_at int8 NULL, -- 删除时间
	remark varchar(512) NULL, -- 备注
	CONSTRAINT mchnt_dept_pk PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.mchnt_dept.id IS 'ID';
COMMENT ON COLUMN public.mchnt_dept.pid IS '上级部门id';
COMMENT ON COLUMN public.mchnt_dept.unit_id IS '组织单位id';
COMMENT ON COLUMN public.mchnt_dept."name" IS '部门名称';
COMMENT ON COLUMN public.mchnt_dept.principal_id IS '负责人id';
COMMENT ON COLUMN public.mchnt_dept.sort IS '排序';
COMMENT ON COLUMN public.mchnt_dept.status IS '状态：0禁用1启用';
COMMENT ON COLUMN public.mchnt_dept.deleted IS '是否删除：0否1是';
COMMENT ON COLUMN public.mchnt_dept.updated_at IS '更新时间';
COMMENT ON COLUMN public.mchnt_dept.deleted_at IS '删除时间';
COMMENT ON COLUMN public.mchnt_dept.remark IS '备注';


-- public.mchnt_menu definition

-- Drop table

-- DROP TABLE public.mchnt_menu;

CREATE TABLE public.mchnt_menu (
	id bpchar(36) NOT NULL, -- 菜单ID
	parent_id varchar(36) NULL, -- 父级菜单ID
	menu_type int4 NOT NULL, -- 菜单类型:0代表菜单、1代表iframe、2代表外链、3代表按钮、4所需额外接口
	title varchar(255) NULL, -- 菜单标题
	"name" varchar(255) NULL, -- 菜单名称
	"path" varchar(255) NULL, -- 路由路径
	component varchar(255) NULL, -- 组件路径
	"rank" int4 NULL, -- 菜单排序
	redirect varchar(255) NULL, -- 重定向地址
	icon varchar(255) NULL, -- 图标
	extra_icon varchar(255) NULL, -- 额外图标
	enter_transition varchar(255) NULL, -- 进入动画
	leave_transition varchar(255) NULL, -- 离开动画
	active_path varchar(255) NULL, -- 激活路径
	auths text NULL, -- 权限标识
	frame_src varchar(255) NULL, -- iframe链接
	frame_loading bool NULL DEFAULT true, -- iframe加载
	keep_alive bool NULL DEFAULT false, -- 缓存页面
	hidden_tag bool NULL DEFAULT false, -- 隐藏标签页
	fixed_tag bool NULL DEFAULT false, -- 固定标签页
	show_link bool NULL DEFAULT true, -- 显示链接
	show_parent bool NULL DEFAULT false, -- 显示父级菜单
	created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	deleted int4 NULL DEFAULT 0, -- 是否删除：0否1是
	clone int4 NULL DEFAULT 0, -- 是否允许克隆：0否1是
	remark varchar NULL DEFAULT ''::character varying, -- 备注
	CONSTRAINT mchnt_menu_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_menu_parent_id ON public.mchnt_menu USING btree (parent_id);
CREATE INDEX idx_menu_path ON public.mchnt_menu USING btree (path);

-- Column comments

COMMENT ON COLUMN public.mchnt_menu.id IS '菜单ID';
COMMENT ON COLUMN public.mchnt_menu.parent_id IS '父级菜单ID';
COMMENT ON COLUMN public.mchnt_menu.menu_type IS '菜单类型:0代表菜单、1代表iframe、2代表外链、3代表按钮、4所需额外接口';
COMMENT ON COLUMN public.mchnt_menu.title IS '菜单标题';
COMMENT ON COLUMN public.mchnt_menu."name" IS '菜单名称';
COMMENT ON COLUMN public.mchnt_menu."path" IS '路由路径';
COMMENT ON COLUMN public.mchnt_menu.component IS '组件路径';
COMMENT ON COLUMN public.mchnt_menu."rank" IS '菜单排序';
COMMENT ON COLUMN public.mchnt_menu.redirect IS '重定向地址';
COMMENT ON COLUMN public.mchnt_menu.icon IS '图标';
COMMENT ON COLUMN public.mchnt_menu.extra_icon IS '额外图标';
COMMENT ON COLUMN public.mchnt_menu.enter_transition IS '进入动画';
COMMENT ON COLUMN public.mchnt_menu.leave_transition IS '离开动画';
COMMENT ON COLUMN public.mchnt_menu.active_path IS '激活路径';
COMMENT ON COLUMN public.mchnt_menu.auths IS '权限标识';
COMMENT ON COLUMN public.mchnt_menu.frame_src IS 'iframe链接';
COMMENT ON COLUMN public.mchnt_menu.frame_loading IS 'iframe加载';
COMMENT ON COLUMN public.mchnt_menu.keep_alive IS '缓存页面';
COMMENT ON COLUMN public.mchnt_menu.hidden_tag IS '隐藏标签页';
COMMENT ON COLUMN public.mchnt_menu.fixed_tag IS '固定标签页';
COMMENT ON COLUMN public.mchnt_menu.show_link IS '显示链接';
COMMENT ON COLUMN public.mchnt_menu.show_parent IS '显示父级菜单';
COMMENT ON COLUMN public.mchnt_menu.deleted IS '是否删除：0否1是';
COMMENT ON COLUMN public.mchnt_menu.clone IS '是否允许克隆：0否1是';
COMMENT ON COLUMN public.mchnt_menu.remark IS '备注';


-- public.mchnt_menu_map definition

-- Drop table

-- DROP TABLE public.mchnt_menu_map;

CREATE TABLE public.mchnt_menu_map (
	id bpchar(36) NOT NULL, -- ID
	unit_id bpchar(36) NOT NULL, -- 组织单位id
	menu_id bpchar(36) NOT NULL, -- 菜单id
	updated_at int8 NULL, -- 更新时间
	deleted int4 NOT NULL DEFAULT 0, -- 是否删除：0否1是
	CONSTRAINT mchnt_menu_map_pk PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.mchnt_menu_map.id IS 'ID';
COMMENT ON COLUMN public.mchnt_menu_map.unit_id IS '组织单位id';
COMMENT ON COLUMN public.mchnt_menu_map.menu_id IS '菜单id';
COMMENT ON COLUMN public.mchnt_menu_map.updated_at IS '更新时间';
COMMENT ON COLUMN public.mchnt_menu_map.deleted IS '是否删除：0否1是';


-- public.mchnt_role definition

-- Drop table

-- DROP TABLE public.mchnt_role;

CREATE TABLE public.mchnt_role (
	id bpchar(36) NOT NULL, -- 角色ID
	unit_id bpchar(36) NOT NULL, -- 组织单位id
	role_name varchar(50) NOT NULL, -- 角色名称
	role_sort int4 NOT NULL, -- 显示顺序
	status int4 NOT NULL DEFAULT 1, -- 角色状态：0停用，1正常
	deleted int4 NULL DEFAULT 0, -- 删除：0否,1是
	created_by varchar(64) NULL DEFAULT ''::character varying, -- 创建者
	created_at int8 NULL, -- 创建时间
	updated_by varchar(64) NULL DEFAULT ''::character varying, -- 更新者
	updated_at int8 NULL, -- 更新时间
	remark varchar(500) NULL DEFAULT NULL::character varying, -- 备注
	CONSTRAINT org_role_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX org_role_org_id_role_name_idx ON public.mchnt_role USING btree (unit_id, role_name);
COMMENT ON TABLE public.mchnt_role IS '组织角色表-Merchant';

-- Column comments

COMMENT ON COLUMN public.mchnt_role.id IS '角色ID';
COMMENT ON COLUMN public.mchnt_role.unit_id IS '组织单位id';
COMMENT ON COLUMN public.mchnt_role.role_name IS '角色名称';
COMMENT ON COLUMN public.mchnt_role.role_sort IS '显示顺序';
COMMENT ON COLUMN public.mchnt_role.status IS '角色状态：0停用，1正常';
COMMENT ON COLUMN public.mchnt_role.deleted IS '删除：0否,1是';
COMMENT ON COLUMN public.mchnt_role.created_by IS '创建者';
COMMENT ON COLUMN public.mchnt_role.created_at IS '创建时间';
COMMENT ON COLUMN public.mchnt_role.updated_by IS '更新者';
COMMENT ON COLUMN public.mchnt_role.updated_at IS '更新时间';
COMMENT ON COLUMN public.mchnt_role.remark IS '备注';


-- public.mchnt_role_classify definition

-- Drop table

-- DROP TABLE public.mchnt_role_classify;

CREATE TABLE public.mchnt_role_classify (
	id bpchar(36) NOT NULL, -- ID
	role_id bpchar(36) NOT NULL, -- 角色id
	"name" varchar(100) NOT NULL, -- 角色分类
	unit_id bpchar(36) NOT NULL DEFAULT ''::bpchar,
	deleted int4 NOT NULL DEFAULT 0,
	CONSTRAINT org_role_classify_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public.mchnt_role_classify IS '组织角色业务分类-Merchant';

-- Column comments

COMMENT ON COLUMN public.mchnt_role_classify.id IS 'ID';
COMMENT ON COLUMN public.mchnt_role_classify.role_id IS '角色id';
COMMENT ON COLUMN public.mchnt_role_classify."name" IS '角色分类';


-- public.mchnt_role_menu definition

-- Drop table

-- DROP TABLE public.mchnt_role_menu;

CREATE TABLE public.mchnt_role_menu (
	id bpchar(36) NOT NULL, -- ID
	role_id bpchar(36) NOT NULL, -- 角色ID
	menu_id bpchar(36) NOT NULL, -- 菜单权限ID
	CONSTRAINT org_role_menu_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public.mchnt_role_menu IS '系统角色菜单权限表-Merchant';

-- Column comments

COMMENT ON COLUMN public.mchnt_role_menu.id IS 'ID';
COMMENT ON COLUMN public.mchnt_role_menu.role_id IS '角色ID';
COMMENT ON COLUMN public.mchnt_role_menu.menu_id IS '菜单权限ID';


-- public.mchnt_user definition

-- Drop table

-- DROP TABLE public.mchnt_user;

CREATE TABLE public.mchnt_user (
	id bpchar(36) NOT NULL, -- ID
	unit_id bpchar(36) NOT NULL, -- 组织单位id
	user_id bpchar(36) NOT NULL, -- 员工id
	deleted int4 NULL DEFAULT 0, -- 删除：0否,1是
	is_default int4 NULL DEFAULT 0, -- 是否默认：0否1是
	"name" varchar(36) NOT NULL, -- 姓名
	phone int8 NOT NULL, -- 手机号
	CONSTRAINT org_staff_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX org_staff_org_id_customer_id_idx ON public.mchnt_user USING btree (unit_id, user_id);
COMMENT ON TABLE public.mchnt_user IS '组织员工关系表-Merchant';

-- Column comments

COMMENT ON COLUMN public.mchnt_user.id IS 'ID';
COMMENT ON COLUMN public.mchnt_user.unit_id IS '组织单位id';
COMMENT ON COLUMN public.mchnt_user.user_id IS '员工id';
COMMENT ON COLUMN public.mchnt_user.deleted IS '删除：0否,1是';
COMMENT ON COLUMN public.mchnt_user.is_default IS '是否默认：0否1是';
COMMENT ON COLUMN public.mchnt_user."name" IS '姓名';
COMMENT ON COLUMN public.mchnt_user.phone IS '手机号';


-- public.mchnt_user_dept definition

-- Drop table

-- DROP TABLE public.mchnt_user_dept;

CREATE TABLE public.mchnt_user_dept (
	id bpchar(36) NOT NULL, -- ID
	user_id bpchar(36) NOT NULL, -- 组织单位用户id
	dept_id bpchar(36) NOT NULL, -- 组织单位部门id
	deleted int4 NULL DEFAULT 0, -- 是否删除：0否1是
	CONSTRAINT mchnt_user_dept_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.mchnt_user_dept.id IS 'ID';
COMMENT ON COLUMN public.mchnt_user_dept.user_id IS '组织单位用户id';
COMMENT ON COLUMN public.mchnt_user_dept.dept_id IS '组织单位部门id';
COMMENT ON COLUMN public.mchnt_user_dept.deleted IS '是否删除：0否1是';


-- public.mchnt_user_profile definition

-- Drop table

-- DROP TABLE public.mchnt_user_profile;

CREATE TABLE public.mchnt_user_profile (
	id bpchar(36) NOT NULL, -- ID
	avatar varchar(255) NULL, -- 头像
	card_type int2 NULL, -- 1大陆身份证2港澳台身份证3护照4军官证5其它
	card_num varchar(100) NULL, -- 证件号码
	card_images varchar(1000) NULL, -- 证件照片
	gender int2 NULL, -- 性别:1男，2女
	birth_date date NULL, -- 出生日期
	constellation varchar(50) NULL, -- 星座
	occupation varchar(50) NULL, -- 职业
	company varchar(500) NULL, -- 所属公司名称
	emergency_name varchar(50) NULL, -- 紧急联系人姓名
	emergency_tel varchar(100) NULL, -- 紧急联系人电话
	address varchar(200) NULL, -- 通讯地址
	email varchar(50) NULL, -- 邮箱
	valid_date_begin timestamp NULL, -- 身份证有效期开始时间
	valid_date_end timestamp NULL, -- 身份证有效期截止时间
	schooling varchar(100) NULL, -- 学历
	degree_number varchar(100) NULL, -- 学位编号
	remark varchar(255) NULL, -- 备注
	professional varchar(100) NULL, -- 专业
	status int4 NOT NULL DEFAULT 1, -- 用户行为状态：1正常，2已注销，组织单位状态：1正常，3禁用，4离职
	created_at int8 NULL, -- 记录创建时间
	updated_at int8 NULL, -- 记录修改时间
	deleted_at int8 NULL, -- 删除时间
	deleted int4 NOT NULL DEFAULT 0, -- 是否删除：0否1是
	graduated_from varchar(100) NULL DEFAULT ''::character varying, -- 毕业院校
	"source" int4 NULL DEFAULT 1, -- 注册来源：1系统录入2微信3web端4app5其它
	CONSTRAINT mchnt_user_profile_pkey PRIMARY KEY (id)
);
CREATE INDEX mchnt_user_profile_card_num_idx ON public.mchnt_user_profile USING btree (card_num);
COMMENT ON TABLE public.mchnt_user_profile IS '商户平台用户信息表';

-- Column comments

COMMENT ON COLUMN public.mchnt_user_profile.id IS 'ID';
COMMENT ON COLUMN public.mchnt_user_profile.avatar IS '头像';
COMMENT ON COLUMN public.mchnt_user_profile.card_type IS '1大陆身份证2港澳台身份证3护照4军官证5其它';
COMMENT ON COLUMN public.mchnt_user_profile.card_num IS '证件号码';
COMMENT ON COLUMN public.mchnt_user_profile.card_images IS '证件照片';
COMMENT ON COLUMN public.mchnt_user_profile.gender IS '性别:1男，2女';
COMMENT ON COLUMN public.mchnt_user_profile.birth_date IS '出生日期';
COMMENT ON COLUMN public.mchnt_user_profile.constellation IS '星座';
COMMENT ON COLUMN public.mchnt_user_profile.occupation IS '职业';
COMMENT ON COLUMN public.mchnt_user_profile.company IS '所属公司名称';
COMMENT ON COLUMN public.mchnt_user_profile.emergency_name IS '紧急联系人姓名';
COMMENT ON COLUMN public.mchnt_user_profile.emergency_tel IS '紧急联系人电话';
COMMENT ON COLUMN public.mchnt_user_profile.address IS '通讯地址';
COMMENT ON COLUMN public.mchnt_user_profile.email IS '邮箱';
COMMENT ON COLUMN public.mchnt_user_profile.valid_date_begin IS '身份证有效期开始时间';
COMMENT ON COLUMN public.mchnt_user_profile.valid_date_end IS '身份证有效期截止时间';
COMMENT ON COLUMN public.mchnt_user_profile.schooling IS '学历';
COMMENT ON COLUMN public.mchnt_user_profile.degree_number IS '学位编号';
COMMENT ON COLUMN public.mchnt_user_profile.remark IS '备注';
COMMENT ON COLUMN public.mchnt_user_profile.professional IS '专业';
COMMENT ON COLUMN public.mchnt_user_profile.status IS '用户行为状态：1正常，2已注销，组织单位状态：1正常，3禁用，4离职';
COMMENT ON COLUMN public.mchnt_user_profile.created_at IS '记录创建时间';
COMMENT ON COLUMN public.mchnt_user_profile.updated_at IS '记录修改时间';
COMMENT ON COLUMN public.mchnt_user_profile.deleted_at IS '删除时间';
COMMENT ON COLUMN public.mchnt_user_profile.deleted IS '是否删除：0否1是';
COMMENT ON COLUMN public.mchnt_user_profile.graduated_from IS '毕业院校';
COMMENT ON COLUMN public.mchnt_user_profile."source" IS '注册来源：1系统录入2微信3web端4app5其它';


-- public.mchnt_user_role definition

-- Drop table

-- DROP TABLE public.mchnt_user_role;

CREATE TABLE public.mchnt_user_role (
	id bpchar(36) NOT NULL, -- ID
	user_id bpchar(36) NOT NULL, -- 员工id
	role_id bpchar(36) NOT NULL, -- 员工角色id
	deleted int4 NULL DEFAULT 0, -- 删除：0否,1是
	CONSTRAINT org_staff_role_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public.mchnt_user_role IS '组织员工角色关系表-Merchant';

-- Column comments

COMMENT ON COLUMN public.mchnt_user_role.id IS 'ID';
COMMENT ON COLUMN public.mchnt_user_role.user_id IS '员工id';
COMMENT ON COLUMN public.mchnt_user_role.role_id IS '员工角色id';
COMMENT ON COLUMN public.mchnt_user_role.deleted IS '删除：0否,1是';


-- public.plat definition

-- Drop table

-- DROP TABLE public.plat;

CREATE TABLE public.plat (
	id bpchar(36) NOT NULL, -- ID
	pid bpchar(36) NOT NULL DEFAULT ''::bpchar, -- 上级组织id
	logo varchar(512) NULL DEFAULT ''::character varying, -- logo
	"name" varchar(100) NOT NULL, -- 单位名称
	code varchar(100) NOT NULL, -- 组织机构代码
	corporation varchar(100) NOT NULL, -- 法人
	license varchar(512) NOT NULL DEFAULT ''::character varying, -- 营业执照
	address varchar(255) NULL DEFAULT ''::character varying, -- 地址
	status int4 NOT NULL DEFAULT 0, -- 0未审核，1审核通过，2审核不通过，3禁用
	plat_status int4 NOT NULL DEFAULT 0, -- 平台审核状态(同上)
	deleted int4 NULL DEFAULT 0, -- 是否删除：0否1是
	created_at int8 NULL, -- 创建时间
	updated_at int8 NULL, -- 更新时间
	deleted_at int8 NULL, -- 删除时间
	created_by bpchar(36) NOT NULL DEFAULT ''::bpchar, -- 创建人
	updated_by bpchar(36) NULL, -- 更新人
	deleted_by bpchar(36) NULL, -- 删除人
	sort int4 NOT NULL DEFAULT 0, -- 排序
	is_official bool NOT NULL DEFAULT false, -- 是否官方平台
	CONSTRAINT plat_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public.plat IS '组织机构表-Platform';

-- Column comments

COMMENT ON COLUMN public.plat.id IS 'ID';
COMMENT ON COLUMN public.plat.pid IS '上级组织id';
COMMENT ON COLUMN public.plat.logo IS 'logo';
COMMENT ON COLUMN public.plat."name" IS '单位名称';
COMMENT ON COLUMN public.plat.code IS '组织机构代码';
COMMENT ON COLUMN public.plat.corporation IS '法人';
COMMENT ON COLUMN public.plat.license IS '营业执照';
COMMENT ON COLUMN public.plat.address IS '地址';
COMMENT ON COLUMN public.plat.status IS '0未审核，1审核通过，2审核不通过，3禁用';
COMMENT ON COLUMN public.plat.plat_status IS '平台审核状态(同上)';
COMMENT ON COLUMN public.plat.deleted IS '是否删除：0否1是';
COMMENT ON COLUMN public.plat.created_at IS '创建时间';
COMMENT ON COLUMN public.plat.updated_at IS '更新时间';
COMMENT ON COLUMN public.plat.deleted_at IS '删除时间';
COMMENT ON COLUMN public.plat.created_by IS '创建人';
COMMENT ON COLUMN public.plat.updated_by IS '更新人';
COMMENT ON COLUMN public.plat.deleted_by IS '删除人';
COMMENT ON COLUMN public.plat.sort IS '排序';
COMMENT ON COLUMN public.plat.is_official IS '是否官方平台';


-- public.plat_api_statistics definition

-- Drop table

-- DROP TABLE public.plat_api_statistics;

CREATE TABLE public.plat_api_statistics (
	id bpchar(36) NOT NULL, -- ID
	perms_id bpchar(36) NULL DEFAULT ''::bpchar, -- menu_perms.id
	uri varchar NOT NULL DEFAULT ''::character varying, -- URI
	pv int8 NOT NULL DEFAULT 0, -- 当日PV
	uv int8 NOT NULL DEFAULT 0, -- 单日UV
	"date" int8 NOT NULL, -- 日期
	unit_id bpchar(36) NOT NULL, -- 组织单位id
	modulename varchar NOT NULL DEFAULT ''::character varying, -- 模块
	CONSTRAINT plat_api_statistics_pk PRIMARY KEY (id)
);
COMMENT ON TABLE public.plat_api_statistics IS 'api统计';

-- Column comments

COMMENT ON COLUMN public.plat_api_statistics.id IS 'ID';
COMMENT ON COLUMN public.plat_api_statistics.perms_id IS 'menu_perms.id';
COMMENT ON COLUMN public.plat_api_statistics.uri IS 'URI';
COMMENT ON COLUMN public.plat_api_statistics.pv IS '当日PV';
COMMENT ON COLUMN public.plat_api_statistics.uv IS '单日UV';
COMMENT ON COLUMN public.plat_api_statistics."date" IS '日期';
COMMENT ON COLUMN public.plat_api_statistics.unit_id IS '组织单位id';
COMMENT ON COLUMN public.plat_api_statistics.modulename IS '模块';


-- public.plat_crontab definition

-- Drop table

-- DROP TABLE public.plat_crontab;

CREATE TABLE public.plat_crontab (
	id bpchar(36) NOT NULL, -- ID
	unit_id bpchar(36) NOT NULL, -- 组织单位id
	"name" varchar(100) NOT NULL, -- 任务名称
	name_en varchar(100) NOT NULL, -- 任务名称-英文
	"group" varchar(100) NOT NULL DEFAULT 'default'::character varying, -- 分组名称
	cron_expr varchar(2048) NOT NULL, -- cron表达式
	status int4 NOT NULL DEFAULT 0, -- 状态：0禁用1启用
	created_by bpchar(36) NOT NULL, -- 创建人
	created_at date NULL, -- 创建时间
	updated_by varchar(36) NULL, -- 更新人
	updated_at date NULL, -- 更新时间
	deleted int4 NOT NULL, -- 是否删除：0否1是
	remark varchar(512) NULL, -- 备注
	CONSTRAINT plat_crontab_pk PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.plat_crontab.id IS 'ID';
COMMENT ON COLUMN public.plat_crontab.unit_id IS '组织单位id';
COMMENT ON COLUMN public.plat_crontab."name" IS '任务名称';
COMMENT ON COLUMN public.plat_crontab.name_en IS '任务名称-英文';
COMMENT ON COLUMN public.plat_crontab."group" IS '分组名称';
COMMENT ON COLUMN public.plat_crontab.cron_expr IS 'cron表达式';
COMMENT ON COLUMN public.plat_crontab.status IS '状态：0禁用1启用';
COMMENT ON COLUMN public.plat_crontab.created_by IS '创建人';
COMMENT ON COLUMN public.plat_crontab.created_at IS '创建时间';
COMMENT ON COLUMN public.plat_crontab.updated_by IS '更新人';
COMMENT ON COLUMN public.plat_crontab.updated_at IS '更新时间';
COMMENT ON COLUMN public.plat_crontab.deleted IS '是否删除：0否1是';
COMMENT ON COLUMN public.plat_crontab.remark IS '备注';


-- public.plat_crontab_log definition

-- Drop table

-- DROP TABLE public.plat_crontab_log;

CREATE TABLE public.plat_crontab_log (
	id bpchar(36) NOT NULL, -- ID
	name_en varchar(36) NOT NULL, -- crontab表id
	created_at int8 NOT NULL, -- 创建时间
	"result" bool NULL, -- 执行结果
	remark varchar NULL -- 备注
);
COMMENT ON TABLE public.plat_crontab_log IS 'crontab日志表';

-- Column comments

COMMENT ON COLUMN public.plat_crontab_log.id IS 'ID';
COMMENT ON COLUMN public.plat_crontab_log.name_en IS 'crontab表id';
COMMENT ON COLUMN public.plat_crontab_log.created_at IS '创建时间';
COMMENT ON COLUMN public.plat_crontab_log."result" IS '执行结果';
COMMENT ON COLUMN public.plat_crontab_log.remark IS '备注';


-- public.plat_dept definition

-- Drop table

-- DROP TABLE public.plat_dept;

CREATE TABLE public.plat_dept (
	id bpchar(36) NOT NULL, -- ID
	pid bpchar(36) NULL, -- 上级部门id
	unit_id bpchar(36) NOT NULL, -- 组织单位id
	"name" varchar(100) NOT NULL, -- 部门名称
	principal_id bpchar(36) NULL, -- 负责人id
	sort int4 NOT NULL DEFAULT 0, -- 排序
	status int4 NOT NULL DEFAULT 0, -- 状态：0禁用1启用
	deleted int4 NOT NULL DEFAULT 0, -- 是否删除：0否1是
	updated_at int8 NOT NULL, -- 更新时间
	deleted_at int8 NULL, -- 删除时间
	remark varchar(512) NULL, -- 备注
	CONSTRAINT plat_dept_pk PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.plat_dept.id IS 'ID';
COMMENT ON COLUMN public.plat_dept.pid IS '上级部门id';
COMMENT ON COLUMN public.plat_dept.unit_id IS '组织单位id';
COMMENT ON COLUMN public.plat_dept."name" IS '部门名称';
COMMENT ON COLUMN public.plat_dept.principal_id IS '负责人id';
COMMENT ON COLUMN public.plat_dept.sort IS '排序';
COMMENT ON COLUMN public.plat_dept.status IS '状态：0禁用1启用';
COMMENT ON COLUMN public.plat_dept.deleted IS '是否删除：0否1是';
COMMENT ON COLUMN public.plat_dept.updated_at IS '更新时间';
COMMENT ON COLUMN public.plat_dept.deleted_at IS '删除时间';
COMMENT ON COLUMN public.plat_dept.remark IS '备注';


-- public.plat_menu definition

-- Drop table

-- DROP TABLE public.plat_menu;

CREATE TABLE public.plat_menu (
	id bpchar(36) NOT NULL, -- 菜单ID
	parent_id varchar(36) NULL, -- 父级菜单ID
	menu_type int4 NOT NULL, -- 菜单类型:0代表菜单、1代表iframe、2代表外链、3代表按钮、4所需额外接口
	title varchar(255) NULL, -- 菜单标题
	"name" varchar(255) NULL, -- 菜单名称
	"path" varchar(255) NULL, -- 路由路径
	component varchar(255) NULL, -- 组件路径
	"rank" int4 NULL, -- 菜单排序
	redirect varchar(255) NULL, -- 重定向地址
	icon varchar(255) NULL, -- 图标
	extra_icon varchar(255) NULL, -- 额外图标
	enter_transition varchar(255) NULL, -- 进入动画
	leave_transition varchar(255) NULL, -- 离开动画
	active_path varchar(255) NULL, -- 激活路径
	auths text NULL, -- 权限标识
	frame_src varchar(255) NULL, -- iframe链接
	frame_loading bool NULL DEFAULT true, -- iframe加载
	keep_alive bool NULL DEFAULT false, -- 缓存页面
	hidden_tag bool NULL DEFAULT false, -- 隐藏标签页
	fixed_tag bool NULL DEFAULT false, -- 固定标签页
	show_link bool NULL DEFAULT true, -- 显示链接
	show_parent bool NULL DEFAULT false, -- 显示父级菜单
	created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	deleted int4 NULL DEFAULT 0, -- 是否删除：0否1是
	clone int4 NULL DEFAULT 0, -- 是否允许克隆：0否1是
	remark varchar NULL DEFAULT ''::character varying, -- 备注
	menu_from varchar(20) DEFAULT 'admin_plat' NOT NULL,
	CONSTRAINT plat_menu_pkey PRIMARY KEY (id)
);
CREATE INDEX plat_idx_menu_parent_id ON public.plat_menu USING btree (parent_id);
CREATE INDEX plat_idx_menu_path ON public.plat_menu USING btree (path);

-- Column comments

COMMENT ON COLUMN public.plat_menu.id IS '菜单ID';
COMMENT ON COLUMN public.plat_menu.parent_id IS '父级菜单ID';
COMMENT ON COLUMN public.plat_menu.menu_type IS '菜单类型:0代表菜单、1代表iframe、2代表外链、3代表按钮、4所需额外接口';
COMMENT ON COLUMN public.plat_menu.title IS '菜单标题';
COMMENT ON COLUMN public.plat_menu."name" IS '菜单名称';
COMMENT ON COLUMN public.plat_menu."path" IS '路由路径';
COMMENT ON COLUMN public.plat_menu.component IS '组件路径';
COMMENT ON COLUMN public.plat_menu."rank" IS '菜单排序';
COMMENT ON COLUMN public.plat_menu.redirect IS '重定向地址';
COMMENT ON COLUMN public.plat_menu.icon IS '图标';
COMMENT ON COLUMN public.plat_menu.extra_icon IS '额外图标';
COMMENT ON COLUMN public.plat_menu.enter_transition IS '进入动画';
COMMENT ON COLUMN public.plat_menu.leave_transition IS '离开动画';
COMMENT ON COLUMN public.plat_menu.active_path IS '激活路径';
COMMENT ON COLUMN public.plat_menu.auths IS '权限标识';
COMMENT ON COLUMN public.plat_menu.frame_src IS 'iframe链接';
COMMENT ON COLUMN public.plat_menu.frame_loading IS 'iframe加载';
COMMENT ON COLUMN public.plat_menu.keep_alive IS '缓存页面';
COMMENT ON COLUMN public.plat_menu.hidden_tag IS '隐藏标签页';
COMMENT ON COLUMN public.plat_menu.fixed_tag IS '固定标签页';
COMMENT ON COLUMN public.plat_menu.show_link IS '显示链接';
COMMENT ON COLUMN public.plat_menu.show_parent IS '显示父级菜单';
COMMENT ON COLUMN public.plat_menu.deleted IS '是否删除：0否1是';
COMMENT ON COLUMN public.plat_menu.clone IS '是否允许克隆：0否1是';
COMMENT ON COLUMN public.plat_menu.remark IS '备注';
COMMENT ON COLUMN public.plat_menu.menu_from IS '系统菜单分类：admin_plat、admin_mchnt；admin_mchnt表示作为商户系统数据管理';


-- public.plat_menu_map definition

-- Drop table

-- DROP TABLE public.plat_menu_map;

CREATE TABLE public.plat_menu_map (
	id bpchar(36) NOT NULL, -- ID
	unit_id bpchar(36) NOT NULL, -- 组织单位id
	menu_id bpchar(36) NOT NULL, -- 菜单id
	updated_at int8 NULL, -- 更新时间
	deleted int4 NOT NULL DEFAULT 0, -- 是否删除：0否1是
	CONSTRAINT plat_menu_map_pk PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.plat_menu_map.id IS 'ID';
COMMENT ON COLUMN public.plat_menu_map.unit_id IS '组织单位id';
COMMENT ON COLUMN public.plat_menu_map.menu_id IS '菜单id';
COMMENT ON COLUMN public.plat_menu_map.updated_at IS '更新时间';
COMMENT ON COLUMN public.plat_menu_map.deleted IS '是否删除：0否1是';


-- public.plat_role definition

-- Drop table

-- DROP TABLE public.plat_role;

CREATE TABLE public.plat_role (
	id bpchar(36) NOT NULL, -- 角色ID
	unit_id bpchar(36) NOT NULL, -- 组织单位id
	role_name varchar(50) NOT NULL, -- 角色名称
	role_sort int4 NOT NULL, -- 显示顺序
	status int4 NOT NULL DEFAULT 1, -- 角色状态：0停用，1正常
	deleted int4 NULL DEFAULT 0, -- 删除：0否,1是
	created_by varchar(64) NULL DEFAULT ''::character varying, -- 创建者
	created_at int8 NULL, -- 创建时间
	updated_by varchar(64) NULL DEFAULT ''::character varying, -- 更新者
	updated_at int8 NULL, -- 更新时间
	remark varchar(500) NULL DEFAULT NULL::character varying, -- 备注
	CONSTRAINT plat_role_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX plat_role_plat_id_role_name_idx ON public.plat_role USING btree (unit_id, role_name);
COMMENT ON TABLE public.plat_role IS '组织角色表-Platform';

-- Column comments

COMMENT ON COLUMN public.plat_role.id IS '角色ID';
COMMENT ON COLUMN public.plat_role.unit_id IS '组织单位id';
COMMENT ON COLUMN public.plat_role.role_name IS '角色名称';
COMMENT ON COLUMN public.plat_role.role_sort IS '显示顺序';
COMMENT ON COLUMN public.plat_role.status IS '角色状态：0停用，1正常';
COMMENT ON COLUMN public.plat_role.deleted IS '删除：0否,1是';
COMMENT ON COLUMN public.plat_role.created_by IS '创建者';
COMMENT ON COLUMN public.plat_role.created_at IS '创建时间';
COMMENT ON COLUMN public.plat_role.updated_by IS '更新者';
COMMENT ON COLUMN public.plat_role.updated_at IS '更新时间';
COMMENT ON COLUMN public.plat_role.remark IS '备注';


-- public.plat_role_classify definition

-- Drop table

-- DROP TABLE public.plat_role_classify;

CREATE TABLE public.plat_role_classify (
	id bpchar(36) NOT NULL, -- ID
	role_id bpchar(36) NOT NULL, -- 角色id
	"name" varchar(100) NOT NULL, -- 角色分类
	unit_id bpchar(36) NOT NULL DEFAULT ''::bpchar, -- 组织单位id
	deleted int4 NOT NULL DEFAULT 0, -- 是否删除：0否1是
	CONSTRAINT plat_role_classify_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public.plat_role_classify IS '组织角色业务分类-Platform';

-- Column comments

COMMENT ON COLUMN public.plat_role_classify.id IS 'ID';
COMMENT ON COLUMN public.plat_role_classify.role_id IS '角色id';
COMMENT ON COLUMN public.plat_role_classify."name" IS '角色分类';
COMMENT ON COLUMN public.plat_role_classify.unit_id IS '组织单位id';
COMMENT ON COLUMN public.plat_role_classify.deleted IS '是否删除：0否1是';


-- public.plat_role_menu definition

-- Drop table

-- DROP TABLE public.plat_role_menu;

CREATE TABLE public.plat_role_menu (
	id bpchar(36) NOT NULL, -- ID
	role_id bpchar(36) NOT NULL, -- 角色ID
	menu_id bpchar(36) NOT NULL, -- 菜单权限ID
	CONSTRAINT plat_role_menu_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public.plat_role_menu IS '系统角色菜单权限表-Platform';

-- Column comments

COMMENT ON COLUMN public.plat_role_menu.id IS 'ID';
COMMENT ON COLUMN public.plat_role_menu.role_id IS '角色ID';
COMMENT ON COLUMN public.plat_role_menu.menu_id IS '菜单权限ID';


-- public.plat_user definition

-- Drop table

-- DROP TABLE public.plat_user;

CREATE TABLE public.plat_user (
	id bpchar(36) NOT NULL, -- ID
	unit_id bpchar(36) NOT NULL, -- 组织单位id
	user_id bpchar(36) NOT NULL, -- 员工id
	deleted int4 NULL DEFAULT 0, -- 删除：0否,1是
	is_default int4 NULL DEFAULT 0, -- 是否默认：0否1是
	"name" varchar(36) NOT NULL DEFAULT ''::character varying, -- 姓名
	phone int8 NOT NULL DEFAULT 0, -- 手机号
	CONSTRAINT plat_staff_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX plat_staff_plat_id_customer_id_idx ON public.plat_user USING btree (unit_id, user_id);
COMMENT ON TABLE public.plat_user IS '组织员工关系表-Platform';

-- Column comments

COMMENT ON COLUMN public.plat_user.id IS 'ID';
COMMENT ON COLUMN public.plat_user.unit_id IS '组织单位id';
COMMENT ON COLUMN public.plat_user.user_id IS '员工id';
COMMENT ON COLUMN public.plat_user.deleted IS '删除：0否,1是';
COMMENT ON COLUMN public.plat_user.is_default IS '是否默认：0否1是';
COMMENT ON COLUMN public.plat_user."name" IS '姓名';
COMMENT ON COLUMN public.plat_user.phone IS '手机号';


-- public.plat_user_dept definition

-- Drop table

-- DROP TABLE public.plat_user_dept;

CREATE TABLE public.plat_user_dept (
	id bpchar(36) NOT NULL, -- ID
	user_id bpchar(36) NOT NULL, -- 组织单位用户id
	dept_id bpchar(36) NOT NULL, -- 组织单位部门id
	deleted int4 NULL DEFAULT 0, -- 是否删除：0否1是
	CONSTRAINT plat_user_dept_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.plat_user_dept.id IS 'ID';
COMMENT ON COLUMN public.plat_user_dept.user_id IS '组织单位用户id';
COMMENT ON COLUMN public.plat_user_dept.dept_id IS '组织单位部门id';
COMMENT ON COLUMN public.plat_user_dept.deleted IS '是否删除：0否1是';


-- public.plat_user_profile definition

-- Drop table

-- DROP TABLE public.plat_user_profile;

CREATE TABLE public.plat_user_profile (
	id bpchar(36) NOT NULL, -- ID
	avatar varchar(255) NULL, -- 头像
	card_type int2 NULL, -- 1大陆身份证2港澳台身份证3护照4军官证5其它
	card_num varchar(100) NULL, -- 证件号码
	card_images varchar(1000) NULL, -- 证件照片
	gender int2 NULL, -- 性别:1男，2女
	birth_date date NULL, -- 出生日期
	constellation varchar(50) NULL, -- 星座
	occupation varchar(50) NULL, -- 职业
	company varchar(500) NULL, -- 所属公司名称
	emergency_name varchar(50) NULL, -- 紧急联系人姓名
	emergency_tel varchar(100) NULL, -- 紧急联系人电话
	address varchar(200) NULL, -- 通讯地址
	email varchar(50) NULL, -- 邮箱
	valid_date_begin timestamp NULL, -- 身份证有效期开始时间
	valid_date_end timestamp NULL, -- 身份证有效期截止时间
	schooling varchar(100) NULL, -- 学历
	degree_number varchar(100) NULL, -- 学位编号
	remark varchar(255) NULL, -- 备注
	professional varchar(100) NULL, -- 专业
	status int4 NOT NULL DEFAULT 1, -- 用户行为状态：1正常，2已注销，组织单位状态：1正常，3禁用，4离职
	created_at int8 NULL, -- 记录创建时间
	updated_at int8 NULL, -- 记录修改时间
	deleted_at int8 NULL, -- 删除时间
	deleted int4 NOT NULL DEFAULT 0, -- 是否删除：0否1是
	graduated_from varchar(100) NULL DEFAULT ''::character varying, -- 毕业院校
	"source" int4 NULL DEFAULT 1, -- 注册来源：1系统录入2微信3web端4app5其它
	CONSTRAINT plat_user_profile_pkey PRIMARY KEY (id)
);
CREATE INDEX plat_user_profile_card_num_idx ON public.plat_user_profile USING btree (card_num);
COMMENT ON TABLE public.plat_user_profile IS '平台用户信息表';

-- Column comments

COMMENT ON COLUMN public.plat_user_profile.id IS 'ID';
COMMENT ON COLUMN public.plat_user_profile.avatar IS '头像';
COMMENT ON COLUMN public.plat_user_profile.card_type IS '1大陆身份证2港澳台身份证3护照4军官证5其它';
COMMENT ON COLUMN public.plat_user_profile.card_num IS '证件号码';
COMMENT ON COLUMN public.plat_user_profile.card_images IS '证件照片';
COMMENT ON COLUMN public.plat_user_profile.gender IS '性别:1男，2女';
COMMENT ON COLUMN public.plat_user_profile.birth_date IS '出生日期';
COMMENT ON COLUMN public.plat_user_profile.constellation IS '星座';
COMMENT ON COLUMN public.plat_user_profile.occupation IS '职业';
COMMENT ON COLUMN public.plat_user_profile.company IS '所属公司名称';
COMMENT ON COLUMN public.plat_user_profile.emergency_name IS '紧急联系人姓名';
COMMENT ON COLUMN public.plat_user_profile.emergency_tel IS '紧急联系人电话';
COMMENT ON COLUMN public.plat_user_profile.address IS '通讯地址';
COMMENT ON COLUMN public.plat_user_profile.email IS '邮箱';
COMMENT ON COLUMN public.plat_user_profile.valid_date_begin IS '身份证有效期开始时间';
COMMENT ON COLUMN public.plat_user_profile.valid_date_end IS '身份证有效期截止时间';
COMMENT ON COLUMN public.plat_user_profile.schooling IS '学历';
COMMENT ON COLUMN public.plat_user_profile.degree_number IS '学位编号';
COMMENT ON COLUMN public.plat_user_profile.remark IS '备注';
COMMENT ON COLUMN public.plat_user_profile.professional IS '专业';
COMMENT ON COLUMN public.plat_user_profile.status IS '用户行为状态：1正常，2已注销，组织单位状态：1正常，3禁用，4离职';
COMMENT ON COLUMN public.plat_user_profile.created_at IS '记录创建时间';
COMMENT ON COLUMN public.plat_user_profile.updated_at IS '记录修改时间';
COMMENT ON COLUMN public.plat_user_profile.deleted_at IS '删除时间';
COMMENT ON COLUMN public.plat_user_profile.deleted IS '是否删除：0否1是';
COMMENT ON COLUMN public.plat_user_profile.graduated_from IS '毕业院校';
COMMENT ON COLUMN public.plat_user_profile."source" IS '注册来源：1系统录入2微信3web端4app5其它';


-- public.plat_user_role definition

-- Drop table

-- DROP TABLE public.plat_user_role;

CREATE TABLE public.plat_user_role (
	id bpchar(36) NOT NULL, -- ID
	user_id bpchar(36) NOT NULL, -- 员工id
	role_id bpchar(36) NOT NULL, -- 员工角色id
	deleted int4 NULL DEFAULT 0, -- 删除：0否,1是
	CONSTRAINT plat_staff_role_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public.plat_user_role IS '组织员工角色关系表-Platform';

-- Column comments

COMMENT ON COLUMN public.plat_user_role.id IS 'ID';
COMMENT ON COLUMN public.plat_user_role.user_id IS '员工id';
COMMENT ON COLUMN public.plat_user_role.role_id IS '员工角色id';
COMMENT ON COLUMN public.plat_user_role.deleted IS '删除：0否,1是';


-- public.queue_dlq_failed_log definition

-- Drop table

-- DROP TABLE public.queue_dlq_failed_log;

CREATE TABLE public.queue_dlq_failed_log (
	task_uuid varchar(36) NOT NULL, -- UUID
	task_name varchar(128) NULL, -- 任务名称
	task_args text NULL, -- 任务参数
	error_msg text NULL, -- 错误信息
	create_time timestamp NULL DEFAULT now(), -- 创建时间
	status int4 NULL DEFAULT 0, -- 状态：0未处理，1已经处理
	deleted varchar NULL DEFAULT 0, -- 是否删除：0否1是
	CONSTRAINT queue_dlq_failed_log_pkey PRIMARY KEY (task_uuid)
);
COMMENT ON TABLE public.queue_dlq_failed_log IS '队列：死信队列消费失败记录';

-- Column comments

COMMENT ON COLUMN public.queue_dlq_failed_log.task_uuid IS 'UUID';
COMMENT ON COLUMN public.queue_dlq_failed_log.task_name IS '任务名称';
COMMENT ON COLUMN public.queue_dlq_failed_log.task_args IS '任务参数';
COMMENT ON COLUMN public.queue_dlq_failed_log.error_msg IS '错误信息';
COMMENT ON COLUMN public.queue_dlq_failed_log.create_time IS '创建时间';
COMMENT ON COLUMN public.queue_dlq_failed_log.status IS '状态：0未处理，1已经处理';
COMMENT ON COLUMN public.queue_dlq_failed_log.deleted IS '是否删除：0否1是';


-- public.queue_dlq_failed_retry definition

-- Drop table

-- DROP TABLE public.queue_dlq_failed_retry;

-- public.queue_dlq_failed_retry definition

-- Drop table

-- DROP TABLE public.queue_dlq_failed_retry;

CREATE TABLE public.queue_dlq_failed_retry (
	id varchar(64) NOT NULL, -- 原任务UUID
	new_id varchar(64) NOT NULL, -- 新任务UUID
	CONSTRAINT queue_dlq_failed_retry_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public.queue_dlq_failed_retry IS '队列：死信队列重新入队记录';

-- Column comments

COMMENT ON COLUMN public.queue_dlq_failed_retry.id IS '原任务UUID';
COMMENT ON COLUMN public.queue_dlq_failed_retry.new_id IS '新任务UUID';


-- public."user" definition

-- Drop table

-- DROP TABLE public."user";

CREATE TABLE public."user" (
	id bpchar(36) NOT NULL, -- ID
	default_unit_id bpchar(36) NOT NULL DEFAULT ''::bpchar, -- 默认组织id
	"name" varchar(32) NOT NULL, -- 员工姓名
	wx_openid varchar(32) NULL, -- 微信openid
	phone int8 NOT NULL, -- 手机号
	"password" varchar(255) NULL DEFAULT ''::character varying, -- 密码
	email varchar(64) NULL, -- 邮箱
	CONSTRAINT sys_staff_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE public."user" IS '系统员工表';

-- Column comments

COMMENT ON COLUMN public."user".id IS 'ID';
COMMENT ON COLUMN public."user".default_unit_id IS '默认组织id';
COMMENT ON COLUMN public."user"."name" IS '员工姓名';
COMMENT ON COLUMN public."user".wx_openid IS '微信openid';
COMMENT ON COLUMN public."user".phone IS '手机号';
COMMENT ON COLUMN public."user"."password" IS '密码';
COMMENT ON COLUMN public."user".email IS '邮箱';


-- public.user_profile definition

-- Drop table

-- DROP TABLE public.user_profile;

CREATE TABLE public.user_profile (
	id bpchar(36) NOT NULL, -- ID
	avatar varchar(255) NULL, -- 头像
	card_type int2 NULL, -- 1大陆身份证2港澳台身份证3护照4军官证5其它
	card_num varchar(100) NULL, -- 证件号码
	card_images varchar(1000) NULL, -- 证件照片
	gender int2 NULL, -- 性别:1男，2女
	birth_date date NULL, -- 出生日期
	constellation varchar(50) NULL, -- 星座
	occupation varchar(50) NULL, -- 职业
	company varchar(500) NULL, -- 所属公司名称
	emergency_name varchar(50) NULL, -- 紧急联系人姓名
	emergency_tel varchar(100) NULL, -- 紧急联系人电话
	address varchar(200) NULL, -- 通讯地址
	email varchar(50) NULL, -- 邮箱
	valid_date_begin timestamp NULL, -- 身份证有效期开始时间
	valid_date_end timestamp NULL, -- 身份证有效期截止时间
	schooling varchar(100) NULL, -- 学历
	degree_number varchar(100) NULL, -- 学位编号
	remark varchar(255) NULL, -- 备注
	professional varchar(100) NULL, -- 专业
	status int4 NOT NULL DEFAULT 1, -- 用户行为状态：1正常，2已注销，平台状态：3禁用
	created_at int8 NULL, -- 记录创建时间
	updated_at int8 NULL, -- 记录修改时间
	deleted_at int8 NULL, -- 删除时间
	deleted int4 NOT NULL DEFAULT 0, -- 是否删除：0否1是
	graduated_from varchar(100) NULL DEFAULT ''::character varying, -- 毕业院校
	"source" int4 NULL DEFAULT 1, -- 注册来源：1系统录入2微信3web端4app5其它
	CONSTRAINT sys_customer_pkey PRIMARY KEY (id)
);
CREATE INDEX sys_customer_card_num_idx ON public.user_profile USING btree (card_num);
COMMENT ON TABLE public.user_profile IS '系统用户信息表';

-- Column comments

COMMENT ON COLUMN public.user_profile.id IS 'ID';
COMMENT ON COLUMN public.user_profile.avatar IS '头像';
COMMENT ON COLUMN public.user_profile.card_type IS '1大陆身份证2港澳台身份证3护照4军官证5其它';
COMMENT ON COLUMN public.user_profile.card_num IS '证件号码';
COMMENT ON COLUMN public.user_profile.card_images IS '证件照片';
COMMENT ON COLUMN public.user_profile.gender IS '性别:1男，2女';
COMMENT ON COLUMN public.user_profile.birth_date IS '出生日期';
COMMENT ON COLUMN public.user_profile.constellation IS '星座';
COMMENT ON COLUMN public.user_profile.occupation IS '职业';
COMMENT ON COLUMN public.user_profile.company IS '所属公司名称';
COMMENT ON COLUMN public.user_profile.emergency_name IS '紧急联系人姓名';
COMMENT ON COLUMN public.user_profile.emergency_tel IS '紧急联系人电话';
COMMENT ON COLUMN public.user_profile.address IS '通讯地址';
COMMENT ON COLUMN public.user_profile.email IS '邮箱';
COMMENT ON COLUMN public.user_profile.valid_date_begin IS '身份证有效期开始时间';
COMMENT ON COLUMN public.user_profile.valid_date_end IS '身份证有效期截止时间';
COMMENT ON COLUMN public.user_profile.schooling IS '学历';
COMMENT ON COLUMN public.user_profile.degree_number IS '学位编号';
COMMENT ON COLUMN public.user_profile.remark IS '备注';
COMMENT ON COLUMN public.user_profile.professional IS '专业';
COMMENT ON COLUMN public.user_profile.status IS '用户行为状态：1正常，2已注销';
COMMENT ON COLUMN public.user_profile.created_at IS '记录创建时间';
COMMENT ON COLUMN public.user_profile.updated_at IS '记录修改时间';
COMMENT ON COLUMN public.user_profile.deleted_at IS '删除时间';
COMMENT ON COLUMN public.user_profile.deleted IS '是否删除：0否1是';
COMMENT ON COLUMN public.user_profile.graduated_from IS '毕业院校';
COMMENT ON COLUMN public.user_profile."source" IS '注册来源：1系统录入2微信3web端4app5其它';

-- ----------------------------
-- 代码生成配置表
-- ----------------------------
DROP TABLE IF EXISTS public.generate_code;
CREATE TABLE public.generate_code (
    id bpchar(36) NOT NULL,
    table_name varchar(100),
    data text,
    create_time timestamptz(6),
    deleted int2 DEFAULT 0,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.generate_code IS '代码生成配置表';
COMMENT ON COLUMN public.generate_code.id IS 'ID';
COMMENT ON COLUMN public.generate_code.table_name IS '表名';
COMMENT ON COLUMN public.generate_code.data IS '数据';
COMMENT ON COLUMN public.generate_code.create_time IS '创建时间';
COMMENT ON COLUMN public.generate_code.deleted IS '是否删除';