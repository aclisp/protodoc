# API Protocol

Table of Contents

* [Service SDK](#service-sdk)
    * [Method SDK.Ready](#method-sdkready)
    * [Method SDK.Allocate](#method-sdkallocate)
    * [Method SDK.Shutdown](#method-sdkshutdown)
    * [Method SDK.Health](#method-sdkhealth)
    * [Method SDK.GetGameServer](#method-sdkgetgameserver)
    * [Method SDK.WatchGameServer](#method-sdkwatchgameserver)
    * [Method SDK.SetLabel](#method-sdksetlabel)
    * [Method SDK.SetAnnotation](#method-sdksetannotation)
* [Enums](#enums)
    * [Enum GameServer.Status.State](#enum-gameserverstatusstate)
* [Objects](#objects)
    * [Object KeyValue](#object-keyvalue)
    * [Object GameServer](#object-gameserver)
    * [Object GameServer.ObjectMeta](#object-gameserverobjectmeta)
    * [Object GameServer.Spec](#object-gameserverspec)
    * [Object GameServer.Spec.Health](#object-gameserverspechealth)
    * [Object GameServer.Status](#object-gameserverstatus)
    * [Object GameServer.Status.Port](#object-gameserverstatusport)




## Service SDK

SDK service to be used in the GameServer SDK to the Pod Sidecar

### Method SDK.Ready

> POST /stable.agones.dev.sdk/SDK/Ready <br/>
> Content-Type: application/json <br/>
> Authorization: Bearer (token) <br/>

Call when the GameServer is ready

Request is empty

Response is empty


### Method SDK.Allocate

> POST /stable.agones.dev.sdk/SDK/Allocate <br/>
> Content-Type: application/json <br/>
> Authorization: Bearer (token) <br/>

Call to self Allocation the GameServer

Request is empty

Response is empty


### Method SDK.Shutdown

> POST /stable.agones.dev.sdk/SDK/Shutdown <br/>
> Content-Type: application/json <br/>
> Authorization: Bearer (token) <br/>

Call when the GameServer is shutting down

Request is empty

Response is empty


### Method SDK.Health

WebSocket client-streaming

> GET /stable.agones.dev.sdk/SDK/Health <br/>

Send a Empty every d Duration to declare that this GameSever is healthy

Request is empty

Response is empty


### Method SDK.GetGameServer

> POST /stable.agones.dev.sdk/SDK/GetGameServer <br/>
> Content-Type: application/json <br/>
> Authorization: Bearer (token) <br/>

Retrieve the current GameServer data

Request is empty

Response parameters

|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
| data | [object GameServer](#object-gameserver) |  |


### Method SDK.WatchGameServer

WebSocket server-streaming

> GET /stable.agones.dev.sdk/SDK/WatchGameServer <br/>

Send GameServer details whenever the GameServer is updated

Request is empty

Response parameters

|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
| data | [object GameServer](#object-gameserver) |  |


### Method SDK.SetLabel

> POST /stable.agones.dev.sdk/SDK/SetLabel <br/>
> Content-Type: application/json <br/>
> Authorization: Bearer (token) <br/>

Apply a Label to the backing GameServer metadata

Request parameters

|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
| kv | [object KeyValue](#object-keyvalue) |  |

Response is empty


### Method SDK.SetAnnotation

> POST /stable.agones.dev.sdk/SDK/SetAnnotation <br/>
> Content-Type: application/json <br/>
> Authorization: Bearer (token) <br/>

Apply a Annotation to the backing GameServer metadata

Request parameters

|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
| kv | [object KeyValue](#object-keyvalue) |  |

Response is empty





## Enums

### enum GameServer.Status.State



Constants

|   Value   |   Name    |  Description |
| --------- | --------- | ------------ |
| 0  | READY | The GameServer is ready to serve |
| 1  | STARTING | The GameServer is starting |
| 2  | SHUTTING_DOWN | The GameServer is shutting down |


## Objects

### object KeyValue

Key, Value entry

Attributes

|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
| key | string |  |
| value | string |  |


### object GameServer

A GameServer Custom Resource Definition object We will only export those resources that make the most sense. Can always expand to more as needed.

Attributes

|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
| object_meta | [object GameServer.ObjectMeta](#object-gameserverobjectmeta) | GameServer meta data |
| spec | [object GameServer.Spec](#object-gameserverspec) | specification |
| status | [object GameServer.Status](#object-gameserverstatus) |  |


### object GameServer.ObjectMeta

representation of the K8s ObjectMeta resource

Attributes

|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
| name | string |  |
| namespace | string |  |
| uid | string |  |
| resource_version | string |  |
| generation | int64 |  |
| creation_timestamp | int64 | timestamp is in Epoch format, unit: seconds |
| deletion_timestamp | int64 | optional deletion timestamp in Epoch format, unit: seconds |


### object GameServer.Spec



Attributes

|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
| health | [object GameServer.Spec.Health](#object-gameserverspechealth) |  |


### object GameServer.Spec.Health



Attributes

|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
| Disabled | bool |  |
| PeriodSeconds | int32 |  |
| FailureThreshold | int32 |  |
| InitialDelaySeconds | int32 |  |


### object GameServer.Status



Attributes

|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
| state | [enum GameServer.Status.State](#enum-gameserverstatusstate) |  |
| address | string |  |
| ports | array of [object GameServer.Status.Port](#object-gameserverstatusport) |  |


### object GameServer.Status.Port



Attributes

|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
| name | string |  |
| port | int32 |  |

