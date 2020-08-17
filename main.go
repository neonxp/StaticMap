/*
Copyright © 2020 Alexander Kiryukhin <a.kiryukhin@mail.ru>
This file is part of StaticMap project.
*/
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/neonxp/StaticMap/pkg/static"
	"log"
	"strconv"
)

func main() {
	e := echo.New()
	e.GET("/map", func(c echo.Context) error {
		lat, err := strconv.ParseFloat(c.QueryParam("lat"), 64)
		if err != nil {
			return err
		}
		lon, err := strconv.ParseFloat(c.QueryParam("lon"), 64)
		if err != nil {
			return err
		}
		w, err := strconv.Atoi(c.QueryParam("w"))
		if err != nil {
			w = 800
		}
		h, err := strconv.Atoi(c.QueryParam("h"))
		if err != nil {
			h = 800
		}
		zoom, err := strconv.Atoi(c.QueryParam("zoom"))
		if err != nil {
			zoom = 16
		}
		if zoom < 1 {
			zoom = 1
		}
		if zoom > 20 {
			zoom = 20
		}
		img, err := static.GetMapImage(lat, lon, zoom, w, h)
		if err != nil {
			return err
		}
		return c.Blob(200, "image/png", img)
	})
	log.Fatal(e.Start(":8000"))
}
