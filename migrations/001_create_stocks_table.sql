-- dirección del archivo: /migrations/001_create_stocks_table.sql
-- Crea la tabla 'stocks' para almacenar recomendaciones y el JSON raw recibido.
-- Si tu CockroachDB no soporta JSONB, cambia JSONB por JSON o BYTES según corresponda.

CREATE TABLE IF NOT EXISTS stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker TEXT NOT NULL,
    company TEXT,
    brokerage TEXT,
    action TEXT,
    rating_from TEXT,
    rating_to TEXT,
    target_from TEXT,
    target_to TEXT,
    raw JSONB,
    source_ts TIMESTAMPTZ DEFAULT now(),
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Índices para búsqueda por ticker y fecha
CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks (ticker);
CREATE INDEX IF NOT EXISTS idx_stocks_source_ts ON stocks (source_ts DESC);
