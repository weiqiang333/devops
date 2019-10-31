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
  "service" VARCHAR(64),
  "status" VARCHAR(10) DEFAULT 'unknown'
);
CREATE UNIQUE index "index_server_on_service" on "service" ("server");