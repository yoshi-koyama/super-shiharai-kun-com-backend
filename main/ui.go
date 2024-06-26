package main

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"net/http"
	"time"
)

// CreateInvoiceRequest は請求書作成 API のリクエストボディを表す構造体です。
// IssueDate と DueDate は "2006-01-02" 形式の文字列を time.Time に変換するためのカスタム型を使用しています。
type CreateInvoiceRequest struct {
	CompanyID     string      `json:"company_id" validate:"required"`
	ClientID      string      `json:"client_id" validate:"required"`
	IssueDate     bindingTime `json:"issue_date" validate:"required"`
	DueDate       bindingTime `json:"due_date" validate:"required"`
	PaymentAmount int         `json:"payment_amount" validate:"required,numeric"`
}

type bindingTime time.Time

func (bt *bindingTime) String() string {
	return time.Time(*bt).Format("2006-01-02")
}

func (bt *bindingTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*bt = bindingTime(t)
	return nil
}

// CreateInvoiceResponse 請求書データ作成 API 用のレスポンスボディを表す構造体です。
type CreateInvoiceResponse struct {
	ID string `json:"id"`
}

// InvoiceHandler 請求書の登録や取得などのロジックを担います。
// 小規模なアプリケーションなので Handler 内にリクエストのバインド、バリデーション業務処理、レスポンスの返却すべてを記述します。
type InvoiceHandler struct {
	db *sqlx.DB
}

func (h *InvoiceHandler) CreateInvoice(c echo.Context) error {
	req := new(CreateInvoiceRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	tx, err := h.db.Begin()
	if err != nil {
		Logger.Error("failed to begin transaction", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	defer tx.Rollback()

	countOfCompanies := 0
	rows := tx.QueryRow("SELECT COUNT(*) FROM company WHERE id = ?", req.CompanyID)
	err = rows.Scan(&countOfCompanies)
	if err != nil {
		Logger.Error("failed to scan a count of companies", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	if countOfCompanies == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "company not found"})
	}

	countOfClients := 0
	rows = tx.QueryRow("SELECT COUNT(*) FROM client WHERE id = ?", req.ClientID)
	err = rows.Scan(&countOfClients)
	if err != nil {
		Logger.Error("failed to scan a count of clients", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	if countOfClients == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "client not found"})
	}

	id := ulid.Make().String()
	invoice := NewInvoice(id, req.CompanyID, req.ClientID, time.Time(req.IssueDate), time.Time(req.DueDate), req.PaymentAmount)

	_, err = tx.Exec(`
		INSERT INTO invoice (
		                      id, 
		                      company_id, 
		                      client_id, 
		                      issue_date, 
		                      payment_amount, 
		                      fee, 
		                      fee_rate, 
		                      tax, 
		                      tax_rate, 
		                      invoice_amount, 
		                      payment_due_date, 
		                      status) 
		VALUES (?, ?, ? ,?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		invoice.ID(),
		invoice.CompanyID(),
		invoice.ClientID(),
		invoice.IssueDate(),
		invoice.PaymentAmount(),
		invoice.Fee(),
		invoice.FeeRate(),
		invoice.Tax(),
		invoice.TaxRate(),
		invoice.InvoiceAmount(),
		invoice.PaymentDueDate(),
		StatusUnprocessed.Code(),
	)
	if err != nil {
		Logger.Error("failed to insert invoice", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	err = tx.Commit()
	if err != nil {
		Logger.Error("failed to commit transaction", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, CreateInvoiceResponse{ID: invoice.ID()})
}
