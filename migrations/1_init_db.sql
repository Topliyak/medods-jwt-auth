DROP TABLE IF EXISTS "users" CASCADE;

CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY,
    "email" VARCHAR(50) UNIQUE,
    "password" VARCHAR(100),
    "refresh" VARCHAR(100),
    "refresh_issued_at" TIME
);
