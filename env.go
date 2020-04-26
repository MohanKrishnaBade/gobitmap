package main

import "os"

func config() {
	os.Setenv("NEWS_SECRET_KEY", "46dc19316c914cb79c56acb8e80c418f")
	os.Setenv("NEW_APP_URL", "https://newsapi.org")
}
