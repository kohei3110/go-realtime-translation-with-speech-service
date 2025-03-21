# 環境変数の設定手順

このドキュメントでは、アプリケーションの環境変数を安全に設定する手順について説明します。

## 1. 環境変数テンプレートのコピー

バックエンドディレクトリで`.env.example`を`.env`にコピーします：

```bash
cd backend
cp .env.example .env
```

## 2. 環境変数の設定

`.env`ファイルを編集して、以下の環境変数を設定します：

```plaintext
# Application Port
PORT=8080

# Azure Speech Service Credentials
AZURE_SPEECH_KEY=your_speech_key_here
AZURE_SPEECH_REGION=your_region_here
```

**重要**: `.env`ファイルは`.gitignore`に追加されており、リポジトリにコミットされません。

## 3. コンテナの起動方法

### 3.1 docker-composeを使用する場合

```bash
docker compose --env-file ./backend/.env up
```

### 3.2 Dockerを直接使用する場合

```bash
docker build -t realtime-translation \
  --build-arg PORT=8080 \
  --build-arg AZURE_SPEECH_KEY=your_key \
  --build-arg AZURE_SPEECH_REGION=your_region \
  ./backend
```

## 4. 環境変数の説明

| 環境変数 | 説明 | 必須 | デフォルト値 |
|----------|------|------|--------------|
| PORT | アプリケーションのポート番号 | はい | 8080 |
| AZURE_SPEECH_KEY | Azure Speech Serviceのキー | はい | なし |
| AZURE_SPEECH_REGION | Azure Speech Serviceのリージョン | はい | なし |

## 5. セキュリティに関する注意事項

1. `.env`ファイルを絶対にGitリポジトリにコミットしないでください
2. 本番環境の環境変数は、安全な方法で管理・配布してください
3. デフォルトの環境変数はDockerfile内で設定されていますが、実運用時は必ず上書きしてください

## 6. トラブルシューティング

環境変数が正しく設定されていない場合、以下の点を確認してください：

1. `.env`ファイルがbackendディレクトリに存在すること
2. `.env`ファイル内の変数が正しい形式で設定されていること
3. docker-composeまたはdocker buildコマンド実行時に正しいパスを指定していること