�;
RestApi�;
	
RestApi"
	interface"Storer"R
interface_docA"?Storer abstracts all required RefData persistence and retrieval"
json_property_separator"-*�
GET /api/{key}/{startTime}�
GET /api/{key}/{startTime}"!
method_name"GetDataWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
DataBH/api/{key}/{startTime}
	startTime	� 
key	� *�
	POST /api�
	POST /api"
method_name"CreateDataSet"
patterns
:
"rest"

middleware"AuthorizeRoot:B
KeyB/apiJ
dsJ

DataSetPayload*�
DELETE /api/{key}/{startTime}�
DELETE /api/{key}/{startTime}"
method_name"
DeleteData"
patterns
:
"rest" 

middleware"AuthorizeDataSet:

returnBH/api/{key}/{startTime}
	startTime	� 
key	� *�
GET /api/{key}/name�
GET /api/{key}/name"
method_name"GetDataSetName"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B	
KeyNameB`/api/{key}/name
key	�  7
time-J$

	
RestApi{queryTime<:string}�  *�
POST /api/admin/{key}/subscribe�
POST /api/admin/{key}/subscribe" 
method_name"PutSubscription"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
SubscriptionB2/api/admin/{key}/subscribe
key	�3 J
sJ

Subscription*�
PUT /api/{key}�
PUT /api/{key}"
method_name	"PutData"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
DataB"
/api/{key}
key	�
 J
dpJ

DataPayload*�
!PUT /api/admin/{key}/restrictions�
!PUT /api/admin/{key}/restrictions"
method_name"PutRestriction"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
RestrictionB5/api/admin/{key}/restrictions
key	�0 J
rJ

Restriction*�
PUT /api/{key}/{startTime}�
PUT /api/{key}/{startTime}"!
method_name"PutDataWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
DataBH/api/{key}/{startTime}
	startTime	� 
key	� J
dpJ

DataPayload*�
GET /api/{key}�
GET /api/{key}"
method_name	"GetData"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
DataB[
/api/{key}
key	� 7
time-J$

	
RestApi{queryTime<:string}� *�
GET /api/{key}/schema�
GET /api/{key}/schema"
method_name"	GetSchema"
patterns
:
"rest" 

middleware"AuthorizeDataSet:
B
SchemaBb/api/{key}/schema
key	� 7
time-J$

	
RestApi{queryTime<:string}� *�
#GET /api/admin/{key}/creation-times�
#GET /api/admin/{key}/creation-times"!
method_name"GetCreationTimes"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
CreationTimesB7/api/admin/{key}/creation-times
key	�+ *�
PUT /api/{key}/schema�
PUT /api/{key}/schema"
method_name"	PutSchema"
patterns
:
"rest" 

middleware"AuthorizeDataSet:
B
SchemaB)/api/{key}/schema
key	� J
spJ

SchemaPayload*�
!PUT /api/{key}/schema/{startTime}�
!PUT /api/{key}/schema/{startTime}"#
method_name"PutSchemaWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:
B
SchemaBO/api/{key}/schema/{startTime}
	startTime	� 
key	� J
spJ

SchemaPayload*�
!GET /api/{key}/schema/{startTime}�
!GET /api/{key}/schema/{startTime}"#
method_name"GetSchemaWithStart"
patterns
:
"rest" 

middleware"AuthorizeDataSet:
B
SchemaBO/api/{key}/schema/{startTime}
	startTime	� 
key	� *�
GET /api�
GET /api"
method_name	"GetKeys"
patterns
:
"rest"

middleware"AuthorizeRoot"

method_doc	"DataSet:B
KeysB/api*�
!GET /api/admin/{key}/restrictions�
!GET /api/admin/{key}/restrictions"
method_name"GetRestriction"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B
RestrictionB5/api/admin/{key}/restrictions
key	�. *�
PUT /api/{key}/name�
PUT /api/{key}/name"
method_name"PutDataSetName"
patterns
:
"rest" 

middleware"AuthorizeDataSet:B	
KeyNameB'/api/{key}/name
key	�" J
npJ

NamePayload*�
!POST /api/admin/{key}/unsubscribe�
!POST /api/admin/{key}/unsubscribe"#
method_name"DeleteSubscription"
patterns
:
"rest" 

middleware"AuthorizeDataSet:

returnB4/api/admin/{key}/unsubscribe
key	�6 J
sJ

Subscription*�
DELETE /api/admin/{key}�
DELETE /api/admin/{key}"
method_name"DeleteDataSet"
patterns
:
"rest" 

middleware"AuthorizeDataSet:

returnB(/api/admin/{key}
key	�% *�
$DELETE /api/{key}/schema/{startTime}�
$DELETE /api/{key}/schema/{startTime}"
method_name"DeleteSchema"
patterns
:
"rest" 

middleware"AuthorizeDataSet:

returnBO/api/{key}/schema/{startTime}
	startTime	� 
key	� *�
 GET /api/admin/{key}/start-times�
 GET /api/admin/{key}/start-times"
method_name"GetStartTimes"
patterns
:
"rest" 

middleware"AuthorizeDataSet:	B
TimesB4/api/admin/{key}/start-times
key	�( 2�
Restriction��

DataFrozenUntil	�I

SchemaFrozenUntil	�H

AdminScopes"
	�L
 
ReadWriteScopes"
	�K


ReadScopes"
	�JB\
docU"SRestriction contains scope access restriction and frozen times for schema and data.2\
KeysT

Keys"
	�OB9
doc2"0Keys is JSON result type for getKeys in REST API2v
NamePayloadg

Name	�jBP
docI"GNamePayload is JSON payload on REST API request to update data set name2�
DataSetPayload�o
/
StartTimeStrB
json"
start-time�f

Name	�e
)

JSONSchemaB
json"schema�gBR
docK"IDataSetPayload is JSON payload on REST API request to create new data set2�
Times�F
+
Data#"!
B
json"
data-times�Y

Schema"
	�ZBN
docG"ETimes contains schema and data times, used to get StartTimes for both2�
UpdateEvent�h

Deleted	�w

Data	�u

	StartTime	�t

Key	�s

Schema	�vBJ
docC"AUpdateEvent holds all information necessary to post to subscribes2�
KeyNamev%

Name	�V

Key	�UBM
docF"DKeyName is JSON result type for get and put dataDetNamre in REST API2u
SchemaPayloadd

Schema	�pBK
docD"BSchemaPayload is JSON payload on REST API request to update schema2[
KeyT

Key	�RB>
doc7"5Key is JSON result type for createDataSet in REST API2�
CreationStartTime�3

CreationTime	�]

	StartTime	�^BY
docR"PCreationStartTime contains start and creation time for a schema or data snapshot2m
DataPayload^

Data	�mBG
doc@">DataPayload is JSON payload on REST API request to update data2�
Subscription�<

URLB
json"url�D

SecreteToken	�EBL
docE"CSubscription holds external endpoint values for change notification2�
Data�Z

CreationTime	�<
%
JSONDataB
json"data�;

	StartTime	�:BJ
docC"AData holds JSON data valid from StartTime created at CreationTime2�
CreationTimes��
i
DataaB
json"data-time-mapJ?

	
RestApiCreationTimes!map of string:CreationStartTime�a
R
SchemaHJ?

	
RestApiCreationTimes!map of string:CreationStartTime�bB]
docV"TCreationTimes contains schema and data times maps, used to StartTime to CreationTims2�
Schema�^

CreationTime	�A
)

JSONSchemaB
json"schema�@

	StartTime	�?BT
docM"KSchema holds JSON Schema to validate Data against, a name for key creation.