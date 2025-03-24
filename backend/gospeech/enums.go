// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license.

// Package enums defines the enumeration types used in the Speech SDK
package gospeech

import "fmt"

// PropertyID represents speech property identifiers
type PropertyID string

// PropertyID constants
const (
	// SpeechServiceConnection properties
	SpeechServiceConnectionKey                    PropertyID = "SpeechServiceConnection_Key"
	SpeechServiceConnectionEndpoint               PropertyID = "SpeechServiceConnection_Endpoint"
	SpeechServiceConnectionRegion                 PropertyID = "SpeechServiceConnection_Region"
	SpeechServiceAuthorizationToken               PropertyID = "SpeechServiceAuthorization_Token"
	SpeechServiceAuthorizationType                PropertyID = "SpeechServiceAuthorization_Type"
	SpeechServiceConnectionEndpointID             PropertyID = "SpeechServiceConnection_EndpointId"
	SpeechServiceConnectionHost                   PropertyID = "SpeechServiceConnection_Host"
	SpeechServiceConnectionProxyHostName          PropertyID = "SpeechServiceConnection_ProxyHostName"
	SpeechServiceConnectionProxyPort              PropertyID = "SpeechServiceConnection_ProxyPort"
	SpeechServiceConnectionProxyUserName          PropertyID = "SpeechServiceConnection_ProxyUserName"
	SpeechServiceConnectionProxyPassword          PropertyID = "SpeechServiceConnection_ProxyPassword"
	SpeechServiceConnectionURL                    PropertyID = "SpeechServiceConnection_Url"
	SpeechServiceConnectionProxyHostBypass        PropertyID = "SpeechServiceConnection_ProxyHostBypass"
	SpeechServiceConnectionTranslationToLanguages PropertyID = "SpeechServiceConnection_TranslationToLanguages"
	SpeechServiceConnectionTranslationVoice       PropertyID = "SpeechServiceConnection_TranslationVoice"
	SpeechServiceConnectionTranslationFeatures    PropertyID = "SpeechServiceConnection_TranslationFeatures"
	SpeechServiceConnectionIntentRegion           PropertyID = "SpeechServiceConnection_IntentRegion"
	SpeechServiceConnectionRecoMode               PropertyID = "SpeechServiceConnection_RecoMode"
	SpeechServiceConnectionRecoLanguage           PropertyID = "SpeechServiceConnection_RecoLanguage"
	SpeechSessionID                               PropertyID = "Speech_SessionId"
	SpeechServiceConnectionUserDefinedQueryParams PropertyID = "SpeechServiceConnection_UserDefinedQueryParameters"
)

// ResultReason defines the reason a result was generated
type ResultReason int

// ResultReason constants
const (
	ResultReasonRecognizedSpeech ResultReason = iota
	ResultReasonNoMatch
	ResultReasonCanceled
	ResultReasonTranslatedSpeech
)

// String returns the string representation of ResultReason
func (r ResultReason) String() string {
	switch r {
	case ResultReasonRecognizedSpeech:
		return "RecognizedSpeech"
	case ResultReasonNoMatch:
		return "NoMatch"
	case ResultReasonCanceled:
		return "Canceled"
	case ResultReasonTranslatedSpeech:
		return "TranslatedSpeech"
	default:
		return fmt.Sprintf("Unknown ResultReason (%d)", r)
	}
}

// CancellationReason defines the reason a recognition was canceled
type CancellationReason int

// CancellationReason constants
const (
	CancellationReasonError CancellationReason = iota
	CancellationReasonEndOfStream
)

// String returns the string representation of CancellationReason
func (r CancellationReason) String() string {
	switch r {
	case CancellationReasonError:
		return "Error"
	case CancellationReasonEndOfStream:
		return "EndOfStream"
	default:
		return fmt.Sprintf("Unknown CancellationReason (%d)", r)
	}
}

// CancellationErrorCode defines specific error codes for cancellation
type CancellationErrorCode int

// CancellationErrorCode constants
const (
	CancellationErrorNoError CancellationErrorCode = iota
	CancellationErrorAuthenticationFailure
	CancellationErrorBadRequest
	CancellationErrorTooManyRequests
	CancellationErrorForbidden
	CancellationErrorConnectionFailure
	CancellationErrorServiceTimeout
	CancellationErrorServiceError
	CancellationErrorServiceUnavailable
	CancellationErrorRuntimeError
)

// String returns the string representation of CancellationErrorCode
func (e CancellationErrorCode) String() string {
	switch e {
	case CancellationErrorNoError:
		return "NoError"
	case CancellationErrorAuthenticationFailure:
		return "AuthenticationFailure"
	case CancellationErrorBadRequest:
		return "BadRequest"
	case CancellationErrorTooManyRequests:
		return "TooManyRequests"
	case CancellationErrorForbidden:
		return "Forbidden"
	case CancellationErrorConnectionFailure:
		return "ConnectionFailure"
	case CancellationErrorServiceTimeout:
		return "ServiceTimeout"
	case CancellationErrorServiceError:
		return "ServiceError"
	case CancellationErrorServiceUnavailable:
		return "ServiceUnavailable"
	case CancellationErrorRuntimeError:
		return "RuntimeError"
	default:
		return fmt.Sprintf("Unknown CancellationErrorCode (%d)", e)
	}
}

// OutputFormat defines different output formats for recognition results
type OutputFormat int

// OutputFormat constants
const (
	OutputFormatSimple OutputFormat = iota
	OutputFormatDetailed
)

// String returns the string representation of OutputFormat
func (f OutputFormat) String() string {
	switch f {
	case OutputFormatSimple:
		return "Simple"
	case OutputFormatDetailed:
		return "Detailed"
	default:
		return fmt.Sprintf("Unknown OutputFormat (%d)", f)
	}
}

// SpeechSynthesisOutputFormat defines the audio output formats supported for synthesis
type SpeechSynthesisOutputFormat int

// SpeechSynthesisOutputFormat constants
const (
	SpeechSynthesisOutputFormatRaw8Khz8BitMonoPCM SpeechSynthesisOutputFormat = iota
	SpeechSynthesisOutputFormatRaw16Khz16BitMonoPCM
	SpeechSynthesisOutputFormatRiff8Khz8BitMonoPCM
	SpeechSynthesisOutputFormatRiff16Khz16BitMonoPCM
)

// String returns the string representation of SpeechSynthesisOutputFormat
func (f SpeechSynthesisOutputFormat) String() string {
	switch f {
	case SpeechSynthesisOutputFormatRaw8Khz8BitMonoPCM:
		return "Raw8Khz8BitMonoPCM"
	case SpeechSynthesisOutputFormatRaw16Khz16BitMonoPCM:
		return "Raw16Khz16BitMonoPCM"
	case SpeechSynthesisOutputFormatRiff8Khz8BitMonoPCM:
		return "Riff8Khz8BitMonoPCM"
	case SpeechSynthesisOutputFormatRiff16Khz16BitMonoPCM:
		return "Riff16Khz16BitMonoPCM"
	default:
		return fmt.Sprintf("Unknown SpeechSynthesisOutputFormat (%d)", f)
	}
}

// ServicePropertyChannel defines the channels used to pass service properties
type ServicePropertyChannel int

// ServicePropertyChannel constants
const (
	ServicePropertyChannelURI ServicePropertyChannel = iota
	ServicePropertyChannelQuery
)

// String returns the string representation of ServicePropertyChannel
func (c ServicePropertyChannel) String() string {
	switch c {
	case ServicePropertyChannelURI:
		return "URI"
	case ServicePropertyChannelQuery:
		return "Query"
	default:
		return fmt.Sprintf("Unknown ServicePropertyChannel (%d)", c)
	}
}
