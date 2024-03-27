package executor

import (
	"context"
	"fmt"
	"io/fs"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/semanser/ai-coder/assets"
	"github.com/semanser/ai-coder/database"
)

func InitBrowser(db *database.Queries) error {
	browserContainerName := BrowserName()
	_, err := SpawnContainer(context.Background(), browserContainerName, "ghcr.io/go-rod/rod", db)

	if err != nil {
		return fmt.Errorf("failed to spawn container: %w", err)
	}

	return nil
}

func Content(url string, query string) (result string, err error) {
	page, err := loadPage(nil)

	if err != nil {
		return "", fmt.Errorf("Error loading page: %w", err)
	}

	pageRouter := page.HijackRequests()

	// Do not load any images or css files
	pageRouter.MustAdd("*", func(ctx *rod.Hijack) {
		// There're a lot of types you can use in this enum, like NetworkResourceTypeScript for javascript files
		// In this case we're using NetworkResourceTypeImage to block images
		if ctx.Request.Type() == proto.NetworkResourceTypeImage ||
			ctx.Request.Type() == proto.NetworkResourceTypeStylesheet ||
			ctx.Request.Type() == proto.NetworkResourceTypeFont ||
			ctx.Request.Type() == proto.NetworkResourceTypeMedia ||
			ctx.Request.Type() == proto.NetworkResourceTypeManifest ||
			ctx.Request.Type() == proto.NetworkResourceTypeOther {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})
	// since we are only hijacking a specific page, even using the "*" won't affect much of the performance
	go pageRouter.Run()

	err = page.Navigate(url)

	if err != nil {
		return "", fmt.Errorf("Error navigating to page: %w", err)
	}

	err = page.WaitDOMStable(time.Second*1, 5)

	if err != nil {
		return "", fmt.Errorf("Error waiting for page to stabilize: %w", err)
	}

	script, err := fs.ReadFile(assets.ScriptTemplates, "scripts/content.js")

	if err != nil {
		return "", fmt.Errorf("Error reading script: %w", err)
	}

	pageText, err := page.Eval(string(script))

	if err != nil {
		return "", fmt.Errorf("Error evaluating script: %w", err)
	}

	return pageText.Value.Str(), nil
}

func Screenshot(url string) ([]byte, error) {
	page, err := loadPage(nil)

	if err != nil {
		return nil, fmt.Errorf("Error loading page: %w", err)
	}

	err = page.Navigate(url)

	if err != nil {
		return nil, fmt.Errorf("Error navigating to page: %w", err)
	}

	screenshot, err := page.Screenshot(true, nil)

	if err != nil {
		return nil, fmt.Errorf("Error taking screenshot: %w", err)
	}

	return screenshot, nil
}

func BrowserName() string {
	return fmt.Sprintf("browser")
}

func loadPage(page *rod.Page) (*rod.Page, error) {
	u := launcher.MustResolveURL("")

	browser := rod.New().ControlURL(u)
	err := browser.Connect()

	if err != nil {
		return nil, fmt.Errorf("Error connecting to browser: %w", err)
	}

	page, err = browser.Page(proto.TargetCreateTarget{})

	if err != nil {
		return nil, fmt.Errorf("Error opening page: %w", err)
	}

	return page, nil
}
