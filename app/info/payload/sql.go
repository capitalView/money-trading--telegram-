package payload

const GetTransactionId = "select transaction_id from payload WHERE (data->>'message_id')::int = $1;"

const InsertPayload = "insert into payload (transaction_id, data) values ($1, $2);"

const UpdatePayload = "UPDATE payload SET data = jsonb_set(data, '{update}', $2) WHERE transaction_id = $1;"
