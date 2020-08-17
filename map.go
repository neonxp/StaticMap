/*
Copyright © 2020 Alexander Kiryukhin <ak@bytechain.ru>
This file is part of OsmStatic project.
*/
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sync"

	"github.com/disintegration/imaging"
)

const (
	tileAddr = `https://a.tile.openstreetmap.org/%d/%d/%d.png`
	tw       = 256
	th       = 256
)

func GetMapImage(lat, lon float64, zoom, width, height int) ([]byte, error) {
	x, y, dx, dy := getCoords(lat, lon, zoom)
	dst := imaging.New(width, height, color.NRGBA{0, 255, 0, 255})
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	cx := width/2 - tw/2
	cy := height/2 - th/2
	sx := width/tw + 2
	sy := height/th + 2
	di := int(dx*float64(tw)) - tw/2
	dj := int(dy*float64(th)) - th/2
	for i := -sx / 2; i <= sx/2; i++ {
		for j := -sy / 2; j <= sy/2; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				tile, err := getTile(x+i, y+j, zoom)
				if err != nil {
					log.Println(err)
					return
				}
				img, err := png.Decode(bytes.NewReader(tile))
				if err != nil {
					log.Println(err)
					return
				}
				mu.Lock()
				defer mu.Unlock()
				tx := cx + i*tw - di
				ty := cy + j*th - dj
				dst = imaging.Paste(dst, img, image.Pt(tx, ty))
			}(int(i), int(j))
		}
	}
	wg.Wait()

	out := bytes.NewBuffer([]byte{})
	if err := imaging.Encode(out, dst, imaging.PNG); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func getCoords(lat, lon float64, zoom int) (int, int, float64, float64) {
	x := (lon + 180.0) / 360.0 * (math.Exp2(float64(zoom)))
	y := (1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(zoom)))
	dx := x - math.Floor(x)
	dy := y - math.Floor(y)
	return int(math.Floor(x)), int(math.Floor(y)), dx, dy
}

func getTile(x, y, z int) ([]byte, error) {
	tile := fmt.Sprintf(tileAddr, z, x, y)
	resp, err := http.DefaultClient.Get(tile)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
