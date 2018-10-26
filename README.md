# ghost
Go Storage hack project

Playing around with distributed systems in Go.  Basic implementation of chain replication object storage. 

To run 
```
go install ./src/cmd/...
docker-compose up
```
to upload or download a file
```
./bin/cp foo.txt doge://cloud/bar.txt
./bin/cp doge://cloud/bar.txt foobar.txt
```
(yes, `doge` is the official protocol for ghost)
