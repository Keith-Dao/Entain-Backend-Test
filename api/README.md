# API Endpoints

## List races

Lists all the races that match the filters provided in request body.

### Request

`POST /v1/list-races`

#### Body

```json
{
  "filter": {
    "meeting_ids": [1, 2, 3],
    "only_show_visible": true
  },
  "sort": [
    {
      "column": "meeting_ids",
      "is_descending": true
    }
  ]
}
```

**filter** - JSON object  
The filters to apply when querying the races. If filter is omitted, no filter is applied.

&emsp; **meeting_ids** - list of ints  
&emsp; List of race meeting ids to be matched.  
&emsp; If left empty or undefined, get all races.

&emsp; **only_show_visible** - bool  
&emsp; When set to `true`, only return races that have `visible` set to true.

---

**sort** - list of JSON objects  
The columns to sort the queries and whether the respective column should be in descending order. The sorting order will follow the order provided in the request. If sort is omitted or empty, sort will be defaulted to `advertised_start_time` in ascending order.

&emsp; **column** - str  
&emsp; The name of the column to sort by in snake case. An error would occur if omitted.

&emsp; **is_descending** - bool  
&emsp; When set to `true`, the rows will sort using this row in descending order.  
&emsp; If omitted or set to `false`, the rows will be sorted in ascending order.

---

### Response:

```json
{
  "races": [
    {
      "id": "2",
      "meetingId": "1",
      "name": "Connecticut griffins",
      "number": "12",
      "visible": true,
      "advertisedStartTime": "2021-03-02T19:16:58Z",
      "status": "CLOSED"
    }
  ]
}
```

**races** - list of JSON objects  
List of all races that match the conditions provided by the filter.

&emsp; **id** - int  
&emsp; The race id.

&emsp; **meetingId** - int  
&emsp; The race's meeting id.

&emsp; **name** - string  
&emsp; The name of the race.

&emsp; **number** - int  
&emsp; The race's number.

&emsp; **visible** - bool  
&emsp; The race's visibility.

&emsp; **advertisedStartTime** - timestamp  
&emsp; The race's advertised start time.

&emsp; **status** - string  
&emsp; The race's status, which will either be `OPEN` if the `advertisedStartTime` is after the current time, other it will be `CLOSED`.

---

## Get race

Gets the race with a given id, if it exists.

### Request

`GET /v1/race/{id}`

**id** - int
The race id.

---

### Response

```json
{
  "id": "2",
  "meetingId": "1",
  "name": "Connecticut griffins",
  "number": "12",
  "visible": true,
  "advertisedStartTime": "2021-03-02T19:16:58Z",
  "status": "CLOSED"
}
```

**id** - int  
The race id.

**meetingId** - int  
The race's meeting id.

**name** - string  
The name of the race.

**number** - int  
The race's number.

**visible** - bool  
The race's visibility.

**advertisedStartTime** - timestamp  
The race's advertised start time.

**status** - string  
The race's status, which will either be `OPEN` if the `advertisedStartTime` is after the current time, other it will be `CLOSED`.

---

## List events

Lists all the sport events that match the filters provided in request body.

### Request

`POST /v1/list-events`

#### Body

```json
{
  "filter": {
    "sports": ["basketball", "tennis"],
    "only_show_visible": true
  },
  "sort": [
    {
      "column": "sports",
      "is_descending": true
    }
  ]
}
```

**filter** - JSON object  
The filters to apply when querying the events. If filter is omitted, no filter is applied.

&emsp; **sports** - list of strings  
&emsp; List of event sports to be matched.  
&emsp; If left empty or undefined, get all events.

&emsp; **only_show_visible** - bool  
&emsp; When set to `true`, only return races that have `visible` set to true.

---

**sort** - list of JSON objects  
The columns to sort the queries and whether the respective column should be in descending order. The sorting order will follow the order provided in the request. If sort is omitted or empty, sort will be defaulted to `advertised_start_time` in ascending order.

&emsp; **column** - str  
&emsp; The name of the column to sort by in snake case. An error would occur if omitted.

&emsp; **is_descending** - bool  
&emsp; When set to `true`, the rows will sort using this row in descending order.  
&emsp; If omitted or set to `false`, the rows will be sorted in ascending order.

---

### Response:

```json
{
  "events": [
    {
      "id": "2",
      "sport": "Tennis",
      "name": "Connecticut griffins",
      "number": "12",
      "visible": true,
      "advertisedStartTime": "2021-03-02T19:16:58Z",
      "status": "CLOSED"
    }
  ]
}
```

**events** - list of JSON objects  
List of all sport events that match the conditions provided by the filter.

&emsp; **id** - int  
&emsp; The event id.

&emsp; **sport** - string  
&emsp; The event's sport type.

&emsp; **name** - string  
&emsp; The name of the event.

&emsp; **number** - int  
&emsp; The event's number.

&emsp; **visible** - bool  
&emsp; The event's visibility.

&emsp; **advertisedStartTime** - timestamp  
&emsp; The event's advertised start time.

&emsp; **status** - string  
&emsp; The event's status, which will either be `OPEN` if the `advertisedStartTime` is after the current time, other it will be `CLOSED`.

---

## Get event

Gets the event with a given id, if it exists.

### Request

`GET /v1/event/{id}`

**id** - int
The event id.

---

### Response

```json
{
  "id": "2",
  "sport": "Tennis",
  "name": "Connecticut griffins",
  "number": "12",
  "visible": true,
  "advertisedStartTime": "2021-03-02T19:16:58Z",
  "status": "CLOSED"
}
```

**id** - int  
The event id.

**sport** - string  
The event's sport type.

**name** - string  
The name of the event.

**number** - int  
The event's number.

**visible** - bool  
The event's visibility.

**advertisedStartTime** - timestamp  
The event's advertised start time.

**status** - string  
The event's status, which will either be `OPEN` if the `advertisedStartTime` is after the current time, other it will be `CLOSED`.
