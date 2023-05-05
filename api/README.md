# API Endpoints

## List races

Lists all the races that match the filters provided in request body.

`POST /v1/list-races`

### Request body:

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
      "advertisedStartTime": "2021-03-02T19:16:58Z"
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

---
