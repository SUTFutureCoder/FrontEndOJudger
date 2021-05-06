package utils

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"math"
	url2 "net/url"
	"strings"
	"time"
)

func ScreenshotTasks(url string, cookies string, imageBuf *interface{}, width int64, height int64, waitBefore int) chromedp.Tasks {
	tasks := make([]chromedp.Action, 0)
	
	// set cookies
	if cookies != "" {
		cookieList := strings.Split(cookies, ";")
		u, _ := url2.Parse(url)
		for _, c := range cookieList {
			c = strings.TrimSpace(c)
			cookieKv := strings.SplitN(c, "=", 2)
			if len(cookieKv) != 2 {
				continue
			}
			tasks = append(tasks, SetCookie(cookieKv[0], cookieKv[1], u.Host, "/", false, false))
		}
	}

	tasks = append(tasks, chromedp.Tasks{
		chromedp.Navigate(url),
		//chromedp.ActionFunc(func(ctx context.Context) (err error) {
		//	*imageBuf, err = page.CaptureScreenshot().WithQuality(90).Do(ctx)
		//	return err
		//}),
		chromedp.EmulateViewport(width, height, chromedp.EmulateScale(1)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).Do(ctx)
			if err != nil {
				return err
			}

			*imageBuf, err = page.CaptureScreenshot().
				WithQuality(100).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)

			if err != nil {
				return err
			}
			return nil
		}),
	})
	if waitBefore != 0 {
		tasks = append(tasks, chromedp.Sleep(time.Duration(waitBefore)*time.Millisecond))
	}
	return tasks
}

func SetCookie(name, value, domain, path string, httpOnly, secure bool) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
		success, err := network.SetCookie(name, value).
			WithExpires(&expr).
			WithDomain(domain).
			WithPath(path).
			WithHTTPOnly(httpOnly).
			WithSecure(secure).
			Do(ctx)
		if err != nil {
			return err
		}
		if !success {
			return fmt.Errorf("could not set cookie %s", name)
		}
		return nil
	})
}
