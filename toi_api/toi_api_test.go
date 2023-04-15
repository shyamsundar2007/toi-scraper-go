package toi_api

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type mockHTTPClient struct{}

type defaultHTTPClient struct{}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func (c *mockHTTPClient) Get(url string) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("Test response")),
	}
	return resp, nil
}

func (c *defaultHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

type mockGoQueryClient struct{}

type defaultGoQueryClient struct{}

func (c *mockGoQueryClient) NewDocumentFromReader(_ io.ReadCloser) (*goquery.Document, error) {
	httpClient := mockHTTPClient{}
	resp, err := httpClient.Get("string")
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(resp.Body)
}

func (c *defaultGoQueryClient) NewDocumentFromReader(resp io.ReadCloser) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(resp)
}

func TestToiMovieApi_GetMovieReviews(t *testing.T) {
	type args struct {
		language Language
	}
	tests := []struct {
		name string
		args args
		want []*MovieReview
	}{
		// TODO: Add test cases.
		{
			"testReturnsNil",
			args{
				language: Tamil,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toiMovieApi := NewToiMovieApi(&mockHTTPClient{}, &mockGoQueryClient{})
			if got, _ := toiMovieApi.GetMovieReviews(tt.args.language); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMovieReviews() = %v, want %v", got, tt.want)
			}
		})
	}
}
