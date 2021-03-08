package caroline

import (
	"context"
	"encoding/json"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)

type screenshotReq struct {
	Url string `json:"url"`
	Height int `json:"height"`
	Width int `json:"width"`
	Scale float64 `json:"scale"`

}

func ScreenShot(w http.ResponseWriter, req *http.Request) {

	var screenReq screenshotReq
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&screenReq)
	// Start Chrome
	// Remove the 2nd param if you don't need debug information logged
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	url := "https://www.bilibili.com/"
	filename := "golangcode.png"

	// Run Tasks
	// List of actions to run in sequence (which also fills our image buffer)
	var imageBuf []byte
	if err := chromedp.Run(ctx, ScreenshotTasks(url, &imageBuf)); err != nil {
		log.Fatal(err)
	}

	// Write our image to file
	if err := ioutil.WriteFile(filename, imageBuf, 0644); err != nil {
		log.Fatal(err)
	}
}


func ScreenshotTasks(url string, imageBuf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		//chromedp.ActionFunc(func(ctx context.Context) (err error) {
		//	*imageBuf, err = page.CaptureScreenshot().WithQuality(90).Do(ctx)
		//	return err
		//}),
		chromedp.EmulateViewport(1400, 0, chromedp.EmulateScale(1)),

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
	}
}