ALTER TABLE beer DROP COLUMN url_name;
ALTER TABLE beer DROP CONSTRAINT unique_url_name;
ALTER TABLE BEER DROP COLUMN page_ipa_ml VARCHAR;