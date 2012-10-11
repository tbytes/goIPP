package main

// import "fmt"
import "ipp"

func main() {
	var c ipp.CupsServer
	c.SetServer("http://www.google.com")
	c.GetPrinters()
}
