package model

import "gorm.io/plugin/soft_delete"

type Repository struct {
	BaseModel
	Name        string                `gorm:"uniqueIndex:idx_url_name;type:varchar(256);not null;column:name;comment:仓库名称"`
	Url         string                `gorm:"uniqueIndex:idx_url_name;type:varchar(256);not null;column:url;comment:仓库地址"`
	Description string                `gorm:"type:longText;column:description;comment:仓库描述"`
	DeletedAt   soft_delete.DeletedAt `gorm:"uniqueIndex:idx_url_name"`
}

func (r Repository) TableName() string {
	return "repository"
}

type RepositoryDo struct {
	Repository
}

type RepositoryImage struct {
	BaseModel
	Name         string                `gorm:"uniqueIndex:idx_repo_name;type:varchar(256);not null;column:name;comment:镜像名称"`
	RepositoryID int64                 `gorm:"uniqueIndex:idx_repo_name;type:int;not null;column:repository_id;comment:仓库id"`
	Description  string                `gorm:"type:longText;column:description;comment:镜像描述"`
	DeletedAt    soft_delete.DeletedAt `gorm:"uniqueIndex:idx_repo_name"`
}

type RepositoryImageDo struct {
	RepositoryImage
}

func (ri RepositoryImage) TableName() string {
	return "repository_image"
}

type RepositoryImageTag struct {
	BaseModel
	Name      string                `gorm:"uniqueIndex:idx_tag_name;type:varchar(256);not null;column:name;comment:标签名称"`
	ImageID   int64                 `gorm:"uniqueIndex:idx_tag_name;type:int;not null;column:image_id;comment:镜像id"`
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:idx_tag_name"`
}

func (rit RepositoryImageTag) TableName() string {
	return "repository_image_tag"
}

type RepositoryImageTagDo struct {
	RepositoryImageTag
}
