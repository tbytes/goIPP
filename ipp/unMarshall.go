package ipp

import (
	"bytes"
	"encoding/binary"
	"log"
	utility "utils"
)

// message
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

//   Each "additional-value" field is encoded as follows:
//
//   -----------------------------------------------
//   |                   value-tag                 |   1 byte // syntax
//   -----------------------------------------------
//   |            name-length  (value is 0x0000)   |   2 bytes
//   -----------------------------------------------
//   |              value-length (value is w)      |   2 bytes
//   -----------------------------------------------
//   |                     value                   |   w bytes
//   -----------------------------------------------
/*
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
*/

func ParseMessage(b []byte) (m Message, err error) {
	err = nil
	i, _ := binary.Varint(b[:1])
	m.majorVer = int8(i) //	version-number            	2 bytes  	- required
	i, _ = binary.Varint(b[1:2])
	m.minorVer = int8(i)
	log.Println(b[1:2])

	ii, _ := binary.Uvarint(b[3:4])
	m.operationIdStatusCode = uint16(ii) //	operation-id (request)		2 bytes  	- required
	m.IsResponse = typeCheck(uint16(ii)) //	or status-code (response)                   		

	iii, _ := binary.Varint(b[4:8]) //	request-id 					4 bytes  	- required	
	m.requestId = int32(iii)
	//et := bytes.IndexByte(b[8:], TAG_END)

	ags := splitAValues(b[8 : ]) //et+9]) //	attribute-group				n bytes 	- 0 or more
	m.attributeGroups = ags

	// m.endAttributeTag = b[et] //	end-of-attributes-tag		1 byte   	- required

	//m.Data = b[et+8:] //	data q bytes  	- optional
	//log.Println("Data: ", m.Data)

	return
}

func splitAValues(b []byte) (ags []attributeGroup) {
	util := utility.NewIterator(b)
	var vTag byte
	var name string
	// ====================================== begin-attribute-group-tag ================================================
	bt, err := util.GetNextOne() // get the first byte which should be 0x01 start operation-attributes	operation-attributes-tag
	if !err {
		log.Println("87: ", bt) // if next byte is a dud
	}
	if bt != 0x01 {
		log.Println("malformed packet", bt) // if first byte is not operation-attributes-tag
	}
	var ag attributeGroup
	ag.beginAttributeGroupTag = bt
	//   ----------------------------------------------------------
	//   |           begin-attribute-group-tag         |  1 byte  |-
	//   ----------------------------------------------------------
	//   |                   attribute                 |  p bytes |- 0 or more
	//   ----------------------------------------------------------
	one := utility.Start(1)
	var v attribute
	// ===================================== attribute-with ====================================================
	for vTag != TAG_END { // parse atribute groups (b []byte) until TAG_END is reached
		
		one.Plus()
		log.Println("Loop number: ", one.Get(), ag.beginAttributeGroupTag, len(ags))
		util.DeBug()
		vTag, err = util.GetNextOne() // get value tag
		_, isDelimitter := checkGroupTag(vTag)
		
		// ================================= New Group vTag is a Delimitter ====================================
		if isDelimitter {
			ag.attributes = append(ag.attributes, v) // add the attributes to the existing group
			ags = append(ags, ag)                    // append the group to the groups
			var agn attributeGroup 
			ag = agn                   // start a new attribute Group
			ag.beginAttributeGroupTag = vTag
			continue
		}

		if !err {
			log.Println("115: ", vTag) // if next byte is a dud 
			break
		}
		//   Each "attribute-with-one-value" field is encoded as follows:
		//   -----------------------------------------------
		//   |                   value-tag                 |   1 byte	Syntax
		//   -----------------------------------------------
		//   |               name-length  (value is u)     |   2 bytes
		//   -----------------------------------------------
		//   |                     name                    |   u bytes
		//   -----------------------------------------------
		//   |              value-length  (value is v)     |   2 bytes
		//   -----------------------------------------------
		//   |                     value                   |   v bytes
		//   -----------------------------------------------
		// =========================================== name length =========================================================
		x, err := util.GetNextN(2) // name-length; if 0 then is additional attribute
		if !err {
			log.Println("142") // there was a problem fetching the name length bytes (2)
			break
		}
		cnt := utility.Start(1) // value counter
		var xx int16
		buf := bytes.NewBuffer(x)
		er := binary.Read(buf, binary.BigEndian, &xx)
		if er != nil {
			log.Println("150: ", er)
		}
		nLength := xx
		// ======================================== name if length not 0 ===================================================
		if nLength != 0 { // attribute is an operation/status not an additional Value to current Attribute	
			if cnt.GTEq(1) { // append v to ag.attributes if not the first time through
				log.Println("158 - append v to ag.attributes: ", v)
				ag.attributes = append(ag.attributes, v)
				var nv attribute
				v = nv // reset the attribute
			}
			n, err := util.GetNextN(int(nLength)) // Start of new Attribute, if cnt > 1 then it is not the first run and we need to append 
			name = string(n)
			cnt.Plus() // increase av counter so we know that the next time nLength != 0 it is a new av and not the first one
			if !err {
				log.Println("166") // there was a problem getting the name bytes
				break
			}
			x, err = util.GetNextN(2) // value-length
			buf := bytes.NewBuffer(x)
			er := binary.Read(buf, binary.BigEndian, &xx)
			vLength := xx
			if er != nil {
				log.Println("177: ", er)
				break
			}
			util.DeBug()
			value, err := util.GetNextN(int(vLength)) // value
			if !err {
				log.Println("182 util.GetNextN(int(vLength)) vLength: ", vLength)
				break
			}
			av, _ := UnMarshallattribute(vTag, value) //returns attributeValue
			av.name = name
			av.nameLength = nLength
			v.appendValue(av)
			continue
		} else if nLength == 0 {
			x, err = util.GetNextN(2)
			buf := bytes.NewBuffer(x)
			er := binary.Read(buf, binary.BigEndian, &xx)
			vLength := xx
			if er != nil {
				log.Println("201: ", er)
				break
			}
			value, err := util.GetNextN(int(vLength)) // value

			if !err {
				log.Println("175")
				break
			}
			av, _ := UnMarshallattribute(vTag, value)
			av.name = ""
			v.appendValue(av)
			continue
		}

		ag.attributes = append(ag.attributes, v)
		ags = append(ags, ag)
	}
	ags = append(ags, ag)
	return
}

// Attribute Group Tags - Delimitters 
func checkGroupTag(b byte) (status string, err bool) {
	status = ""
	err = false
	switch b {
	case 0x01: // "operation-attributes-tag"
		status = "TAG_OPERATION"
		err = true
	case 0x02: // "job-attributes-tag"
		status = "TAG_JOB"
		err = true
	case 0x03: // "end-of-attributes-tag"
		status = "TAG_END"
		err = true
	case 0x04:
		status = "TAG_PRINTER" // "printer-attributes-tag"
		err = true
	case 0x05:
		status = "TAG_UNSUPPORTED_GROUP" // "unsupported-attributes-tag"
		err = true
	case 0x06:
		status = "TAG_SUBSCRIPTION"
		err = true
	case 0x07:
		status = "TAG_EVENT_NOTIFICATION"
		err = true
	case 0x08:
		status = "0x08" // Reserved (ipp-get-resources)
		err = true
	case 0x09:
		status = "TAG_DOCUMENT_ATTRIBUTES" //
		err = true
	}
	return status, err
}

// Attribute Syntaxes
//	bi = value tag as byte; bts = value as []byte
func UnMarshallattribute(bi byte, bts []byte) (attributeValue, error) {
	var a attributeValue
	// a.value = bts
	switch bi {
	case 0x21:
		a.valueTag = TAG_INTEGER // integer
		a.valueTagStr = "TAG_INTEGER"
		a.Marshal = (func() ([]byte, error) { b := a.value.(integer); return b.MarshalIPP() })
		a.UnMarshal = (func([]byte) error { b := a.value.(integer); return b.UnMarshalIPP(bts) })
		a.Length = (func() uint16 { return uint16(4) })
	case 0x22:
		a.valueTag = TAG_BOOLEAN // boolean
		a.valueTagStr = "TAG_BOOLEAN"
		a.Marshal = (func() ([]byte, error) { b := a.value.(ippBoolean); return b.MarshalIPP() })
		a.UnMarshal = (func(bts []byte) error { var b ippBoolean; b.UnMarshalIPP(bts); a.value = b; return nil})
		a.Length = (func() uint16 {return uint16(1) })
		a.String = (func() string {x := a.value.(ippBoolean); return x.String()})
	case 0x23:
		a.valueTag = TAG_ENUM // enum
		a.valueTagStr = "TAG_ENUM"
		a.Marshal = (func() ([]byte, error) { b := a.value.(enum); return b.MarshalIPP()})
		a.UnMarshal = (func(bts []byte) error { var b enum; b.UnMarshalIPP(bts); a.value = b; return nil})
		a.String = (func() string {x := a.value.(enum); return x.String()})
		a.Length = (func() uint16 {return uint16(4)})				
	case 0x30:
		a.valueTag = TAG_STRING // octetString with an  unspecified format
		a.valueTagStr = "TAG_STRING"
		a.Marshal = (func() ([]byte, error) {b := a.value.(octetString); return b.MarshalIPP()})
		a.UnMarshal = (func([]byte) error {b := a.value.(octetString); return b.UnMarshalIPP(bts)})
		a.Length = (func() uint16 {b := a.value.(octetString); return b.len()})
	case 0x31:
		a.valueTag = TAG_DATE // dateTime
		a.valueTagStr = "TAG_DATE"
		a.Marshal = (func() ([]byte, error) { b := a.value.(dateTime); return b.MarshalIPP() })
		a.Length = (func() uint16 { return uint16(9) })
		a.UnMarshal = (func([]byte) error { b := a.value.(dateTime); return b.UnMarshalIPP(bts) })
	case 0x32:
		a.valueTag = TAG_RESOLUTION // resolution
		a.valueTagStr = "TAG_RESOLUTION"
		a.Marshal = (func() ([]byte, error) { b := a.value.(resolution); return b.MarshalIPP() })
		a.UnMarshal = (func([]byte) error { b := a.value.(resolution); return b.UnMarshalIPP(bts) })
		a.Length = (func() uint16 { return uint16(11) })
	case 0x33:
		a.valueTag = TAG_RANGE // rangeOfInteger
		a.valueTagStr = "TAG_RANGE"
		a.Marshal = (func() ([]byte, error) { b := a.value.(rangeOfInteger); return b.MarshalIPP() })
		a.UnMarshal = (func(bts []byte) error { var b rangeOfInteger; b.UnMarshalIPP(bts); a.value = b; return nil})
		a.String = (func() string {x := a.value.(rangeOfInteger); return x.String()})
		a.Length = (func() uint16 { return uint16(8) })
	// case 0x34:
	//	a.valueTag = TAG_BEGIN_COLLECTION	// reserved for definition in a future IETF standards track document
	case 0x35:
		a.valueTag = TAG_TEXTLANG // textWithLanguage
		a.valueTagStr = "TAG_TEXTLANG"
		a.Marshal = (func() ([]byte, error) { b := a.value.(textWithLanguage); return b.MarshalIPP() })
		a.UnMarshal = (func([]byte) error { b := a.value.(textWithLanguage); return b.UnMarshalIPP(bts) })
		a.Length = (func() uint16 { b := a.value.(textWithLanguage); return b.length() })
	case 0x36:
		a.valueTag = TAG_NAMELANG // nameWithLanguage
		a.valueTagStr = "TAG_NAMELANG"
		a.Marshal = (func() ([]byte, error) { b := a.value.(nameWithLanguage); return b.MarshalIPP() })
		a.UnMarshal = (func([]byte) error { b := a.value.(nameWithLanguage); return b.UnMarshalIPP(bts) })
		a.Length = (func() uint16 { b := a.value.(nameWithLanguage); return b.length() })
	case 0x47:
		a.valueTag = TAG_CHARSET
		a.valueTagStr = "TAG_CHARSET"
		a.Marshal = (func() ([]byte, error) {b := a.value.(charset); return b.MarshalIPP()})
		a.UnMarshal = (func(bts []byte) error {var b charset; b.UnMarshalIPP(bts); a.value = b; return nil})
		a.Length = (func() uint16 {b := a.value.(charset); return uint16(b.len())})
		a.String = (func() string {x := a.value.(charset); return x.String()})
	case 0x48:
		a.valueTag = TAG_LANGUAGE // textWithLanguage
		a.valueTagStr = "TAG_LANGUAGE"
		a.Marshal = (func() ([]byte, error) {b := a.value.(naturalLanguage); return b.MarshalIPP()})
		// a.UnMarshal = (func([]byte) error {b := a.value.(naturalLanguage); return b.UnMarshalIPP(bts)})
		a.UnMarshal = (func(bts []byte) error {var b naturalLanguage; b.UnMarshalIPP(bts); a.value = b; return nil})
		a.Length = (func() uint16 {b := a.value.(naturalLanguage); return uint16(b.len())})
		a.String = (func() string {x := a.value.(naturalLanguage); return x.String()})
	case 0x44:
		a.valueTag = TAG_KEYWORD // textWithLanguage
		a.valueTagStr = "TAG_KEYWORD"
		a.Marshal = (func() ([]byte, error) { b := a.value.(keyword); return b.MarshalIPP() })
		a.UnMarshal = (func(bts []byte) error {var b keyword; b.UnMarshalIPP(bts); a.value = b; return nil })
		a.Length = (func() uint16 { b := a.value.(keyword); return uint16(b.len()) })
		a.String = (func() string {x := a.value.(keyword); return x.String()})
		// case 0x37:
		//	a.valueTag = TAG_END_COLLECTION     
		// 0x24-0x2F         reserved for integer types for definition in future IETF standards track documents 
	case 0x49:
		a.valueTag = TAG_MIMETYPE // textWithLanguage
		a.valueTagStr = "TAG_MIMETYPE"
		a.Marshal = (func() ([]byte, error) { b := a.value.(mimeMediaType); return b.MarshalIPP() })
		a.UnMarshal = (func(bts []byte) error {var b mimeMediaType; b.UnMarshalIPP(bts); a.value = b; return nil })
		a.Length = (func() uint16 { b := a.value.(mimeMediaType); return uint16(b.len()) })
		a.String = (func() string {x := a.value.(mimeMediaType); return x.String()})
	
	}

	a.UnMarshal(bts)
	return a, nil
}
