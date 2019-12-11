CREATE TABLE "server_list" (
  "id" SERIAL NOT NULL PRIMARY KEY,
  "server" VARCHAR NOT NULL,
  "name" VARCHAR(20) DEFAULT '' ,
  "app" VARCHAR(64) DEFAULT '',
  "pillar" VARCHAR(20) DEFAULT '',
  "status" BOOLEAN NOT NULL DEFAULT FALSE,
  "uptime" TIMESTAMP DEFAULT (now() at time zone 'utc')
);
CREATE UNIQUE index "index_server_on_server_list" on "server_list" ("server");


CREATE TABLE "service" (
  "id" SERIAL NOT NULL PRIMARY KEY,
  "server" VARCHAR NOT NULL,
  "service" VARCHAR,
  "status" VARCHAR(10) DEFAULT 'unknown'
);


CREATE TABLE "google_auth" (
  "id" SERIAL NOT NULL PRIMARY KEY,
  "name" VARCHAR(30) NOT NULL,
  "secret"  VARCHAR(32) NOT NULL,
  "created_at" TIMESTAMP DEFAULT (now() at time zone 'utc'),
  "updated_at" TIMESTAMP DEFAULT (now() at time zone 'utc')
);
CREATE UNIQUE index "index_name_on_google_auth" on "google_auth" ("name");
