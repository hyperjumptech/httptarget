# HttpTarget

**HttpTarget** is a very simple, small and lightweight HTTP server that would be helpful for http client development tool.
Simply start the server and it'll be ready to accept incoming http requests. It can easily simulate heavy server side load
by implementing random delay range, or simulate any kind of http response code.  
 
It has the following features:

1. Multiple URL path support, each of them is configurable for their response code, response time, response body and header.
2. Add, modify, remove URL path in realtime without the need to restart the server.
3. Simple API (equiped with OpenAPI 2.0 Swagger) for you app to integrate straight to the server for their testing purpose.

## Get HttpTarget

You can obtain **HttpTarget** binary in 2 ways. By simply download the released binary for 
your platform or you can always build the binary your self.

### Download the HttpTarget binary

Your can download **HttpTarget** binary from the [release page on Github](https://github.com/hyperjumptech/httptarget/releases).
There you download the binary for your platform:

- Windows : httptarget-windows.tar.gz
- Linux : httptarget-linux.tar.gz
- MacOS : httptarget-macos.tar.gz

You can then extract the downloaded tarball (.tar.gz) for the executable. The executable will run straight-away without needing any dependency. 

### Build HttpTarget binary

To build the binary, you to install the following apps.

- Git client
- Golang 1.16 Compiler

Once done, assuming you know basic GIT, you can do the following steps.

1. Clone the httptarget project from [github repository](https://github.com/hyperjumptech/httptarget).
2. From the clone directory you can simply execute `make build-linux` to build linux binary, `make build-windows` to build windows binary or `make build-macos` to build mac os binary. Or simply build them all using `make build-all`
3. Once done, a new directory `build` is created contains your compressed binary.
4. Extract the compressed binary (tarball .tar.gz) to get your executable.

If you don't have `gnumake` application to run `make`, you can build the binary straight away using `go build -o . ./...`

## Binary Usage

Now you have the executable binary. To run the server you can simply run the executable

```bash
$ httptarget.app
INFO[0000] Added test endpoint on [/], code 200, minDelay 0 ms, maxDelay 200 ms 
INFO[0000] Server listening at 0.0.0.0:51423 
```

By default, the server will listen on port `51423`. To stop the server simply hit `ctrl+c`

There are few argument available when you start the server. The argument is usefull to 
configure the server for:

- Listening on speciffic port
- Bind to speciffic host IP (network interface)
- The initial __endpoint__ to service, including its response code, body, header
- The initial responseTime range (to simulate slow response due to load or network problem)

To see all possible argument, simply put `-help` argument.

```bash
$ httptarget.app -help
Usage of httptarget:
  -body string
        HTTP response body (default "OK")
  -code int
        Response code (default 200)
  -h string
        Bind host (default "0.0.0.0")
  -help
        Display this help message
  -maxDelay int
        Maximum Delay Millisecond (default 200)
  -minDelay int
        Minimum Delay Millisecond
  -p int
        Listen port (default 51423)
  -path string
        Base path (default "/")
```

## API to Manage or to Integrate to the Server.

**HttpTarget** is implementing OpenAPI standard. The OpenAPI Swagger documentation is immediately available
as soon as you start the server.

```bash
$ ./httptarget.app 
INFO[0000] Added test endpoint on [/], code 200, minDelay 0 ms, maxDelay 200 ms 
INFO[0000] Server listening at 0.0.0.0:51423 
```

There you notice that it's starting on the server that bind to any interface (0.0.0.0), on port `51423`. Thus you can open your
favourite web-browser and go to `http://localhost:51423/docs/index.html`. It will open the swagger API spec page where 
you can configure your server or integrate your testing with. The __OpenAPI Swagger__ page will tell you everything about
__API Endpoints__ available for you to use, what __method__, __parameters__, __URL__, etc.

## All configuration stuff is too complicated, Let me test straight away !!

**HttpTarget** configuration (swagger etc) is indeed not for every body ;). However, for you who don't want to sweat the configuration, 
you can straight away uses the following 2 built-in endpoint groups. (oh yeah, they accept any `method`)

**Delay Endpoints**

They have the following path patterns

- `/delay/{delay}` This endpoint will simply delay it-self for __{delay}__ milliseconds and return http response code `200`
- `/delay/{minDelay}/{maxDelay}` This endpoint will simply delay it-self for between __{minDelay}__ and __{maxDelay}__ milliseconds and return http response code `200`

for example :

```text
>> GET http://localhost:51423/delay/1000/2000
<< 200 OK
<< OK, delayed between 1000ms and 2000ms : 1120ms
```

```text
>> GET http://localhost:51423/delay/1000
<< 200 OK
<< OK, delayed between 1000ms and 1000ms : 1000ms
```

**Code Endpoints**

They have the following path pattern

- `/code/{responsecode}` This endpoint will simply http response code __{responsecode}__
- `/code/{responsecode}/{delay}` This endpoint will simply delay it-self for __{delay}__ milliseconds and return http response code __{responsecode}__
- `/code/{responsecode}/{minDelay}/{maxDelay}` This endpoint will simply delay it-self for between __{minDelay}__ and __{maxDelay}__ milliseconds and return http response code __{responsecode}__

for example :

```text
>> GET http://localhost:51423/code/202/2000
<< 202 Accepted
<< OK, delayed between 2000ms and 2000ms : 2000ms
```

```text
>> GET http://localhost:51423/code/500/2000/3000
<< 500 Internal Server Error
<< OK, delayed between 2000ms and 3000ms : 2630ms
```


Finally, Happy Testing from Hyperjump team !!!