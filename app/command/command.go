package command

import (
	"fmt"
	"github.com/mymmrac/telego"
	"main/info/payload"
	"main/info/transaction"
)

func (m *ResponseMapper) Save(message *telego.Message) (string, error) {
	res, parseErr := transaction.ParseText(message.Text, m.rate)
	if parseErr != nil {
		return "", parseErr
	}
	tr := transaction.NewTransaction(m.db, m.rate)
	id, err := tr.Insert(res)
	if err != nil {
		return "", err
	}
	pay := payload.NewPayload(m.db, message)
	pay.SavePayload(id)
	return tr.GetAll(), nil
}

func (m *ResponseMapper) Balances() string {
	return transaction.NewTransaction(m.db, m.rate).GetAll()
}

func (m *ResponseMapper) Edit(message *telego.Message) string {
	pay := payload.NewPayload(m.db, message)
	id, err := pay.GetTransactionId(message.MessageID)
	if err != nil {
		return fmt.Sprintf("%x", err)
	}
	result, err := transaction.ParseText(message.Text, m.rate)
	if err != nil {
		return fmt.Sprintf("%x", err)
	}
	tr := transaction.NewTransaction(m.db, m.rate)
	errorUpdate := tr.Update(id, result)

	if errorUpdate != nil {
		return fmt.Sprintf("%x", errorUpdate)
	}
	pay.UpdatePayload(id)
	return tr.GetAll()
}
