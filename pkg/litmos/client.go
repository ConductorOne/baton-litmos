package litmos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/davecgh/go-spew/spew"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const pageSize = 100

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
	l := ctxzap.Extract(ctx)
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
	l.Debug("sending request", zap.String("url", url.String()), zap.String("method", method))
	resp, err := c.wrapper.Do(req, uhttp.WithXMLResponse(response))
	if err != nil && resp != nil {
		// Retry 503s & 504s because the Litmos API is flaky
		if resp.StatusCode == http.StatusGatewayTimeout || resp.StatusCode == http.StatusServiceUnavailable {
			return resp, status.Error(codes.Unavailable, resp.Status)
		}
	}
	return resp, err
}

type Team struct {
	Id                    string `xml:"Id"`
	Name                  string `xml:"Name"`
	TeamCodeForBulkImport string `xml:"TeamCodeForBulkImport"`
	ParentTeamId          string `xml:"ParentTeamId"`
}

type TeamsResp struct {
	Teams []Team `xml:"Team"`
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

type UsersResp struct {
	XMLName xml.Name `xml:"Users"`
	Users   []User   `xml:"User"`
}

func pageTokenToQuery(pToken *pagination.Token) *url.Values {
	query := &url.Values{}
	query.Add("limit", strconv.Itoa(pageSize))

	if pToken == nil || pToken.Token == "" {
		return query
	}

	_, err := strconv.Atoi(pToken.Token)
	if err != nil {
		fmt.Printf("error converting token %s to int: %v\n", pToken.Token, err)
		return query
	}
	query.Add("start", pToken.Token)

	return query
}

func getNextPageToken(pToken *pagination.Token, numItems int) string {
	if pToken == nil {
		return ""
	}

	if numItems < pageSize {
		// no more pages
		return ""
	}

	if pToken.Token == "" {
		return strconv.Itoa(numItems)
	}

	start, err := strconv.Atoi(pToken.Token)
	if err != nil {
		fmt.Printf("error converting token %s to int: %v\n", pToken.Token, err)
		return ""
	}

	return strconv.Itoa(start + numItems)
}

func (c *Client) ListUsers(ctx context.Context, pToken *pagination.Token) ([]User, string, error) {
	usersResp := UsersResp{}
	query := pageTokenToQuery(pToken)
	res, err := c.Do(ctx, "GET", "/v1.svc/users", query, &usersResp)
	if err != nil {
		spew.Dump(res.Body)
		return nil, pToken.Token, err
	}

	spew.Dump(usersResp)
	nextPageToken := getNextPageToken(pToken, len(usersResp.Users))
	return usersResp.Users, nextPageToken, nil
}

func (c *Client) ListTeams(ctx context.Context, pToken *pagination.Token) ([]Team, string, error) {
	teamsResp := TeamsResp{}
	query := pageTokenToQuery(pToken)
	res, err := c.Do(ctx, "GET", "/v1.svc/teams", query, &teamsResp)
	if err != nil {
		spew.Dump(res.Body)
		return nil, pToken.Token, err
	}

	spew.Dump(teamsResp)
	nextPageToken := getNextPageToken(pToken, len(teamsResp.Teams))
	return teamsResp.Teams, nextPageToken, nil
}
