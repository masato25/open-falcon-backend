---
category: Screen_Graph
apiurl: '/api/v1/dashboard/graph_clone'
title: "Clone Graph by ID"
type: 'POST'
sample_doc: ''
layout: default
---

* [Session](#/authentication) Required
* id
  * graph id
  * int64 [required]
* name
  * string
  * new name of copied screen

### Request

```
{
  "id": 4626,
  "name": "CPU_copy"
}
```

### Response

```Status: 200```
```
{
  "id": 4993,
  "title": "CPU_copy",
  "hosts": "agent227|agent228|agent230",
  "counters": "cpu.guest|cpu.idle|cpu.iowait|cpu.irq|cpu.nice|cpu.softirq|cpu.steal|cpu.switches|cpu.system|cpu.user",
  "screen_id": 955,
  "timespan": 3600,
  "graph_type": "h",
  "method": "sum",
  "position": 4626,
  "falcon_tags": "",
  "creator": "root"
}
```
