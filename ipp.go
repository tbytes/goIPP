package ipp

import (
	// "bytes"
	// "encoding/binary"
	"fmt"
	//"log"
	"io/ioutil"
	"net/http"
	"os"
//	"sync/atomic"
)

type CupsServer struct {
	uri            string
	username       string
	password       string
	requestCounter int32
}

func (c *CupsServer) SetServer(server string) {
	c.uri = server
}

func (c *CupsServer) CreateRequest(operationId uint16) Message {
	m := newMessage(operationId)
	return m
}

func (c *CupsServer) GetPrinters() {
	m := c.CreateRequest(CUPS_GET_PRINTERS)
	m.addAttribute(TAG_CHARSET, "attributes-charset", charset("utf-8"))
	m.addAttribute(TAG_LANGUAGE, "attributes-natural-language", naturalLanguage("en-us"))
	m.addAttribute(TAG_KEYWORD, "requested-attributes", keyword("printer-name"))
	m.addAttribute(TAG_ENUM, "printer-type", enum(0))
	m.addAttribute(TAG_ENUM, "printer-type-mask", enum(1))
	c.DoRequest(m)
	
}

func (c *CupsServer) DoRequest(m Message) {

    fii, _ := os.Create("/Users/tmartino/aC")
	defer fii.Close()
	s := m.marshallMsg()
	fii.Write(s.Bytes())
    fmt.Println(s)
	resp, err := http.Post("http://localhost:631", "application/ipp", s)
	if err != nil {
	fmt.Println("err: ",err)
		// handle error
	}
  body, errr := ioutil.ReadAll(resp.Body)
	if errr != nil {
		// handle error
		fmt.Println("errr:   ", errr)
	}

  fmt.Println("Response Body: ", string(body))

}
