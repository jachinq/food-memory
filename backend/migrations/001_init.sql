-- 食忆 MVP 初始表结构（PostgreSQL）

CREATE TABLE IF NOT EXISTS dishes (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    cover_image_url TEXT,
    source_url TEXT,
    source_platform VARCHAR(50),
    description TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'want_to_cook',
    rating NUMERIC(2,1),
    difficulty SMALLINT,
    cook_time_minutes INTEGER,
    main_ingredients TEXT,
    taste VARCHAR(80),
    scene VARCHAR(80),
    note TEXT,
    last_cooked_at TIMESTAMPTZ,
    cook_count INTEGER NOT NULL DEFAULT 0,
    is_favorite BOOLEAN NOT NULL DEFAULT FALSE,
    extra JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_dishes_status ON dishes(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_dishes_last_cooked_at ON dishes(last_cooked_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_dishes_rating ON dishes(rating) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_dishes_name ON dishes(name) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS cook_records (
    id BIGSERIAL PRIMARY KEY,
    dish_id BIGINT NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    cooked_at TIMESTAMPTZ NOT NULL,
    result VARCHAR(30) NOT NULL DEFAULT 'normal',
    rating NUMERIC(2,1),
    notes TEXT,
    changes TEXT,
    failure_reason TEXT,
    next_improvement TEXT,
    extra JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_cook_records_dish_id ON cook_records(dish_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_cook_records_cooked_at ON cook_records(cooked_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_cook_records_result ON cook_records(result) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(30) NOT NULL DEFAULT 'custom',
    usage_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uk_tags_name_type UNIQUE(name, type)
);

CREATE INDEX IF NOT EXISTS idx_tags_type ON tags(type) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tags_usage_count ON tags(usage_count) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS dish_tags (
    dish_id BIGINT NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (dish_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_dish_tags_tag_id ON dish_tags(tag_id);

CREATE TABLE IF NOT EXISTS attachments (
    id BIGSERIAL PRIMARY KEY,
    biz_type VARCHAR(50) NOT NULL,
    biz_id BIGINT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_url TEXT NOT NULL,
    thumbnail_url TEXT,
    mime_type VARCHAR(80),
    file_size BIGINT,
    width INTEGER,
    height INTEGER,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_attachments_biz ON attachments(biz_type, biz_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS recook_plans (
    id BIGSERIAL PRIMARY KEY,
    dish_id BIGINT NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    planned_date DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'active',
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_recook_plans_status ON recook_plans(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_recook_plans_planned_date ON recook_plans(planned_date) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS app_settings (
    id BIGSERIAL PRIMARY KEY,
    setting_key VARCHAR(100) NOT NULL UNIQUE,
    setting_value JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL UNIQUE,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
