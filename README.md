# DynDNS.com CheckIP in golang by Josimar Andrade

For test: http://67.205.183.189:82

DynDNS.com runs a simple web service for detecting a user's remote IP address called
You can demo the web service at http://checkip.dyndns.com/

The goal of this challenge is to implement a new working CheckIP application server and in
addition include the location and time-zone.

* Requirements
    Written in GO;
    Return an HTML output the IP, time-zone, and location;
    Use User Agents X-Forwarded-For to get the IP, otherwise session connection IP;

* Workflow
```
          The browser                         The Aplication                     The browser
                                      +                          +
                                      |                          |
                                      |                          |
+------------+      +-------------+   |    +--------------+      |  +---------------+      +-------------+
|            |      |             |   |    |              |      |  |Send HTML      |      |             |
|   START    |      |Access the   |   |    | Get the IP   |      |  |Response with  |      |             |
|            +----->+web service  +------->+              |      |  |the information+----->+     END     |
|            |      |end point    |   |    |              |      |  |               |      |             |
+------------+      +-------------+   |    +-------+------+      |  +-------+-------+      +-------------+
                                      |            |             |          ^
                                      |            |             |          |
                                      |    +-------v------+      |          |
                                      |    |              |      |          |
                                      |    | Retrive IP   |      |          |
                                      |    | Geolocation  +-----------------+
                                      |    | (ipapi.co)   |      |
                                      |    +--------------+      |
                                      |                          |
                                      +                          +
```
NOTE: All the development was made on Ubuntu 18.04.3 LTS x86_64

## File directory
```
└── checkIP
    ├── challenge               --> challenge description
    ├── checkIP                 --> Binary executable
    ├── checkIP.go              --> program code
    ├── checkIP_test.go         --> test code
    ├── docker-compose.yml      --> (docker configuration)
    ├── Dockerfile              --> (docker configuration)
    ├── k8s-deployment.yaml     --> (Kubernetes configuration)
    ├── k8s-pod.yml             --> (Kubernetes configuration)
    ├── k8s-service.yml         --> (Kubernetes configuration)
    └── README                  --> documentation
```
1 directory, 10 files

## Tests
```
All                                         --> go test -run ''
    Health server status                    --> go test -run webHealth
    Internet connectivity to ipapi.co       --> go test -run ipapiAPIconectivity
    Command line test (Unit test)           --> go test -run cmd
    ipapi.co API test (API test)            --> go test -run ipInfo
    Web request test  (Integration test)    --> go test -run webHome
```
## Installation (build and run) single server.
```
Build:
    go build checkIP.go
run:
    ./checkIP        (start the server)
    ./checkIP <IP>   (get a geolocation for a sigle IP
```
## Usage

change server port by environment variable (default: 8080):
    export PORT=<DESIRED_PORT_NUMBER>
Server health check (should return "up" if the server is up):
    http://localhost:8080/health
Browser:
    http://localhost:8080
Curl test:
    curl --header "X-Forwarded-For: 8.8.8.8" http://localhost:8080

## Scalability with Kubernetes

-- Docker
Build (creat the executable):
    GOOS=linux GOARCH=amd64 go build
Build a docker image:
    docker build -t checkip:1 .
Build the service:
    docker-compose build
Start container:
    docker-compose up -d

-- Kubernetes
Create a kubernetes Pod:
    kubectl create -f k8s-pod.yml
Access:
    kubectl port-forward checkip 8080:8080
Deployment our pods (nodes):
    kubectl create -f k8s-deployment.yaml
Create a kubernetes Services to handle the nodes:
    kubectl create -f k8s-service.yml
    kubectl get services (info)


-- Scale the nodes (horizontal scaling, add 10 more nodes):
    kubectl scale deployment checkip --replica=10
dashboard of pods:
    kubectl get pods -w

## High performance Test

In my computer one node can handle with 65 concurrent request.

* Host info:
    OS: Ubuntu 18.04.3 LTS x86_64
    CPU: Intel i7-4702MQ (8) @ 3.200GHz
    Memory: 3974MiB / 7867MiB

* Measuring performance

    jandrade@JAFA-PC:~$ ab -n 65 -c 65 -H "X-Forwarded-For:8.8.8.8" http://192.168.99.100:30518/
    This is ApacheBench, Version 2.3 <$Revision: 1807734 $>
    Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
    Licensed to The Apache Software Foundation, http://www.apache.org/

    Benchmarking localhost (be patient).....done


    Server Software:
    Server Hostname:        localhost
    Server Port:            8080

    Document Path:          /
    Document Length:        183 bytes

    Concurrency Level:      65
    Time taken for tests:   0.706 seconds
    Complete requests:      65
    Failed requests:        0
    Total transferred:      12000 bytes
    HTML transferred:       7320 bytes
    Requests per second:    56.69 [#/sec] (mean)
    Time per request:       705.543 [ms] (mean)
    Time per request:       17.639 [ms] (mean, across all concurrent requests)
    Transfer rate:          16.61 [Kbytes/sec] received

    Connection Times (ms)
                  min  mean[+/-sd] median   max
    Connect:        1    2   0.3      2       2
    Processing:   510  647  44.2    666     703
    Waiting:      510  647  44.2    666     703
    Total:        512  648  44.2    668     705

    Percentage of the requests served within a certain time (ms)
      50%    668
      66%    676
      75%    680
      80%    682
      90%    691
      95%    699
      98%    705
      99%    705
     100%    705 (longest request)
    jandrade@JAFA-PC:~$


## License
Copyright (c) 2019. Josimar Andrade, No Rights Reserved
