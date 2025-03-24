// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kohei3110/gospeech"
)

func main() {
	// Create a translation config with your Azure Speech service subscription
	config, err := gospeech.SpeechTranslationConfigFromSubscription(
		os.Getenv("SPEECH_KEY"),
		os.Getenv("SPEECH_REGION"),
	)
	if err != nil {
		fmt.Printf("Error creating config: %v\n", err)
		return
	}

	// Set speech recognition language
	config.SetSpeechRecognitionLanguage("en-US")

	// Add translation target languages
	config.AddTargetLanguage("ja")
	config.AddTargetLanguage("es")
	config.AddTargetLanguage("fr")

	// Create an audio config from the default microphone
	audioConfig, err := gospeech.NewAudioConfigFromDefaultMicrophone()
	if err != nil {
		fmt.Printf("Error creating audio config: %v\n", err)
		return
	}

	// Create a translation recognizer
	recognizer, err := gospeech.NewTranslationRecognizer(config, audioConfig)
	if err != nil {
		fmt.Printf("Error creating recognizer: %v\n", err)
		return
	}
	defer recognizer.Close()

	// Set up event handlers
	recognizer.Recognizing().Connect(func(e interface{}) {
		if args, ok := e.(*gospeech.TranslationRecognitionEventArgs); ok {
			fmt.Printf("RECOGNIZING: %s\n", args.Result.Text)
			for lang, translation := range args.Result.Translations {
				fmt.Printf("  %s: %s\n", lang, translation)
			}
		}
	})

	recognizer.Recognized().Connect(func(e interface{}) {
		if args, ok := e.(*gospeech.TranslationRecognitionEventArgs); ok {
			fmt.Printf("RECOGNIZED: %s\n", args.Result.Text)
			for lang, translation := range args.Result.Translations {
				fmt.Printf("  %s: %s\n", lang, translation)
			}
		}
	})

	recognizer.Canceled().Connect(func(e interface{}) {
		if args, ok := e.(*gospeech.TranslationRecognitionCanceledEventArgs); ok {
			fmt.Printf("CANCELED: Reason=%v\n", args.CancellationDetails.Reason)
			if args.CancellationDetails.Reason == gospeech.CancellationReasonError {
				fmt.Printf("CANCELED: ErrorCode=%v\n", args.CancellationDetails.ErrorCode)
				fmt.Printf("CANCELED: ErrorDetails=%s\n", args.CancellationDetails.ErrorDetails)
			}
		}
	})

	recognizer.SessionStarted().Connect(func(e interface{}) {
		if args, ok := e.(*gospeech.SessionEventArgs); ok {
			fmt.Printf("SESSION STARTED: SessionId=%s\n", args.SessionID)
		}
	})

	recognizer.SessionStopped().Connect(func(e interface{}) {
		if args, ok := e.(*gospeech.SessionEventArgs); ok {
			fmt.Printf("SESSION STOPPED: SessionId=%s\n", args.SessionID)
		}
	})

	// Start continuous recognition
	ctx := context.Background()
	err = recognizer.StartContinuousRecognition(ctx)
	if err != nil {
		fmt.Printf("Error starting continuous recognition: %v\n", err)
		return
	}

	// Wait for Ctrl+C to stop
	fmt.Println("Speak into your microphone. Press Ctrl+C to stop.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Stop recognition
	err = recognizer.StopContinuousRecognition()
	if err != nil {
		fmt.Printf("Error stopping continuous recognition: %v\n", err)
		return
	}

	// Wait a moment for all events to be processed
	time.Sleep(2 * time.Second)
}
