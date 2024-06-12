package litmos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/davecgh/go-spew/spew"
)

type Client struct {
	wrapper *uhttp.BaseHttpClient
	apiKey  string
	source  string
}

func NewClient(ctx context.Context, apiKey, source string) (*Client, error) {
	options := []uhttp.Option{uhttp.WithLogger(true, nil)}

	httpClient, err := uhttp.NewClient(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("creating HTTP client failed: %w", err)
	}
	wrapper := uhttp.NewBaseHttpClient(httpClient)

	return &Client{
		wrapper: wrapper,
		apiKey:  apiKey,
		source:  source,
	}, nil
}

func (c *Client) Do(ctx context.Context, method string, path string, query *url.Values, response interface{}, options ...uhttp.RequestOption) (*http.Response, error) {
	options = append(options,
		uhttp.WithHeader("apikey", c.apiKey), uhttp.WithAcceptXMLHeader(),
	)

	rawQuery := ""
	if query != nil {
		rawQuery = query.Encode()
	}
	url := &url.URL{
		Scheme:   "https",
		Host:     "api.litmos.com",
		Path:     path,
		RawQuery: rawQuery,
	}
	q := url.Query()
	q.Add("source", c.source)
	url.RawQuery = q.Encode()

	req, err := c.wrapper.NewRequest(ctx, method, url, options...)
	if err != nil {
		return nil, err
	}
	return c.wrapper.Do(req, uhttp.WithXMLResponse(response))
}

type User struct {
	Id          string `xml:"Id"`
	UserName    string `xml:"UserName"`
	FirstName   string `xml:"FirstName"`
	LastName    string `xml:"LastName"`
	Active      bool   `xml:"Active"`
	Email       string `xml:"Email"`
	AccessLevel string `xml:"AccessLevel"`
	Brand       string `xml:"Brand"`
}

type PaginationInfo struct {
	BatchParam string `xml:"BatchParam"`
	BatchSize  int    `xml:"BatchSize"`
	Start      int    `xml:"Start"`
	TotalCount int    `xml:"TotalCount"`
}

type UserItems struct {
	Users []User `xml:"User"`
}
type UserCollection struct {
	Pagination PaginationInfo `xml:"Pagination"`
	Items      UserItems      `xml:"Items"`
}

func (c *Client) ListUsers(ctx context.Context, pToken *pagination.Token) ([]User, error) {
	// TODO: figure out query args for pagination
	userCollection := UserCollection{}
	res, err := c.Do(ctx, "GET", "/v1.svc/users/paginated", nil, &userCollection)
	if err != nil {
		return nil, err
	}

	spew.Dump(res.Body)
	spew.Dump(userCollection)
	return nil, nil
}
