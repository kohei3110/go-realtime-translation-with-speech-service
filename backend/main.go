package main

import (
	"context"
	"log"

	translatortext "go-realtime-translation-with-speech-service/backend/translatortext"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	// 1. 認証情報の取得
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("認証情報の取得に失敗しました: %v", err)
	}

	// 2. Speech Serviceのエンドポイント設定
	endpoint := "https://api.cognitive.microsofttranslator.com/"

	// 3. NewTranslatorClientの作成
	client, err := translatortext.NewTranslatorClient(endpoint, cred, nil)
	if err != nil {
		log.Fatalf("NewTranslatorClientの作成に失敗しました: %v", err)
	}

	// 4. 翻訳リクエストの作成
	text := "Hello, how are you?"
	textParam := []*translatortext.TranslateTextInput{
		{
			Text: &text,
		},
	}

	// 5. 翻訳の実行
	ctx := context.Background()
	// リクエストのデバッグ
	log.Printf("ctx: %s", ctx)

	result, err := client.Translate(ctx, []string{"ja"}, textParam, nil)
	if err != nil {
		log.Fatalf("翻訳の実行に失敗しました: %v", err)
	}

	// デバッグ用に構造を出力
	log.Printf("Translation response structure: %+v", result)

	// 翻訳結果を見やすく出力
	log.Printf("=== 翻訳結果 ===")
	if result.TranslateResultAllItemArray != nil {
		for i, item := range result.TranslateResultAllItemArray {
			log.Printf("翻訳 %d:", i+1)

			// 検出された言語情報
			if item.DetectedLanguage != nil {
				log.Printf("検出された言語: %s (信頼度: %.2f)",
					*item.DetectedLanguage.Language,
					item.DetectedLanguage.Score)
			}

			// 翻訳結果
			if item.Translations != nil {
				for j, translation := range item.Translations {
					log.Printf("  翻訳 %d.%d:", i+1, j+1)
					log.Printf("  対象言語: %s", translation.To)
					log.Printf("  翻訳文: %s", *translation.Text)

					// transliterationがある場合
					if translation.Transliteration != nil {
						log.Printf("    音訳: %s (スクリプト: %s)",
							*translation.Transliteration.Text,
							translation.Transliteration.Script)
					}
				}
			}

			// 元のテキスト（ソースの再表示）
			if i < len(textParam) && textParam[i].Text != nil {
				log.Printf("  元のテキスト: %s", *textParam[i].Text)
			}

			log.Printf("---")
		}
	} else {
		log.Printf("翻訳結果がありません")
	}
}
