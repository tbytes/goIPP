package main

// import "fmt"
import "ipp"
  
func main() {
	var c ipp.CupsServer
	c.SetServer("192.168.1.8")
	c.GetPrinterAttributes()

}
