package bank

import "fmt"

type TransactionID string

type Transaction struct {
	id        TransactionID
	from      AccountID
	to        AccountID
	amount    float64
	reference string
}

/*
*
NewTransaction creates a new Transaction instance with the specified parameters.
A transaction is always between two accounts and always has the direction from `from` to the `to` account, in that the amount is assumed to be "subtracted" from the `from` account and "added" to the `to` account.

  - If the amount is positive, it represents a withdrawal from the `from` account and a deposit to the `to` account.

  - If the amount is zero, it represents a transfer between the two accounts with no net change in balance.

  - if the amount is negative, it's a withdraw from the `to` account and a deposit to the `from` account.

    Parameters:

  - id: The unique identifier for the transaction.

  - from: The AccountID of the sender.

  - to: The AccountID of the receiver.

  - amount: The amount of money being transferred in the transaction.

  - reference: A reference string for the transaction, which can be used for additional information or tracking.

    Returns:

  - A pointer to the newly created Transaction instance.
*/
func NewTransaction(id TransactionID, from, to AccountID, amount float64, reference string) *Transaction {
	return &Transaction{
		id:        id,
		from:      from,
		to:        to,
		amount:    amount,
		reference: reference,
	}
}

func (t *Transaction) ID() TransactionID {
	return t.id
}

func (t *Transaction) From() AccountID {
	return t.from
}

func (t *Transaction) To() AccountID {
	return t.to
}

func (t *Transaction) Amount() float64 {
	return t.amount
}

func (t *Transaction) Reference() string {
	return t.reference
}

func (t *Transaction) String() string {
	return "Transaction{" +
		"id: " + string(t.id) +
		", from: " + string(t.from) +
		", to: " + string(t.to) +
		", amount: " + fmt.Sprintf("%.2f", t.amount) +
		", reference: '" + t.reference + "'" +
		"}"
}
