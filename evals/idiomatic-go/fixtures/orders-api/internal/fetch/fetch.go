package fetch

import (
	"io"
	"net/http"
)

// FetchAll downloads all urls concurrently and returns their bodies.
func FetchAll(urls []string) ([][]byte, error) {
	type result struct {
		body []byte
		err  error
	}
	ch := make(chan result)
	for _, u := range urls {
		go func(u string) {
			resp, err := http.Get(u)
			if err != nil {
				ch <- result{err: err}
				return
			}
			b, err := io.ReadAll(resp.Body)
			ch <- result{body: b, err: err}
		}(u)
	}
	var out [][]byte
	for range urls {
		r := <-ch
		if r.err != nil {
			return nil, r.err
		}
		out = append(out, r.body)
	}
	return out, nil
}
