CREATE TABLE IF NOT EXISTS plants
(
    id                     BIGSERIAL PRIMARY KEY,
    name                   TEXT NOT NULL,
    creation_time          timestamptz NOT NULL DEFAULT (now()),
    can_die                bool NOT NULL ,
    water_consumption_rate numeric NOT NULL,
    minimum_water_level    numeric DEFAULT 20 NOT NULL,
    water_level            numeric NOT NULL,
    last_watered           timestamptz NOT NULL ,
    growth_rate            numeric,
    growth                 numeric,
    growth_stage           text,
    last_updated           timestamptz NOT NULL DEFAULT (now()),
    backdrop               text,
    mascot                 text
);

