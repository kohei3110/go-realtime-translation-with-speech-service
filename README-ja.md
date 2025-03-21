# リアルタイム翻訳サービス

Azure Speech Serviceを使用したリアルタイム翻訳サービスのバックエンドAPIです。

## 必要要件

- Go 1.19以上
- Azure Speech Serviceのアカウント
- C++ コンパイラ (gcc または clang)
- Azure Speech SDK for C/C++

### macOSでの依存関係のインストール

```bash
# Homebrewを使用してC++コンパイラをインストール
brew install gcc

# Azure Speech SDK for C/C++のインストール
curl -L https://aka.ms/csspeech/macosbinary -o speechsdk.tar.gz
tar -xzf speechsdk.tar.gz
sudo mkdir -p /usr/local/include
sudo cp SpeechSDK-macOS/include/* /usr/local/include/
sudo cp SpeechSDK-macOS/lib/libMicrosoft.CognitiveServices.Speech.core.dylib /usr/local/lib/
rm -rf speechsdk.tar.gz SpeechSDK-macOS
```

## 環境変数の設定

以下の環境変数を設定してください：

```bash
export PORT=8080  # APIサーバーのポート（省略可、デフォルト: 8080）
export AZURE_SPEECH_KEY=your_key_here  # Azure Speech Serviceのキー
export AZURE_SPEECH_REGION=your_region_here  # Azure Speech Serviceのリージョン

# Azure Speech SDKのライブラリパスを設定
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -lMicrosoft.CognitiveServices.Speech.core"
export DYLD_LIBRARY_PATH="/usr/local/lib:$DYLD_LIBRARY_PATH"
```

## セットアップ手順

1. リポジトリのクローン
```bash
git clone [repository-url]
cd go-realtime-translation-with-speech-service
```

2. 依存関係のインストール
```bash
cd backend
go mod tidy
```

## Dockerでの実行方法

1. .envファイルの作成
```bash
AZURE_SPEECH_KEY=your_key_here
AZURE_SPEECH_REGION=your_region_here
```

2. Dockerコンテナのビルドと起動
```bash
docker compose up --build
```

コンテナを停止するには以下のコマンドを実行してください：
```bash
docker compose down
```

## APIの起動方法（ローカル環境）

1. バックエンドディレクトリに移動
```bash
cd backend  # もし既にbackendディレクトリにいない場合
```

2. APIサーバーの起動
```bash
go run cmd/api/main.go
```

サーバーが正常に起動すると、以下のようなメッセージが表示されます：
```
Server is running on port 8080
```

## APIの終了方法

サーバーを終了するには、ターミナルで `Ctrl+C` を押してください。グレースフルシャットダウンが実行されます。

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