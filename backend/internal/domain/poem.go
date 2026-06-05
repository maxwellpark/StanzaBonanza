package domain

import (
	"time"

	"github.com/google/uuid"
)

type PoemFormat string

const (
	FormatFreeVerse          PoemFormat = "free_verse"
	FormatHaiku              PoemFormat = "haiku"
	FormatSonnet             PoemFormat = "sonnet"
	FormatLimerick           PoemFormat = "limerick"
	FormatIambicPentameter   PoemFormat = "iambic_pentameter"
	FormatRhymingCouplets    PoemFormat = "rhyming_couplets"
	FormatCustom             PoemFormat = "custom"
)

type ApprovalMode string

const (
	ApprovalOpen     ApprovalMode = "open"
	ApprovalRequired ApprovalMode = "approval_required"
	ApprovalClosed   ApprovalMode = "closed"
)

type Poem struct {
	ID              uuid.UUID    `json:"id"`
	AuthorID        uuid.UUID    `json:"authorId"`
	Author          *User        `json:"author,omitempty"`
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	Format          PoemFormat   `json:"format"`
	FormatRulesJSON string       `json:"formatRules,omitempty"`
	ApprovalMode    ApprovalMode `json:"approvalMode"`
	MaxStanzas      *int         `json:"maxStanzas"`
	IsHallOfFame    bool         `json:"isHallOfFame"`
	LikeCount       int          `json:"likeCount"`
	StanzaCount     int          `json:"stanzaCount"`
	CommentCount    int          `json:"commentCount"`
	Stanzas         []Stanza     `json:"stanzas,omitempty"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
}

type StanzaStatus string

const (
	StanzaApproved StanzaStatus = "approved"
	StanzaPending  StanzaStatus = "pending"
	StanzaRejected StanzaStatus = "rejected"
)

type Stanza struct {
	ID             uuid.UUID    `json:"id"`
	PoemID         uuid.UUID    `json:"poemId"`
	AuthorID       uuid.UUID    `json:"authorId"`
	Author         *User        `json:"author,omitempty"`
	Text           string       `json:"text"`
	Position       int          `json:"position"`
	LiteraryDevice string       `json:"literaryDevice,omitempty"`
	Status         StanzaStatus `json:"status"`
	CreatedAt      time.Time    `json:"createdAt"`
}
