package ipp

import (

)
/* 
3.9 (Attribute) Value

   The syntax types (specified by the "value-tag" field) and most of the
   details of the representation of attribute values are defined in the
   IPP model document. The table below augments the information in the
   model document, and defines the syntax types from the model document
   in terms of the 5 basic types defined in section 3, "Encoding of the
   Operation Layer". The 5 types are US-ASCII-STRING, LOCALIZED-STRING,
   SIGNED-INTEGER, SIGNED-SHORT, SIGNED-BYTE, and OCTET-STRING.

  Syntax of Attribute   Encoding
  Value

  textWithoutLanguage,  LOCALIZED-STRING.
  nameWithoutLanguage

  textWithLanguage      OCTET-STRING consisting of 4 fields:
                          a. a SIGNED-SHORT which is the number of
                             octets in the following field
                          b. a value of type natural-language,
                          c. a SIGNED-SHORT which is the number of
                             octets in the following field,
                          d. a value of type textWithoutLanguage.
                        The length of a textWithLanguage value MUST be
                        4 + the value of field a + the value of field c.

  nameWithLanguage      OCTET-STRING consisting of 4 fields:
                          a. a SIGNED-SHORT which is the number of
                             octets in the following field
                          b. a value of type natural-language,
                          c. a SIGNED-SHORT which is the number of
                             octets in the following field
                          d. a value of type nameWithoutLanguage.
                        The length of a nameWithLanguage value MUST be
                        4 + the value of field a + the value of field c.

  charset,              US-ASCII-STRING.
  naturalLanguage,
  mimeMediaType,
  keyword, uri, and
  uriScheme

  boolean               SIGNED-BYTE  where 0x00 is 'false' and 0x01 is 'true'.

  integer and enum      a SIGNED-INTEGER.

  dateTime              OCTET-STRING consisting of eleven octets whose
                        contents are defined by "DateAndTime" in RFC
                        1903 [RFC1903].

  resolution            OCTET-STRING consisting of nine octets of  2
                        SIGNED-INTEGERs followed by a SIGNED-BYTE. The
                        first SIGNED-INTEGER contains the value of
                        cross feed direction resolution. The second
                        SIGNED-INTEGER contains the value of feed
                        direction resolution. The SIGNED-BYTE contains
                        the units

  rangeOfInteger        Eight octets consisting of 2 SIGNED-INTEGERs.
                        The first SIGNED-INTEGER contains the lower
                        bound and the second SIGNED-INTEGER contains
                        the upper bound.

  1setOf  X             Encoding according to the rules for an
                        attribute with more than 1 value.  Each value
                        X is encoded according to the rules for
                        encoding its type.

  octetString           OCTET-STRING

   The attribute syntax type of the value determines its encoding and
   the value of its "value-tag"
*/

/*
   Every operation request contains the following REQUIRED parameters:
      - a "version-number",
      - an "operation-id",
      - a "request-id", and
      - the attributes that are REQUIRED for that type of request.

   Every operation response contains the following REQUIRED parameters:
      - a "version-number",
      - a "status-code",
      - the "request-id" that was supplied in the corresponding request,
        and
      - the attributes that are REQUIRED for that type of response.
*/
 
 
/* Example Protocol Request

13.1 Print-Job Request

   The following is an example of a Print-Job request with job-name,
   copies, and sides specified. The "ipp-attribute-fidelity" attribute
   is set to 'true' so that the print request will fail if the "copies"
   or the "sides" attribute are not supported or their values are not
   supported.

  Octets          Symbolic Value                		Protocol field
  
  0x0101          					1.1                           		version-number
  0x0002          					Print-Job                     		operation-id
  0x00000001      					1                             		request-id
  0x01            					start operation-attributes    		operation-attributes-tag
  
  
  0x47            					charset type                  		value-tag
  0x0012                                        						name-length
  attributes-charset     			attributes-charset     				name 
  0x0008                     					                   		value-length
  us-ascii        					US-ASCII                      		value
  
  0x48            					natural-language type         		value-tag
  0x001B                                        						name-length
  attributes-natural-language		attributes-natural-language			name
  0x0005                                        						value-length
  en-us           					en-US			                    value
  
  
  0x45            					uri type                      		value-tag
  0x000B                                        						name-length
  printer-uri     					printer-uri                   		name
  0x0015                                        						value-length
  ipp://forest/pinetree   			printer pinetree              		value
    
  0x42            					nameWithoutLanguage type      		value-tag
  0x0008                                        						name-length
  job-name        					job-name                      		name
  0x0006                                        						value-length
  foobar          					foobar                        		value
  
  0x22            					boolean type                  		value-tag
  0x0016                                        						name-length
  ipp-attribute-fidelity	  		ipp-attribute-fidelity        		name
  0x0001                                        						value-length
  0x01            					true                          		value
  
  0x02            					start job-attributes          		job-attributes-tag
  
  0x21            					integer type                  		value-tag
  0x0006                                        						name-length
  copies          					copies                        		name
  0x0004                                        						value-length
  0x00000014      					20                            		value
  
  0x44            					keyword type                  		value-tag
  0x0005                                        						name-length
  sides           					sides                         		name
  0x0013                                        						value-length
  two-sided-long-edge      			two-sided-long-edge           		value
  
  0x03            					end-of-attributes             		end-of-attributes-tag
  %!PS...         					<PostScript>                  		data


13.2 Print-Job Response (successful)

   Here is an example of a successful Print-Job response to the previous
   Print-Job request.  The printer supported the "copies" and "sides"
   attributes and their supplied values.  The status code returned is
   'successful-ok'.

  Octets            				Symbolic Value              		Protocol field

  0x0101            				1.1                         		version-number
  0x0000            				successful-ok               		status-code
  0x00000001        				1                           		request-id
  0x01              				start operation-attributes  		operation-attributes-tag
  
  0x47              				charset type                		value-tag
  0x0012                                        						name-length
  attributes-charset		      	attributes-charset          		name
  0x0008                                        						value-length
  us-ascii          				US-ASCII                    		value
  
  0x48              				natural-language type       		value-tag
  0x001B                                        						name-length
  attributes-natural-language		attributes-natural-language			name
  0x0005                                        						value-length
  en-us             				en-US         						value
  
  0x41              				textWithoutLanguage type    		value-tag
  0x000E                                        						name-length
  status-message    				status-message              		name
  0x000D                                        						value-length
  successful-ok     				successful-ok               		value
  
  0x02              				start job-attributes        		job-attributes-tag
  
  0x21              				integer                     		value-tag
  0x0006                                        						name-length
  job-id            				job-id                      		name
  0x0004                                        						value-length
  147               				147                         		value
  
  0x45              				uri type                    		value-tag
  0x0007                                        						name-length
  job-uri           				job-uri                     		name
  0x0019                                        						value-length
  ipp://forest/pinetree/123     	job 123 on pinetree         		value
    
  0x23              				enum type                   		value-tag
  0x0009                                        						name-length
  job-state         				job-state                   		name
  0x0004                                        						value-length
  0x0003            				pending                     		value
  0x03              				end-of-attributes           		end-of-attributes-tag


13.3 Print-Job Response (failure)

   Here is an example of an unsuccessful Print-Job response to the
   previous Print-Job request. It fails because, in this case, the
   printer does not support the "sides" attribute and because the value
   '20' for the "copies" attribute is not supported. Therefore, no job
   is created, and neither a "job-id" nor a "job-uri" operation
   attribute is returned. The error code returned is 'client-error-
   attributes-or-values-not-supported' (0x040B).

  0x0101        1.1                           version-number
  0x040B        client-error-attributes-or-   status-code
                values-not-supported
  0x00000001    1                             request-id
  0x01          start operation-attributes    operation-attributes tag
  0x47          charset type                  value-tag
  0x0012                                      name-length
  attributes-   attributes-charset            name
  charset
  0x0008                                      value-length
  us-ascii      US-ASCII                      value
  0x48          natural-language type         value-tag
  0x001B                                      name-length
  attributes-   attributes-natural-language   name
  natural-
  language
  0x0005                                      value-length
  en-us         en-US                         value
  0x41          textWithoutLanguage type      value-tag
  0x000E                                      name-length
  status-       status-message                name
  message
  0x002F                                      value-length
  client-error-                               value
  attributes-   values-not-supported
  or-values-    client-error-attributes-or-
  not-supported
  0x05          start unsupported-attributes  unsupported-attributes tag
  0x21          integer type                  value-tag
  0x0006                                      name-length
  copies        copies                        name
  0x0004                                      value-length
  0x00000014    20                            value
  0x10          unsupported  (type)           value-tag
  0x0005                                      name-length
  sides         sides                         name
  0x0000                                      value-length
  0x03          end-of-attributes             end-of-attributes-tag
  */