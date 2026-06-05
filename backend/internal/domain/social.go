package domain

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	UserID    uuid.UUID `json:"userId"`
	PoemID    uuid.UUID `json:"poemId"`
	CreatedAt time.Time `json:"createdAt"`
}

type Follow struct {
	FollowerID uuid.UUID `json:"followerId"`
	FollowedID uuid.UUID `json:"followedId"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Comment struct {
	ID        uuid.UUID  `json:"id"`
	PoemID    uuid.UUID  `json:"poemId"`
	AuthorID  uuid.UUID  `json:"authorId"`
	Author    *User      `json:"author,omitempty"`
	ParentID  *uuid.UUID `json:"parentId,omitempty"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type NotificationType string

const (
	NotifLike           NotificationType = "like"
	NotifComment        NotificationType = "comment"
	NotifFollow         NotificationType = "follow"
	NotifStanzaSubmit   NotificationType = "stanza_submitted"
	NotifStanzaApproved NotificationType = "stanza_approved"
	NotifStanzaRejected NotificationType = "stanza_rejected"
	NotifPoemFeatured   NotificationType = "poem_featured"
)

type Notification struct {
	ID          uuid.UUID        `json:"id"`
	RecipientID uuid.UUID        `json:"recipientId"`
	ActorID     *uuid.UUID       `json:"actorId,omitempty"`
	Actor       *User            `json:"actor,omitempty"`
	Type        NotificationType `json:"type"`
	PoemID      *uuid.UUID       `json:"poemId,omitempty"`
	Poem        *Poem            `json:"poem,omitempty"`
	Read        bool             `json:"read"`
	CreatedAt   time.Time        `json:"createdAt"`
}

type Tutorial struct {
	ID           uuid.UUID  `json:"id"`
	Title        string     `json:"title"`
	Slug         string     `json:"slug"`
	Format       PoemFormat `json:"format"`
	ContentMD    string     `json:"contentMd"`
	Difficulty   string     `json:"difficulty"`
	DisplayOrder int        `json:"displayOrder"`
	CreatedAt    time.Time  `json:"createdAt"`
}
