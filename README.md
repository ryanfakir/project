# Nmap Practice

### Installation

nmap-lookup is UI, practice is backend.
In order to run the application, you need NPM installed, cd to nmap-lookup folder.

```sh
$ cd nmap-lookup
$ npm install 
$ npm start
```
In practice folder, If you want to check code and log. put the folder under your **$GOPATH/src**. 
Before start application, application assume you install **nmap** commmand line tool and using local mysql db store entry, you can change **config.go** to push your credentials. 
Default will be username: root, password : root, after it start, it will create ip database for saving IP ports. 
After that you can run application:

```sh
$ cd practice
$ ./main
```
you can build application after modify it.
```sh
cd practice
go build main.go
```
you can also install application if you like.
```sh
go install
```
### Data Flow

UI(POST) ----> Web Server ----> Client(cache with ttl) ---> Service(look up nmap command tool) ---> DB(transaction update)

