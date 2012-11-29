package ipp

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
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
/*
 Octets           Symbolic Value               Protocol field

 0x0100           1.0                          version-number
 0x000A           Get-Jobs                     operation-id
 0x00000123       0x123                        request-id
 0x01             start operation-attributes   operation-attributes-tag
 0x47             charset type                 value-tag
 */

func (c *CupsServer) GetPrinters() {
	m := c.CreateRequest(CUPS_GET_PRINTERS)
	m.AddAttribute(TAG_CHARSET, "attributes-charset", charset("utf-8"))
	m.AddAttribute(TAG_LANGUAGE, "attributes-natural-language", naturalLanguage("en-us"))
	m.AddAttribute(TAG_KEYWORD, "requested-attributes", keyword("printer-name"))
	m.AddAttribute(TAG_ENUM, "printer-type", enum(0))
	m.AddAttribute(TAG_ENUM, "printer-type-mask", enum(1))
	c.DoRequest(m)
	
}

func (c *CupsServer) GetPrinterAttributes() {
	m := c.CreateRequest(GET_PRINTER_ATTRIBUTES)
	m.AddAttribute(TAG_CHARSET, "attributes-charset", charset("utf-8"))
	m.AddAttribute(TAG_LANGUAGE, "attributes-natural-language", naturalLanguage("en-us"))
	m.AddAttribute(TAG_URI, "printer-uri", uri("ipp://192.168.1.1:631/ipp/"))
	
	a := NewAttribute()
	a.AddValue(TAG_KEYWORD, "requested-attributes", keyword("copies-supported"))
	a.AddValue(TAG_KEYWORD, "", keyword("document-format-supported"))
	a.AddValue(TAG_KEYWORD, "", keyword("printer-is-accepting-jobs"))
	a.AddValue(TAG_KEYWORD, "", keyword("printer-state"))
	a.AddValue(TAG_KEYWORD, "", keyword("printer-state-message"))
	a.AddValue(TAG_KEYWORD, "", keyword("printer-state-reasons"))
	
	m.AppendAttribute(a)
	
	c.DoRequest(m)
	
}

func (c *CupsServer) PrintTestPage() {
	m := c.CreateRequest(PAUSE_PRINTER)
	m.AddAttribute(TAG_CHARSET, "attributes-charset", charset("utf-8"))
	m.AddAttribute(TAG_LANGUAGE, "attributes-natural-language", naturalLanguage("en-us"))
	c.DoRequest(m)
	
}

func (c *CupsServer) DoRequest(m Message) {

    fii, _ := os.Create("/Users/tmartino/aC")
	defer fii.Close()
	s := m.marshallMsg()
	fii.Write(s.Bytes())
	// "http://192.168.1.8:631/ipp/printer" "application/ipp"
	
	resp, err := http.Post("http://192.168.1.8:631/ipp/printer", "application/ipp", s)
	if err != nil {
		fmt.Println("err: ",err)
	}
  body, errr := ioutil.ReadAll(resp.Body)
	if errr != nil {
		fmt.Println("errr:   ", errr)
	}
  fmt.Println("Response Body: ", string(body))
  fmt.Println("Response Body bytes[]: ", body)
  fmt.Println("End Tag: ", TAG_END)
  fmt.Println("Header: ", resp.Header)
  x, eerr := ParseMessage(body)
  fmt.Println("Message: ", x.attributeGroups[1].attributes[0], len(x.attributeGroups[1].attributes[0].values), "eerr: ", eerr)

}
