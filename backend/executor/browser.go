package executor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/semanser/ai-coder/assets"
	"github.com/semanser/ai-coder/database"
	"github.com/semanser/ai-coder/templates"
)

var (
	browser *rod.Browser
)

const port = "9222"

func InitBrowser(db *database.Queries) error {
	browserContainerName := BrowserName()
	portBinding := nat.Port(fmt.Sprintf("%s/tcp", port))

	_, err := SpawnContainer(context.Background(), browserContainerName, &container.Config{
		Image: "ghcr.io/go-rod/rod",
		ExposedPorts: nat.PortSet{
			portBinding: struct{}{},
		},
		Cmd: []string{"chrome", "--headless", "--no-sandbox", fmt.Sprintf("--remote-debugging-port=%s", port), "--remote-debugging-address=0.0.0.0"},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			portBinding: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
				},
			},
		},
	}, db)

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

	script, err := templates.Render(assets.ScriptTemplates, "scripts/content.js", nil)

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
	return fmt.Sprintf("codel-browser")
}

func loadPage(page *rod.Page) (*rod.Page, error) {
	u, err := launcher.ResolveURL("")

	if err != nil {
		return nil, fmt.Errorf("Error resolving url: %w", err)
	}

	browser := rod.New().ControlURL(u)
	err = browser.Connect()

	version, err := browser.Version()

	if err != nil {
		return nil, fmt.Errorf("Error getting browser version: %w", err)
	}
	log.Println("Connected to browser %s", version.Product)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to browser: %w", err)
	}

	page, err = browser.Page(proto.TargetCreateTarget{})

	if err != nil {
		return nil, fmt.Errorf("Error opening page: %w", err)
	}

	return page, nil
}
