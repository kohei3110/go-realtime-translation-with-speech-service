package main

import (
	"context"
	"log"

	speechtotext "go-realtime-translation-with-speech-service/backend/speechtotext"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	// 1. 認証情報の取得
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("認証情報の取得に失敗しました: %v", err)
	}

	// 2. Speech Serviceのエンドポイント設定
	endpoint := "https://eastus.api.cognitive.microsoft.com/"

	// 3. TranscriptionsClientの作成
	client, err := speechtotext.NewTranscriptionsClient(endpoint, cred, nil)
	if err != nil {
		log.Fatalf("TranscriptionsClientの作成に失敗しました: %v", err)
	}

	// 4. 翻訳リクエストの作成
	displayName := "task1"
	locale := "ja-JP"
	audioURL := "https://stspeechsdkwithgodemoeu1.blob.core.windows.net/samplewav/AESOP_JPN.wav"
	transcription := speechtotext.Transcription{
		DisplayName: &displayName,
		Locale:      &locale,
		ContentUrls: []*string{&audioURL},
	}

	// 5. 翻訳の実行
	ctx := context.Background()
	result, err := client.Create(ctx, transcription, nil)
	if err != nil {
		log.Fatalf("翻訳の実行に失敗しました: %v", err)
	}

	log.Printf("翻訳タスクが作成されました: %v", *result.Self)

	// 6. 状態の確認
	taskID := *result.Self
	status, err := client.Get(ctx, taskID, nil)
	if err != nil {
		log.Fatalf("状態の取得に失敗しました: %v", err)
	}

	log.Printf("翻訳タスクの状態: %v", *status.Status)
}
