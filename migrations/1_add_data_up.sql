CREATE TABLE IF NOT EXISTS "results"
(
    "id" SERIAL UNIQUE NOT NULL,
    "url" text UNIQUE NOT NULL,
    "statuscode" integer NOT NULL,
    "text" text NOT NULL,
    CONSTRAINT "Results_pkey" PRIMARY KEY ("id")
);


