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

CREATE TABLE "rds_rsync_workorder_logs" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "workorderid" INT NOT NULL,
    "username" VARCHAR(30) NOT NULL,
    "created_at" TIMESTAMP DEFAULT (now() at time zone 'utc'),
    "get_snapshot_at" TIMESTAMP DEFAULT (now() at time zone 'utc'),
    "delete_at" TIMESTAMP DEFAULT (now() at time zone 'utc'),
    "restore_at" TIMESTAMP DEFAULT (now() at time zone 'utc'),
    "modify_config_at" TIMESTAMP DEFAULT (now() at time zone 'utc'),
    "execute_sql_at" TIMESTAMP DEFAULT (now() at time zone 'utc'),
    "status" VARCHAR(30) NOT NULL
);
CREATE UNIQUE index "index_workorderid_on_rds_rsync_workorder_logs" on "rds_rsync_workorder_logs" ("workorderid");


CREATE TABLE "aws_instance_types" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "instance_type" VARCHAR(30) NOT NULL,
    "vcpus" INT NOT NULL,
    "memory" INT NOT NULL,
    "update_at" TIMESTAMP DEFAULT (now() at time zone 'utc')
)
CREATE UNIQUE index "index_instance_type_on_aws_instance_types" on "aws_instance_types" ("instance_type");

CREATE TABLE "aws_volumes" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "disk_id"   VARCHAR NOT NULL,
    "create_at" TIMESTAMP,
    "size"  INT NOT NULL,
    "iops"  INT,
    "state"   VARCHAR,
    "type"  VARCHAR,
    "tag_name"  VARCHAR,
    "tag_app"   VARCHAR,
    "tag_env"   VARCHAR,
    "tag_pillar"    VARCHAR,
    "update_at" TIMESTAMP DEFAULT (now() at time zone 'utc')
)
CREATE UNIQUE index "index_disk_id_on_aws_volumes" on "aws_volumes" ("disk_id");