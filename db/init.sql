CREATE TABLE addresses (
  address TEXT PRIMARY KEY,
  value NUMERIC NOT NULL DEFAULT '0',
  nonce BIGINT NOT NULL DEFAULT '0'
);

CREATE TABLE blocks (
  "number" BIGSERIAL PRIMARY KEY,
  hash TEXT NOT NULL UNIQUE,
  confirmations NUMERIC NOT NULL DEFAULT '0',
  parentHash TEXT NOT NULL DEFAULT '',
  nonce TEXT NOT NULL DEFAULT '',
  sha3Uncles TEXT NOT NULL DEFAULT '',
  logsBloom TEXT NOT NULL DEFAULT '',
  transactionsRoot TEXT NOT NULL DEFAULT '',
  stateRoot TEXT NOT NULL DEFAULT '',
  receiptsRoot TEXT NOT NULL,
  miner TEXT NOT NULL DEFAULT '',
  difficulty NUMERIC NOT NULL DEFAULT 0,
  size BIGINT NOT NULL DEFAULT 0,
  extraData TEXT NOT NULL DEFAULT '',
  gasLimit BIGINT NOT NULL DEFAULT 0,
  gasUsed BIGINT NOT NULL DEFAULT 0,
  "timestamp" BIGINT NOT NULL DEFAULT '0',
  mixhash TEXT DEFAULT ''
);

CREATE TABLE transactions (
  hash TEXT PRIMARY KEY,
  nonce BIGINT,
  blockHash TEXT NOT NULL REFERENCES blocks(hash) ON DELETE CASCADE ON UPDATE CASCADE,
  blockNumber BIGINT NOT NULL REFERENCES blocks("number") ON DELETE CASCADE ON UPDATE CASCADE,
  transactionIndex BIGINT NOT NULL,
  "from" TEXT NOT NULL,
  "to" TEXT NOT NULL,
  "value" NUMERIC NOT NULL,
  gas BIGINT NOT NULL DEFAULT 0,
  gasPrice NUMERIC NOT NULL DEFAULT 0,
  "input" TEXT DEFAULT '',
  v TEXT DEFAULT '',
  r TEXT DEFAULT '',
  s TEXT DEFAULT ''
);


CREATE VIEW view_last_block
AS
SELECT b.number
FROM blocks b
WHERE b.number = (SELECT max(b2.number) FROM blocks b2);

CREATE INDEX idx_transactions_from ON transactions("from");
CREATE INDEX idx_transactions_to ON transactions("to");
