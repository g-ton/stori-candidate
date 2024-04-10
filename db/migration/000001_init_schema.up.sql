CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "card_number" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "date" varchar NOT NULL,
  "transaction" float NOT NULL
);

COMMENT ON COLUMN "transactions"."transaction" IS 'can be negative ( - Debit) or possitive ( + Credit)';

ALTER TABLE "transactions" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");