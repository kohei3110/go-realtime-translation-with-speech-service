# リアルタイム翻訳アプリケーション フロントエンド

音声のリアルタイム翻訳を行うWebアプリケーションのフロントエンドです。React、TypeScript、Viteを使用して構築されています。

## 機能

- テキスト翻訳
- リアルタイム音声翻訳
- Material-UIベースのモダンなUI

## 必要条件

- Node.js (v18以上)
- npm または yarn

## セットアップ手順

1. 依存パッケージのインストール:
```bash
npm install
# または
yarn
```

2. 開発サーバーの起動:
```bash
npm run dev
# または
yarn dev
```

アプリケーションは デフォルトで http://localhost:5173 で起動します。

## 利用可能なスクリプト

- `npm run dev`: 開発サーバーを起動します
- `npm run build`: プロダクション用にアプリケーションをビルドします
- `npm run preview`: ビルドしたアプリケーションをプレビューします
- `npm run lint`: ESLintでコードを検証します

## 開発

このアプリケーションは以下の主要な技術を使用しています：

- React v19
- TypeScript
- Material-UI v6
- Vite v6
- Axios (APIクライアント)

バックエンドAPIとの通信には、デフォルトで `http://localhost:8080` を使用します。

## 本番環境へのデプロイ

ビルドを実行してプロダクション用のファイルを生成：

```bash
npm run build
# または
yarn build
```

ビルドされたファイルは `dist` ディレクトリに生成されます。

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type-aware lint rules:

```js
export default tseslint.config({
  extends: [
    // Remove ...tseslint.configs.recommended and replace with this
    ...tseslint.configs.recommendedTypeChecked,
    // Alternatively, use this for stricter rules
    ...tseslint.configs.strictTypeChecked,
    // Optionally, add this for stylistic rules
    ...tseslint.configs.stylisticTypeChecked,
  ],
  languageOptions: {
    // other options...
    parserOptions: {
      project: ['./tsconfig.node.json', './tsconfig.app.json'],
      tsconfigRootDir: import.meta.dirname,
    },
  },
})
```

You can also install [eslint-plugin-react-x](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-x) and [eslint-plugin-react-dom](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-dom) for React-specific lint rules:

```js
// eslint.config.js
import reactX from 'eslint-plugin-react-x'
import reactDom from 'eslint-plugin-react-dom'

export default tseslint.config({
  plugins: {
    // Add the react-x and react-dom plugins
    'react-x': reactX,
    'react-dom': reactDom,
  },
  rules: {
    // other rules...
    // Enable its recommended typescript rules
    ...reactX.configs['recommended-typescript'].rules,
    ...reactDom.configs.recommended.rules,
  },
})
```