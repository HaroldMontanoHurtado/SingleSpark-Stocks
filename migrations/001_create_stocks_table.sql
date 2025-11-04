-- dirección del archivo: /migrations/001_create_recommendations.sql
-- Crea la tabla 'stocks' para almacenar recomendaciones y el JSON raw recibido.
-- Si tu CockroachDB no soporta JSONB, cambia JSONB por JSON o BYTES según corresponda.

CREATE TABLE IF NOT EXISTS stocks (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    ticker STRING NOT NULL,
    company STRING,
    brokerage STRING,
    action STRING,
    rating_from STRING,
    rating_to STRING,
    target_from STRING,
    target_to STRING,
    raw JSONB,
    source_ts TIMESTAMPTZ DEFAULT now(),
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Índices para búsqueda por ticker y fecha
CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks (ticker);
CREATE INDEX IF NOT EXISTS idx_stocks_source_ts ON stocks (source_ts DESC);
