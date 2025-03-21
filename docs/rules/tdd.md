# テスト駆動開発 (TDD) ガイドライン

## 1. TDDの基本サイクル

リアルタイム翻訳 API 開発プロジェクトでは、以下のTDDサイクルを採用します：

### Red-Green-Refactor サイクル

1. **Red**: 失敗するテストを書く
- 実装したい機能を明確にする
- テストは実装前に記述し、必ず失敗することを確認する
- テストは機能要件を明確に表現したものであること

2. **Green**: 最小限のコードで成功させる
- テストが通るための最小限（必要十分）のコードを実装する
- この段階ではパフォーマンスやエレガントさより機能の正しさを優先する

3. **Refactor**: リファクタリングする
- コードの品質を高めるための改善を行う
- テストは引き続き成功する状態を維持する
- コードの重複を排除し、可読性を向上させる

## 2. テストの種類と役割

### 2.1 ユニットテスト

- 対象: 関数、メソッド、小さなクラス
- 目的: コードの最小単位の正確性を検証
- 特徴: 高速に実行可能、外部依存をモック化
- ツール: Testify、GoMock、httptest
- 場所: 各機能モジュール内の `tests` ディレクトリ

```go
package speech

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestJapaneseToEnglishTranslation(t *testing.T) {
    // 準備
    translator := NewTranslator()
    japaneseText := "こんにちは、元気ですか？"
    
    // 実行
    result, err := translator.Translate(japaneseText, "ja", "en")
    
    // 検証 - require は失敗時にテスト終了
    require.NoError(t, err, "翻訳中にエラーが発生しました")
    
    // assert は失敗を記録するが続行
    assert.Equal(t, "en", result.TargetLanguage, "ターゲット言語が一致しません")
    assert.Contains(t, result.TranslatedText, "Hello", "英語の挨拶が含まれていません")
    assert.Greater(t, result.ConfidenceScore, 0.8, "信頼度スコアが低すぎます")
}
```

### 2.2 統合テスト

- 対象: 複数のコンポーネントの相互作用
- 目的: コンポーネント間の連携が正しく機能することを確認
- 特徴: やや実行に時間がかかる、一部の外部依存を実際に使用
- 場所: `tests/integration` ディレクトリ

```go
package integration

import (
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/yourusername/yourproject/internal/file"
    "github.com/yourusername/yourproject/internal/audio"
    "github.com/yourusername/yourproject/internal/storage"
)

// TestFileUploadAndConversionFlow はファイル処理と音声変換の統合テスト
func TestFileUploadAndConversionFlow(t *testing.T) {
    t.Parallel() // 可能な場合は並列実行

    // テスト用の一時ディレクトリを作成
    tempDir, err := os.MkdirTemp("", "integration-test")
    require.NoError(t, err, "一時ディレクトリの作成に失敗しました")
    defer os.RemoveAll(tempDir) // テスト終了後にクリーンアップ

    // テスト用ファイルの準備
    testFilePath := filepath.Join("testdata", "sample.mp4")
    
    // 1. テスト用ファイルをアップロード
    fileService := file.NewFileProcessingService(tempDir)
    fileID, err := fileService.Upload(testFilePath)
    require.NoError(t, err, "ファイルのアップロードに失敗しました")
    
    // 2. 音声変換を実行
    converterService := audio.NewAudioConversionService()
    result, err := converterService.ProcessFile(fileID)
    require.NoError(t, err, "音声変換に失敗しました")
    
    // 3. 検証
    assert.Equal(t, "completed", result.Status, "処理状態が不正です")
    
    // ストレージサービスで結果ファイルを確認
    storageService := storage.NewStorageService()
    exists, err := storageService.Exists(filepath.Join(fileID, "converted.mp3"))
    require.NoError(t, err, "ファイル存在確認中にエラーが発生しました")
    assert.True(t, exists, "変換後のファイルが存在しません")
    
    // 必要に応じて変換されたファイルの内容を検証
    // 例: ファイルサイズやメタデータの確認
    fileInfo, err := storageService.GetFileInfo(filepath.Join(fileID, "converted.mp3"))
    require.NoError(t, err)
    assert.Greater(t, fileInfo.Size, int64(0), "ファイルサイズが0です")
}

// 複数コンポーネントを使用した別の統合テスト例
func TestRealTimeTranslationFlow(t *testing.T) {
    // 音声認識、翻訳、テキスト出力を統合的にテスト
    recognizer := speech.NewSpeechRecognizer()
    translator := translation.NewTranslator()
    outputService := output.NewOutputService()
    
    // テスト用音声ファイル
    audioData, err := os.ReadFile("testdata/japanese_sample.wav")
    require.NoError(t, err)
    
    // 音声認識
    recognizedText, err := recognizer.Recognize(audioData)
    require.NoError(t, err)
    assert.NotEmpty(t, recognizedText, "認識テキストが空です")
    
    // テキスト翻訳
    translatedText, err := translator.Translate(recognizedText, "ja", "en")
    require.NoError(t, err)
    assert.NotEmpty(t, translatedText, "翻訳テキストが空です")
    
    // 出力処理
    outputID, err := outputService.SaveTranslation(translatedText)
    require.NoError(t, err)
    
    // 保存された翻訳を確認
    savedTranslation, err := outputService.GetTranslation(outputID)
    require.NoError(t, err)
    assert.Equal(t, translatedText, savedTranslation, "保存された翻訳が一致しません")
}
```

### 2.3 エンドツーエンド (E2E) テスト

- 対象: システム全体の機能
- 目的: ユーザーの視点からのシナリオ検証
- 特徴: 実行に時間がかかる、実際のサービスと連携
- 場所: `tests/e2e` ディレクトリ

```go
package integration

import (
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/yourusername/yourproject/internal/file"
    "github.com/yourusername/yourproject/internal/audio"
    "github.com/yourusername/yourproject/internal/storage"
)

// TestFileUploadAndConversionFlow はファイル処理と音声変換の統合テスト
func TestFileUploadAndConversionFlow(t *testing.T) {
    t.Parallel() // 可能な場合は並列実行

    // テスト用の一時ディレクトリを作成
    tempDir, err := os.MkdirTemp("", "integration-test")
    require.NoError(t, err, "一時ディレクトリの作成に失敗しました")
    defer os.RemoveAll(tempDir) // テスト終了後にクリーンアップ

    // テスト用ファイルの準備
    testFilePath := filepath.Join("testdata", "sample.mp4")
    
    // 1. テスト用ファイルをアップロード
    fileService := file.NewFileProcessingService(tempDir)
    fileID, err := fileService.Upload(testFilePath)
    require.NoError(t, err, "ファイルのアップロードに失敗しました")
    
    // 2. 音声変換を実行
    converterService := audio.NewAudioConversionService()
    result, err := converterService.ProcessFile(fileID)
    require.NoError(t, err, "音声変換に失敗しました")
    
    // 3. 検証
    assert.Equal(t, "completed", result.Status, "処理状態が不正です")
    
    // ストレージサービスで結果ファイルを確認
    storageService := storage.NewStorageService()
    exists, err := storageService.Exists(filepath.Join(fileID, "converted.mp3"))
    require.NoError(t, err, "ファイル存在確認中にエラーが発生しました")
    assert.True(t, exists, "変換後のファイルが存在しません")
    
    // 必要に応じて変換されたファイルの内容を検証
    // 例: ファイルサイズやメタデータの確認
    fileInfo, err := storageService.GetFileInfo(filepath.Join(fileID, "converted.mp3"))
    require.NoError(t, err)
    assert.Greater(t, fileInfo.Size, int64(0), "ファイルサイズが0です")
}

// 複数コンポーネントを使用した別の統合テスト例
func TestRealTimeTranslationFlow(t *testing.T) {
    // 音声認識、翻訳、テキスト出力を統合的にテスト
    recognizer := speech.NewSpeechRecognizer()
    translator := translation.NewTranslator()
    outputService := output.NewOutputService()
    
    // テスト用音声ファイル
    audioData, err := os.ReadFile("testdata/japanese_sample.wav")
    require.NoError(t, err)
    
    // 音声認識
    recognizedText, err := recognizer.Recognize(audioData)
    require.NoError(t, err)
    assert.NotEmpty(t, recognizedText, "認識テキストが空です")
    
    // テキスト翻訳
    translatedText, err := translator.Translate(recognizedText, "ja", "en")
    require.NoError(t, err)
    assert.NotEmpty(t, translatedText, "翻訳テキストが空です")
    
    // 出力処理
    outputID, err := outputService.SaveTranslation(translatedText)
    require.NoError(t, err)
    
    // 保存された翻訳を確認
    savedTranslation, err := outputService.GetTranslation(outputID)
    require.NoError(t, err)
    assert.Equal(t, translatedText, savedTranslation, "保存された翻訳が一致しません")
}
```

## 3. テストカバレッジ方針

- **目標カバレッジ**: コードベース全体で80%以上
- **重要コンポーネント**: 核となるビジネスロジックは90%以上
- **測定**: pytestとCoverageを使用
- **レポート**: CI/CDパイプラインで自動生成

## 4. モックとスタブの使用ガイドライン

### 4.1 モックを使用するケース

- 外部サービス（Azure Speech等）の呼び出し
- ファイルシステムアクセス
- 時間依存の処理

### 4.2 モック実装例

```go
package speech

// SpeechService は音声認識サービスのインターフェース
type SpeechService interface {
    Transcribe(audioData []byte) (TranscriptionResult, error)
}

// TranscriptionResult は音声認識の結果
type TranscriptionResult struct {
    Text       string
    Speaker    string
    Confidence float64
}
```

### 4.3 テストダブル選択指針

- **Stub**: 単純な戻り値が必要な場合
- **Mock**: 呼び出し回数や引数の検証が必要な場合
- **Fake**: 軽量な代替実装が必要な場合（インメモリDBなど）
- **Spy**: 実際の処理を行いつつ呼び出し情報も記録する場合

## 5. テストデータ管理

- テストデータは `tests/fixtures` ディレクトリに格納
- 大容量ファイルはGitに含めず、CIパイプラインで取得
- 機密データはテスト用に匿名化したものを使用
- フィクスチャとファクトリパターンを活用

```go

```

## 6. テスト自動化

### 6.1 ローカル開発環境

- コミット前に自動実行するプリコミットフック
- 変更されたコードに関連するテストのみ実行する機能

### 6.2 CI/CD パイプライン

- プルリクエスト時に全テストを自動実行
- テストカバレッジレポートの自動生成
- 失敗テストの通知とレポート

### 6.3 実行コマンド

```bash
# ユニットテストのみ実行
go test ./internal/...

# 統合テストを実行
go test ./tests/integration/...
# または統合テストタグを使用
go test -tags=integration ./...

# E2Eテストを実行
go test ./tests/e2e/...
# またはE2Eテストタグを使用
go test -tags=e2e ./...

# カバレッジレポート付きで全テスト実行
go test -coverprofile=coverage.out ./...
# HTML形式でカバレッジレポートを表示
go tool cover -html=coverage.out -o coverage.html
# コンソールにカバレッジサマリーを表示
go tool cover -func=coverage.out
```

## 7. TDDのベストプラクティス

### 7.1 テストファーストの原則

- 必ず実装前にテストを記述する
- テストが意味のある失敗を示すことを確認してから実装に進む

### 7.2 FIRST原則

- **Fast**: テストは高速に実行できること
- **Independent**: テスト間に依存関係がないこと
- **Repeatable**: 何度実行しても同じ結果が得られること
- **Self-validating**: テストは自己検証可能であること
- **Timely**: テストは実装前に書くこと

### 7.3 テストの表現力

- テスト名は検証内容を明確に表現する
- Given-When-Then パターンで条件、操作、期待結果を明確に
- ヘルパー関数を使って複雑なセットアップを抽象化

```go
def test_when_large_file_uploaded_then_chunked_processing_used():
    # Given: 大容量ファイルの準備
    large_file = create_test_file(size_mb=500)
    
    # When: ファイル処理サービスに渡す
    processor = FileProcessingService()
    result = processor.process(large_file)
    
    # Then: チャンク処理が適用されていることを確認
    assert result.processing_strategy == "chunked"
    assert len(result.chunks) > 1
```

## 8. 特定のテスト戦略

### 8.1 非同期コードのテスト

```go
func TestAsyncFileProcessing(t *testing.T) {
    // コンテキストを作成（タイムアウト付き）
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    processor := NewAsyncFileProcessor()
    
    // 非同期処理の実行
    resultChan := make(chan ProcessResult)
    go func() {
        result, err := processor.ProcessFile(ctx, "test.mp4")
        require.NoError(t, err)
        resultChan <- result
    }()
    
    // 結果を待機
    select {
    case result := <-resultChan:
        assert.Equal(t, "completed", result.Status)
    case <-ctx.Done():
        t.Fatal("処理がタイムアウトしました")
    }
}
```

### 8.2 例外とエラー処理のテスト

```go
func TestInvalidFileFormatReturnsError(t *testing.T) {
    converter := NewAudioConverter()
    
    // Goではエラーを返り値として検証
    result, err := converter.Convert("test.txt", "mp3")
    
    // エラーが期待通り発生したか確認
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "unsupported format")
    
    // もしくはエラータイプの検証
    var invalidFormatErr *InvalidFormatError
    assert.True(t, errors.As(err, &invalidFormatErr))
}
```

### 8.3 分岐条件のテスト

```go
func TestProcessingStrategySelection(t *testing.T) {
    processor := NewFileProcessor()
    
    // テストケースを定義
    testCases := []struct {
        name           string
        fileSize       int64
        expectedStrategy string
    }{
        {
            name:             "小サイズファイルは単一処理",
            fileSize:         1024 * 1024 * 10, // 10MB
            expectedStrategy: "single",
        },
        {
            name:             "大サイズファイルはチャンク処理",
            fileSize:         1024 * 1024 * 500, // 500MB
            expectedStrategy: "chunked",
        },
    }
    
    // 各ケースを実行
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            file := createTestFile(tc.fileSize)
            result, err := processor.Process(file)
            
            assert.NoError(t, err)
            assert.Equal(t, tc.expectedStrategy, result.ProcessingStrategy)
            
            // チャンク処理の場合は追加検証
            if (tc.expectedStrategy == "chunked") {
                assert.Greater(t, len(result.Chunks), 1)
            }
        })
    }
}
```

## 9. コードレビューとテスト

コードレビューでは以下の点を重点的に確認します：

1. テストが機能要件を適切にカバーしているか
2. テスト自体の品質は適切か
3. エッジケースや例外パターンがテストされているか
4. テストの可読性と保守性は確保されているか

## 10. 継続的改善

- テスト戦略を定期的に見直し、改善する
- 重要な障害発生時には、該当するケースのテストを追加
- チーム内でテスト技術の共有と学習を促進