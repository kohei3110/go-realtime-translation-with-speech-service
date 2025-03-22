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
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// DatasetsClient contains the methods for the Datasets group.
// Don't use this type directly, use a constructor function instead.
type DatasetsClient struct {
	internal *azcore.Client
	endpoint string
}

// CommitBlocks - Commit block list to complete the upload of the dataset.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the dataset.
//   - blockList - The list of blocks that compile the dataset.
//   - options - DatasetsClientCommitBlocksOptions contains the optional parameters for the DatasetsClient.CommitBlocks method.
func (client *DatasetsClient) CommitBlocks(ctx context.Context, id string, blockList []*CommitBlocksEntry, options *DatasetsClientCommitBlocksOptions) (DatasetsClientCommitBlocksResponse, error) {
	var err error
	req, err := client.commitBlocksCreateRequest(ctx, id, blockList, options)
	if err != nil {
		return DatasetsClientCommitBlocksResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatasetsClientCommitBlocksResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DatasetsClientCommitBlocksResponse{}, err
	}
	return DatasetsClientCommitBlocksResponse{}, nil
}

// commitBlocksCreateRequest creates the CommitBlocks request.
func (client *DatasetsClient) commitBlocksCreateRequest(ctx context.Context, id string, blockList []*CommitBlocksEntry, _ *DatasetsClientCommitBlocksOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets/{id}/blocks:commit"
	if id == "" {
		return nil, errors.New("parameter id cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{id}", url.PathEscape(id))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, blockList); err != nil {
	return nil, err
}
;	return req, nil
}

// Create - Uploads and creates a new dataset by getting the data from a specified URL or starts waiting for data blocks to
// be uploaded.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - dataset - Definition for the new dataset.
//   - options - DatasetsClientCreateOptions contains the optional parameters for the DatasetsClient.Create method.
func (client *DatasetsClient) Create(ctx context.Context, dataset Dataset, options *DatasetsClientCreateOptions) (DatasetsClientCreateResponse, error) {
	var err error
	req, err := client.createCreateRequest(ctx, dataset, options)
	if err != nil {
		return DatasetsClientCreateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatasetsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return DatasetsClientCreateResponse{}, err
	}
	resp, err := client.createHandleResponse(httpResp)
	return resp, err
}

// createCreateRequest creates the Create request.
func (client *DatasetsClient) createCreateRequest(ctx context.Context, dataset Dataset, _ *DatasetsClientCreateOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets"
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, dataset); err != nil {
	return nil, err
}
;	return req, nil
}

// createHandleResponse handles the Create response.
func (client *DatasetsClient) createHandleResponse(resp *http.Response) (DatasetsClientCreateResponse, error) {
	result := DatasetsClientCreateResponse{}
	if val := resp.Header.Get("Location"); val != "" {
		result.Location = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Dataset); err != nil {
		return DatasetsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the specified dataset.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the dataset.
//   - options - DatasetsClientDeleteOptions contains the optional parameters for the DatasetsClient.Delete method.
func (client *DatasetsClient) Delete(ctx context.Context, id string, options *DatasetsClientDeleteOptions) (DatasetsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, id, options)
	if err != nil {
		return DatasetsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatasetsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return DatasetsClientDeleteResponse{}, err
	}
	return DatasetsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *DatasetsClient) deleteCreateRequest(ctx context.Context, id string, _ *DatasetsClientDeleteOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets/{id}"
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

// Get - Gets the dataset identified by the given ID.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the dataset.
//   - options - DatasetsClientGetOptions contains the optional parameters for the DatasetsClient.Get method.
func (client *DatasetsClient) Get(ctx context.Context, id string, options *DatasetsClientGetOptions) (DatasetsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, id, options)
	if err != nil {
		return DatasetsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatasetsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DatasetsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *DatasetsClient) getCreateRequest(ctx context.Context, id string, _ *DatasetsClientGetOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets/{id}"
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
func (client *DatasetsClient) getHandleResponse(resp *http.Response) (DatasetsClientGetResponse, error) {
	result := DatasetsClientGetResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return DatasetsClientGetResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Dataset); err != nil {
		return DatasetsClientGetResponse{}, err
	}
	return result, nil
}

// GetBlocks - Gets the list of uploaded blocks for this dataset.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the dataset.
//   - options - DatasetsClientGetBlocksOptions contains the optional parameters for the DatasetsClient.GetBlocks method.
func (client *DatasetsClient) GetBlocks(ctx context.Context, id string, options *DatasetsClientGetBlocksOptions) (DatasetsClientGetBlocksResponse, error) {
	var err error
	req, err := client.getBlocksCreateRequest(ctx, id, options)
	if err != nil {
		return DatasetsClientGetBlocksResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatasetsClientGetBlocksResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DatasetsClientGetBlocksResponse{}, err
	}
	resp, err := client.getBlocksHandleResponse(httpResp)
	return resp, err
}

// getBlocksCreateRequest creates the GetBlocks request.
func (client *DatasetsClient) getBlocksCreateRequest(ctx context.Context, id string, _ *DatasetsClientGetBlocksOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets/{id}/blocks"
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

// getBlocksHandleResponse handles the GetBlocks response.
func (client *DatasetsClient) getBlocksHandleResponse(resp *http.Response) (DatasetsClientGetBlocksResponse, error) {
	result := DatasetsClientGetBlocksResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return DatasetsClientGetBlocksResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.UploadedBlocks); err != nil {
		return DatasetsClientGetBlocksResponse{}, err
	}
	return result, nil
}

// GetFile - Gets one specific file (identified with fileId) from a dataset (identified with id).
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the dataset.
//   - fileID - The identifier of the file.
//   - options - DatasetsClientGetFileOptions contains the optional parameters for the DatasetsClient.GetFile method.
func (client *DatasetsClient) GetFile(ctx context.Context, id string, fileID string, options *DatasetsClientGetFileOptions) (DatasetsClientGetFileResponse, error) {
	var err error
	req, err := client.getFileCreateRequest(ctx, id, fileID, options)
	if err != nil {
		return DatasetsClientGetFileResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatasetsClientGetFileResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DatasetsClientGetFileResponse{}, err
	}
	resp, err := client.getFileHandleResponse(httpResp)
	return resp, err
}

// getFileCreateRequest creates the GetFile request.
func (client *DatasetsClient) getFileCreateRequest(ctx context.Context, id string, fileID string, options *DatasetsClientGetFileOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets/{id}/files/{fileId}"
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
func (client *DatasetsClient) getFileHandleResponse(resp *http.Response) (DatasetsClientGetFileResponse, error) {
	result := DatasetsClientGetFileResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.File); err != nil {
		return DatasetsClientGetFileResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets a list of datasets for the authenticated subscription.
//
// Generated from API version 3.2
//   - options - DatasetsClientListOptions contains the optional parameters for the DatasetsClient.NewListPager method.
func (client *DatasetsClient) NewListPager(options *DatasetsClientListOptions) (*runtime.Pager[DatasetsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[DatasetsClientListResponse]{
		More: func(page DatasetsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *DatasetsClientListResponse) (DatasetsClientListResponse, error) {
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return DatasetsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
			},
	})
}

// listCreateRequest creates the List request.
func (client *DatasetsClient) listCreateRequest(ctx context.Context, options *DatasetsClientListOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets"
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
func (client *DatasetsClient) listHandleResponse(resp *http.Response) (DatasetsClientListResponse, error) {
	result := DatasetsClientListResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return DatasetsClientListResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.PaginatedDatasets); err != nil {
		return DatasetsClientListResponse{}, err
	}
	return result, nil
}

// NewListFilesPager - Gets the files of the dataset identified by the given ID.
//
// Generated from API version 3.2
//   - id - The identifier of the dataset.
//   - options - DatasetsClientListFilesOptions contains the optional parameters for the DatasetsClient.NewListFilesPager method.
func (client *DatasetsClient) NewListFilesPager(id string, options *DatasetsClientListFilesOptions) (*runtime.Pager[DatasetsClientListFilesResponse]) {
	return runtime.NewPager(runtime.PagingHandler[DatasetsClientListFilesResponse]{
		More: func(page DatasetsClientListFilesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *DatasetsClientListFilesResponse) (DatasetsClientListFilesResponse, error) {
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listFilesCreateRequest(ctx, id, options)
			}, nil)
			if err != nil {
				return DatasetsClientListFilesResponse{}, err
			}
			return client.listFilesHandleResponse(resp)
			},
	})
}

// listFilesCreateRequest creates the ListFiles request.
func (client *DatasetsClient) listFilesCreateRequest(ctx context.Context, id string, options *DatasetsClientListFilesOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets/{id}/files"
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
func (client *DatasetsClient) listFilesHandleResponse(resp *http.Response) (DatasetsClientListFilesResponse, error) {
	result := DatasetsClientListFilesResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return DatasetsClientListFilesResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.PaginatedFiles); err != nil {
		return DatasetsClientListFilesResponse{}, err
	}
	return result, nil
}

// ListSupportedLocales - Gets a list of supported locales for datasets.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - options - DatasetsClientListSupportedLocalesOptions contains the optional parameters for the DatasetsClient.ListSupportedLocales
//     method.
func (client *DatasetsClient) ListSupportedLocales(ctx context.Context, options *DatasetsClientListSupportedLocalesOptions) (DatasetsClientListSupportedLocalesResponse, error) {
	var err error
	req, err := client.listSupportedLocalesCreateRequest(ctx, options)
	if err != nil {
		return DatasetsClientListSupportedLocalesResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatasetsClientListSupportedLocalesResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DatasetsClientListSupportedLocalesResponse{}, err
	}
	resp, err := client.listSupportedLocalesHandleResponse(httpResp)
	return resp, err
}

// listSupportedLocalesCreateRequest creates the ListSupportedLocales request.
func (client *DatasetsClient) listSupportedLocalesCreateRequest(ctx context.Context, _ *DatasetsClientListSupportedLocalesOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets/locales"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listSupportedLocalesHandleResponse handles the ListSupportedLocales response.
func (client *DatasetsClient) listSupportedLocalesHandleResponse(resp *http.Response) (DatasetsClientListSupportedLocalesResponse, error) {
	result := DatasetsClientListSupportedLocalesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Value); err != nil {
		return DatasetsClientListSupportedLocalesResponse{}, err
	}
	return result, nil
}

// Update - Updates the mutable details of the dataset identified by its ID.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the dataset.
//   - datasetUpdate - The updated values for the dataset.
//   - options - DatasetsClientUpdateOptions contains the optional parameters for the DatasetsClient.Update method.
func (client *DatasetsClient) Update(ctx context.Context, id string, datasetUpdate DatasetUpdate, options *DatasetsClientUpdateOptions) (DatasetsClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, id, datasetUpdate, options)
	if err != nil {
		return DatasetsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatasetsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DatasetsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *DatasetsClient) updateCreateRequest(ctx context.Context, id string, datasetUpdate DatasetUpdate, _ *DatasetsClientUpdateOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets/{id}"
	if id == "" {
		return nil, errors.New("parameter id cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{id}", url.PathEscape(id))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, datasetUpdate); err != nil {
	return nil, err
}
;	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *DatasetsClient) updateHandleResponse(resp *http.Response) (DatasetsClientUpdateResponse, error) {
	result := DatasetsClientUpdateResponse{}
	if val := resp.Header.Get("Retry-After"); val != "" {
		retryAfter32, err := strconv.ParseInt(val, 10, 32)
		retryAfter := int32(retryAfter32)
		if err != nil {
			return DatasetsClientUpdateResponse{}, err
		}
		result.RetryAfter = &retryAfter
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Dataset); err != nil {
		return DatasetsClientUpdateResponse{}, err
	}
	return result, nil
}

// Upload - Uploads data and creates a new dataset.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - displayName - The name of this dataset.
//   - locale - The locale of this dataset.
//   - kind - The kind of the dataset. Possible values are "Language", "Acoustic", "Pronunciation", "AudioFiles", "LanguageMarkdown",
//     "OutputFormatting".
//   - options - DatasetsClientUploadOptions contains the optional parameters for the DatasetsClient.Upload method.
func (client *DatasetsClient) Upload(ctx context.Context, displayName string, locale string, kind string, options *DatasetsClientUploadOptions) (DatasetsClientUploadResponse, error) {
	var err error
	req, err := client.uploadCreateRequest(ctx, displayName, locale, kind, options)
	if err != nil {
		return DatasetsClientUploadResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatasetsClientUploadResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return DatasetsClientUploadResponse{}, err
	}
	resp, err := client.uploadHandleResponse(httpResp)
	return resp, err
}

// uploadCreateRequest creates the Upload request.
func (client *DatasetsClient) uploadCreateRequest(ctx context.Context, displayName string, locale string, kind string, options *DatasetsClientUploadOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets/upload"
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	formData := map[string]any{}
	if options != nil && options.Project != nil {
	formData["Project"] = *options.Project
	}
	formData["displayName"] = displayName
	if options != nil && options.Description != nil {
	formData["Description"] = *options.Description
	}
	formData["locale"] = locale
	formData["kind"] = kind
	if options != nil && options.CustomProperties != nil {
	formData["CustomProperties"] = *options.CustomProperties
	}
	if options != nil && options.Data != nil {
	formData["Data"] = options.Data
	}
	if options != nil && options.Email != nil {
	formData["Email"] = *options.Email
	}
	if err := runtime.SetMultipartFormData(req, formData); err != nil {
		return nil, err
	}
	return req, nil
}

// uploadHandleResponse handles the Upload response.
func (client *DatasetsClient) uploadHandleResponse(resp *http.Response) (DatasetsClientUploadResponse, error) {
	result := DatasetsClientUploadResponse{}
	if val := resp.Header.Get("Location"); val != "" {
		result.Location = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Dataset); err != nil {
		return DatasetsClientUploadResponse{}, err
	}
	return result, nil
}

// UploadBlock - Upload a block of data for the dataset. The maximum size of the block is 8MiB.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 3.2
//   - id - The identifier of the dataset.
//   - blockid - A valid Base64 string value that identifies the block. Prior to encoding, the string must be less than or equal
//     to 64 bytes in size. For a given blob, the length of the value specified for the blockid
//     parameter must be the same size for each block. Note that the Base64 string must be URL-encoded.
//   - options - DatasetsClientUploadBlockOptions contains the optional parameters for the DatasetsClient.UploadBlock method.
func (client *DatasetsClient) UploadBlock(ctx context.Context, id string, blockid string, body io.ReadSeekCloser, options *DatasetsClientUploadBlockOptions) (DatasetsClientUploadBlockResponse, error) {
	var err error
	req, err := client.uploadBlockCreateRequest(ctx, id, blockid, body, options)
	if err != nil {
		return DatasetsClientUploadBlockResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatasetsClientUploadBlockResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return DatasetsClientUploadBlockResponse{}, err
	}
	return DatasetsClientUploadBlockResponse{}, nil
}

// uploadBlockCreateRequest creates the UploadBlock request.
func (client *DatasetsClient) uploadBlockCreateRequest(ctx context.Context, id string, blockid string, body io.ReadSeekCloser, _ *DatasetsClientUploadBlockOptions) (*policy.Request, error) {
	host := "{endpoint}/speechtotext/v3.2"
	host = strings.ReplaceAll(host, "{endpoint}", client.endpoint)
	urlPath := "/datasets/{id}/blocks"
	if id == "" {
		return nil, errors.New("parameter id cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{id}", url.PathEscape(id))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("blockid", blockid)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := req.SetBody(body, "application/octet-stream"); err != nil {
	return nil, err
}
;	return req, nil
}

