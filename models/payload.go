package models

type SourceForm struct {
	Source string `form:"src" binding:"required"`
}
