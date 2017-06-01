---
category: Screen
apiurl: '/api/v1/dashboard/screen_clone'
title: "Clone Screen & graphs of this Screen"
type: 'PUT'
sample_doc: ''
layout: default
---

* [Session](#/authentication) Required
* id
  * screen_id [copy base on]
  * required

### Request

```
{"id": 965}
```

### Response

```Status: 200```
```
{
  "graph_names": [
    "net.if.total.bytes",
    "net.if.total.bytes"
  ],
  "id": 1276,
  "pid": 0,
  "name": "net.if.total.bytes_copy",
  "creator": "root"
}
```
