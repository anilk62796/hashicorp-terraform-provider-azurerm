package configurations

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByServerResponse struct {
	HttpResponse *http.Response
	Model        *[]ServerConfiguration

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByServerResponse, error)
}

type ListByServerCompleteResult struct {
	Items []ServerConfiguration
}

func (r ListByServerResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByServerResponse) LoadMore(ctx context.Context) (resp ListByServerResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByServer ...
func (c ConfigurationsClient) ListByServer(ctx context.Context, id ServerId) (resp ListByServerResponse, err error) {
	req, err := c.preparerForListByServer(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServer", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServer", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByServer(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServer", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByServerComplete retrieves all of the results into a single object
func (c ConfigurationsClient) ListByServerComplete(ctx context.Context, id ServerId) (ListByServerCompleteResult, error) {
	return c.ListByServerCompleteMatchingPredicate(ctx, id, ServerConfigurationPredicate{})
}

// ListByServerCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ConfigurationsClient) ListByServerCompleteMatchingPredicate(ctx context.Context, id ServerId, predicate ServerConfigurationPredicate) (resp ListByServerCompleteResult, err error) {
	items := make([]ServerConfiguration, 0)

	page, err := c.ListByServer(ctx, id)
	if err != nil {
		err = fmt.Errorf("loading the initial page: %+v", err)
		return
	}
	if page.Model != nil {
		for _, v := range *page.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	for page.HasMore() {
		page, err = page.LoadMore(ctx)
		if err != nil {
			err = fmt.Errorf("loading the next page: %+v", err)
			return
		}

		if page.Model != nil {
			for _, v := range *page.Model {
				if predicate.Matches(v) {
					items = append(items, v)
				}
			}
		}
	}

	out := ListByServerCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByServer prepares the ListByServer request.
func (c ConfigurationsClient) preparerForListByServer(ctx context.Context, id ServerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/configurations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByServerWithNextLink prepares the ListByServer request with the given nextLink token.
func (c ConfigurationsClient) preparerForListByServerWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
	uri, err := url.Parse(nextLink)
	if err != nil {
		return nil, fmt.Errorf("parsing nextLink %q: %+v", nextLink, err)
	}
	queryParameters := map[string]interface{}{}
	for k, v := range uri.Query() {
		if len(v) == 0 {
			continue
		}
		val := v[0]
		val = autorest.Encode("query", val)
		queryParameters[k] = val
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByServer handles the response to the ListByServer request. The method always
// closes the http.Response Body.
func (c ConfigurationsClient) responderForListByServer(resp *http.Response) (result ListByServerResponse, err error) {
	type page struct {
		Values   []ServerConfiguration `json:"value"`
		NextLink *string               `json:"nextLink"`
	}
	var respObj page
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	result.Model = &respObj.Values
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByServerResponse, err error) {
			req, err := c.preparerForListByServerWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServer", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServer", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByServer(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "ListByServer", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
