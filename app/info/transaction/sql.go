package transaction

const InsertTransaction = "insert into money (amount, currency, type, name, description, usd_rate) values ($1,$2, $3,$4, $5, $6) returning id;"

const SelectTransactions = "select amount, currency, usd_rate from money;"

const UpdateTransaction = "update money set amount = $1, currency = $2, type = $3, name = $4, description = $5, usd_rate = $6, updated_at = now() where id = $7;"
