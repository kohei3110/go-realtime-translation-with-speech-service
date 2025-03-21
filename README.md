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