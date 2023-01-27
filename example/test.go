package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bastianccm/future"
)

func fetch(ctx context.Context, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("unable to fetch %q: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %w", err)
	}

	return string(body), nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	de, en, fr, err := future.Resolve3(
		future.Promise(ctx, func(ctx context.Context) (string, error) {
			return fetch(ctx, "https://de.wikipedia.org/wiki/Go_(Programmiersprache)")
		}),
		future.Promise(ctx, func(ctx context.Context) (string, error) {
			return fetch(ctx, "https://en.wikipedia.org/wiki/Go_(programming_language)")
		}),
		future.Promise(ctx, func(ctx context.Context) (string, error) {
			return fetch(ctx, "https://fr.wikipedia.org/wiki/Go_(langage)")
		}),
	)

	if err != nil {
		fmt.Printf("error fetching golang article: %s\n", err)
		return
	}

	fmt.Printf("Fetched all in %s:\n", time.Since(start))
	fmt.Printf("DE: %s\n", de[strings.Index(de, "<title>")+7:strings.Index(de, "</title>")])
	fmt.Printf("EN: %s\n", en[strings.Index(en, "<title>")+7:strings.Index(en, "</title>")])
	fmt.Printf("FR: %s\n", fr[strings.Index(fr, "<title>")+7:strings.Index(fr, "</title>")])
}
