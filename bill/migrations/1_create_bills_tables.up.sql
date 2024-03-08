CREATE TABLE bills (
  id          BIGSERIAL PRIMARY KEY,
  customer_id INTEGER   NOT NULL,
  time_period INTEGER   NOT NULL,
  created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
  closed_at   TIMESTAMP
);

CREATE TABLE bill_charges (
  id         BIGSERIAL PRIMARY KEY,
  bill_id    INTEGER   NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
  amount     DECIMAL   NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
)
