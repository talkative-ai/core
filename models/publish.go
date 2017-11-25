package models

type PublishStatus int

const (
	PublishStatusNotPublished PublishStatus = iota
	PublishStatusPublishing
	PublishStatusPublished
	PublishStatusProblem
)
