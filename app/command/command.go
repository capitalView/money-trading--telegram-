package command

import (
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	"main/info/payload"
	"main/info/transaction"
)

func (m *ResponseMapper) Save(ctx context.Context, message *telego.Message) (string, error) {
	res, parseErr := transaction.ParseText(message.Text, m.rate)
	if parseErr != nil {
		return "", parseErr
	}

	tr := transaction.NewTransaction(m.db, m.rate)
	if err := tr.Insert(ctx, res); err != nil {
		return "", err
	}

	pay := payload.NewPayload(m.db, message)
	if err := pay.SavePayload(ctx, tr.GetID()); err != nil {
		return "", err
	}

	return tr.GetAll(ctx), nil
}

func (m *ResponseMapper) Balances(ctx context.Context) string {
	tr := transaction.NewTransaction(m.db, m.rate)
	return tr.GetAll(ctx)
}

func (m *ResponseMapper) Edit(ctx context.Context, message *telego.Message) string {
	pay := payload.NewPayload(m.db, message)
	id, err := pay.GetTransactionId(ctx, message.MessageID)
	if err != nil {
		return fmt.Sprintf("%x", err)
	}
	result, err := transaction.ParseText(message.Text, m.rate)
	if err != nil {
		return fmt.Sprintf("%x", err)
	}
	tr := transaction.NewTransaction(m.db, m.rate)
	errorUpdate := tr.Update(ctx, id, result)

	if errorUpdate != nil {
		return fmt.Sprintf("%x", errorUpdate)
	}
	pay.UpdatePayload(ctx, id)
	return tr.GetAll(ctx)
}
