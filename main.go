package main

import (
	"fmt"
	"log"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"binjuice/bin"
	"binjuice/notification"
)

var notifyHours = []int{8}

func main() {

	app := initialiseApp()
	go notificationLooper()
	log.Fatal(app.Listen(":8000"))
}

func initialiseApp() *fiber.App {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", handleIndex)
	app.Static("/", "./public")
	return app
}

func notificationLooper() {
	for range time.Tick(1 * time.Hour) {
		currentHour := time.Now().Hour()
		if !slices.Contains(notifyHours, currentHour) {
			fmt.Println("Not a notification hour - skipping")
			continue
		}
		fmt.Println("Starting notification looper")
		ncData := bin.NextCollectionData()
		ncDate := bin.NextCollectionDate(time.Now())
		daysToNc := math.Ceil(time.Until(ncDate).Hours() / 24)

		if daysToNc >= 2 {
			continue
		}

		fmt.Println("Sending notification about next collection.")
		var whichBins []string
		for _, bin := range ncData {
			if bin.Collected {
				whichBins = append(whichBins, bin.FriendlyName)
			}
		}
		whichBinsString := strings.Join(whichBins, ", ")

		daysString := fmt.Sprintf("Bin Alert! (%v day(s))", daysToNc)
		if daysToNc <= 0 {
			daysString = "Bin Alert! (today)"
		}

		msg := notification.Message{
			Title: daysString,
			Body:  whichBinsString,
			Tags:  "wastebasket",
		}

		if err := notification.Send(msg); err != nil {
			fmt.Println("Error sending notification")
		}

	}
}

func handleIndex(c *fiber.Ctx) error {
	ncData := bin.NextCollectionData()
	ncDate := bin.NextCollectionDate(time.Now())
	daysToNc := math.Ceil(time.Until(ncDate).Hours() / 24)
	var daysToNcString string
	if daysToNc == 0 {
		daysToNcString = "TODAY"
	} else if daysToNc == 1 {
		daysToNcString = "tomorrow"
	} else {
		daysToNcString = fmt.Sprintf("in %v days", daysToNc)
	}

	return c.Render("index", fiber.Map{"data": ncData, "nextCollectionDate": ncDate, "daysToNextCollection": daysToNcString})
}
