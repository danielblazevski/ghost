# ghost
Go Storage hack project

Playing around with distributed systems in Go.  Basic implementation of chain replication object storage. 

To run 
```
go install ./src/cmd/...
docker-compose up
```
to upload a file
```
./bin/cp cp foo.txt doge://cloud/bar.txt
```
(yes, `doge` is the official protocol for ghost)
