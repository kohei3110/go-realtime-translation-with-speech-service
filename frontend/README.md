# Real-time Translation Application Frontend

This is the frontend of a web application that performs real-time voice translation. It is built using React, TypeScript, and Vite.

## Features

- Text translation
- Real-time voice translation
- Modern UI based on Material-UI

## Requirements

- Node.js (v18 or higher)
- npm or yarn

## Setup Instructions

1. Install dependencies:
```bash
npm install
# or
yarn
```

2. Start development server:
```bash
npm run dev
# or
yarn dev
```

The application will start at http://localhost:5173 by default.

## Available Scripts

- `npm run dev`: Start development server
- `npm run build`: Build application for production
- `npm run preview`: Preview built application
- `npm run lint`: Validate code with ESLint

## Development

This application uses the following key technologies:

- React v19
- TypeScript
- Material-UI v6
- Vite v6
- Axios (API client)

For backend API communication, it uses `http://localhost:8080` by default.

## Production Deployment

Generate production files by running the build:

```bash
npm run build
# or
yarn build
```

Built files will be generated in the `dist` directory.

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
