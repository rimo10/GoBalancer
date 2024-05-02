# Golang Load Balancer
This project implements a load balancer which is a necessicity in todays internet era. Modern webistes might have to deal with hundreds of thousands of requests and thus it's important to route traffic to the suitable server so that the underlying infrastructure can handle the traffic and keep on its buisness ongoing. 
<br>
A load balancer is a tool that is commonly used to coordinate the volume of traffic between the available servers.
A load balancer sits in front of a group of servers which is referred to as Serverpool, and the load distribution is performed based on various load balancing algorithms. 
Our projects implements only RoundRobin scheduling algorithm
<br>
The input is given through the configuration file `config.yaml`. You can change it as per your requirements
``` yaml
//port in which server will be running
lb_port: 3000
algorithm: roundroubin // for now roundrobin works by default
//available backend servers
backends:
  - "http://localhost:5100"
  - "http://localhost:5200"
  - "http://localhost:5300"
  - "http://localhost:5400"
  - "http://localhost:5500"

```
## Future Works
* Implement the various other scheduling algorithms
* Add some security measures 
## References
* ([load_balancer](https://github.com/leonardo5621/golang-load-balancer/tree/master))
* ([simplelb](https://github.com/kasvith/simplelb))
