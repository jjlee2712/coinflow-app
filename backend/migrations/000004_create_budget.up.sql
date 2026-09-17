CREATE TABLE budgets (
      id          BIGSERIAL PRIMARY KEY,
      user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
      category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
      amount      NUMERIC(12,2) NOT NULL,
      created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      UNIQUE(user_id, category_id)
  );