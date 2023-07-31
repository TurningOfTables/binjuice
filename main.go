package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {

	localIp := getLocalIP()
	app := initialiseApp()
	app.Listen(localIp + ":8000")
}

func initialiseApp() *fiber.App {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", handleIndex)
	app.Static("/", "./public")
	return app
}

func handleIndex(c *fiber.Ctx) error {
	ncData := NextCollectionData()
	ncDate := NextCollectionDate(time.Now())
	daysToNc := time.Until(ncDate).Hours() / 24
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
