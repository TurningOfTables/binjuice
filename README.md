# Bin Juice

## Description

Uses Go with the Fiber framework and Bootstrap to generate a basic html page showing what bins are next being collected in my area. Contains unnecessary CSS animations.

![Screenshot](readme_screenshot.png)

## Notes

Run with `go run main.go binlib.go`

By default binds to the current IP address on port 8000

If you set an env variable of `BIN_JUICE_NTFY_TOPIC` with a [ntfy.sh](https://ntfy.sh) topic then you can get notifications as the bin collection day nears. See the ntfy.sh website for more info.

## Attributions

[Fiber web framework](https://gofiber.io)

[Bootstrap](https://getbootstrap.com)

[Blue recycling bin icon created by Smashicons - Freepik](https://www.freepik.com/icon/bin_10509062)

[Garden waste bin icon created by Smashicons - Freepik](https://www.freepik.com/icon/bin_6303004)

[Grey bin bag icon created by Smashicons - Freepik](https://www.freepik.com/icon/bag_10722122)

[Food waste icon created by Smashicons - Flaticon](https://www.flaticon.com/free-icons/food-waste)

[Red cross icon by Pixel Perfect - Flaticon](https://www.flaticon.com/free-icons/delete)

[Green tick icon by kliriw art - Flaticon](https://www.flaticon.com/free-icons/correct)

[Optional push notifications by ntfy.sh](https://ntfy.sh/)