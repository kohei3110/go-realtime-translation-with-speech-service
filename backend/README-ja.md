# リアルタイム音声翻訳サービス バックエンド

## 概要

このバックエンドサービスは、Azure Translator Serviceを利用して、テキスト翻訳およびリアルタイム音声ストリーミング翻訳機能を提供するRESTful APIです。Goで実装され、Ginフレームワークを使用しています。

## システムアーキテクチャ

システムは以下のコンポーネントで構成されています：

- **Gin Webサーバー**: HTTPリクエストを処理し、各種エンドポイントを提供
- **Azure Translator クライアント**: Azure Translator Text APIと通信するためのクライアント
- **セッション管理**: ストリーミング翻訳セッションを管理するためのインメモリストレージ
- **音声処理**: Base64エンコードされた音声データを処理するためのモジュール

```
+----------------+        +-------------------+
|                |        |                   |
| クライアント     +------->+ Gin Webサーバー   |
|                |        |                   |
+----------------+        +--------+----------+
                                  |
                                  v
                          +----------------+         +------------------+
                          |                |         |                  |
                          | TranslatorClient+-------->+ Azure Translator |
                          |                |         |                  |
                          +----------------+         +------------------+
```

## API エンドポイント

### ヘルスチェック

```
GET /api/v1/health
```

サーバーの状態を確認するためのエンドポイント。

**レスポンス例**:
```json
{
  "status": "ok"
}
```

### テキスト翻訳

```
POST /api/v1/translate
```

テキストを指定した言語に翻訳します。

**リクエスト例**:
```json
{
  "text": "こんにちは",
  "targetLanguage": "en",
  "sourceLanguage": "ja"
}
```

**レスポンス例**:
```json
{
  "originalText": "こんにちは",
  "translatedText": "Hello",
  "sourceLanguage": "ja",
  "targetLanguage": "en",
  "confidence": 0.98
}
```

### ストリーミング翻訳セッション開始

```
POST /api/v1/streaming/start
```

ストリーミング翻訳セッションを開始します。

**リクエスト例**:
```json
{
  "sourceLanguage": "ja",
  "targetLanguage": "en",
  "audioFormat": "wav"
}
```

**レスポンス例**:
```json
{
  "sessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

### 音声データ処理

```
POST /api/v1/streaming/process
```

Base64エンコードされた音声チャンクを送信して処理します。

**リクエスト例**:
```json
{
  "sessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "audioChunk": "UklGRjoAAABXQVZFZm10IBIAAAAHAAEAQB8AAEAfAAABAAgAAABMSVNUHAAAAElORk9JU0ZUDQAAAExhdmY1OC4yOS4xMDDA/w=="
}
```

**レスポンス例**:
```json
[
  {
    "sourceLanguage": "ja",
    "targetLanguage": "en",
    "translatedText": "Hello, how are you?",
    "isFinal": true,
    "segmentId": "f7e8d9c0-b1a2-3456-7890-abcdef123456"
  }
]
```

### ストリーミングセッション終了

```
POST /api/v1/streaming/close
```

ストリーミングセッションを終了します。

**リクエスト例**:
```json
{
  "sessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**レスポンス例**:
```json
{
  "status": "セッションを終了しました"
}
```

## 音声データ要件

- サポートされているフォーマット: WAV
- サンプリングレート: 16kHz推奨
- ビット深度: 16bit
- チャンネル: モノラル
- Base64エンコード: 音声データはBase64エンコードして送信する必要があります

## サービスプリンシパルの作成

Azure CLIを使用してサービスプリンシパルを作成します。

```bash
az ad sp create-for-rbac --name "go-translation-service" --role contributor --scopes /subscriptions/{subscription-id}/resourceGroups/{resource-group}
```

コマンド実行後、以下の情報が表示されます：
- appId (AZURE_CLIENT_ID)
- password (AZURE_CLIENT_SECRET)
- tenant (AZURE_TENANT_ID)

## 権限を設定

- 簡単のため、リソースグループスコープで `共同作成者` を付与。
- 本番環境では、最小権限の原則に従い、必要な権限のみを付与することをお勧めします。

## 環境変数の設定

- `.env.example` ファイルをコピーし、`.env` ファイルを作成。

```bash
cp .env.example .env
```

- `.env` ファイルに以下の環境変数を設定。

| 環境変数 | 説明 |
|----------|------|
| AZURE_CLIENT_ID | サービスプリンシパルのクライアントID |
| AZURE_CLIENT_SECRET | サービスプリンシパルのシークレット |
| AZURE_TENANT_ID | Entra IDのテナントID |
| TRANSLATOR_SUBSCRIPTION_KEY | Azure Translator リソースのサブスクリプションキー |
| TRANSLATOR_SUBSCRIPTION_REGION | Azure Translator リソースのリージョン（例: japaneast） |
| PORT | サーバーが使用するポート（デフォルト: 8080） |

## ローカル開発

### 必要条件

- Go 1.16以上
- Azure サブスクリプション
- Azure Translator リソース

### ローカル実行

```bash
go run main.go
```

## Dockerでの実行

```bash
# Dockerイメージをビルド
docker build -t go-translation-service .

# コンテナを実行
docker run --env-file .env -p 8080:8080 go-translation-service
```

## エラーハンドリング

サービスは以下のHTTPステータスコードを返します：

- 200 OK: リクエストが成功
- 400 Bad Request: リクエストパラメータが無効
- 401 Unauthorized: 認証に失敗
- 404 Not Found: リソースが見つからない
- 500 Internal Server Error: サーバー内部エラー

## パフォーマンスに関する考慮事項

- ストリーミングセッションはインメモリで管理されるため、サーバーの再起動時にすべてのセッションが失われます
- 大規模な環境では、Redisなどの外部キャッシュを使用してセッション状態を保存することを検討してください
- 長時間のアイドル状態のセッションを自動的に削除するタイムアウトメカニズムの実装を検討してください

## サポートされている言語

サポートされている言語のリストは、Azure Translator Serviceのドキュメントを参照してください。現在、100以上の言語がサポートされています。