package model

type ImageInfo struct {
	BaseModel
	ImageName    string `gorm:"column:image_name;type:varchar(255);NOT NULL" json:"image_name"`
	BuildType    int    `gorm:"column:build_type;type:tinyint(4);default:0;comment:构建方式 1:from self 2docker file 3 harbor;NOT NULL" json:"build_type"`
	ImageExt     string `gorm:"column:image_ext;type:text;NOT NULL" json:"image_ext"`
	ImageUserFor string `gorm:"column:image_user_for;type:varchar(1000);NOT NULL" json:"image_user_for"`
	CreateUser   int64  `gorm:"column:create_user;type:bigint(20);default:0;NOT NULL" json:"create_user"`
	UpdateUser   int64  `gorm:"column:update_user;type:bigint(20);default:0;NOT NULL" json:"update_user"`
}

func (m *ImageInfo) TableName() string {
	return "image_info"
}

type ImageBuild struct {
	BaseModel
	ImageID     int64  `gorm:"column:image_id;type:bigint(20);default:0;NOT NULL" json:"image_id"`
	BuildStatus int    `gorm:"column:build_status;type:tinyint(4);default:0;NOT NULL" json:"build_status"`
	BuildLog    string `gorm:"column:build_log;type:longtext;NOT NULL" json:"build_log"`
	BuildUser   int64  `gorm:"column:build_user;type:bigint(20);default:0;NOT NULL" json:"build_user"`
}

func (m *ImageBuild) TableName() string {
	return "image_build"
}

type ImageTag struct {
	BaseModel
	TagName    string `gorm:"column:tag_name;type:varchar(255);NOT NULL" json:"tag_name"`
	TagDesc    string `gorm:"column:tag_desc;type:varchar(2000);NOT NULL" json:"tag_desc"`
	ImageId    int64  `gorm:"column:image_id;type:bigint(20);default:0;NOT NULL" json:"image_id"`
	CreateUser int64  `gorm:"column:create_user;type:bigint(20);default:0;NOT NULL" json:"create_user"`
}

func (m *ImageTag) TableName() string {
	return "image_tag"
}

type ImageDo struct {
	ImageInfo
}
