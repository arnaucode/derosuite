- blockchain is stored in a boltdb, in tmp/derod_database.db

- the seeds are hardcoded in p2p/controller.go, line 61:
```go
// add hard-coded seeds
end_point_list = append(end_point_list, "212.8.242.60:18090")
```
And also can be setted in the parameter:
```
derod --add-exclusive-node
```

- mainnetwork is defined in config/config.go line 91
	- NETWORK_ID
        - Line 92 can be substitued by:
        ```go
        Network_ID:                       uuid.FromStringOrNil("afdb5368-641d-41a2-8bc2-d2e825e25c46"),
        ```
        Where the string is a uuid string
	- Genesis_Block_Hash
	- Genesis_Tx

- blocktime: config/config.go line 24