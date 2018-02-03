package models

import (
	"github.com/artificial-universe-maker/core/common"
	uuid "github.com/artificial-universe-maker/go.uuid"
	"github.com/go-gorp/gorp"
)

type ProjectReviewResult int

const (
	ProjectReviewResultApprove ProjectReviewResult = iota
	ProjectReviewResultReject
)

type ProjectReview struct {
	ProjectID       uuid.UUID
	Version         int64
	Reviewer        string
	Result          ProjectReviewResult
	SeriousProblems []ReviewSeriousProblem
	MinorProblems   []ReviewMinorProblem
	Dialogues       []common.StringArray
	ReviewedAt      gorp.NullTime
}

type ReviewSeriousProblem int64

const (
	ReviewSeriousProblemSexuallyExplicit ReviewSeriousProblem = iota
	ChildEndangerment
	ViolenceDangerousActivities
	BullyingAndHarassment
	HateSpeech
	SensitiveEvent
	Gambling
	IllegalActivities
	RecreationalDrugs
	Health
	Language
	MatureContent
)

type ReviewMinorProblem int64

const (
	ConversationHangingOpenReviewMinorProblem ReviewMinorProblem = iota
)
