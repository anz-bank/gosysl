Ê;
RestApi⁄;
	
RestApi"
	interface"Storer"R
interface_docA"?Storer abstracts all required RefData persistence and retrieval"
json_property_separator"-*Ï
GET /api/{key}/{startTime}Õ
GET /api/{key}/{startTime}"!
method_name"GetDataWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
DataBH/api/{key}/{startTime}
	startTime	ö 
key	ö *°
	POST /apiì
	POST /api"
method_name"CreateDataSet"
patterns
:
"rest"

middleware"AuthorizeRoot:B
KeyB/apiJ
dsJ

DataSetPayload*Ó
DELETE /api/{key}/{startTime}Ã
DELETE /api/{key}/{startTime}"
method_name"
DeleteData"
patterns
:
"rest" 

middleware"AuthorizeDataSet:

returnBH/api/{key}/{startTime}
	startTime	ö 
key	ö *˜
GET /api/{key}/nameﬂ
GET /api/{key}/name"
method_name"GetDataSetName"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B	
KeyNameB`/api/{key}/name
key	ö  7
time-J$

	
RestApi{queryTime<:string}ö  *Ç
POST /api/admin/{key}/subscribeﬁ
POST /api/admin/{key}/subscribe" 
method_name"PutSubscription"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
SubscriptionB2/api/admin/{key}/subscribe
key	ö3 J
sJ

Subscription*¿
PUT /api/{key}≠
PUT /api/{key}"
method_name	"PutData"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
DataB"
/api/{key}
key	ö
 J
dpJ

DataPayload*Ü
!PUT /api/admin/{key}/restrictions‡
!PUT /api/admin/{key}/restrictions"
method_name"PutRestriction"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
RestrictionB5/api/admin/{key}/restrictions
key	ö0 J
rJ

Restriction*á
PUT /api/{key}/{startTime}Ë
PUT /api/{key}/{startTime}"!
method_name"PutDataWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
DataBH/api/{key}/{startTime}
	startTime	ö 
key	ö J
dpJ

DataPayload*ﬁ
GET /api/{key}À
GET /api/{key}"
method_name	"GetData"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
DataB[
/api/{key}
key	ö 7
time-J$

	
RestApi{queryTime<:string}ö *˜
GET /api/{key}/schema›
GET /api/{key}/schema"
method_name"	GetSchema"
patterns
:
"rest" 

middleware"AuthorizeDataSet:
B
SchemaBb/api/{key}/schema
key	ö 7
time-J$

	
RestApi{queryTime<:string}ö *ˆ
#GET /api/admin/{key}/creation-timesŒ
#GET /api/admin/{key}/creation-times"!
method_name"GetCreationTimes"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
CreationTimesB7/api/admin/{key}/creation-times
key	ö+ *€
PUT /api/{key}/schema¡
PUT /api/{key}/schema"
method_name"	PutSchema"
patterns
:
"rest" 

middleware"AuthorizeDataSet:
B
SchemaB)/api/{key}/schema
key	ö J
spJ

SchemaPayload*¢
!PUT /api/{key}/schema/{startTime}¸
!PUT /api/{key}/schema/{startTime}"#
method_name"PutSchemaWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:
B
SchemaBO/api/{key}/schema/{startTime}
	startTime	ö 
key	ö J
spJ

SchemaPayload*Ö
!GET /api/{key}/schema/{startTime}ﬂ
!GET /api/{key}/schema/{startTime}"#
method_name"GetSchemaWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:
B
SchemaBO/api/{key}/schema/{startTime}
	startTime	ö 
key	ö *ï
GET /apià
GET /api"
method_name	"GetKeys"
patterns
:
"rest"

middleware"AuthorizeRoot"

method_doc	"DataSet:B
KeysB/api*Ï
!GET /api/admin/{key}/restrictions∆
!GET /api/admin/{key}/restrictions"
method_name"GetRestriction"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
RestrictionB5/api/admin/{key}/restrictions
key	ö. *Ÿ
PUT /api/{key}/name¡
PUT /api/{key}/name"
method_name"PutDataSetName"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B	
KeyNameB'/api/{key}/name
key	ö" J
npJ

NamePayload*Ö
!POST /api/admin/{key}/unsubscribeﬂ
!POST /api/admin/{key}/unsubscribe"#
method_name"DeleteSubscription"
patterns
:
"rest" 

middleware"AuthorizeDataSet:

returnB4/api/admin/{key}/unsubscribe
key	ö6 J
sJ

Subscription*≈
DELETE /api/admin/{key}©
DELETE /api/admin/{key}"
method_name"DeleteDataSet"
patterns
:
"rest" 

middleware"AuthorizeDataSet:

returnB(/api/admin/{key}
key	ö% *Ö
$DELETE /api/{key}/schema/{startTime}‹
$DELETE /api/{key}/schema/{startTime}"
method_name"DeleteSchema"
patterns
:
"rest" 

middleware"AuthorizeDataSet:

returnBO/api/{key}/schema/{startTime}
	startTime	ö 
key	ö *‚
 GET /api/admin/{key}/start-timesΩ
 GET /api/admin/{key}/start-times"
method_name"GetStartTimes"
patterns
:
"rest" 

middleware"AuthorizeDataSet:	B
TimesB4/api/admin/{key}/start-times
key	ö( 2å
Restriction¸õ

DataFrozenUntil	öI

SchemaFrozenUntil	öH

AdminScopes"
	öL
 
ReadWriteScopes"
	öK


ReadScopes"
	öJB\
docU"SRestriction contains scope access restriction and frozen times for schema and data.2\
KeysT

Keys"
	öOB9
doc2"0Keys is JSON result type for getKeys in REST API2v
NamePayloadg

Name	öjBP
docI"GNamePayload is JSON payload on REST API request to update data set name2ÿ
DataSetPayload≈o
/
StartTimeStrB
json"
start-timeöf

Name	öe
)

JSONSchemaB
json"schemaögBR
docK"IDataSetPayload is JSON payload on REST API request to create new data set2¢
TimesòF
+
Data#"!
B
json"
data-timesöY

Schema"
	öZBN
docG"ETimes contains schema and data times, used to get StartTimes for both2∆
UpdateEvent∂h

Deleted	öw

Data	öu

	StartTime	öt

Key	ös

Schema	övBJ
docC"AUpdateEvent holds all information necessary to post to subscribes2Å
KeyNamev%

Name	öV

Key	öUBM
docF"DKeyName is JSON result type for get and put dataDetNamre in REST API2u
SchemaPayloadd

Schema	öpBK
docD"BSchemaPayload is JSON payload on REST API request to update schema2[
KeyT

Key	öRB>
doc7"5Key is JSON result type for createDataSet in REST API2¶
CreationStartTimeê3

CreationTime	ö]

	StartTime	ö^BY
docR"PCreationStartTime contains start and creation time for a schema or data snapshot2m
DataPayload^

Data	ömBG
doc@">DataPayload is JSON payload on REST API request to update data2ù
Subscriptionå<

URLB
json"urlöD

SecreteToken	öEBL
docE"CSubscription holds external endpoint values for change notification2±
Data®Z

CreationTime	ö<
%
JSONDataB
json"dataö;

	StartTime	ö:BJ
docC"AData holds JSON data valid from StartTime created at CreationTime2≥
CreationTimes°ø
i
DataaB
json"data-time-mapJ?

	
RestApiCreationTimes!map of string:CreationStartTimeöa
R
SchemaHJ?

	
RestApiCreationTimes!map of string:CreationStartTimeöbB]
docV"TCreationTimes contains schema and data times maps, used to StartTime to CreationTims2¡
Schema∂^

CreationTime	öA
)

JSONSchemaB
json"schemaö@

	StartTime	ö?BT
docM"KSchema holds JSON Schema to validate Data against, a name for key creation.