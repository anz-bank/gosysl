«<
RestApiŸ<
	
RestApi"
	interface"Storer"R
interface_docA"?Storer abstracts all required RefData persistence and retrieval"
json_property_separator"-*ì
GET /api/{key}/{startTime}Í
GET /api/{key}/{startTime}"!
method_name"GetDataWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
DataBH/api/{key}/{startTime}
	startTime	š 
key	š *¡
	POST /api“
	POST /api"
method_name"CreateDataSet"
patterns
:
"rest"

middleware"AuthorizeRoot:B
KeyB/apiJ
dsJ

DataSetPayload*î
DELETE /api/{key}/{startTime}Ì
DELETE /api/{key}/{startTime}"
method_name"
DeleteData"
patterns
:
"rest" 

middleware"AuthorizeDataSet:

returnBH/api/{key}/{startTime}
	startTime	š 
key	š *÷
GET /api/{key}/nameß
GET /api/{key}/name"
method_name"GetDataSetName"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B	
KeyNameB`/api/{key}/name
key	š 7
time-J$

	
RestApi{queryTime<:string}š *‚
POST /api/admin/{key}/subscribeŞ
POST /api/admin/{key}/subscribe" 
method_name"PutSubscription"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
SubscriptionB2/api/admin/{key}/subscribe
key	š3 J
sJ

Subscription*À
PUT /api/{key}­
PUT /api/{key}"
method_name	"PutData"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
DataB"
/api/{key}
key	š J
dpJ

DataPayload*†
!PUT /api/admin/{key}/restrictionsà
!PUT /api/admin/{key}/restrictions"
method_name"PutRestriction"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
RestrictionB5/api/admin/{key}/restrictions
key	š0 J
rJ

Restriction*‡
PUT /api/{key}/{startTime}è
PUT /api/{key}/{startTime}"!
method_name"PutDataWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
DataBH/api/{key}/{startTime}
	startTime	š 
key	š J
dpJ

DataPayload*ô
GET /api/{key}á
GET /api/{key}"
method_name	"GetData"
patterns
:
"rest" 

middleware"AuthorizeDataSet"

method_doc"Data:B
DataB[
/api/{key}
key	š 7
time-J$

	
RestApi{queryTime<:string}š *
GET /api/{key}/schemaõ
GET /api/{key}/schema"
method_name"	GetSchema"
patterns
:
"rest" 

middleware"AuthorizeDataSet"

method_doc"Schema:
B
SchemaBb/api/{key}/schema
key	š 7
time-J$

	
RestApi{queryTime<:string}š *ö
#GET /api/admin/{key}/creation-timesÎ
#GET /api/admin/{key}/creation-times"!
method_name"GetCreationTimes"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
CreationTimesB7/api/admin/{key}/creation-times
key	š+ *Û
PUT /api/{key}/schemaÁ
PUT /api/{key}/schema"
method_name"	PutSchema"
patterns
:
"rest" 

middleware"AuthorizeDataSet:
B
SchemaB)/api/{key}/schema
key	š J
spJ

SchemaPayload*Ù
PUT /api/{key}/nameÁ
PUT /api/{key}/name"
method_name"PutDataSetName"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B	
KeyNameB'/api/{key}/name
key	š
 J
npJ

NamePayload*…
!GET /api/{key}/schema/{startTime}ß
!GET /api/{key}/schema/{startTime}"#
method_name"GetSchemaWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:
B
SchemaBO/api/{key}/schema/{startTime}
	startTime	š 
key	š *•
GET /apiˆ
GET /api"
method_name	"GetKeys"
patterns
:
"rest"

middleware"AuthorizeRoot"

method_doc	"DataSet:B
KeysB/api*ì
!GET /api/admin/{key}/restrictionsÆ
!GET /api/admin/{key}/restrictions"
method_name"GetRestriction"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
RestrictionB5/api/admin/{key}/restrictions
key	š. *¢
!PUT /api/{key}/schema/{startTime}ü
!PUT /api/{key}/schema/{startTime}"#
method_name"PutSchemaWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:
B
SchemaBO/api/{key}/schema/{startTime}
	startTime	š  
key	š  J
spJ

SchemaPayload*…
!POST /api/admin/{key}/unsubscribeß
!POST /api/admin/{key}/unsubscribe"#
method_name"DeleteSubscription"
patterns
:
"rest" 

middleware"AuthorizeDataSet:

returnB4/api/admin/{key}/unsubscribe
key	š6 J
sJ

Subscription*Ü
DELETE /api/admin/{key}À
DELETE /api/admin/{key}"
method_name"DeleteDataSet"
patterns
:
"rest" 

middleware"AuthorizeDataSet"

method_doc"Admin:

returnB(/api/admin/{key}
key	š% *…
$DELETE /api/{key}/schema/{startTime}Ü
$DELETE /api/{key}/schema/{startTime}"
method_name"DeleteSchema"
patterns
:
"rest" 

middleware"AuthorizeDataSet:

returnBO/api/{key}/schema/{startTime}
	startTime	š" 
key	š" *â
 GET /api/admin/{key}/start-times½
 GET /api/admin/{key}/start-times"
method_name"GetStartTimes"
patterns
:
"rest" 

middleware"AuthorizeDataSet:	B
TimesB4/api/admin/{key}/start-times
key	š( 2Œ
Restrictionü›

DataFrozenUntil	šI

SchemaFrozenUntil	šH

AdminScopes"
	šL
 
ReadWriteScopes"
	šK


ReadScopes"
	šJB\
docU"SRestriction contains scope access restriction and frozen times for schema and data.2\
KeysT

Keys"
	šOB9
doc2"0Keys is JSON result type for getKeys in REST API2v
NamePayloadg

Name	šjBP
docI"GNamePayload is JSON payload on REST API request to update data set name2Ø
DataSetPayloadÅo
/
StartTimeStrB
json"
start-timešf

Name	še
)

JSONSchemaB
json"schemašgBR
docK"IDataSetPayload is JSON payload on REST API request to create new data set2¢
Times˜F
+
Data#"!
B
json"
data-timesšY

Schema"
	šZBN
docG"ETimes contains schema and data times, used to get StartTimes for both2Æ
UpdateEvent¶h

Deleted	šw

Data	šu

	StartTime	št

Key	šs

Schema	švBJ
docC"AUpdateEvent holds all information necessary to post to subscribes2
KeyNamev%

Name	šV

Key	šUBM
docF"DKeyName is JSON result type for get and put dataDetNamre in REST API2u
SchemaPayloadd

Schema	špBK
docD"BSchemaPayload is JSON payload on REST API request to update schema2[
KeyT

Key	šRB>
doc7"5Key is JSON result type for createDataSet in REST API2¦
CreationStartTime3

CreationTime	š]

	StartTime	š^BY
docR"PCreationStartTime contains start and creation time for a schema or data snapshot2m
DataPayload^

Data	šmBG
doc@">DataPayload is JSON payload on REST API request to update data2
SubscriptionŒ<

URLB
json"urlšD

SecreteToken	šEBL
docE"CSubscription holds external endpoint values for change notification2±
Data¨Z

CreationTime	š<
%
JSONDataB
json"dataš;

	StartTime	š:BJ
docC"AData holds JSON data valid from StartTime created at CreationTime2³
CreationTimes¡¿
i
DataaB
json"data-time-mapJ?

	
RestApiCreationTimes!map of string:CreationStartTimeša
R
SchemaHJ?

	
RestApiCreationTimes!map of string:CreationStartTimešbB]
docV"TCreationTimes contains schema and data times maps, used to StartTime to CreationTims2Á
Schema¶^

CreationTime	šA
)

JSONSchemaB
json"schemaš@

	StartTime	š?BT
docM"KSchema holds JSON Schema to validate Data against, a name for key creation.