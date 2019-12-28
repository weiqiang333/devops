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


CREATE TABLE "ldap_pwd_expired" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "name" VARCHAR(30) NOT NULL,
    "pwd_last_set"  TIMESTAMP,
    "pwd_expired"   TIMESTAMP
);
CREATE UNIQUE index "index_name_on_ldap_pwd_expired" on "ldap_pwd_expired" ("name");


CREATE TABLE "rds_rsync_order" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "database" VARCHAR NOT NULL,
    "priority" INT DEFAULT 10,
    "authorized_user" VARCHAR(30) NOT NULL
);

CREATE TABLE "rds_rsync_workorder" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "database" VARCHAR NOT NULL,
    "username" VARCHAR(30) NOT NULL,
    "created_at" TIMESTAMP DEFAULT (now() at time zone 'utc'),
    "pass_at" TIMESTAMP DEFAULT (now() at time zone 'utc'),
    "order_status" VARCHAR(20) DEFAULT 'review'
);

CREATE TABLE "rds_rsync_order_logs" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "workorderid" INT NOT NULL,
    "orderid" INT NOT NULL,
    "status" BOOLEAN,
    "created_at" TIMESTAMP DEFAULT (now() at time zone 'utc')
);
