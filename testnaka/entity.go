package testnaka

import (
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
	LastUpdatedDescription null.String `json:"last_updated_description"`
}

// Define specific structures for each event_code
type DueBillsMessage struct {
	AccountNumber   null.Int64             `json:"account_number"`
	AccountSequence null.Int64             `json:"account_sequence"`
	PrincipalAmount null.Dec2              `json:"principal_amount"`
	InterestAmount  null.Dec2              `json:"interest_amount"`
	PenaltyAmount   null.Dec2              `json:"penalty_amount"`
	VatAmount       null.Dec2              `json:"vat_amount"`
	EffectiveDate   null.String            `json:"effective_date"`
	ChannelPostDate null.String            `json:"channel_post_date"`
	CurrencyCode    null.String            `json:"currency_code"`
	ServiceBranch   null.Int64             `json:"service_branch"`
	OtherProperties map[string]interface{} `json:"other_properties"`
}

// type DueBillsMessageOtherProperties struct {
// 	Bills   []Bill    `json:"bills"`
// 	Penalty []Penalty `json:"penalty"`
// }

type FeeMessage struct {
	FeeAmount       null.Dec2              `json:"fee_amount"`
	ServiceBranch   null.Int64             `json:"service_branch"`
	OtherProperties map[string]interface{} `json:"other_properties"`
}

// type FeeMessageOtherProperties struct {
// 	Fee []Fee `json:"fee"`
// }

type OthersMessage struct {
	AccountNumber   null.Int64             `json:"account_number"`
	AccountSequence null.Int64             `json:"account_sequence"`
	PrincipalAmount null.Dec2              `json:"principal_amount"`
	InterestAmount  null.Dec2              `json:"interest_amount"`
	PenaltyAmount   null.Dec2              `json:"penalty_amount"`
	VatAmount       null.Dec2              `json:"vat_amount"`
	EffectiveDate   null.String            `json:"effective_date"`
	ChannelPostDate null.String            `json:"channel_post_date"`
	CurrencyCode    null.String            `json:"currency_code"`
	ServiceBranch   null.Int64             `json:"service_branch"`
	OtherProperties map[string]interface{} `json:"other_properties"`
}

// type OthersMessageOtherProperties struct {
// 	Penalties      []Penalty      `json:"penalties"`
// 	AdvancePayment AdvancePayment `json:"advance_payment"`
// }

// Struct for "bills"
type Bill struct {
	BillSequence          null.Int64  `json:"bill_sequence"`
	BillDueDate           null.String `json:"bill_due_date"`
	PrincipalAmount       null.Dec2   `json:"principal_amount"`
	InterestAmount        null.Dec2   `json:"interest_amount"`
	PenaltyAmount         null.Dec2   `json:"penalty_amount"`
	VatAmount             null.Dec2   `json:"vat_amount"`
	UnpaidPrincipalAmount null.Dec2   `json:"unpaid_principal_amount"`
	UnpaidInterestAmount  null.Dec2   `json:"unpaid_interest_amount"`
	UnpaidPenaltyAmount   null.Dec2   `json:"unpaid_penalty_amount"`
	UnpaidVatAmount       null.Dec2   `json:"unpaid_vat_amount"`
}

// Struct for "penalties"
type Penalty struct {
	BillSequence  null.Int64  `json:"bill_sequence"`
	BillDueDate   null.String `json:"bill_due_date"`
	PenaltyAmount null.Dec2   `json:"penalty_amount"`
}

// Struct for "fee"
type Fee struct {
	LoanDueDate  null.String `json:"loan_due_date"`
	BillSequence null.Int64  `json:"bill_sequence"`
	FeeAmount    null.Dec2   `json:"fee_amount"`
}

// Struct for "advance_payment"
type AdvancePayment struct {
	PrincipalAmount null.Dec2 `json:"principal_amount"`
	InterestAmount  null.Dec2 `json:"interest_amount"`
	PenaltyAmount   null.Dec2 `json:"penalty_amount"`
}
