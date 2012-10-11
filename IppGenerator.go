package ipp

import (
	"bytes"
	"encoding/binary"
	//"log"
)

//   -----------------------------------------------
//   |                  version-number..           |   2 bytes  - required
//   -----------------------------------------------
//   |               operation-id (request)..      |
//   |                      or                     |   2 bytes  - required
//   |               status-code (response)..      |
//   -----------------------------------------------
//   |                   request-id..              |   4 bytes  - required
//   -----------------------------------------------
//   |                 attribute-group             |   n bytes - 0 or more
//   -----------------------------------------------
//   |              end-of-attributes-tag..        |   1 byte   - required
//   -----------------------------------------------
//   |                     data                    |   q bytes  - optional
//   -----------------------------------------------

func (im *Message) marshallMsg() *bytes.Buffer {

	b := new(bytes.Buffer)
	x := im.marshallAtrib()
	binary.Write(b, binary.BigEndian, im.majorVer)
	binary.Write(b, binary.BigEndian, im.minorVer)
	binary.Write(b, binary.BigEndian, im.operationIdStatusCode)
	binary.Write(b, binary.BigEndian, im.requestId)
	binary.Write(b, binary.BigEndian, x.Bytes())
	binary.Write(b, binary.BigEndian, uint8(3))
	binary.Write(b, binary.BigEndian, im.Data)
	return b
}

//   Each "attribute-group" field is encoded as follows:
//
//   -----------------------------------------------
//   |           begin-attribute-group-tag         |  1 byte
//   ----------------------------------------------------------
//   |                   attribute                 |  p bytes |- 0 or more
//   ----------------------------------------------------------
//
//   When an attribute is single valued (e.g. "copies" with value of 10) or multi-valued with one value
//   (e.g. "sides-supported" with just the value 'one-sided') it is encoded with just an "attribute-with-one-value"
//   field. When an attribute is multi-valued with n values (e.g."sides-supported" with the values 'one-sided' and 
//	 'two-sided-long-edge') , it is encoded with an "attribute-with-one-value" field followed by n-1 "additional-value" fields.

//      The "value-tag" field specifies the attribute syntax, e.g. 0x44 for the attribute syntax 'keyword'.
//      The "name-length" field specifies the length of the "name" field in bytes, e.g. u in the above 
//			diagram or 15 for the name "sides-supported".
//      The "name" field contains the textual name of the attribute, e.g."sides-supported".
//      The "value-length" field specifies the length of the "value" field in bytes, e.g. v in the above 
//			diagram or 9 for the (keyword) value 'one-sided'.
//      The "value" field contains the value of the attribute, e.g. the textual value 'one-sided'.

//   Each "attribute-with-one-value" field is encoded as follows:
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
//   Each "additional-value" field is encoded as follows:
//   -----------------------------------------------
//   |                   value-tag                 |   1 byte
//   -----------------------------------------------
//   |            name-length  (value is 0x0000)   |   2 bytes
//   -----------------------------------------------
//   |              value-length (value is w)      |   2 bytes
//   -----------------------------------------------
//   |                     value                   |   w bytes
//   -----------------------------------------------

func (im *Message) marshallAtrib() *bytes.Buffer {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, uint8(TAG_OPERATION))
	for _, ag := range im.attributeGroups {
		for _, a := range ag.attributes {
			//  The "value-tag" field specifies the attribute syntax, e.g. 0x44 for the attribute syntax 'keyword'.
			//		valueTag    byte
			//	The "name-length" specifies the length of the "name" field in bytes
			//	IF the field has the value of 0 it signifies that this is an "additional-value". 
			//	The value of the "name-length" field distinguishes an "additional-value" field
			//	from an "attribute-with-one-value" field ("name-length" is not 0).
			//		nameLength  int16
			//	The "name" field contains the textual name of the attribute, e.g."sides-supported"
			//		name        string
			//	The "value-length" field specifies the length of the "value" field in bytes, e.g. w in the above diagram or 19 for the (keyword) value 'two-sided-long-edge'.
			//		valueLength int16
			//	The "value" field contains the value of the attribute, e.g. the textual value 'one-sided'.
			//		value       []byte
			for iii, v := range a.values {

				if iii == 0 {
					binary.Write(b, binary.BigEndian, v.valueTag)
					binary.Write(b, binary.BigEndian, v.nameLength)
					binary.Write(b, binary.BigEndian, []byte(v.name))
					binary.Write(b, binary.BigEndian, v.valueLength)
					if v.value == nil {
						binary.Write(b, binary.BigEndian, 0x0000)
					} else {
						binary.Write(b, binary.BigEndian, v.value)
					}
				} else {
					binary.Write(b, binary.BigEndian, v.valueTag)
					binary.Write(b, binary.BigEndian, uint16(0))
					binary.Write(b, binary.BigEndian, v.valueLength)
					binary.Write(b, binary.BigEndian, v.value)
				}
			}
		}
	}

	return b
}
