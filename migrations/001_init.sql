CREATE TABLE IF NOT EXISTS route_origin_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  observer_id TEXT NOT NULL,
  peer_address TEXT NOT NULL,
  prefix CIDR NOT NULL,
  operation TEXT NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS route_origin_events_cursor ON route_origin_events(id);
CREATE TABLE IF NOT EXISTS route_origin_checkpoints (
  observer_id TEXT PRIMARY KEY,
  sequence BIGINT NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
