---
category: Screen_Graph
apiurl: '/api/v1/dashboard/graph/:id'
title: "Update Graph by id"
type: 'PUT'
sample_doc: ''
layout: default
---

* [Session](#/authentication) Required
* id
  * graph id
  * int64 [required] - in url path
* title
  * string
  * name of graph
* endpoints
  * []string
* counters
  * []string
* timespan
  * int64
  * 时间区段 (秒)
    * default 3600
* graph_type
  * string
  * 视角: h (endpoint view), k (counter view), a (combo view)
  * accept values:
    * h
    * a
    * k
* method
  * string
  * accept values:
    * 'sum'
    * '' (空值)
* position
  * int64
  * 排序
    * 预设值為0
* falcon_tags
  * string
  * owl-light not this concept, keep empty. (open-falcon only)

### Request

```
/api/v1/dashboard/graph/4246
{
  "title": "testtt",
  "screen_id": 955,
  "endpoints": ["a1","a2"],
  "counters": ["c1","c2"],
  "timespan":9999,
  "graph_type":"a",
  "method": "sum",
  "graph_type": "h",
  "position": 22,
  "falcon_tags": "a=1,b=2"
}
```

### Response

```Status: 200```
```
{
  "graph_id": 4626,
  "title": "testtt",
  "screen_id": 955,
  "endpoints": [
    "a1",
    "a2"
  ],
  "counters": [
    "c1",
    "c2"
  ],
  "timespan": 9999,
  "graph_type": "a",
  "method": "sum",
  "position": 22,
  "falcon_tags": "a=1,b=2"
}
```
