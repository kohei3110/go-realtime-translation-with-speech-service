// Code generated by Microsoft (R) AutoRest Code Generator (autorest: 3.10.4, generator: @autorest/go@4.0.0-preview.70)
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// Code generated by @autorest/go. DO NOT EDIT.

package speechclient

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// TranscriptionsClient contains the methods for the Transcriptions group.
// Don't use this type directly, use a constructor function instead.
type TranscriptionsClient struct {
	internal *azcore.Client
	endpoint string
}

// NewTranscriptionsClient creates a new instance of TranscriptionsClient with the specified values.
//   - endpoint - The endpoint of your speech service resource.
//   - credential - Used to authorize requests. Usually a credential from azidentity.
//   - options - Pass nil to accept the default values.
func NewTranscriptionsClient(endpoint string, credential azcore.TokenCredential, options *azcore.ClientOptions) (*TranscriptionsClient, error) {
	if endpoint == "" {
		return nil, errors.New("parameter endpoint cannot be empty")
	}
	if credential == nil {
		return nil, errors.New("parameter credential cannot be nil")
	}

	// Set up authentication policy
	pipelineOptions := runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewBearerTokenPolicy(credential, []string{"https://cognitiveservices.azure.com/.default"}, nil),
		},
	}

	cl, err := azcore.NewClient("speechclient.TranscriptionsClient", "v3.2.0", pipelineOptions, options)
	if err != nil {
		return nil, err
	}
	client := &TranscriptionsClient{
		internal: cl,
		endpoint: endpoint,
	}
	return client, nil
}

// Create - Creates a new transcription.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - transcription - The details of the new transcription.
//   - options - TranscriptionsClientCreateOptions contains the optional parameters for the TranscriptionsClient.Create method.
func (client *TranscriptionsClient) Create(ctx context.Context, transcription Transcription, options *TranscriptionsClientCreateOptions) (TranscriptionsClientCreateResponse, error) {
	var err error
	req, err := client.createCreateRequest(ctx, transcription, options)
	if err != nil {
		return TranscriptionsClientCreateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TranscriptionsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return TranscriptionsClientCreateResponse{}, err
	}
	resp, err := client.createHandleResponse(httpResp)
	return resp, err
}

// createCreateRequest creates the Create request.
func (client *TranscriptionsClient) createCreateRequest(ctx context.Context, transcription Transcription, _ *TranscriptionsClientCreateOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/transcriptions"
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	// TODO: FIXME
	req.Raw().Header["Ocp-Apim-Subscription-Key"] = []string{"REPLACE_WITH_YOUR_SUBSCRIPTION_KEY"}
	if err := runtime.MarshalAsJSON(req, transcription); err != nil {
	return nil, err
}
;	return req, nil
}

// createHandleResponse handles the Create response.
func (client *TranscriptionsClient) createHandleResponse(resp *http.Response) (TranscriptionsClientCreateResponse, error) {
	result := TranscriptionsClientCreateResponse{}
	if val := resp.Header.Get("Location"); val != "" {
		result.Location = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Transcription); err != nil {
		return TranscriptionsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the specified transcription task.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the transcription.
//   - options - TranscriptionsClientDeleteOptions contains the optional parameters for the TranscriptionsClient.Delete method.
func (client *TranscriptionsClient) Delete(ctx context.Context, id string, options *TranscriptionsClientDeleteOptions) (TranscriptionsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, id, options)
	if err != nil {
		return TranscriptionsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TranscriptionsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return TranscriptionsClientDeleteResponse{}, err
	}
	return TranscriptionsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *TranscriptionsClient) deleteCreateRequest(ctx context.Context, id string, _ *TranscriptionsClientDeleteOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/transcriptions/{id}"
	if id == "" {
		return nil, errors.New("parameter id cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{id}", url.PathEscape(id))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the transcription identified by the given ID.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the transcription.
//   - options - TranscriptionsClientGetOptions contains the optional parameters for the TranscriptionsClient.Get method.
func (client *TranscriptionsClient) Get(ctx context.Context, id string, options *TranscriptionsClientGetOptions) (TranscriptionsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, id, options)
	if err != nil {
		return TranscriptionsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TranscriptionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return TranscriptionsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *TranscriptionsClient) getCreateRequest(ctx context.Context, id string, _ *TranscriptionsClientGetOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/transcriptions/{id}"
	if id == "" {
		return nil, errors.New("parameter id cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{id}", url.PathEscape(id))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	// TODO: FIXME
	req.Raw().Header["Ocp-Apim-Subscription-Key"] = []string{"REPLACE_WITH_YOUR_SUBSCRIPTION_KEY"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *TranscriptionsClient) getHandleResponse(resp *http.Response) (TranscriptionsClientGetResponse, error) {
	result := TranscriptionsClientGetResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return TranscriptionsClientGetResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Transcription); err != nil {
		return TranscriptionsClientGetResponse{}, err
	}
	return result, nil
}

// GetFile - Gets one specific file (identified with fileId) from a transcription (identified with id).
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the transcription.
//   - fileID - The identifier of the file.
//   - options - TranscriptionsClientGetFileOptions contains the optional parameters for the TranscriptionsClient.GetFile method.
func (client *TranscriptionsClient) GetFile(ctx context.Context, id string, fileID string, options *TranscriptionsClientGetFileOptions) (TranscriptionsClientGetFileResponse, error) {
	var err error
	req, err := client.getFileCreateRequest(ctx, id, fileID, options)
	if err != nil {
		return TranscriptionsClientGetFileResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TranscriptionsClientGetFileResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return TranscriptionsClientGetFileResponse{}, err
	}
	resp, err := client.getFileHandleResponse(httpResp)
	return resp, err
}

// getFileCreateRequest creates the GetFile request.
func (client *TranscriptionsClient) getFileCreateRequest(ctx context.Context, id string, fileID string, options *TranscriptionsClientGetFileOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/transcriptions/{id}/files/{fileId}"
	if id == "" {
		return nil, errors.New("parameter id cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{id}", url.PathEscape(id))
	if fileID == "" {
		return nil, errors.New("parameter fileID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{fileId}", url.PathEscape(fileID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.SasValidityInSeconds != nil {
		reqQP.Set("sasValidityInSeconds", strconv.FormatInt(int64(*options.SasValidityInSeconds), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getFileHandleResponse handles the GetFile response.
func (client *TranscriptionsClient) getFileHandleResponse(resp *http.Response) (TranscriptionsClientGetFileResponse, error) {
	result := TranscriptionsClientGetFileResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.File); err != nil {
		return TranscriptionsClientGetFileResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets a list of transcriptions for the authenticated subscription.
//
// Generated from API version 3.2
//   - options - TranscriptionsClientListOptions contains the optional parameters for the TranscriptionsClient.NewListPager method.
func (client *TranscriptionsClient) NewListPager(options *TranscriptionsClientListOptions) (*runtime.Pager[TranscriptionsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[TranscriptionsClientListResponse]{
		More: func(page TranscriptionsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *TranscriptionsClientListResponse) (TranscriptionsClientListResponse, error) {
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return TranscriptionsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
			},
	})
}

// listCreateRequest creates the List request.
func (client *TranscriptionsClient) listCreateRequest(ctx context.Context, options *TranscriptionsClientListOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/transcriptions"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("filter", *options.Filter)
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Top != nil {
		reqQP.Set("top", strconv.FormatInt(int64(*options.Top), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *TranscriptionsClient) listHandleResponse(resp *http.Response) (TranscriptionsClientListResponse, error) {
	result := TranscriptionsClientListResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return TranscriptionsClientListResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.PaginatedTranscriptions); err != nil {
		return TranscriptionsClientListResponse{}, err
	}
	return result, nil
}

// NewListFilesPager - Gets the files of the transcription identified by the given ID.
//
// Generated from API version 3.2
//   - id - The identifier of the transcription.
//   - options - TranscriptionsClientListFilesOptions contains the optional parameters for the TranscriptionsClient.NewListFilesPager
//     method.
func (client *TranscriptionsClient) NewListFilesPager(id string, options *TranscriptionsClientListFilesOptions) (*runtime.Pager[TranscriptionsClientListFilesResponse]) {
	return runtime.NewPager(runtime.PagingHandler[TranscriptionsClientListFilesResponse]{
		More: func(page TranscriptionsClientListFilesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *TranscriptionsClientListFilesResponse) (TranscriptionsClientListFilesResponse, error) {
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listFilesCreateRequest(ctx, id, options)
			}, nil)
			if err != nil {
				return TranscriptionsClientListFilesResponse{}, err
			}
			return client.listFilesHandleResponse(resp)
			},
	})
}

// listFilesCreateRequest creates the ListFiles request.
func (client *TranscriptionsClient) listFilesCreateRequest(ctx context.Context, id string, options *TranscriptionsClientListFilesOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/transcriptions/{id}/files"
	if id == "" {
		return nil, errors.New("parameter id cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{id}", url.PathEscape(id))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("filter", *options.Filter)
	}
	if options != nil && options.SasValidityInSeconds != nil {
		reqQP.Set("sasValidityInSeconds", strconv.FormatInt(int64(*options.SasValidityInSeconds), 10))
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Top != nil {
		reqQP.Set("top", strconv.FormatInt(int64(*options.Top), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listFilesHandleResponse handles the ListFiles response.
func (client *TranscriptionsClient) listFilesHandleResponse(resp *http.Response) (TranscriptionsClientListFilesResponse, error) {
	result := TranscriptionsClientListFilesResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return TranscriptionsClientListFilesResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.PaginatedFiles); err != nil {
		return TranscriptionsClientListFilesResponse{}, err
	}
	return result, nil
}

// ListSupportedLocales - Gets a list of supported locales for offline transcriptions.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - options - TranscriptionsClientListSupportedLocalesOptions contains the optional parameters for the TranscriptionsClient.ListSupportedLocales
//     method.
func (client *TranscriptionsClient) ListSupportedLocales(ctx context.Context, options *TranscriptionsClientListSupportedLocalesOptions) (TranscriptionsClientListSupportedLocalesResponse, error) {
	var err error
	req, err := client.listSupportedLocalesCreateRequest(ctx, options)
	if err != nil {
		return TranscriptionsClientListSupportedLocalesResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TranscriptionsClientListSupportedLocalesResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return TranscriptionsClientListSupportedLocalesResponse{}, err
	}
	resp, err := client.listSupportedLocalesHandleResponse(httpResp)
	return resp, err
}

// listSupportedLocalesCreateRequest creates the ListSupportedLocales request.
func (client *TranscriptionsClient) listSupportedLocalesCreateRequest(ctx context.Context, _ *TranscriptionsClientListSupportedLocalesOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/transcriptions/locales"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listSupportedLocalesHandleResponse handles the ListSupportedLocales response.
func (client *TranscriptionsClient) listSupportedLocalesHandleResponse(resp *http.Response) (TranscriptionsClientListSupportedLocalesResponse, error) {
	result := TranscriptionsClientListSupportedLocalesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StringArray); err != nil {
		return TranscriptionsClientListSupportedLocalesResponse{}, err
	}
	return result, nil
}

// Update - Updates the mutable details of the transcription identified by its ID.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the transcription.
//   - transcriptionUpdate - The updated values for the transcription.
//   - options - TranscriptionsClientUpdateOptions contains the optional parameters for the TranscriptionsClient.Update method.
func (client *TranscriptionsClient) Update(ctx context.Context, id string, transcriptionUpdate TranscriptionUpdate, options *TranscriptionsClientUpdateOptions) (TranscriptionsClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, id, transcriptionUpdate, options)
	if err != nil {
		return TranscriptionsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TranscriptionsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return TranscriptionsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *TranscriptionsClient) updateCreateRequest(ctx context.Context, id string, transcriptionUpdate TranscriptionUpdate, _ *TranscriptionsClientUpdateOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/transcriptions/{id}"
	if id == "" {
		return nil, errors.New("parameter id cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{id}", url.PathEscape(id))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, transcriptionUpdate); err != nil {
	return nil, err
}
;	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *TranscriptionsClient) updateHandleResponse(resp *http.Response) (TranscriptionsClientUpdateResponse, error) {
	result := TranscriptionsClientUpdateResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return TranscriptionsClientUpdateResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Transcription); err != nil {
		return TranscriptionsClientUpdateResponse{}, err
	}
	return result, nil
}

