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

var jst = time.FixedZone("JST", 9*60*60)

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
		chromedp.WaitVisible(`#new_user`, chromedp.ByID),
		chromedp.SendKeys(`#user_email`, email, chromedp.ByID),
		chromedp.SendKeys(`#user_password`, password, chromedp.ByID),
		chromedp.Submit(`#new_user`, chromedp.ByID),
		chromedp.WaitVisible(`#adit-button-push`, chromedp.ByID),
	}

	hour := time.Now().In(jst).Hour()
	if hour >= 0 && hour < 5 {
		actions = append(actions, chromedp.Click(`#is_yakin`, chromedp.ByID))
	}

	actions = append(
		actions,
		chromedp.Click(`#adit-button-push`, chromedp.ByID),
		chromedp.WaitVisible(`#adit-button-wait`, chromedp.ByID),
		chromedp.WaitNotVisible(`#adit-button-wait`, chromedp.ByID),
	)

	if err := runActions(ctx, actions); err != nil {
		log.Println(err)
		return 1
	}

	return 0
}

func runActions(ctx context.Context, targets []chromedp.Action) error {
	actions := make([]chromedp.Action, 0)
	for _, target := range targets {
		actions = append(actions, target)
		actions = append(actions, chromedp.Sleep(3*time.Second))
	}
	return chromedp.Run(ctx, actions...)
}
