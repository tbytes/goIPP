package ipp

const (
	UUID_SOURCE = "/proc/sys/kernel/random/uuid"

	MAJOR_VERSION = 0x02
	MINOR_VERSION = 0x01

	PORT = 631

	MAX_NAME   = 256
	MAX_VALUES = 8

	/*
	   3.5.1 Delimiter Tags

	   The following table specifies the values for the delimiter tags:

	   		When a "begin-attribute-group-tag" field occurs in the protocol, it means that zero or 
	   more following attributes up to the next delimiter tag MUST be attributes belonging to 
	   the attribute group specified by the value of the "begin-attribute-group-tag". For 
	   example, if the value of "begin-attribute-group-tag" is 0x01, the following attributes 
	   MUST be members of the Operations Attributes group.
	   		The "end-of-attributes-tag" (value 0x03) MUST occur exactly once in an operation.  
	   It MUST be the last "delimiter-tag". If the operation has a document-content group, 
	   the document data in that group MUST follow the "end-of-attributes-tag".
	   		The order and presence of "attribute-group" fields (whose beginning is marked by the 
	   "begin-attribute-group-tag" subfield) for each operation request and each operation 
	   response MUST be that defined in the model document. For further details, see section 
	   3.7 "(Attribute) Name" and 13 "Appendix A: Protocol Examples".
	   		A Printer MUST treat a "delimiter-tag" (values from 0x00 through 0x0F) differently from 
	   a "value-tag" (values from 0x10 through 0xFF) so that the Printer knows that there is an 
	   entire attribute group that it doesn't understand as opposed to a single value that it doesn't understand.
	*/

	TAG_ZERO               = 0x00 // reserved for definition in a future IETF standards track document
	TAG_OPERATION          = 0x01 // "operation-attributes-tag"
	TAG_JOB                = 0x02 // "job-attributes-tag"
	TAG_END                = 0x03 // "end-of-attributes-tag"
	TAG_PRINTER            = 0x04 // "printer-attributes-tag"
	TAG_UNSUPPORTED_GROUP  = 0x05 // "unsupported-attributes-tag"
	TAG_SUBSCRIPTION       = 0x06
	TAG_EVENT_NOTIFICATION = 0x07
	//	0x08	Reserved (ipp-get-resources)	
	TAG_DOCUMENT_ATTRIBUTES = 0x09 //	document-attributes-tag	[PWG5100.5]
	//	0x0A-0x0F	Unassigned
	/*
	   	The remaining tables show values for the "value-tag" field, which is the first octet of an 
	   	attribute. The "value-tag" field specifies the type of the value of the attribute.

	      	The following table specifies the "out-of-band" values for the "value-tag" field.
	*/

	TAG_UNSUPPORTED_VALUE = 0x10 // unsupported
	TAG_DEFAULT           = 0x11 // reserved for 'default' for definition in a future IETF standards track document
	TAG_UNKNOWN           = 0x12 // unknown
	TAG_NOVALUE           = 0x13 // no-value
	//	0x14-0x1F        reserved for "out-of-band" values in future IETF standards track documents.

	// integer values for the "value-tag"
	TAG_NOTSETTABLE = 0x15
	TAG_DELETEATTR  = 0x16
	TAG_ADMINDEFINE = 0x17
	// 0x20              reserved for definition in a future IETF standards track document
	TAG_INTEGER = 0x21 // integer
	TAG_BOOLEAN = 0x22 // boolean
	TAG_ENUM    = 0x23 // enum
	// 0x24-0x2F         reserved for integer types for definition in future IETF standards track documents

	/*
	   octetString values for the "value-tag" field:
	*/

	TAG_STRING           = 0x30 // octetString with an  unspecified format
	TAG_DATE             = 0x31 // dateTime
	TAG_RESOLUTION       = 0x32 // resolution
	TAG_RANGE            = 0x33 // rangeOfInteger
	TAG_BEGIN_COLLECTION = 0x34 // reserved for definition in a future IETF standards track document
	TAG_TEXTLANG         = 0x35 // textWithLanguage
	TAG_NAMELANG         = 0x36 // nameWithLanguage
	TAG_END_COLLECTION   = 0x37

	//	character-string values for the "value-tag" field:

	//	NOTE:  an attribute value always has a type, which is explicitly
	//	specified by its tag; one such tag value is "nameWithoutLanguage".
	//	An attribute's name has an implicit type, which is keyword.

	//		0x40              reserved for definition in a future IETF standards track document
	//		NOTE: 0x40 is reserved for "generic character-string" if it should ever be needed.
	TAG_TEXT = 0x41 // textWithoutLanguage
	TAG_NAME = 0x42 // nameWithoutLanguage
	//    	0x43              reserved for definition in a future IETF standards track document
	TAG_KEYWORD    = 0x44
	TAG_URI        = 0x45
	TAG_URISCHEME  = 0x46
	TAG_CHARSET    = 0x47
	TAG_LANGUAGE   = 0x48 // naturalLanguage
	TAG_MIMETYPE   = 0x49 // mimeMediaType
	TAG_MEMBERNAME = 0x4a
	TAG_MASK       = 0x7fffffff
	TAG_COPY       = -0x7fffffff - 1

	RES_PER_INCH = 3
	RES_PER_CM   = 4

	// finishings
	FINISHINGS_NONE                = 3
	FINISHINGS_STAPLE              = 4
	FINISHINGS_PUNCH               = 5
	FINISHINGS_COVER               = 6
	FINISHINGS_BIND                = 7
	FINISHINGS_SADDLE_STITCH       = 8
	FINISHINGS_EDGE_STITCH         = 9
	FINISHINGS_FOLD                = 10
	FINISHINGS_TRIM                = 11
	FINISHINGS_BALE                = 12
	FINISHINGS_BOOKLET_MAKER       = 13
	FINISHINGS_JOB_OFFSET          = 14
	FINISHINGS_STAPLE_TOP_LEFT     = 20
	FINISHINGS_STAPLE_BOTTOM_LEFT  = 21
	FINISHINGS_STAPLE_TOP_RIGHT    = 22
	FINISHINGS_STAPLE_BOTTOM_RIGHT = 23
	FINISHINGS_EDGE_STITCH_LEFT    = 24
	FINISHINGS_EDGE_STITCH_TOP     = 25
	FINISHINGS_EDGE_STITCH_RIGHT   = 26
	FINISHINGS_EDGE_STITCH_BOTTOM  = 27
	FINISHINGS_STAPLE_DUAL_LEFT    = 28
	FINISHINGS_STAPLE_DUAL_TOP     = 29
	FINISHINGS_STAPLE_DUAL_RIGHT   = 30
	FINISHINGS_STAPLE_DUAL_BOTTOM  = 31
	FINISHINGS_BIND_LEFT           = 50
	FINISHINGS_BIND_TOP            = 51
	FINISHINGS_BIND_RIGHT          = 52
	FINISHINGS_BIND_BOTTO          = 53

	PORTRAIT          = 3
	LANDSCAPE         = 4
	REVERSE_LANDSCAPE = 5
	REVERSE_PORTRAIT  = 6

	QUALITY_DRAFT  = 3
	QUALITY_NORMAL = 4
	QUALITY_HIGH   = 5

	// job-state
	JOB_PENDING    = 3
	JOB_HELD       = 4
	JOB_PROCESSING = 5
	JOB_STOPPED    = 6
	JOB_CANCELLED  = 7
	JOB_ABORTED    = 8
	JOB_COMPLETE   = 9

	// document-state
	PRINTER_IDLE       = 3
	PRINTER_PROCESSING = 4
	PRINTER_STOPPED    = 5

	ERROR     = -1
	IDLE      = 0
	HEADER    = 1
	ATTRIBUTE = 2
	DATA      = 3

	//	============ Operation-Ids =============
	PRINT_JOB              = 0x0002
	PRINT_URI              = 0x0003
	VALIDATE_JOB           = 0x0004
	CREATE_JOB             = 0x0005
	SEND_DOCUMENT          = 0x0006
	SEND_URI               = 0x0007
	CANCEL_JOB             = 0x0008
	GET_JOB_ATTRIBUTES     = 0x0009
	GET_JOBS               = 0x000a
	GET_PRINTER_ATTRIBUTES = 0x000b
	HOLD_JOB               = 0x000c
	RELEASE_JOB            = 0x000d
	RESTART_JOB            = 0x000e
	//	0x000F              reserved for a future operation	
	PAUSE_PRINTER                   = 0x0010
	RESUME_PRINTER                  = 0x0011
	PURGE_JOBS                      = 0x0012
	SET_PRINTER_ATTRIBUTES          = 0x0013
	SET_JOB_ATTRIBUTES              = 0x0014
	GET_PRINTER_SUPPORTED_VALUES    = 0x0015
	CREATE_PRINTER_SUBSCRIPTION     = 0x0016
	CREATE_JOB_SUBSCRIPTION         = 0x0017
	GET_SUBSCRIPTION_ATTRIBUTES     = 0x0018
	GET_SUBSCRIPTIONS               = 0x0019
	RENEW_SUBSCRIPTION              = 0x001a
	CANCEL_SUBSCRIPTION             = 0x001b
	GET_NOTIFICATIONS               = 0x001c
	SEND_NOTIFICATIONS              = 0x001d
	GET_PRINT_SUPPORT_FILES         = 0x0021
	ENABLE_PRINTER                  = 0x0022
	DISABLE_PRINTER                 = 0x0023
	PAUSE_PRINTER_AFTER_CURRENT_JOB = 0x0024
	HOLD_NEW_JOBS                   = 0x0025
	RELEASE_HELD_NEW_JOBS           = 0x0026
	DEACTIVATE_PRINTER              = 0x0027
	ACTIVATE_PRINTER                = 0x0028
	RESTART_PRINTER                 = 0x0029
	SHUTDOWN_PRINTER                = 0x002a
	STARTUP_PRINTER                 = 0x002b
	REPROCESS_JOB                   = 0x002c
	CANCEL_CURRENT_JOB              = 0x002d
	SUSPEND_CURRENT_JOB             = 0x002e
	RESUME_JOB                      = 0x002f
	PROMOTE_JOB                     = 0x0030
	SCHEDULE_JOB_AFTER              = 0x0031
	PRIVATE                         = 0x4000

	//	============ CUPS ============	
	CUPS_PRINTER_LOCAL         = 0x0000
	CUPS_PRINTER_CLASS         = 0x0001
	CUPS_PRINTER_REMOTE        = 0x0002
	CUPS_PRINTER_BW            = 0x0004
	CUPS_PRINTER_COLOR         = 0x0008
	CUPS_PRINTER_DUPLEX        = 0x0010
	CUPS_PRINTER_STAPLE        = 0x0020
	CUPS_PRINTER_COPIES        = 0x0040
	CUPS_PRINTER_COLLATE       = 0x0080
	CUPS_PRINTER_PUNCH         = 0x0100
	CUPS_PRINTER_COVER         = 0x0200
	CUPS_PRINTER_BIND          = 0x0400
	CUPS_PRINTER_SORT          = 0x0800
	CUPS_PRINTER_SMALL         = 0x1000
	CUPS_PRINTER_MEDIUM        = 0x2000
	CUPS_PRINTER_LARGE         = 0x4000
	CUPS_PRINTER_VARIABLE      = 0x8000
	CUPS_PRINTER_IMPLICIT      = 0x1000
	CUPS_PRINTER_DEFAULT       = 0x2000
	CUPS_PRINTER_FAX           = 0x4000
	CUPS_PRINTER_REJECTING     = 0x8000
	CUPS_PRINTER_DELETE        = 0x1000
	CUPS_PRINTER_NOT_SHARED    = 0x2000
	CUPS_PRINTER_AUTHENTICATED = 0x4000
	CUPS_PRINTER_COMMANDS      = 0x8000
	CUPS_PRINTER_OPTIONS       = 0xe6ff

	CUPS_GET_DEFAULT      = 0x4001
	CUPS_GET_PRINTERS     = 0x4002
	CUPS_ADD_PRINTER      = 0x4003
	CUPS_DELETE_PRINTER   = 0x4004
	CUPS_GET_CLASSES      = 0x4005
	CUPS_ADD_CLASS        = 0x4006
	CUPS_DELETE_CLASS     = 0x4007
	CUPS_ACCEPT_JOBS      = 0x4008
	CUPS_REJECT_JOBS      = 0x4009
	CUPS_SET_DEFAULT      = 0x400a
	CUPS_GET_DEVICES      = 0x400b
	CUPS_GET_PPDS         = 0x400c
	CUPS_MOVE_JOB         = 0x400d
	CUPS_AUTHENTICATE_JOB = 0x400e

	//	============ Status Codes ===========
	OK                           = 0x0000
	OK_SUBST                     = 0x0001
	OK_CONFLICT                  = 0x0002
	OK_IGNORED_SUBSCRIPTIONS     = 0x0003
	OK_IGNORED_NOTIFICATIONS     = 0x0004
	OK_TOO_MANY_EVENTS           = 0x0005
	OK_BUT_CANCEL_SUBSCRIPTION   = 0x0006
	REDIRECTION_OTHER_SITE       = 0x0300
	BAD_REQUEST                  = 0x0400
	FORBIDDEN                    = 0x0401
	NOT_AUTHENTICATED            = 0x0402
	NOT_AUTHORIZED               = 0x0403
	NOT_POSSIBLE                 = 0x0404
	TIMEOUT                      = 0x0405
	NOT_FOUND                    = 0x0406
	GONE                         = 0x0407
	REQUEST_ENTITY               = 0x0408
	REQUEST_VALUE                = 0x0409
	DOCUMENT_FORMAT              = 0x040a
	ATTRIBUTES                   = 0x040b
	URI_SCHEME                   = 0x040c
	CHARSET                      = 0x040d
	CONFLICT                     = 0x040e
	COMPRESSION_NOT_SUPPORTED    = 0x040f
	COMPRESSION_ERROR            = 0x0410
	DOCUMENT_FORMAT_ERROR        = 0x0411
	DOCUMENT_ACCESS_ERROR        = 0x0412
	ATTRIBUTES_NOT_SETTABLE      = 0x0413
	IGNORED_ALL_SUBSCRIPTIONS    = 0x0414
	TOO_MANY_SUBSCRIPTIONS       = 0x0415
	IGNORED_ALL_NOTIFICATIONS    = 0x0416
	PRINT_SUPPORT_FILE_NOT_FOUND = 0x0417

	INTERNAL_ERROR              = 0x0500
	OPERATION_NOT_SUPPORTED     = 0x0501
	SERVICE_UNAVAILABLE         = 0x0502
	VERSION_NOT_SUPPORTED       = 0x0503
	DEVICE_ERROR                = 0x0504
	TEMPORARY_ERROR             = 0x0505
	NOT_ACCEPTING               = 0x0506
	PRINTER_BUSY                = 0x0507
	ERROR_JOB_CANCELLED         = 0x0508
	MULTIPLE_JOBS_NOT_SUPPORTED = 0x0509
	PRINTER_IS_DEACTIVATED      = 0x50a
)
