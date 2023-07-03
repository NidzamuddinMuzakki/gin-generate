package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"
)

func main() {
	ConfigRuntime()
	StartGin()
}

// ConfigRuntime sets the number of operating system threads.
func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

// StartWorkers start starsWorker by goroutine.

// StartGin starts gin web server with setting router.
func StartGin() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/generate/signature", func(c *gin.Context) {
		jsonData, _ := ioutil.ReadAll(c.Request.Body)
		rawPayload := string(jsonData)
		merchantKey := c.Request.Header.Get("merchant-key")
		fmt.Println(rawPayload, merchantKey)
		h := sha512.New()
		h.Write([]byte(rawPayload + merchantKey))
		bs := h.Sum(nil)
		ourSignature := hex.EncodeToString(bs)

		data := make(map[string]string)
		data["signature"] = ourSignature
		c.JSON(200, data)
	})
	port := "8080"
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
