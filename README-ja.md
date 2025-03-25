# リアルタイム翻訳サービス

Azure Speech ServiceとAzure Translatorを使用したリアルタイム翻訳サービスです。

## AutoRest を使った Go Client Library の生成

[AutoRest](https://github.com/Azure/autorest)に API 仕様書を読み込ませ、任意の言語でクライアントライブラリを生成させるアプローチを採用。
今回は、Translator API 仕様書を読み込ませ、Go のクライアントライブラリを生成。

```bash
autorest --go --input-file=https://raw.githubusercontent.com/Azure/azure-rest-api-specs/refs/heads/master/specification/cognitiveservices/data-plane/TranslatorText/stable/v3.0/TranslatorText.json --output-folder=./translatortext --namespace=translatortext
```

生成されたモジュールをプロジェクトに追加するだけで、Azure リソースに対する API リクエストを生成可能。

- [Translator API Spec v3.0](https://learn.microsoft.com/en-us/azure/ai-services/translator/text-translation/reference/v3/reference)

### Issues

- 翻訳結果のレスポンス受け取り時に、以下のエラーが発生する。対処法は以下のとおり。

```
unmarshalling type *[]*translatortext.TranslateResultAllItem: unmarshalling type *translatortext.TranslateResultAllItem: struct field DetectedLanguage: unmarshalling type *translatortext.TranslateResultAllItemDetectedLanguage: struct field Score: json: cannot unmarshal number 1.0 into Go value of type int32
→
// 以下のような構造体を探します
type TranslateResultAllItemDetectedLanguage struct {
    Language string `json:"language,omitempty"`
    Score    int32  `json:"score,omitempty"` // このフィールドの型が問題
}
 
// 以下のように変更します
type TranslateResultAllItemDetectedLanguage struct {
    Language string  `json:"language,omitempty"`
    Score    float64 `json:"score,omitempty"` // int32からfloat64に変更
}
```

## セットアップ手順

- [バックエンド API 起動手順書](./backend/README-ja.md)

サーバーが正常に起動すると、以下のようなメッセージが表示されます：
```
Speech Recognition and Translation Server is running on port 8080
```

- [フロントエンド起動手順書](./frontend/README-ja.md)

## 環境変数

アプリケーションは以下の環境変数を必要とします：

### バックエンド：
- `AZURE_CLIENT_ID` - サービスプリンシパルのクライアントID
- `AZURE_CLIENT_SECRET` - サービスプリンシパルのシークレット
- `AZURE_TENANT_ID` - Entra IDのテナントID
- `SPEECH_SERVICE_KEY` - Azure Speech Serviceのサブスクリプションキー
- `SPEECH_SERVICE_REGION` - Azure Speech Serviceのリージョン（例：japaneast）
- `PORT` - サーバーが使用するポート（デフォルト：8080）

## APIの終了方法

サーバーを終了するには、ターミナルで `Ctrl+C` を押してください。グレースフルシャットダウンが実行されます。

## curlを使用したAPI利用例

APIとは以下のcurlコマンドを使って対話できます：

### テキスト翻訳

テキストをある言語から別の言語に翻訳する：

```bash
curl -X POST http://localhost:8080/api/v1/translate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, how are you?",
    "sourceLanguage": "en",
    "targetLanguage": "ja"
  }'
```

### ストリーミング翻訳

#### 1. ストリーミングセッションの開始

```bash
curl -X POST http://localhost:8080/api/v1/streaming/start \
  -H "Content-Type: application/json" \
  -d '{
    "sourceLanguage": "en",
    "targetLanguage": "ja",
    "audioFormat": "wav"
  }'
```

レスポンスには、後続のリクエストに必要な `sessionId` と `webSocketURL` が含まれます：

```json
{
  "sessionId": "12345678-1234-1234-1234-123456789abc",
  "webSocketURL": "/api/v1/streaming/ws/12345678-1234-1234-1234-123456789abc",
  "sourceLanguage": "en",
  "targetLanguage": "ja"
}
```

#### 2. リアルタイム翻訳のためのWebSocket接続

リアルタイム翻訳のためには、startエンドポイントのレスポンスで提供されたWebSocket URLに接続します。
接続後は、バックエンドのREADMEに記載されているWebSocketプロトコルに従ってください。

#### 3. オーディオチャンクの処理（WebSocketの代替手段）

```bash
curl -X POST http://localhost:8080/api/v1/streaming/process \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "12345678-1234-1234-1234-123456789abc",
    "audioChunk": "BASE64でエンコードされたオーディオデータ"
  }'
```

#### 4. ストリーミングセッションの終了

```bash
curl -X POST http://localhost:8080/api/v1/streaming/close \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "12345678-1234-1234-1234-123456789abc"
  }'
```

### ヘルスチェック

APIサーバーが実行中かどうかを確認する：

```bash
curl http://localhost:8080/api/v1/health
```

期待されるレスポンス：

```json
{
  "status": "ok"
}
```