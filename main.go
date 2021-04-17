package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	email := os.Getenv("JOBCAN_EMAIL")
	if len(email) == 0 {
		panic("no JOBCAN_EMAIL specified")
	}
	password := os.Getenv("JOBCAN_PASSWORD")
	if len(password) == 0 {
		panic("no JOBCAN_PASSWORD specified")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	actions := []chromedp.Action{
		chromedp.Navigate("https://id.jobcan.jp/users/sign_in?app_key=atd&redirect_to=https://ssl.jobcan.jp/jbcoauth/callback"),
		chromedp.WaitVisible(`#new_user`, chromedp.ByID),
		chromedp.SendKeys(`#user_email`, email, chromedp.ByID),
		chromedp.SendKeys(`#user_password`, password, chromedp.ByID),
		chromedp.Submit(`#new_user`, chromedp.ByID),
		chromedp.WaitVisible(`#adit-button-push`, chromedp.ByID),
		chromedp.Click(`#adit-button-push`, chromedp.ByID),
		chromedp.WaitVisible(`#adit-button-wait`, chromedp.ByID),
		chromedp.WaitNotVisible(`#adit-button-wait`, chromedp.ByID),
	}

	if err := chromedp.Run(ctx, actions...); err != nil {
		panic(err)
	}
}
