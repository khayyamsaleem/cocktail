package main

import (
	"strconv"
	"time"

	"github.com/warthog618/gpiod"
	// "github.com/warthog618/gpiod/device/rpi"
	"log"

	"github.com/gin-gonic/gin"
)

const GPIO_CHIP = "gpiochip0"

const HIGH = 1
const LOW = 0

const (
	LIQUOR0 int = 22
	LIQUOR1 int = 17
	LIQUOR2 int = 18
	LIQUOR3 int = 23
	LIQUOR4 int = 5
	WATER   int = 6
	DRAIN   int = 12
	JUICE   int = 16
	X1      int = 19
	X2      int = 20
	X3      int = 21
	X4      int = 26
)

var PINS = []int{
	LIQUOR0,
	LIQUOR1,
	LIQUOR2,
	LIQUOR3,
	LIQUOR4,
	WATER,
	DRAIN,
	JUICE,
}

func makePin(chip *gpiod.Chip, pinNum int) *gpiod.Line {
	pin, err := chip.RequestLine(pinNum, gpiod.AsOutput(LOW))
	if err != nil {
		panic(err)
	}
	return pin
}

func main() {
	chip, err := gpiod.NewChip(GPIO_CHIP, gpiod.WithConsumer("cocktail"))
	if err != nil {
		log.Fatalf("Error creating new chip interface: %s", err)
	}
	defer chip.Close()

	liquor0 := makePin(chip, LIQUOR0)
	liquor1 := makePin(chip, LIQUOR1)
	liquor2 := makePin(chip, LIQUOR2)
	liquor3 := makePin(chip, LIQUOR3)
	liquor4 := makePin(chip, LIQUOR4)
	water := makePin(chip, WATER)
	drain := makePin(chip, DRAIN)
	juice := makePin(chip, JUICE)
	x1 := makePin(chip, X1)
	x2 := makePin(chip, X2)
	x3 := makePin(chip, X3)
	x4 := makePin(chip, X4)

	r := gin.Default()
	r.StaticFile("/", "./public")
	r.POST("/pump/:id/:value", func(c *gin.Context) {
		id := c.Param("id")
		value, _ := strconv.Atoi(c.Param("value"))
		if value != 0 && value != 1 {
			c.JSON(400, gin.H{
				"error": "value can only be 0 or 1",
			})
		}

		switch id {
		case "liquor0":
			liquor0.SetValue(value)
		case "liquor1":
			liquor1.SetValue(value)
		case "liquor2":
			liquor2.SetValue(value)
		case "liquor3":
			liquor3.SetValue(value)
		case "liquor4":
			liquor4.SetValue(value)
		case "water":
			water.SetValue(value)
		case "drain":
			drain.SetValue(value)
		case "juice":
			juice.SetValue(value)
		case "x1":
			x1.SetValue(value)
		case "x2":
			x2.SetValue(value)
		case "x3":
			x3.SetValue(value)
		case "x4":
			x4.SetValue(value)
		}

		c.JSON(200, gin.H{
			"message": "ok",
			"id":      id,
			"value":   value,
		})
	})
	r.POST("/shake/:sleeptime/:iter", func(c *gin.Context) {
		sleeptime, _ := strconv.ParseInt(c.Param("sleeptime"), 10, 64)
		iterations, _ := strconv.Atoi(c.Param("iter"))

		for i := 0; i < iterations; i++ {
			x2.SetValue(1)
			time.Sleep(time.Duration(sleeptime) * time.Millisecond)
			x2.SetValue(0)
			time.Sleep(time.Duration(sleeptime) * time.Millisecond)
		}

		c.JSON(200, gin.H{})
	})

	log.Println("http://192.168.1.45:8080")
	r.Run() // listen and serve on 0.0.0.0:8080
}
