package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestInvoice_InvoiceAmount は InvoiceAmount() のテストです。
// InvoiceAmount() は支払金額の計算という重要なルールを内包しているので単体テストを作成しています。
func TestInvoice_InvoiceAmount(t *testing.T) {
	type fields struct {
		id             string
		companyID      string
		clientID       string
		issueDate      time.Time
		paymentAmount  int
		paymentDueDate time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "正常系: paymentAmount が 10,000 円の場合に InvoiceAmount() が 10,440 円を返す",
			fields: fields{
				id:             "test_id",
				companyID:      "test_company_id",
				clientID:       "test_client_id",
				issueDate:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				paymentAmount:  10000,
				paymentDueDate: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			want: 10440,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewInvoice(
				tt.fields.id,
				tt.fields.companyID,
				tt.fields.clientID,
				tt.fields.issueDate,
				tt.fields.paymentDueDate,
				tt.fields.paymentAmount,
			)
			assert.Equalf(t, tt.want, i.InvoiceAmount(), "InvoiceAmount()")
		})
	}
}
