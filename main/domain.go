package main

import "time"

// Invoice は請求書を表す構造体です。
// feeRate / taxRate は今後変化する可能性があるため、インスタンスごとに設定できるようにする。
type Invoice struct {
	id             string
	companyID      string
	clientID       string
	issueDate      time.Time
	paymentAmount  int
	feeRate        float32
	taxRate        float32
	paymentDueDate time.Time
	status         Status
}

func NewInvoice(
	id string,
	companyID string,
	clientID string,
	issueDate time.Time,
	paymentDueDate time.Time,
	paymentAmount int,
) *Invoice {
	return &Invoice{
		id:             id,
		companyID:      companyID,
		clientID:       clientID,
		issueDate:      issueDate,
		paymentAmount:  paymentAmount,
		feeRate:        FeeRate,
		taxRate:        TaxRate,
		paymentDueDate: paymentDueDate,
		status:         StatusUnprocessed,
	}
}

func (i *Invoice) ID() string {
	return i.id
}

func (i *Invoice) CompanyID() string {
	return i.companyID
}

func (i *Invoice) ClientID() string {
	return i.clientID
}

func (i *Invoice) IssueDate() time.Time {
	return i.issueDate
}

func (i *Invoice) PaymentAmount() int {
	return i.paymentAmount
}

func (i *Invoice) FeeRate() float32 {
	return i.feeRate
}

func (i *Invoice) TaxRate() float32 {
	return i.taxRate
}

func (i *Invoice) Fee() int {
	return int(float32(i.paymentAmount) * i.feeRate)
}

func (i *Invoice) Tax() int {
	return int(float32(i.paymentAmount) * i.feeRate * i.taxRate)
}

func (i *Invoice) InvoiceAmount() int {
	return i.paymentAmount + i.Fee() + i.Tax()
}

func (i *Invoice) PaymentDueDate() time.Time {
	return i.paymentDueDate
}

func (i *Invoice) Status() Status {
	return i.status
}

const FeeRate = 0.04
const TaxRate = 0.1

// Status 支払い状況を表すステータス
type Status struct {
	code int
	name string
}

func (s Status) Name() string {
	return s.name
}

func (s Status) Code() int {
	return s.code
}

var StatusUnprocessed = Status{
	code: 1,
	name: "unprocessed",
}

var StatusProcessing = Status{
	code: 2,
	name: "processing",
}

var StatusPaid = Status{
	code: 3,
	name: "paid",
}

var StatusError = Status{
	code: 4,
	name: "error",
}
