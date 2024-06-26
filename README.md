# 目的

とある企業の API コーディングテストの成果物として作成。

コーディングテストの内容は以下の通り。

https://github.com/upsidr/coding-test/blob/main/web-api-language-agnostic/README.ja.md

# API コーディングテストの進め方

完成までは 8 時間を想定しているが、 3 時間以内で終了する。

つまり、コーディングテスト記載の内容すべてを終わらせる必要はなく最初の 3 時間で何をするかということを考えて進める。

コーディングテストの評価基準をある程度満たせるように進める。

実装前に、どのような API を実装するかをドキュメントにする。

ドキュメントの形式は Google やメルカリで採用されている Design Doc を参考にする。

Design Doc について

https://engineering.mercari.com/blog/entry/20220225-design-docs-by-mercari-shops/

https://www.industrialempathy.com/posts/design-docs-at-google/

ただし、時間も限られているので簡易的なものに抑える。

できればこの Design Doc を元に残りの 5 時間で私以外でも実装を進められるところまで書けると Good。

---

# Design Doc of スーパー支払い君.com API

## 概要

スーパー支払い君.com はユーザーが未来の支払期日の請求書のデータを登録しておくと、期日に残高がなくとも自動的に銀行振り込みを行うことができ、現金の支出を最大一ヶ月遅らせることができるため、ユーザーにとって便利なウェブサービスである。

この Design Doc では、スーパー支払い君.com のバックエンドの REST API を開発することについて記述する。

## 背景

下記を参照すること。

https://github.com/upsidr/coding-test/blob/main/web-api-language-agnostic/README.ja.md

## Goals

下記の 2 つの機能をもつ API を提供する。

- ユーザーとして請求書データを新規に作成する API
- ユーザーとして、指定期間内に支払いが発生する請求書データの一覧を取得する API

## Non-Goals

- ローカル環境で API が動作すればよいので、インフラの構築はしない

## アーキテクチャ

ローカル環境で以下の構成で API を提供する。

- API サーバーを Go で実装
- データベースには MySQL を使用
- データベースの構築に Docker を使用
- 認証に Firebase Authentication を使用

まずは認証機能無しで API 実装する。  
<img width="780" alt="スクリーンショット 2024-06-26 15 41 06" src="https://github.com/yoshi-koyama/super-shiharai-kun-com-backend/assets/62045457/ad8dfcb5-68dd-44a3-8ce8-f1eef2c64e44">

API 実装後に認証機能を付加する。

<img width="810" alt="スクリーンショット 2024-06-26 15 41 44" src="https://github.com/yoshi-koyama/super-shiharai-kun-com-backend/assets/62045457/394e8aee-9ecc-4fdd-a1f4-45a0de08086d">

## API 仕様

### 請求書データ作成 API POST /api/invoices

新しい請求書データを作成する。

請求金額を下記のロジックで自動計算して自動で登録する。

```
支払金額 に手数料 4% を加えたものに更に手数料の消費税を加えたものを請求金額とする
例) 支払金額 10,000 の場合は請求金額は 10,000 + (10,000 * 0.04 * 1.10) = 10,440
```

登録時に請求書のステータスは未処理で登録される。

#### ヘッダー

Content-Type: application/json

#### リクエストボディ

- 企業 ID
    - 必須 string
- 取引先 ID
    - 必須 string
- 発行日
    - 必須 string (yyyy-mm-dd)
- 支払金額
    - 必須 int
- 支払期日
    - 必須 string (yyyy-mm-dd)

#### HTTP ステータスコード

200 登録成功
400 リクエストボディが不正
404 企業 ID または取引先 ID が存在しない
500 サーバーエラー

#### レスポンスボディ

- 請求書 ID string

### 請求書一覧データ取得 API GET /api/invoices

指定期間内に支払いが発生するつまりステータスが未処理または処理中の請求書データの一覧を取得する。

#### ヘッダー

Content-Type: application/json

#### クエリパラメータ

- 開始日
    - 必須 string (yyyy-mm-dd)
- 終了日
    - 必須 string (yyyy-mm-dd)

#### HTTP ステータスコード

200 取得成功
400 リクエストボディが不正

#### レスポンスボディ

- 請求書 ID array of strings

## データベース

### MySQL

#### company

| 列名                  | データ型         | 制約          | 備考   |
|---------------------|--------------|-------------|------|
| id                  | CHAR(26)     | PRIMARY KEY | ULID |
| company_name        | VARCHAR(255) | NOT NULL    |      |
| representative_name | VARCHAR(255) | NOT NULL    |      |
| phone_number        | VARCHAR(20)  | NOT NULL    |      |
| postal_code         | VARCHAR(10)  | NOT NULL    |      |
| address             | VARCHAR(255) | NOT NULL    |      |

### user

メールアドレスとパスワードは Firebase Authentication にて保管する。
このテーブルはメールアドレスに紐付くユーザー名と企業 ID を保管する。

| 列名         | データ型       | 制約                                                        | 備考   |
|------------| -------------- |-----------------------------------------------------------| ------ |
| id         | CHAR(26)       | PRIMARY KEY                                               | ULID   |
| company_id | CHAR(26)       | NOT NULL, FOREIGN KEY (company_id) REFERENCES company(id) |     |
| email_address | VARCHAR(255) | NOT NULL                                                  |     |
| full_name  | VARCHAR(255)   | NOT NULL                                          |     |

#### client

| 列名                  | データ型         | 制約                                                        | 備考   |
|---------------------|--------------|-----------------------------------------------------------|------|
| id                  | CHAR(26)     | PRIMARY KEY                                               | ULID |
| company_id          | CHAR(26)     | NOT NULL, FOREIGN KEY (company_id) REFERENCES company(id) |      |
| company_name        | VARCHAR(255) | NOT NULL                                                  |      |
| representative_name | VARCHAR(255) | NOT NULL                                                  |      |
| phone_number        | VARCHAR(20)  | NOT NULL                                                  |      |
| postal_code         | VARCHAR(10)  | NOT NULL                                                  |      |
| address             | VARCHAR(255) | NOT NULL                                                  |      |

#### client_bank_account

| 列名             | データ型         | 制約                                                      | 備考   |
|----------------|--------------|---------------------------------------------------------|------|
| id             | CHAR(26)     | PRIMARY KEY                                             | ULID |
| client_id      | CHAR(26)     | NOT NULL, FOREIGN KEY (client_id) REFERENCES client(id) |      |
| bank_name      | VARCHAR(255) | NOT NULL                                                |      |
| branch_name    | VARCHAR(255) | NOT NULL                                                |      |
| account_number | VARCHAR(20)  | NOT NULL                                                |      |
| account_name   | VARCHAR(255) | NOT NULL                                                |      |

#### invoice

| 列名               | データ型                                               | 制約                                                        | 備考   |
|------------------|----------------------------------------------------|-----------------------------------------------------------| ------|
| id               | CHAR(26)                                           | PRIMARY KEY                                               | ULID |
| company_id       | CHAR(26)                                           | NOT NULL, FOREIGN KEY (company_id) REFERENCES company(id) |     |
| client_id        | CHAR(26)                                           | NOT NULL, FOREIGN KEY (client_id) REFERENCES client(id)   |    |
| issue_date       | DATE                                               | NOT NULL                                                  |    |
| payment_amount   | DECIMAL(10, 2)                                     | NOT NULL                                                  |   |
| fee              | DECIMAL(10, 2)                                     |                                                           |  |
| fee_rate         | DECIMAL(5, 2)                                      |                                                           | |
| tax              | DECIMAL(10, 2)                                     |                                                           | |
| tax_rate         | DECIMAL(5, 2)                                      |                                                           | |
| invoice_amount   | DECIMAL(10, 2)                                     | NOT NULL                                                  | |
| payment_due_date | DATE                                               | NOT NULL, MUL                                             | |
| status           | ENUM('unprocessed', 'processing', 'paid', 'error') | NOT NULL                                                  | |

### Firebase Authentication

- ユーザーの認証に Firebase Authentication を使用
- メールアドレスとパスワードでの認証を使用

# Alternative Solution

特になし

# Milestones

8 時間での完成を目指すので特になし

# Concerns

- 手数料や消費税の端数は切り捨てか？
- 請求額の金額の上限はあるか？

# Logs

- JSON 形式で出力する
- 標準出力する
- アクセスログを出力する
- 500 エラー時にエラーレベルのログを出力する
- アプリケーション起動失敗時にエラーレベルのログを出力する
- アプリケーションの規模が小さいのでスタックトレースはひとまず出力しない

# Security

- リクエストに含まれるパラメータをログ出力する場合、 ID を除いた個人情報をマスクする
- SQL インジェクション対策としてプレースホルダーを使用して SQL を組み立てる
- パスワードは Firebase Authentication にて保管する

Firebase Authentication では独自のハッシュ関数を使用してパスワードをハッシュ化して保存する
ref: https://firebaseopensource.com/projects/firebase/scrypt/

# Observability

監視システムの導入はスコープ外なので特になし

# Glossary

```
企業: company
法人名: company_name
代表者名: representative_name
電話番号: phone_number
郵便番号: postal_code
住所: address
ユーザー: user
企業ID: company_id
氏名: full_name
メールアドレス: email_address
パスワード: password
取引先: client
取引先銀行口座: client_bank_account
銀行名: bank_name
支店名: branch_name
口座番号: account_number
口座名: account_name
請求書: invoice
発行日: issue_date
支払金額: payment_amount
手数料: fee
手数料率: fee_rate
消費税: tax
消費税率: tax_rate
請求金額: invoice_amount
支払期日: payment_due_date
ステータス: status
```

---
