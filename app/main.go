package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/playwright-community/playwright-go"
)

func assertErrorToNilf(message string, err error) {
	if err != nil {
		log.Fatalf("%v: %v", message, err)
	}
}

type Theme struct {
	BarWidth int 
	Green    string
	Blue     string
	Yellow   string
	Magenta  string
	White    string
}
func initTheme() *Theme {
	return &Theme{
		BarWidth: 50,
		Green: "\033[32m",
		Blue: "\033[34m",
		Yellow: "\033[33m",
		Magenta: "\033[35m",
		White: "\033[97m",
	}
}

var theme *Theme = initTheme()

type Options struct {
	Headless        bool
	UserAgent       string
	SleepDuration   int
	ActionTimeout   int
	DownloadTimeout int
	MaxWorkers      int    
	ConvertorUrl    string
	DownloadDir     string
	SpotifyUrl      string 
}
func initOptions() *Options {
	var options Options 
	flag.BoolVar(
		&options.Headless,
		"headless",
		true,
		"hide browser windows or not",
	)
	flag.StringVar(
		&options.UserAgent,
		"uagent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"provide custom user agent",
	)
	flag.IntVar(
		&options.SleepDuration,
		"sleep",
		1000,
		"time in ms for the program to sleep between repeating actions",
	)
	flag.IntVar(
		&options.ActionTimeout,
		"atimeout",
		3000,
		"time in ms to wait for an action preformed on a locator",
	)
	flag.IntVar(
		&options.DownloadTimeout,
		"dtimeout",
		5000,
		"time in ms to wait for a download before retry",
	)
	flag.IntVar(
		&options.MaxWorkers,
		"workers",
		5,
		"maximum number of browser contexts attempting download simultaneously",
	)
	flag.StringVar(
		&options.ConvertorUrl,
		"convertor",
		"https://spotidown.app/",
		"url of a spotify to audio converotor whitout captcha",
	)
	flag.StringVar(
		&options.DownloadDir,
		"dir",
		"./",
		"directory on this pc to download the track/s to",
	)
	flag.StringVar(
		&options.SpotifyUrl,
		"download",
		"",
		"spotify url of spotify track or album/playlist",
	)
	flag.Parse()

	if options.SpotifyUrl == "" {
		log.Fatalf("no spotify url provided set -download flag or -help for help")
	} else if !strings.Contains(options.SpotifyUrl, "https://open.spotify.com/") {
		log.Fatalf("invalid spotify url provided: %v", options.SpotifyUrl)
	}
	return &options
}
func (options *Options) SetDownloadDirName(name string) {
	options.DownloadDir = filepath.Join(options.DownloadDir, name)
}

func displayProgressBar(progress, total int) {
	decimal := float32(progress) / float32(total)
	barLen := int(decimal * float32(theme.BarWidth))
	var bar strings.Builder
	bar.WriteRune('\r')
	bar.WriteString(theme.Blue)
	bar.WriteString(strings.Repeat("=", barLen))
	bar.WriteString(strings.Repeat(" ", theme.BarWidth-barLen))
	bar.WriteString(fmt.Sprintf(
		" %v%v/%v %v(%.1f%%)%v", 
		theme.Green,
		progress, 
		total, 
		theme.Magenta,
		decimal*100,
		theme.White,
	))
	fmt.Print(bar.String())
}

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func initBrowser(options *Options) (browser playwright.Browser, close func()) {
	pw, err := playwright.Run()
	assertErrorToNilf("failed to start playwright", err)
	browser, err = pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(options.Headless),
			Args: []string{
				"--disable-background-timer-throttling",
				"--disable-backgrounding-occluded-windows",
				"--disable-gpu",
				"--no-sandbox",
			},
		},
	)
	assertErrorToNilf("failed to launch browser", err)
	close = func() {
		browser.Close()
		pw.Stop()
	}
	return
}

func newContextPage(browser playwright.Browser, options *Options, url string) (page playwright.Page, close func()) {
	ctxt, err := browser.NewContext(
		playwright.BrowserNewContextOptions{
			Viewport: &playwright.Size{
				Width:  960,
				Height: 720,
			},
			UserAgent: playwright.String(options.UserAgent),
		},
	)
	assertErrorToNilf("failed to launch new browserContext", err)
	page, err = ctxt.NewPage()
	assertErrorToNilf("failed to create page", err)
	_, err = page.Goto(url)
	assertErrorToNilf("failed to find url", err)
	close = func() {
		page.Close()
		ctxt.Close()
	}
	return
}

func scrapeTracks(browser playwright.Browser, options *Options) (trackLinks []string) {
	page, closeContext := newContextPage(browser, options, options.SpotifyUrl)
	defer closeContext()
	notFound := page.Locator("h1")
	if text, _:=notFound.InnerText(); text == "Page not available" {
		log.Fatalf("invalid spotify album/playlist url: %v", options.SpotifyUrl)
	}
	cookiesBtn := page.Locator("#onetrust-accept-btn-handler")
	cookiesBtn.Click()

	listName, _ := page.Locator(".encore-text-headline-large").InnerText()
	options.SetDownloadDirName(listName)
	fmt.Printf("Getting track urls from %v%v%v\n", theme.Yellow, listName, theme.White)

	tracksStr, _ := page.Locator("div.GI8QLntnaSCh2ONX_y2c > span:nth-child(1)").InnerText()
	tracksCount, _ := strconv.Atoi(strings.TrimSuffix(tracksStr, " songs"))

	var (
		trackElements []playwright.Locator
		loadedLen     int
	)
	for {
		trackElements, _ = page.Locator(`div:nth-child(2) > div[role="presentation"]:nth-child(2) > div a.btE2c3IKaOXZ4VNAb8WQ`).All()
		for _, element := range trackElements {
			trackLink, _ := element.GetAttribute("href")
			if !slices.Contains(trackLinks, trackLink) {
				trackLinks = append(trackLinks, trackLink)
			}
		}

		displayProgressBar(len(trackLinks), tracksCount)
		if len(trackLinks) >= tracksCount {
			fmt.Println()
			return
		}

		loadedLen = len(trackElements)
		if loadedLen > 0 {
			trackElements[loadedLen-1].ScrollIntoViewIfNeeded()
			sleep(options.SleepDuration)
		}
	}
}

func adBlock(page playwright.Page, options *Options) {
	page.OnDialog(func(dialog playwright.Dialog) { dialog.Dismiss() })
	page.Route("**/*.{jpg,png,gif,css}", func(route playwright.Route) { route.Abort() })
	btns := []playwright.Locator{
		page.Locator(`iframe[name="aswift_10"]`).
			ContentFrame().
			Locator(`iframe[name="ad_iframe"]`).
			ContentFrame().
			GetByRole(`button`, playwright.FrameLocatorGetByRoleOptions{Name: `Close ad`}),
		page.Locator(`iframe[name="aswift_9"]`).
			ContentFrame().
			Locator(`iframe[name="ad_iframe"]`).
			ContentFrame().
			GetByRole(`button`, playwright.FrameLocatorGetByRoleOptions{Name: `Close ad`}),
		page.Locator(`iframe[name="aswift_10"]`).
			ContentFrame().
			GetByRole("button", playwright.FrameLocatorGetByRoleOptions{Name: `Close ad`}),
		page.Locator(`iframe[name="aswift_9"]`).
			ContentFrame().
			GetByRole(`button`, playwright.FrameLocatorGetByRoleOptions{Name: `Close ad`}),
		page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Close"}),
	}
	for !page.IsClosed() {
		for i := range btns {
			if count, _ := btns[i].Count(); count > 0 {
				btns[i].Click()
			}
		}
		sleep(options.SleepDuration)
	}
}
func saveTrack(downloadEvent playwright.Download, options *Options) error {
	suggestedName := downloadEvent.SuggestedFilename()
	if length := len(suggestedName); length >= 4 && suggestedName[length-4:] != ".mp3" {
		defer downloadEvent.Cancel()
		return fmt.Errorf("unsupported file format %s", suggestedName)
	}
	fileName := strings.TrimPrefix(suggestedName, "SpotiDown.App - ")
	filePath := filepath.Join(options.DownloadDir, fileName)
	return downloadEvent.SaveAs(filePath)
}

func download(browser playwright.Browser, options *Options, trackUrl string) {
	page, closeContext := newContextPage(browser, options, options.ConvertorUrl)
	defer closeContext()
	go adBlock(page, options)
	page.BringToFront()

	consentBtn := page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Do not consent"})
	urlInput := page.Locator("#url")
	submitBtn := page.Locator("#send")
	notFound := page.Locator("#alert")
	downloadBtn := page.GetByRole("link", playwright.PageGetByRoleOptions{Name: "Download Mp3"})
	tryDownload:
		consentBtn.Click(playwright.LocatorClickOptions{Timeout: playwright.Float(float64(options.ActionTimeout))})
		urlInput.Fill(trackUrl, playwright.LocatorFillOptions{Timeout: playwright.Float(float64(options.ActionTimeout))})
		submitBtn.Click(playwright.LocatorClickOptions{Timeout: playwright.Float(float64(options.ActionTimeout))})
		if visible, _ := notFound.IsVisible(); visible {
			log.Printf(`invalid track url: "%v"`, trackUrl)
		}
		// page.Pause()

		downloadEvent, err := page.ExpectDownload(
			func() error {
				err := downloadBtn.Click(playwright.LocatorClickOptions{Timeout: playwright.Float(float64(options.ActionTimeout))})
				return err
			},
			playwright.PageExpectDownloadOptions{Timeout: playwright.Float(float64(options.DownloadTimeout))},
		)
		if err != nil || saveTrack(downloadEvent, options) != nil {
			page.Goto(options.ConvertorUrl)
			goto tryDownload
		}
		// page.Pause()
}

func resolveTrack(browser playwright.Browser, options *Options) {
	fmt.Println("Downloading track:")
	displayProgressBar(0, 1)
	download(browser, options, options.SpotifyUrl)
	displayProgressBar(1, 1)
}

func resolveList(browser playwright.Browser, options *Options) {
	links := scrapeTracks(browser, options)
	fmt.Println("Downloading tracks:")
	sem := make(chan struct{}, options.MaxWorkers)
	total := len(links)
	awaiting := total
	displayProgressBar(0, total)
	var wg sync.WaitGroup
	wg.Add(total)
	for _, link := range links {
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			download(browser, options, "https://open.spotify.com"+link)
			<-sem
			awaiting -= 1
			displayProgressBar(total-awaiting, total)
		}()
	}
	wg.Wait()
}

func main() {
	options := initOptions()

	browser, closeBrowser := initBrowser(options)
	defer closeBrowser()

	var funcToCall func(playwright.Browser, *Options)
	if strings.Contains(options.SpotifyUrl, "/track") {
		funcToCall = resolveTrack
	} else if strings.Contains(options.SpotifyUrl, "/album") || strings.Contains(options.SpotifyUrl, "/playlist") {
		funcToCall = resolveList
	} else {
		log.Fatalf("invalid spotify url!")
	}
	startTime := time.Now()
	funcToCall(browser, options)
	fmt.Printf("\nResolved in %.3fs\n", time.Since(startTime).Seconds())
}