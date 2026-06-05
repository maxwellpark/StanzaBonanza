BEGIN;

-- Enums

CREATE TYPE poem_format AS ENUM (
    'free_verse',
    'haiku',
    'sonnet',
    'limerick',
    'iambic_pentameter',
    'rhyming_couplets',
    'custom'
);

CREATE TYPE approval_mode AS ENUM (
    'open',
    'approval_required',
    'closed'
);

CREATE TYPE notification_type AS ENUM (
    'like',
    'comment',
    'follow',
    'stanza_submitted',
    'stanza_approved',
    'stanza_rejected',
    'poem_featured'
);

-- 1. users

CREATE TABLE users (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    display_name    VARCHAR(50) NOT NULL,
    email           VARCHAR(254) NOT NULL UNIQUE,
    bio             TEXT,
    avatar_url      TEXT,
    is_verified     BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 2. webauthn_credentials

CREATE TABLE webauthn_credentials (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    credential_id   BYTEA       NOT NULL UNIQUE,
    public_key      BYTEA       NOT NULL,
    sign_count      BIGINT      NOT NULL DEFAULT 0,
    transports      TEXT[],
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 3. magic_links

CREATE TABLE magic_links (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(254)    NOT NULL,
    token_hash      VARCHAR(128)    NOT NULL UNIQUE,
    expires_at      TIMESTAMPTZ     NOT NULL,
    used_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now()
);

-- 4. sessions

CREATE TABLE sessions (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash      VARCHAR(128)    NOT NULL UNIQUE,
    expires_at      TIMESTAMPTZ     NOT NULL,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now()
);

-- 5. poems

CREATE TABLE poems (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id           UUID            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title               VARCHAR(200)    NOT NULL,
    description         TEXT,
    format              poem_format     NOT NULL DEFAULT 'free_verse',
    format_rules_json   JSONB           NOT NULL DEFAULT '{}',
    approval_mode       approval_mode   NOT NULL DEFAULT 'open',
    max_stanzas         INT,
    is_hall_of_fame     BOOLEAN         NOT NULL DEFAULT FALSE,
    like_count          INT             NOT NULL DEFAULT 0,
    stanza_count        INT             NOT NULL DEFAULT 0,
    comment_count       INT             NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT now()
);

-- 6. stanzas

CREATE TABLE stanzas (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    poem_id         UUID            NOT NULL REFERENCES poems(id) ON DELETE CASCADE,
    author_id       UUID            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text            TEXT            NOT NULL,
    position        INT             NOT NULL,
    literary_device VARCHAR(50),
    status          VARCHAR(20)     NOT NULL DEFAULT 'approved',
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_stanzas_poem_position ON stanzas (poem_id, position);

-- 7. likes

CREATE TABLE likes (
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    poem_id     UUID        NOT NULL REFERENCES poems(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, poem_id)
);

-- 8. follows

CREATE TABLE follows (
    follower_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followed_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (follower_id, followed_id),
    CHECK (follower_id != followed_id)
);

-- 9. comments

CREATE TABLE comments (
    id          UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    poem_id     UUID            NOT NULL REFERENCES poems(id) ON DELETE CASCADE,
    author_id   UUID            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id   UUID            REFERENCES comments(id) ON DELETE CASCADE,
    text        TEXT            NOT NULL,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT now()
);

-- 10. notifications

CREATE TABLE notifications (
    id              UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_id    UUID                NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    actor_id        UUID                REFERENCES users(id) ON DELETE SET NULL,
    type            notification_type   NOT NULL,
    poem_id         UUID                REFERENCES poems(id) ON DELETE CASCADE,
    read            BOOLEAN             NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ         NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_unread
    ON notifications (recipient_id, created_at DESC)
    WHERE NOT read;

-- 11. tutorials

CREATE TABLE tutorials (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(200)    NOT NULL,
    slug            VARCHAR(200)    NOT NULL UNIQUE,
    format          poem_format     NOT NULL,
    content_md      TEXT            NOT NULL,
    difficulty      VARCHAR(20)     NOT NULL DEFAULT 'beginner',
    display_order   INT             NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now()
);

-- 12. user_tiers

CREATE TABLE user_tiers (
    user_id     UUID        PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    tier        VARCHAR(20) NOT NULL DEFAULT 'free',
    expires_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Additional indexes

CREATE INDEX idx_poems_author_id   ON poems (author_id);
CREATE INDEX idx_poems_created_at  ON poems (created_at DESC);
CREATE INDEX idx_poems_hall_of_fame ON poems (created_at DESC) WHERE is_hall_of_fame;

CREATE INDEX idx_comments_poem_created ON comments (poem_id, created_at);
CREATE INDEX idx_follows_followed_id   ON follows (followed_id);
CREATE INDEX idx_likes_poem_id         ON likes (poem_id);

COMMIT;
