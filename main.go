package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	os.Exit(run())
}

func run() int {
	email := os.Getenv("JOBCAN_EMAIL")
	if len(email) == 0 {
		log.Println("no JOBCAN_EMAIL specified")
		return 1
	}
	password := os.Getenv("JOBCAN_PASSWORD")
	if len(password) == 0 {
		log.Println("no JOBCAN_PASSWORD specified")
		return 1
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	actions := []chromedp.Action{
		chromedp.Navigate("https://id.jobcan.jp/users/sign_in?app_key=atd&redirect_to=https://ssl.jobcan.jp/jbcoauth/callback"),
		chromedp.Sleep(3 * time.Second),
		chromedp.WaitVisible(`#new_user`, chromedp.ByID),
		chromedp.Sleep(3 * time.Second),
		chromedp.SendKeys(`#user_email`, email, chromedp.ByID),
		chromedp.Sleep(3 * time.Second),
		chromedp.SendKeys(`#user_password`, password, chromedp.ByID),
		chromedp.Sleep(3 * time.Second),
		chromedp.Submit(`#new_user`, chromedp.ByID),
		chromedp.Sleep(3 * time.Second),
		chromedp.WaitVisible(`#adit-button-push`, chromedp.ByID),
		chromedp.Sleep(3 * time.Second),
		chromedp.Click(`#adit-button-push`, chromedp.ByID),
		chromedp.Sleep(3 * time.Second),
		chromedp.WaitVisible(`#adit-button-wait`, chromedp.ByID),
		chromedp.Sleep(3 * time.Second),
		chromedp.WaitNotVisible(`#adit-button-wait`, chromedp.ByID),
		chromedp.Sleep(3 * time.Second),
	}

	if err := chromedp.Run(ctx, actions...); err != nil {
		log.Println(err)
		return 1
	}

	return 0
}
