package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {

	localIp := getLocalIP()
	app := initialiseApp()

	ntfyKey := os.Getenv("BIN_JUICE_NTFY_KEY")
	if ntfyKey != "" {
		fmt.Println("Found ntfy key, starting notification looper")
		go notificationLooper(ntfyKey)
	} else {
		fmt.Println("No ntfy key found, ignoring notification looper")
	}
	app.Listen(localIp + ":8000")
}

func initialiseApp() *fiber.App {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", handleIndex)
	app.Static("/", "./public")
	return app
}

func notificationLooper(ntfyKey string) {
	for range time.Tick(12 * time.Second) {
		fmt.Println("Starting notification looper")
		ncData := NextCollectionData()
		ncDate := NextCollectionDate(time.Now())
		daysToNc := math.Ceil(time.Until(ncDate).Hours() / 24)

		if daysToNc <= 2 {
			fmt.Println("Sending notification about next collection.")
			var whichBins []string
			for _, bin := range ncData {
				if bin.Collected {
					whichBins = append(whichBins, bin.FriendlyName)
				}
			}
			whichBinsString := strings.Join(whichBins, ", ")
			daysString := fmt.Sprintf("Bin Alert! (%v day(s))", daysToNc)
			req, _ := http.NewRequest("POST", "https://ntfy.sh/"+ntfyKey, strings.NewReader(whichBinsString))
			req.Header.Set("Title", daysString)
			req.Header.Set("Tags", "wastebasket")
			http.DefaultClient.Do(req)
		} else {
			fmt.Printf("Not close enough to collection date (%v days away) to notify. Skipping.", daysToNc)
		}
	}
}

func handleIndex(c *fiber.Ctx) error {
	ncData := NextCollectionData()
	ncDate := NextCollectionDate(time.Now())
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

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Print("Problem getting network interfaces")
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "n/a"
}
