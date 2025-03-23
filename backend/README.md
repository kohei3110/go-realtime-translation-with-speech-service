## サービスプリンシパルの作成

## 権限を設定

- 簡単のため、リソースグループスコープで `共同作成者` を付与。

## 環境変数の設定

- `.env.example` ファイルをコピーし、`.env` ファイルを作成。

```
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

```
docker run --env-file .env -p 8080:8080 your-image-name
```