package main

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 請求書データ作成 API の結合テスト
// TODO: データベースの初期化、テストデータの投入、テストデータの削除を行う
// TODO: 登録 API の実行後のデータベースの状態を確認する
// TODO: テーブル駆動テストに書き換える
// TODO: 正常系 1 パターンしかテストしていないので、異常系も含めてテストする
func TestCreateInvoiceAPI(t *testing.T) {
	// setup
	db, err := sqlx.Open("mysql", "user:password@tcp(127.0.0.1:3306)/shiharai_com_db")
	if err != nil {
		assert.Fail(t, "failed to open db")
	}

	e := NewEchoServer(db)
	testServer := httptest.NewServer(e)
	defer testServer.Close()

	// make request
	body := `{
		"company_id": "01J13SNTHMBQ8A5039QM3QB3JC",
		"client_id": "01J13SR5J68KDGNQ8Y6HECCSA5",
		"issue_date": "2024-01-01",
		"due_date": "2024-02-01",
		"payment_amount": 100000
	}`

	req, err := http.NewRequest("POST", testServer.URL+"/api/invoices", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		assert.Fail(t, "failed to create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		assert.Fail(t, "failed to send request")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		assert.Fail(t, "failed to read response body")
	}

	// assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	actual := new(CreateInvoiceResponse)
	err = json.Unmarshal(respBody, actual)
	if err != nil {
		assert.Fail(t, "failed to unmarshal response body")
	}
	assert.Regexp(t, "^[0-9A-Z]{26}$", actual.ID)

}
