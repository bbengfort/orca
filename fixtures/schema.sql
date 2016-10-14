/**
 * orca_schema.sql
 * Copyright 2016 University of Maryland
 *
 * Author:  Benjamin Bengfort <benjamin@bengfort.com>
 * Created: Fri Oct 14 16:11:30 2016 -0400
 */

-------------------------------------------------------------------------
-- Ensure transaction security by placing all CREATE and ALTER statements
-- inside of `BEGIN` and `COMMIT` statements.
-------------------------------------------------------------------------

BEGIN;

/**
 *  CREATE ENTITY TABLES
 */

-------------------------------------------------------------------------
-- devices Table
-------------------------------------------------------------------------

-- DROP TABLE IF EXISTS "devices";

CREATE TABLE "devices"
(
    "id" INTEGER PRIMARY KEY,
    "name" TEXT NOT NULL UNIQUE,
    "ipaddr" TEXT,
    "domain" TEXT,
    "created" DATETIME,
    "updated" DATETIME
);

-------------------------------------------------------------------------
-- locations Table
-------------------------------------------------------------------------

-- DROP TABLE IF EXISTS "locations";

CREATE TABLE "locations"
(
    "id" INTEGER PRIMARY KEY,
    "ipaddr" TEXT NOT NULL UNIQUE,
    "latitude" REAL,
    "longitude" REAL,
    "city" TEXT,
    "postcode" TEXT,
    "country" TEXT,
    "organization",
    "domain" TEXT,
    "created" DATETIME,
    "updated" DATETIME
);

-------------------------------------------------------------------------
-- pings Table
-------------------------------------------------------------------------

-- DROP TABLE IF EXISTS "pings";

CREATE TABLE "pings"
(
    "id" INTEGER PRIMARY KEY,
    "source_id" INTEGER NOT NULL,
    "target_id" INTEGER NOT NULL,
    "location_id" INTEGER,
    "request" INTEGER NOT NULL,
    "response" INTEGER,
    "sent" DATETIME NOT NULL,
    "recv" DATETIME,
    FOREIGN KEY ("source_id") REFERENCES devices("id"),
    FOREIGN KEY ("target_id") REFERENCES devices("id"),
    FOREIGN KEY ("location_id") REFERENCES locations("id")
);

 /**
  *  CREATE INDICIES
  */

 COMMIT;

 -------------------------------------------------------------------------
 -- No CREATE or ALTER statements should be outside of the `COMMIT`.
 -------------------------------------------------------------------------
