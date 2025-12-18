package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Kullanım: go run main.go <hedef-url>")
		return
	}
	targetURL := os.Args[1]

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal("[-] Hatalı URL")
	}
	hostName := parsedURL.Host
	if hostName == "" {
		hostName = strings.Split(parsedURL.Path, "/")[0]
	}
	cleanFileName := strings.ReplaceAll(hostName, ".", "_")

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Flag("headless", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var htmlContent string
	var imageBuffer []byte

	fmt.Printf("[+] %s kaydediliyor\n", targetURL)

	err = chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(4*time.Second),
		chromedp.OuterHTML(`html`, &htmlContent),
		chromedp.FullScreenshot(&imageBuffer, 90),
	)

	if err != nil {
		log.Fatalf("[-] Bir hata oluştu: %v", err)
	}

	htmlOutput := cleanFileName + ".html"
	imgOutput := cleanFileName + "_screenshot.png"

	_ = os.WriteFile(htmlOutput, []byte(htmlContent), 0644)
	_ = os.WriteFile(imgOutput, imageBuffer, 0644)

	fmt.Printf("kaydedildi:\n - %s\n - %s\n", htmlOutput, imgOutput)
}
