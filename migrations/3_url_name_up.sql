ALTER TABLE beer ADD COLUMN url_name VARCHAR NOT NULL;
ALTER TABLE beer ADD CONSTRAINT unique_url_name UNIQUE(url_name);
ALTER TABLE BEER ADD COLUMN page_ipa_ml VARCHAR NOT NULL;