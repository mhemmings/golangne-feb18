package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var twClient *twitter.Client
var gopherImage image.Image

func main() {
	gopherFile, err := os.Open("gopher.png")
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}
	gopherImage, err = png.Decode(gopherFile)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer gopherFile.Close()

	config := oauth1.NewConfig(os.Getenv("TW_CONSUMER_KEY"), os.Getenv("TW_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TW_ACCESS_TOKEN"), os.Getenv("TW_ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	twClient = twitter.NewClient(httpClient)

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Printf("Listening on port %s", port)

	// Fire up server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Get user from Twitter
	user, _, err := twClient.Users.Show(&twitter.UserShowParams{
		ScreenName: r.URL.Path[1:],
	})
	if err != nil {
		fmt.Fprintf(w, "failed to get user: %s [%v]", r.URL.Path[1:], err)
		return
	}

	// Get fullsize twitter image.
	// TODO: Edgecase where this will break if filename has "_normal" in
	imgUrl := strings.Replace(user.ProfileImageURLHttps, "_normal", "", -1)

	resp, err := http.Get(imgUrl)
	if err != nil {
		fmt.Fprintf(w, "failed to get user image: %s [%v]", imgUrl, err)
		return
	}
	defer resp.Body.Close()

	var inputImage image.Image
	if strings.ToLower(path.Ext(imgUrl)) == ".png" {
		inputImage, err = png.Decode(resp.Body)
	} else {
		// Just fallback to jpeg encoding
		inputImage, err = jpeg.Decode(resp.Body)
	}
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}

	b := inputImage.Bounds()
	offset := image.Pt(b.Max.X-170, b.Max.Y-110)
	outputImage := image.NewRGBA(b)
	draw.Draw(outputImage, b, inputImage, image.ZP, draw.Src)
	draw.Draw(outputImage, gopherImage.Bounds().Add(offset), gopherImage, image.ZP, draw.Over)

	// Send to client
	// TODO: Improve headers etc
	jpeg.Encode(w, outputImage, &jpeg.Options{jpeg.DefaultQuality})
}
