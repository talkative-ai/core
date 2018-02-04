package models

import (
	"database/sql/driver"
	"fmt"
	"strings"

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
	Dialogues     common.StringArray2D
	ReviewedAt    gorp.NullTime
}

type ReviewMajorProblem int64

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

type ReviewMinorProblem int64

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
	v := []string{}
	for _, a := range *arr {
		v = append(v, fmt.Sprintf("%v", a))
	}

	s := strings.Join(v, ",")
	return fmt.Sprintf("{%v}", s), nil
}

func (arr *ReviewMajorProblemArray) Value() (driver.Value, error) {
	v := []string{}
	for _, a := range *arr {
		v = append(v, fmt.Sprintf("%v", a))
	}

	s := strings.Join(v, ",")
	return fmt.Sprintf("{%v}", s), nil
}
