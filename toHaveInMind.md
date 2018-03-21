- blockchain is stored in a boltdb, in tmp/derod_database.db

- the seeds are hardcoded in derosuite/p2p/controller.go, line 61:
```go
// add hard-coded seeds
end_point_list = append(end_point_list, "212.8.242.60:18090")
```
And also can be setted in the parameter:
```
derod --add-exclusive-node
```
