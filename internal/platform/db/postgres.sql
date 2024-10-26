create table if not exists accounts
(
    id
    bigserial
    primary
    key,
    created_at
    timestamp
    with
    time
    zone,
    updated_at
    timestamp
    with
    time
    zone,
    deleted_at
    timestamp
    with
    time
    zone,
    name
    text,
    last_name
    text,
    email
    text
);


create index if not exists idx_accounts_deleted_at
    on accounts (deleted_at);

create table if not exists transactions
(
    id
    bigserial
    primary
    key,
    created_at
    timestamp
    with
    time
    zone,
    updated_at
    timestamp
    with
    time
    zone,
    deleted_at
    timestamp
    with
    time
    zone,
    external_id
    text
    UNIQUE,
    date
    timestamp
    with
    time
    zone,
    amount
    bigint,
    type
    text,
    account_id
    bigint
    constraint
    fk_accounts_transactions
    references
    accounts
);

create index if not exists idx_transactions_deleted_at
    on transactions (deleted_at);


CREATE
MATERIALIZED VIEW IF NOT EXISTS monthly_balances AS
SELECT account_id,
       SUM(CASE WHEN type = 'credit' THEN amount ELSE 0 END) +
       SUM(CASE WHEN type = 'debit' THEN amount ELSE 0 END)                   AS total_balance,
       COUNT(*)                                                               AS transaction_count,
       CAST(SUM(CASE WHEN type = 'credit' THEN amount ELSE 0 END) /
            NULLIF(COUNT(CASE WHEN type = 'credit' THEN 1 END), 0) AS bigint) AS avg_credit_amount,
       CAST(SUM(CASE WHEN type = 'debit' THEN amount ELSE 0 END) /
            NULLIF(COUNT(CASE WHEN type = 'debit' THEN 1 END), 0) AS bigint)  AS avg_debit_amount,
       TO_CHAR(date, 'MM') AS month
FROM transactions
GROUP BY account_id, TO_CHAR(date, 'MM');

CREATE UNIQUE INDEX IF NOT EXISTS idx_monthly_balances_account_month ON monthly_balances (account_id, month);

CREATE
MATERIALIZED VIEW IF NOT EXISTS balances AS
SELECT account_id,
       SUM(CASE WHEN type = 'credit' THEN amount ELSE 0 END) +
       SUM(CASE WHEN type = 'debit' THEN amount ELSE 0 END)                   AS total_balance,
       COUNT(*)                                                               AS transaction_count,
       CAST(SUM(CASE WHEN type = 'credit' THEN amount ELSE 0 END) /
            NULLIF(COUNT(CASE WHEN type = 'credit' THEN 1 END), 0) AS bigint) AS avg_credit_amount,
       CAST(SUM(CASE WHEN type = 'debit' THEN amount ELSE 0 END) /
            NULLIF(COUNT(CASE WHEN type = 'debit' THEN 1 END), 0) AS bigint)  AS avg_debit_amount
FROM transactions
GROUP BY account_id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_balances_account ON balances (account_id);

