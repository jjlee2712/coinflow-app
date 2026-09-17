CREATE TABLE expenses (
      id          BIGSERIAL PRIMARY KEY,
      user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
      category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
      amount      NUMERIC(12,2) NOT NULL,
      note        TEXT,
      date        DATE NOT NULL,
      created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );