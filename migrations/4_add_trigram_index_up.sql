CREATE extension pg_trgm;
CREATE INDEX trgm_idx ON beer USING GIST ("name" gist_trgm_ops);