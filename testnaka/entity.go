package testnaka

import (
	"encoding/json"

	"github.com/TN-INCORPORATION/kit/v2/null"
)

type Body struct {
	ReqBody []Transaction `json:"rs_body,omitempty"`
}

// Define the main structure
type Transaction struct {
	TransactionDate null.String `json:"transaction_date"`
	ChronoSequence  null.String `json:"chrono_sequence"`
	JobID           null.String `json:"job_id"`
	AccountNumber   null.Int64  `json:"account_number"`
	AccountSequence null.Int64  `json:"account_sequence"`
	EventCode       null.String `json:"event_code"`
	Message         null.String `json:"message"`
}

// Define specific structures for each event_code
type DueBillsMessage struct {
	AccountNumber   null.Int64      `json:"account_number"`
	AccountSequence null.Int64      `json:"account_sequence"`
	PrincipalAmount null.Dec2       `json:"principal_amount"`
	InterestAmount  null.Dec2       `json:"null.Int64erest_amount"`
	PenaltyAmount   null.Dec2       `json:"penalty_amount"`
	VatAmount       null.Dec2       `json:"vat_amount"`
	EffectiveDate   null.String     `json:"effective_date"`
	ChannelPostDate null.String     `json:"channel_post_date"`
	CurrencyCode    null.String     `json:"currency_code"`
	ServiceBranch   null.Int64      `json:"service_branch"`
	OtherProperties json.RawMessage `json:"other_properties"`
}

type FeeMessage struct {
	FeeAmount       null.Dec2       `json:"fee_amount"`
	ServiceBranch   null.Int64      `json:"service_branch"`
	OtherProperties json.RawMessage `json:"other_properties"`
}

type OthersMessage struct {
	AccountNumber   null.Int64      `json:"account_number"`
	AccountSequence null.Int64      `json:"account_sequence"`
	PrincipalAmount null.Dec2       `json:"principal_amount"`
	InterestAmount  null.Dec2       `json:"null.Int64erest_amount"`
	PenaltyAmount   null.Dec2       `json:"penalty_amount"`
	VatAmount       null.Dec2       `json:"vat_amount"`
	EffectiveDate   null.String     `json:"effective_date"`
	ChannelPostDate null.String     `json:"channel_post_date"`
	CurrencyCode    null.String     `json:"currency_code"`
	ServiceBranch   null.Int64      `json:"service_branch"`
	OtherProperties json.RawMessage `json:"other_properties"`
}
