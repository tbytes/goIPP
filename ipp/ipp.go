package ipp

import (

)
//   Every operation request contains the following REQUIRED parameters:
//
//      - a "version-number",
//      - an "operation-id",
//      - a "request-id", and
//      - the attributes that are REQUIRED for that type of request.	

func NewRequest(idStatusCode uint16) Message {
	return newMessage(idStatusCode)
}

//   Every operation response contains the following REQUIRED parameters:
//
//      - a "version-number",
//      - a "status-code",
//      - the "request-id" that was supplied in the corresponding request,
//        and
//      - the attributes that are REQUIRED for that type of response.

func NewResponse(idStatusCode uint16) Message {
	return newMessage(idStatusCode)
}