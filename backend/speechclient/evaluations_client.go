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

// EvaluationsClient contains the methods for the Evaluations group.
// Don't use this type directly, use a constructor function instead.
type EvaluationsClient struct {
	internal *azcore.Client
	endpoint string
}

// Create - Creates a new evaluation.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - evaluation - The details of the new evaluation.
//   - options - EvaluationsClientCreateOptions contains the optional parameters for the EvaluationsClient.Create method.
func (client *EvaluationsClient) Create(ctx context.Context, evaluation Evaluation, options *EvaluationsClientCreateOptions) (EvaluationsClientCreateResponse, error) {
	var err error
	req, err := client.createCreateRequest(ctx, evaluation, options)
	if err != nil {
		return EvaluationsClientCreateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EvaluationsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return EvaluationsClientCreateResponse{}, err
	}
	resp, err := client.createHandleResponse(httpResp)
	return resp, err
}

// createCreateRequest creates the Create request.
func (client *EvaluationsClient) createCreateRequest(ctx context.Context, evaluation Evaluation, _ *EvaluationsClientCreateOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/evaluations"
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, evaluation); err != nil {
	return nil, err
}
;	return req, nil
}

// createHandleResponse handles the Create response.
func (client *EvaluationsClient) createHandleResponse(resp *http.Response) (EvaluationsClientCreateResponse, error) {
	result := EvaluationsClientCreateResponse{}
	if val := resp.Header.Get("Location"); val != "" {
		result.Location = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Evaluation); err != nil {
		return EvaluationsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the evaluation identified by the given ID.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the evaluation.
//   - options - EvaluationsClientDeleteOptions contains the optional parameters for the EvaluationsClient.Delete method.
func (client *EvaluationsClient) Delete(ctx context.Context, id string, options *EvaluationsClientDeleteOptions) (EvaluationsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, id, options)
	if err != nil {
		return EvaluationsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EvaluationsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return EvaluationsClientDeleteResponse{}, err
	}
	return EvaluationsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *EvaluationsClient) deleteCreateRequest(ctx context.Context, id string, _ *EvaluationsClientDeleteOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/evaluations/{id}"
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

// Get - Gets the evaluation identified by the given ID.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the evaluation.
//   - options - EvaluationsClientGetOptions contains the optional parameters for the EvaluationsClient.Get method.
func (client *EvaluationsClient) Get(ctx context.Context, id string, options *EvaluationsClientGetOptions) (EvaluationsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, id, options)
	if err != nil {
		return EvaluationsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EvaluationsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EvaluationsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *EvaluationsClient) getCreateRequest(ctx context.Context, id string, _ *EvaluationsClientGetOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/evaluations/{id}"
	if id == "" {
		return nil, errors.New("parameter id cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{id}", url.PathEscape(id))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *EvaluationsClient) getHandleResponse(resp *http.Response) (EvaluationsClientGetResponse, error) {
	result := EvaluationsClientGetResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return EvaluationsClientGetResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Evaluation); err != nil {
		return EvaluationsClientGetResponse{}, err
	}
	return result, nil
}

// GetFile - Gets one specific file (identified with fileId) from an evaluation (identified with id).
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the evaluation.
//   - fileID - The identifier of the file.
//   - options - EvaluationsClientGetFileOptions contains the optional parameters for the EvaluationsClient.GetFile method.
func (client *EvaluationsClient) GetFile(ctx context.Context, id string, fileID string, options *EvaluationsClientGetFileOptions) (EvaluationsClientGetFileResponse, error) {
	var err error
	req, err := client.getFileCreateRequest(ctx, id, fileID, options)
	if err != nil {
		return EvaluationsClientGetFileResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EvaluationsClientGetFileResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EvaluationsClientGetFileResponse{}, err
	}
	resp, err := client.getFileHandleResponse(httpResp)
	return resp, err
}

// getFileCreateRequest creates the GetFile request.
func (client *EvaluationsClient) getFileCreateRequest(ctx context.Context, id string, fileID string, options *EvaluationsClientGetFileOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/evaluations/{id}/files/{fileId}"
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
func (client *EvaluationsClient) getFileHandleResponse(resp *http.Response) (EvaluationsClientGetFileResponse, error) {
	result := EvaluationsClientGetFileResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.File); err != nil {
		return EvaluationsClientGetFileResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets the list of evaluations for the authenticated subscription.
//
// Generated from API version 3.2
//   - options - EvaluationsClientListOptions contains the optional parameters for the EvaluationsClient.NewListPager method.
func (client *EvaluationsClient) NewListPager(options *EvaluationsClientListOptions) (*runtime.Pager[EvaluationsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[EvaluationsClientListResponse]{
		More: func(page EvaluationsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *EvaluationsClientListResponse) (EvaluationsClientListResponse, error) {
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return EvaluationsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
			},
	})
}

// listCreateRequest creates the List request.
func (client *EvaluationsClient) listCreateRequest(ctx context.Context, options *EvaluationsClientListOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/evaluations"
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
func (client *EvaluationsClient) listHandleResponse(resp *http.Response) (EvaluationsClientListResponse, error) {
	result := EvaluationsClientListResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return EvaluationsClientListResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.PaginatedEvaluations); err != nil {
		return EvaluationsClientListResponse{}, err
	}
	return result, nil
}

// NewListFilesPager - Gets the files of the evaluation identified by the given ID.
//
// Generated from API version 3.2
//   - id - The identifier of the evaluation.
//   - options - EvaluationsClientListFilesOptions contains the optional parameters for the EvaluationsClient.NewListFilesPager
//     method.
func (client *EvaluationsClient) NewListFilesPager(id string, options *EvaluationsClientListFilesOptions) (*runtime.Pager[EvaluationsClientListFilesResponse]) {
	return runtime.NewPager(runtime.PagingHandler[EvaluationsClientListFilesResponse]{
		More: func(page EvaluationsClientListFilesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *EvaluationsClientListFilesResponse) (EvaluationsClientListFilesResponse, error) {
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listFilesCreateRequest(ctx, id, options)
			}, nil)
			if err != nil {
				return EvaluationsClientListFilesResponse{}, err
			}
			return client.listFilesHandleResponse(resp)
			},
	})
}

// listFilesCreateRequest creates the ListFiles request.
func (client *EvaluationsClient) listFilesCreateRequest(ctx context.Context, id string, options *EvaluationsClientListFilesOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/evaluations/{id}/files"
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
func (client *EvaluationsClient) listFilesHandleResponse(resp *http.Response) (EvaluationsClientListFilesResponse, error) {
	result := EvaluationsClientListFilesResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return EvaluationsClientListFilesResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.PaginatedFiles); err != nil {
		return EvaluationsClientListFilesResponse{}, err
	}
	return result, nil
}

// ListSupportedLocales - Gets a list of supported locales for evaluations.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - options - EvaluationsClientListSupportedLocalesOptions contains the optional parameters for the EvaluationsClient.ListSupportedLocales
//     method.
func (client *EvaluationsClient) ListSupportedLocales(ctx context.Context, options *EvaluationsClientListSupportedLocalesOptions) (EvaluationsClientListSupportedLocalesResponse, error) {
	var err error
	req, err := client.listSupportedLocalesCreateRequest(ctx, options)
	if err != nil {
		return EvaluationsClientListSupportedLocalesResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EvaluationsClientListSupportedLocalesResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EvaluationsClientListSupportedLocalesResponse{}, err
	}
	resp, err := client.listSupportedLocalesHandleResponse(httpResp)
	return resp, err
}

// listSupportedLocalesCreateRequest creates the ListSupportedLocales request.
func (client *EvaluationsClient) listSupportedLocalesCreateRequest(ctx context.Context, _ *EvaluationsClientListSupportedLocalesOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/evaluations/locales"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listSupportedLocalesHandleResponse handles the ListSupportedLocales response.
func (client *EvaluationsClient) listSupportedLocalesHandleResponse(resp *http.Response) (EvaluationsClientListSupportedLocalesResponse, error) {
	result := EvaluationsClientListSupportedLocalesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StringArray); err != nil {
		return EvaluationsClientListSupportedLocalesResponse{}, err
	}
	return result, nil
}

// Update - Updates the mutable details of the evaluation identified by its id.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the evaluation.
//   - evaluationUpdate - The object containing the updated fields of the evaluation.
//   - options - EvaluationsClientUpdateOptions contains the optional parameters for the EvaluationsClient.Update method.
func (client *EvaluationsClient) Update(ctx context.Context, id string, evaluationUpdate EvaluationUpdate, options *EvaluationsClientUpdateOptions) (EvaluationsClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, id, evaluationUpdate, options)
	if err != nil {
		return EvaluationsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EvaluationsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EvaluationsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *EvaluationsClient) updateCreateRequest(ctx context.Context, id string, evaluationUpdate EvaluationUpdate, _ *EvaluationsClientUpdateOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/evaluations/{id}"
	if id == "" {
		return nil, errors.New("parameter id cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{id}", url.PathEscape(id))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, evaluationUpdate); err != nil {
	return nil, err
}
;	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *EvaluationsClient) updateHandleResponse(resp *http.Response) (EvaluationsClientUpdateResponse, error) {
	result := EvaluationsClientUpdateResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return EvaluationsClientUpdateResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Evaluation); err != nil {
		return EvaluationsClientUpdateResponse{}, err
	}
	return result, nil
}

