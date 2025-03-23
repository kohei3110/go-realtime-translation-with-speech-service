# リアルタイム翻訳サービス

Azure Speech Serviceを使用したリアルタイム翻訳サービスです。

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

1. リポジトリのクローン
```bash
git clone [repository-url]
cd go-realtime-translation-with-speech-service
```

1. バックエンドディレクトリに移動
```bash
cd backend
```

2. APIサーバーの起動
```bash
go run main.go
```

サーバーが正常に起動すると、以下のようなメッセージが表示されます：
```
Server is running on port 8080
```

## APIの終了方法

サーバーを終了するには、ターミナルで `Ctrl+C` を押してください。グレースフルシャットダウンが実行されます。

## Web APIの使い方

APIは以下のエンドポイントを提供しており、REST APIとストリーミング翻訳の両方の機能があります：

### ヘルスチェックAPI
- **エンドポイント**: `GET /api/v1/health`
- **説明**: APIサーバーが実行中かどうかを確認します
- **使用例**:
  ```bash
  curl http://localhost:8080/api/v1/health
  ```
- **レスポンス**:
  ```json
  {
    "status": "ok"
  }
  ```

### テキスト翻訳API
- **エンドポイント**: `POST /api/v1/translate`
- **説明**: テキストをある言語から別の言語に翻訳します
- **リクエストボディ**:
  ```json
  {
    "text": "こんにちは",
    "targetLanguage": "en", 
    "sourceLanguage": "ja"  // オプション：指定しない場合は自動検出されます
  }
  ```
- **使用例**:
  ```bash
  curl -X POST http://localhost:8080/api/v1/translate \
    -H "Content-Type: application/json" \
    -d '{"text": "こんにちは", "targetLanguage": "en"}'
  ```
- **レスポンス**:
  ```json
  {
    "originalText": "こんにちは",
    "translatedText": "Hello",
    "sourceLanguage": "ja",
    "targetLanguage": "en",
    "confidence": 0.98
  }
  ```

### ストリーミング翻訳API

#### 1. ストリーミングセッション開始
- **エンドポイント**: `POST /api/v1/streaming/start`
- **説明**: 新しいストリーミング翻訳セッションを開始します
- **リクエストボディ**:
  ```json
  {
    "sourceLanguage": "ja",
    "targetLanguage": "en",
    "audioFormat": "wav"
  }
  ```
- **使用例**:
  ```bash
  curl -X POST http://localhost:8080/api/v1/streaming/start \
    -H "Content-Type: application/json" \
    -d '{"sourceLanguage": "ja", "targetLanguage": "en", "audioFormat": "wav"}'
  ```
- **レスポンス**:
  ```json
  {
    "sessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }
  ```

#### 2. 音声チャンク処理
- **エンドポイント**: `POST /api/v1/streaming/process`
- **説明**: アクティブなセッションで翻訳のための音声チャンクを処理します
- **リクエストボディ**:
  ```json
  {
    "sessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "audioChunk": "base64エンコードされた音声データ"
  }
  ```
- **使用例**:
  ```bash
  curl -X POST http://localhost:8080/api/v1/streaming/process \
    -H "Content-Type: application/json" \
    -d '{"sessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890", "audioChunk": "base64エンコードされた音声データ"}'
  ```
- **レスポンス**:
  ```json
  [
    {
      "sourceLanguage": "ja",
      "targetLanguage": "en",
      "translatedText": "こんにちは、お元気ですか？",
      "isFinal": false,
      "segmentId": "segment-123"
    }
  ]
  ```

#### 3. ストリーミングセッション終了
- **エンドポイント**: `POST /api/v1/streaming/close`
- **説明**: アクティブなストリーミング翻訳セッションを終了します
- **リクエストボディ**:
  ```json
  {
    "sessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }
  ```
- **使用例**:
  ```bash
  curl -X POST http://localhost:8080/api/v1/streaming/close \
    -H "Content-Type: application/json" \
    -d '{"sessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"}'
  ```
- **レスポンス**:
  ```json
  {
    "status": "セッションを終了しました"
  }
  ```

### エラーレスポンス
すべてのAPIエンドポイントは適切なHTTPステータスコードを返します：
- `400 Bad Request`: 無効な入力パラメータ
- `401 Unauthorized`: 認証失敗
- `404 Not Found`: リソースが見つからない
- `500 Internal Server Error`: サーバー側のエラー

エラーレスポンスはJSON形式でフォーマットされます：
```json
{
  "error": "エラーメッセージの詳細"
}
```

## 仕様書

### バックエンド仕様

#### 技術スタック
- 言語: Go 1.19+
- フレームワーク: 標準ライブラリ + Gorilla WebSocket
- 外部サービス: Azure Speech Service
- インフラ: Docker

#### API エンドポイント

1. **WebSocket接続エンドポイント**
   - パス: `/ws`
   - メソッド: GET (WebSocket Upgrade)
   - 機能: 音声ストリームの送受信およびリアルタイム翻訳のための双方向通信

2. **ヘルスチェックエンドポイント**
   - パス: `/health`
   - メソッド: GET
   - レスポンス: `{"status": "ok"}`
   - 機能: サービスの稼働状態確認

#### WebSocketメッセージ形式

**クライアントからサーバーへ:**
```json
{
  "type": "start_translation",
  "sourceLanguage": "ja-JP",
  "targetLanguage": "en-US",
  "audioFormat": "audio/wav"
}
```

```json
{
  "type": "audio_data",
  "data": "base64エンコードされた音声データ"
}
```

```json
{
  "type": "stop_translation"
}
```

**サーバーからクライアントへ:**
```json
{
  "type": "translation_result",
  "sourceText": "こんにちは",
  "translatedText": "Hello",
  "isFinal": true
}
```

```json
{
  "type": "error",
  "message": "エラーメッセージ"
}
```

#### エラーハンドリング
- 全てのエラーはログに記録
- クライアントにはJSON形式でエラーメッセージを返却
- 接続エラーが発生した場合は自動的に再接続を試行

#### パフォーマンス要件
- 最大同時接続数: 100
- レイテンシ: 音声入力から翻訳結果表示まで1秒以内
- CPU使用率: 平均60%以下
- メモリ使用量: 最大512MB

### フロントエンド仕様

#### 技術スタック
- 言語: TypeScript
- フレームワーク: React
- スタイリング: CSS Modules または Tailwind CSS
- ビルドツール: Vite

#### 機能要件

1. **ユーザーインターフェース**
   - シンプルで直感的なUI
   - レスポンシブデザイン（モバイル、タブレット、デスクトップ対応）
   - ダークモード/ライトモード切り替え

2. **音声入力**
   - マイク音声の録音と送信
   - 音声レベルインジケーター表示
   - 無音検出による自動一時停止

3. **翻訳表示**
   - 元の言語テキストと翻訳テキストの同時表示
   - 翻訳履歴の保存と表示
   - テキストのコピー機能

4. **設定**
   - 言語ペアの選択（源言語と目標言語）
   - 音声入力感度の調整
   - フォントサイズ調整

5. **状態表示**
   - 接続状態インジケーター
   - エラーメッセージの表示
   - 音声認識状態の表示

#### 非機能要件
- 初期読み込み時間: 2秒以内
- オフライン機能: 基本UIの表示とエラーメッセージ
- アクセシビリティ: WCAG 2.1 AAレベル準拠
- モバイルデバイスの電池消費最適化

#### ユーザーフロー
1. アプリケーションにアクセス
2. 言語ペアを選択
3. マイクへのアクセス許可を付与
4. 開始ボタンをクリック
5. 話し始める
6. リアルタイムで翻訳結果を確認
7. 必要に応じて停止/再開
8. 翻訳履歴を確認またはエクスポート

#### デザイン要件
- モダンでクリーンなインターフェース
- 視覚的フィードバックの提供
- 色のコントラスト比: 4.5:1以上
- アイコンと操作ボタンのサイズ: 最小44px×44px（タッチデバイス用）