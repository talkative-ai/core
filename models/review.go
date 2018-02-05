package models

import (
	"database/sql/driver"
	"encoding/json"

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
	ProjectID     uuid.UUID
	Version       int64
	Reviewer      string
	BadTitle      bool
	Result        ProjectReviewResult
	MajorProblems ReviewMajorProblemArray
	MinorProblems ReviewMinorProblemArray
	ProblemWith   ReviewProblemWith
	Dialogues     common.StringArray2DJSON
	ReviewedAt    gorp.NullTime
}

type ProjectReviewPublic struct {
	BadTitle      bool
	MajorProblems ReviewMajorProblemArray
	MinorProblems ReviewMinorProblemArray
	ProblemWith   ReviewProblemWith
	Dialogues     common.StringArray2DJSON
}

type ReviewMajorProblem int

const (
	ReviewMajorProblemSexuallyExplicit ReviewMajorProblem = iota
	ReviewMajorProblemChildEndangerment
	ReviewMajorProblemViolenceDangerousActivities
	ReviewMajorProblemBullyingAndHarassment
	ReviewMajorProblemHateSpeech
	ReviewMajorProblemSensitiveEvent
	ReviewMajorProblemGambling
	ReviewMajorProblemIllegalActivities
	ReviewMajorProblemRecreationalDrugs
	ReviewMajorProblemHealth
	ReviewMajorProblemLanguage
	ReviewMajorProblemMatureContent
)

type ReviewMinorProblem int

const (
	ReviewMinorProblemConversationHangingOpen ReviewMinorProblem = iota
	ReviewMinorProblemZoneEntryHangingOpen
)

type ReviewProblemWith int

const (
	ReviewProblemWithDialog ReviewProblemWith = iota
	ReviewProblemWithZoneIntroductionTrigger
)

type ReviewMajorProblemArray []ReviewMajorProblem
type ReviewMinorProblemArray []ReviewMinorProblem

func (arr *ReviewMinorProblemArray) Value() (driver.Value, error) {
	return json.Marshal(*arr)
}

func (arr *ReviewMajorProblemArray) Value() (driver.Value, error) {
	return json.Marshal(*arr)
}

func (a *ReviewMajorProblemArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}

func (a *ReviewMinorProblemArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}
