# golangAPI-example
I intend this to be a simple example of a golang db API with simple caching for beginners and people new to golang to refer to if they so choose. The patterns and concepts are there and should be interchagable with other caching and databases.

NB: Go-cache should be inter changagable with something more suited to distributed computing like redis. Go-Cache is only suited to single VM applications and should be switched out if this is to be used in a horizontally scalable application. I intended to create this with as little extra set up as possible. 

## Set up in Ubuntu

### Installing Golang
Download the archive into /usr/local/
https://golang.org/dl/

`tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz`

Add /usr/local/go/bin to the PATH environment variable

`export PATH=$PATH:/usr/local/go/bin`

### Installting mySQL
For my root password for mysql I just used "root" for simplicity, you can change this in the code in cabTrips.go manually if theres discrepecies with your existing password.
```
sudo apt-get update
sudo apt-get install mysql-server
```
## Running the application
###Install the application dependencies
```
go get -u github.com/go-sql-driver/mysql
go get github.com/patrickmn/go-cache
```
### Setting up the database
Log into the Database, sql runs after set up so use the password you used for the root user here
```
mysql -u root -p
Enter password:************
mysql> CHREATE DATABASE cabdata
mysql> USE cabdata
mysql> exit
```
Now load the data into the DB
```
cd /golangAPI-example
mysql -u root -p cabdata < cabData.sql
Enterpassword: ***************

```
This will take a while, just wait for it to return to bash cli
### Run the application
```
cd /src
go run *
```
The go server is now serving on port 8080.
## Making Calls to the application
### /healthcheck
Checks if service is still active
### /flushCache
Clears the in cache memory
###


