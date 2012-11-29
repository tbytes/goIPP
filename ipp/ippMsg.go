package ipp

import (
	// "encoding/binary"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"strconv"
)

//   -----------------------------------------------
//   |                  version-number             |   2 bytes  - required
//   -----------------------------------------------
//   |               operation-id (request)        |
//   |                      or                     |   2 bytes  - required
//   |               status-code (response)        |
//   -----------------------------------------------
//   |                   request-id                |   4 bytes  - required
//   -----------------------------------------------
//   |                 attribute-group             |   n bytes - 0 or more
//   -----------------------------------------------
//   |              end-of-attributes-tag          |   1 byte   - required
//   -----------------------------------------------
//   |                     data                    |   q bytes  - optional
//   -----------------------------------------------

/*
   "request-id":	every invocation of an operation is identified by a
   "request-id" value. For each request, the client chooses the
   "request-id" which MUST be an integer (possibly unique depending on
   client requirements) in the range from 1 to 2**31 - 1 (inclusive).
   This "request-id" allows clients to manage multiple outstanding
   requests. The receiving IPP object copies all 32-bits of the client-
   supplied "request-id" attribute into the response so that the client
   can match the response with the correct outstanding request, even if
   the "request-id" is out of range.  If the request is terminated
   before the complete "request-id" is received, the IPP object rejects
   the request and returns a response with a "request-id" of 0.
*/

type Message struct {
	majorVer              	int8
	minorVer              	int8
	operationIdStatusCode 	uint16
	operationOrStatusCode 	byte
	requestId             	int32
	attributeGroups       	[]attributeGroup
	endAttributeTag       	byte
	Data                  	[]byte
	IsResponse				bool
}

func NewMessage(idStatusCode uint16) Message {
	return newMessage(idStatusCode)
}

func typeCheck(i uint16) bool {
	if i >= 0x0000 && i <= 0x7FFF {
		return true
	}
	return false
}

func newMessage(idStatusCode uint16) Message {
	
	var m Message
	// version-number
	m.majorVer = int8(MAJOR_VERSION)
	m.minorVer = int8(MINOR_VERSION)
	// operation-id (request) or status-code (response)	
	m.operationIdStatusCode = idStatusCode
	// request-id
	m.requestId = newID()
	// end-of-attributes-tag
	m.endAttributeTag = TAG_END
	return m
}

//   Each "attribute-group" field is encoded as follows:
//
//   -----------------------------------------------
//   |           begin-attribute-group-tag         |  1 byte
//   ----------------------------------------------------------
//   |                   attribute                 |  p bytes |- 0 or more
//   ----------------------------------------------------------
//
//   When an attribute is single valued (e.g. "copies" with value of 10)
//   or multi-valued with one value (e.g. "sides-supported" with just the
//   value 'one-sided') it is encoded with just an "attribute-with-one-
//   value" field. When an attribute is multi-valued with n values (e.g.
//   "sides-supported" with the values 'one-sided' and 'two-sided-long-
//   edge'), it is encoded with an "attribute-with-one-value" field
//   followed by n-1 "additional-value" fields.
//   
//      The "value-tag" field specifies the attribute syntax, 
//      e.g. 0x44 for the attribute syntax 'keyword'.
//
//      The "name-length" field specifies the length of the "name" field in bytes, 
//      e.g. u in the above diagram or 15 for the name "sides-supported".
//
//      The "name" field contains the textual name of the attribute, 
//      e.g."sides-supported".
//
//      The "value-length" field specifies the length of the "value" field in bytes, 
//      e.g. v in the above diagram or 9 for the (keyword) value 'one-sided'.
//
//      The "value" field contains the value of the attribute, 
//      e.g. the textual value 'one-sided'.

type attributeGroup struct {
	beginAttributeGroupTag byte
	attributes             []attribute
}

func newAg(a attribute) attributeGroup {
	var x attributeGroup
	x.beginAttributeGroupTag = uint8(TAG_OPERATION)
	x.attributes = append(x.attributes, a)
	return x
}

func (im *Message) AddAttribute(tag byte, name string, value interface{}) {
	im.addAttribute(tag, name, value)
	return
}

func (im *Message) AppendAttribute(attrib attribute) {
	im.attributeGroups = append(im.attributeGroups, newAg(attrib))
	return
}

func (im *Message) addAttribute(tag byte, name string, value interface{}) {
	var attrib attribute
	attrib.addValue(tag, name, value)
	im.attributeGroups = append(im.attributeGroups, newAg(attrib))
	return
}

//   Each "attribute-with-one-value" field is encoded as follows:
//
//   -----------------------------------------------
//   |                   value-tag                 |   1 byte
//   -----------------------------------------------
//   |               name-length  (value is u)     |   2 bytes
//   -----------------------------------------------
//   |                     name                    |   u bytes
//   -----------------------------------------------
//   |              value-length  (value is v)     |   2 bytes
//   -----------------------------------------------
//   |                     value                   |   v bytes
//   -----------------------------------------------
type attribute struct {
	values []attributeValue
}

type attributeValue struct {
	//  The "value-tag" field specifies the attribute syntax, 
	//	e.g. 0x44 for the attribute syntax 'keyword'.
	valueTag byte
	valueTagStr string
	//	The "name-length" specifies the length of the "name" field in bytes
	//	IF the field has the value of 0 it signifies that this is an "additional-value". 
	//	The value of the "name-length" field distinguishes an "additional-value" field
	//	from an "attribute-with-one-value" field ("name-length" is not 0).
	nameLength int16
	//	The "name" field contains the textual name of the attribute, 
	//	e.g."sides-supported"
	name string
	//	The "value-length" field specifies the length of the "value" field in bytes, 
	//	e.g. w in the above diagram or 19 for the (keyword) value 'two-sided-long-edge'.	
	valueLength uint16
	//	The "value" field contains the value of the attribute, 
	//	e.g. the textual value 'one-sided'.
	value interface{}
	Marshal	func() ([]byte, error)
	UnMarshal	func([]byte) (error)
	Length	func() (uint16)
	String func() (string)
}

func NewAttribute() attribute {
	var a attribute
	return a
}

func (i *attribute) AddValue(tag byte, name string, value interface{}) {
	i.addValue(tag, name, value)
	return
}
func (a *attributeValue) Str() (s string) {
	s = a.valueTagStr + " - " + a.name + " - " 
	return
}

func (i *attribute) appendValue(av attributeValue) {
	i.values = append(i.values, av)
	return
}
// 0x20 to 0xFF

//  The "value-tag" (tag byte) field specifies the attribute syntax, 
//	e.g. 0x44 for the attribute syntax 'keyword'.
//	The "name" field (name string) contains the textual name of the attribute, 
//	e.g."sides-supported"
//	The "value" (value []byte) field contains the value of the attribute, 
//	e.g. the textual value 'one-sided'.
func (i *attribute) addValue(tag byte, name string, value interface{}) {

	// if name == "" && len(i.values) <= 0 {return} //error: the first value of an attribute must be named  
	// if name == "" && len(i.values) <= 0 {return} //error: additional values of an attribute cannot be named  
	if len(i.values) > 0 {
		var v attributeValue
		
		v.valueTag = tag
		v.nameLength = 0x0000
		v.value = value
		v.refer()
		v.valueLength = v.Length()
		i.values = append(i.values, v)

		return
	}
	
	var vv attributeValue
	vv.valueTag = tag
	vv.name = name
	vv.nameLength = int16(len(name))
	vv.value = value
	vv.refer()
	vv.valueLength = vv.Length()
	i.values = append(i.values, vv)
	return
}

//   Each "additional-value" field is encoded as follows:
//
//   -----------------------------------------------
//   |                   value-tag                 |   1 byte
//   -----------------------------------------------
//   |            name-length  (value is 0x0000)   |   2 bytes
//   -----------------------------------------------
//   |              value-length (value is w)      |   2 bytes
//   -----------------------------------------------
//   |                     value                   |   w bytes
//   -----------------------------------------------

func newID() int32 {
	uid := int32(1) //, _ := UUID4()
	return uid
}

// https://github.com/serverhorror/uuid/blob/master/uuid4.go
// UUID4() string generates a random version 4 uuid and returns an
// appropriate string representation
func UUID4() (int32, error) {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		log.Fatal(err)
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	str := fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
	uid, err := strconv.Atoi(str)
	return int32(uid), err
}
