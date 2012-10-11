package ipp

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
)

/*
type                                 size in bytes

byte, uint8, int8                     1
uint16, int16                         2
uint32, int32, float32                4
uint64, int64, float64, complex64     8
complex128                           16


uint8       the set of all unsigned  8-bit integers (0 to 255)
uint16      the set of all unsigned 16-bit integers (0 to 65535)
uint32      the set of all unsigned 32-bit integers (0 to 4294967295)
uint64      the set of all unsigned 64-bit integers (0 to 18446744073709551615)

int8        the set of all signed  8-bit integers (-128 to 127)
int16       the set of all signed 16-bit integers (-32768 to 32767)
int32       the set of all signed 32-bit integers (-2147483648 to 2147483647)
int64       the set of all signed 64-bit integers (-9223372036854775808 to 9223372036854775807)

float32     the set of all IEEE-754 32-bit floating-point numbers
float64     the set of all IEEE-754 64-bit floating-point numbers

complex64   the set of all complex numbers with float32 real and imaginary parts
complex128  the set of all complex numbers with float64 real and imaginary parts

byte        alias for uint8
rune        alias for int32
*/

/*
   From the standpoint of a parser that performs an action based on a
   "tag" value, the encoding consists of:

   -----------------------------------------------
   |                  version-number             |   2 bytes  - required
   -----------------------------------------------
   |               operation-id (request)        |
   |                      or                     |   2 bytes  - required
   |               status-code (response)        |
   -----------------------------------------------
   |                   request-id                |   4 bytes  - required
   -----------------------------------------------------------
   |        tag (delimiter-tag or value-tag)     |   1 byte  |
   -----------------------------------------------           |-0 or more
   |           empty or rest of attribute        |   x bytes |
   -----------------------------------------------------------
   |              end-of-attributes-tag          |   1 byte   - required
   -----------------------------------------------
   |                     data                    |   y bytes  - optional
   -----------------------------------------------

   The following show what fields the parser would expect after each
   type of  "tag":

      -  "begin-attribute-group-tag": expect zero or more "attribute" fields
      -  "value-tag": expect the remainder of an "attribute-with-one-value" or  an "additional-value".
      -  "end-of-attributes-tag": expect that "attribute" fields are complete and there is optional "data"

*/
const (
	ippTrue  = 0x01
	ippFalse = 0x00
)

// Marshaler is the interface implemented by objects that
// can marshal themselves into valid IPP Message.
type Marshler interface {
	MarshalIPP() ([]byte, error)
}

type textWithoutLanguage []byte

func (i *textWithoutLanguage) bytes() []byte {
	return []byte(*i)
}

func (i *textWithoutLanguage) len() int {
	return len(i.bytes())
}

type nameWithoutLanguage []byte

func (i *nameWithoutLanguage) bytes() []byte {
	return []byte(*i)
}
func (i *nameWithoutLanguage) len() int {
	return len(i.bytes())
}

type signedShort int16

func (i *signedShort) bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, i)
	return buf.Bytes()
}

func (s *signedShort) len() int {
	return len(s.bytes())
}

type signedByte int8

func (i *signedByte) bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, i)
	return buf.Bytes()
}

type charset []byte                // US-ASCII-STRING.
func (i *charset) bytes() []byte { // US-ASCII-STRING.
	return []byte(*i)
}
func (i *charset) len() int { // US-ASCII-STRING.
	return len(i.bytes())
}

type naturalLanguage []byte // US-ASCII-STRING. 

func (i *naturalLanguage) bytes() []byte { // US-ASCII-STRING. 
	return []byte(*i)
}
func (i *naturalLanguage) len() int {
	return len(i.bytes())
}
func (i *naturalLanguage) MarshalIPP() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, []byte(*i)...)
	return buf, nil
}

type mimeMediaType []byte // US-ASCII-STRING.

func (i *mimeMediaType) bytes() []byte { // US-ASCII-STRING.
	return []byte(*i)
}

type keyword []byte // US-ASCII-STRING.

func (i *keyword) bytes() []byte { // US-ASCII-STRING.
	return []byte(*i)
}

func (i *keyword) len() int {
	return len(i.bytes())
}

func (i *keyword) MarshalIPP() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, []byte(*i)...)
	return buf, nil
}

type uri []byte // US-ASCII-STRING.

func (i *uri) bytes() []byte { // US-ASCII-STRING.
	return []byte(*i)
}

type uriScheme []byte // US-ASCII-STRING.

func (i *uriScheme) bytes() []byte { // US-ASCII-STRING. 
	return []byte(*i)
}

type signedInteger int32 // SIGNED-INTEGER

func (i *signedInteger) bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, i)
	return buf.Bytes()
}

// The length of a textWithLanguage value MUST be
// 4 + the value of field a + the value of field c.
type octetString struct {
	nameLength  signedShort         // a. number of octets in the following field
	name        naturalLanguage     // b. type natural-language,
	valueLength signedShort         // c. the number of octets in the following field,
	value       textWithoutLanguage // d. type textWithoutLanguage.
}

func (o *octetString) MarshalIPP() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, o.nameLength.bytes()...)
	buf = append(buf, o.name.bytes()...)
	buf = append(buf, o.valueLength.bytes()...)
	buf = append(buf, o.value.bytes()...)
	//s3 := append(s2, s0...)

	return buf, nil
}

func (o *octetString) length() uint16 {
	x := uint16(o.nameLength.len() + o.valueLength.len() + o.name.len() + o.value.len())
	log.Println(x)
	return x
}

func (o *octetString) validate() bool {
	x := 4+o.nameLength.len()+o.valueLength.len() == o.nameLength.len()
	log.Println(x)
	return x
}

// OCTET-STRING consisting of 4 fields
type textWithLanguage octetString

func (o *textWithLanguage) length() uint16 {
	return uint16(4 + o.nameLength + o.valueLength)
}

func (t *textWithLanguage) MarshalIPP() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, t.nameLength.bytes()...)
	buf = append(buf, t.name.bytes()...)
	buf = append(buf, t.valueLength.bytes()...)
	buf = append(buf, t.value.bytes()...)

	return buf, nil
}

// OCTET-STRING consisting of 4 fields: []byte
type nameWithLanguage octetString

func (o *nameWithLanguage) length() uint16 {
	return uint16(4 + o.nameLength + o.valueLength)
}

func (t *nameWithLanguage) MarshalIPP() ([]byte, error) {
	return t.MarshalIPP()
}

// SIGNED-BYTE  where 0x00 is 'false' and 0x01 is 'true'.
type ippBoolean signedByte

func (i *ippBoolean) MarshalIPP() ([]byte, error) {
	return i.MarshalIPP()
}

type integer signedInteger

func (i *integer) MarshalIPP() ([]byte, error) {
	return i.MarshalIPP()
}

type enum signedInteger

func (e *enum) MarshalIPP() ([]byte, error) {
	return e.MarshalIPP()
}

// OCTET-STRING consisting of eleven octets whose
// contents are defined by "DateAndTime" in [RFC1903]:
//	DateAndTime ::= TEXTUAL-CONVENTION
//	DISPLAY-HINT "2d-1d-1d,1d:1d:1d.1d,1a1d:1d"
//	STATUS:       current
//	DESCRIPTION: "A date-time specification.							
type dateTime struct { //	field  octets  contents                range
	year    signedShort //	1      1-2   year                      0..65536
	month   signedByte  //	2       3    month                     1..12
	day     signedByte  //	3       4    day                       1..31
	hour    signedByte  //	4       5    hour                      0..23
	minutes signedByte  //	5       6    minutes                   0..59
	seconds signedByte  //	6       7    seconds                   0..60
	//	             (use 60 for leap-second)
	deciSeconds signedByte //	7       8    deci-seconds              0..9
	UTC         signedByte //	8       9    direction from UTC        '+' / '-'
	hoursFrUTC  signedByte //	9      10    hours from UTC            0..11	
}

func NewDateTime(d time.Time) dateTime {
	dt := dateTime{}
	dt.year = signedShort(d.Year())
	dt.month = signedByte(d.Month())
	dt.day = signedByte(d.Day())
	dt.hour = signedByte(d.Hour())
	dt.minutes = signedByte(d.Minute())
	dt.seconds = signedByte(d.Second())
	dt.deciSeconds = signedByte(d.Nanosecond())
	dt.UTC = signedByte(0x002d)
	_, frUtc := d.Zone()
	dt.hoursFrUTC = signedByte(frUtc / 60)
	return dt
}
func (o *dateTime) MarshalIPP() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, o.year.bytes()...)
	buf = append(buf, o.month.bytes()...)
	buf = append(buf, o.day.bytes()...)
	buf = append(buf, o.hour.bytes()...)
	buf = append(buf, o.minutes.bytes()...)
	buf = append(buf, o.seconds.bytes()...)
	buf = append(buf, o.deciSeconds.bytes()...)
	buf = append(buf, o.UTC.bytes()...)
	buf = append(buf, o.hoursFrUTC.bytes()...)

	return buf, nil
}

type resolution struct { // OCTET-STRING consisting of nine octets of 2 SIGNED-INTEGERs followed by a SIGNED-BYTE.
	crossFeedDirection, feedDirection signedInteger // The first SIGNED-INTEGER: cross feed direction resolution. The second SIGNED-INTEGER: value of feed direction resolution.
	units                             signedByte    // The SIGNED-BYTE contains the units
}

func NewResolution(crossFdDirection, feedDirection int, unit int8) (resolution, error) {
	res := resolution{}
	res.crossFeedDirection = signedInteger(crossFdDirection)
	res.feedDirection = signedInteger(feedDirection)
	res.units = signedByte(unit)
	return res, nil
}

func (o *resolution) MarshalIPP() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, o.crossFeedDirection.bytes()...)
	buf = append(buf, o.feedDirection.bytes()...)
	buf = append(buf, o.units.bytes()...)

	return buf, nil
}

type rangeOfInteger struct { // Eight octets consisting of 2 SIGNED-INTEGERs.
	lowerBound, upperBound signedInteger // The first SIGNED-INTEGER contains the lower
	// bound and the second SIGNED-INTEGER contains the upper bound.
}

func (o *rangeOfInteger) MarshalIPP() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, o.lowerBound.bytes()...)
	buf = append(buf, o.upperBound.bytes()...)

	return buf, nil
}

// 1setOfX        		// Encoding according to the rules for an attribute with more than 1 value.
// Each value X is encoded according to the rules for encoding its type.

// function refer() uses the TAG value to set the MarshalIPP and Length Functions 
// (eventually the UnMarshall as well)
func (a *attributeValue) refer() {
	switch a.valueTag {
	case TAG_STRING: // octetString with an  unspecified format
		a.Marshal = (func() ([]byte, error) { b := a.value.(octetString); return b.MarshalIPP() })
		a.Length = (func() uint16 { b := a.value.(octetString); return b.length() })
	case TAG_DATE: // dateTime
		a.Marshal = (func() ([]byte, error) { b := a.value.(dateTime); return b.MarshalIPP() })
		a.Length = (func() uint16 { return uint16(9) })
	case TAG_RESOLUTION: // resolution
		a.Marshal = (func() ([]byte, error) { b := a.value.(resolution); return b.MarshalIPP() })
		a.Length = (func() uint16 { return uint16(11) })
	case TAG_RANGE: // rangeOfInteger
		a.Marshal = (func() ([]byte, error) { b := a.value.(rangeOfInteger); return b.MarshalIPP() })
		a.Length = (func() uint16 { return uint16(8) })
	// case TAG_BEGIN_COLLECTION: 	// reserved for definition in a future IETF standards track document
	//	a.Marshal = (func() ([]byte, error) {b := a.value.(); return b.MarshalIPP()})
	case TAG_TEXTLANG: // textWithLanguage
		a.Marshal = (func() ([]byte, error) { b := a.value.(textWithLanguage); return b.MarshalIPP() })
		a.Length = (func() uint16 { b := a.value.(textWithLanguage); return b.length() })
	case TAG_LANGUAGE: // textWithLanguage
		a.Marshal = (func() ([]byte, error) { b := a.value.(naturalLanguage); return b.MarshalIPP() })
		a.Length = (func() uint16 { b := a.value.(naturalLanguage); return uint16(b.len()) })
	case TAG_KEYWORD: // textWithLanguage
		a.Marshal = (func() ([]byte, error) { b := a.value.(keyword); return b.MarshalIPP() })
		a.Length = (func() uint16 { b := a.value.(keyword); return uint16(b.len()) })
	case TAG_NAMELANG: // nameWithLanguage
		a.Marshal = (func() ([]byte, error) { b := a.value.(nameWithLanguage); return b.MarshalIPP() })
		a.Length = (func() uint16 { b := a.value.(nameWithLanguage); return b.length() })
	case TAG_INTEGER: // integer
		a.Marshal = (func() ([]byte, error) { b := a.value.(integer); return b.MarshalIPP() })
		a.Length = (func() uint16 { return uint16(4) })
	case TAG_BOOLEAN: // boolean
		a.Marshal = (func() ([]byte, error) { b := a.value.(ippBoolean); return b.MarshalIPP() })
		a.Length = (func() uint16 { return uint16(1) })
	case TAG_ENUM:
		a.Marshal = (func() ([]byte, error) { b := a.value.(enum); return b.MarshalIPP() })
		a.Length = (func() uint16 { return uint16(4) })
	case TAG_CHARSET:
		a.Marshal = (func() ([]byte, error) { b := a.value.(charset); return b.bytes(), nil })
		a.Length = (func() uint16 { b := a.value.(charset); return uint16(b.len()) })
	}
}
