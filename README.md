# Ballot app
This is ballot micro-service sample
## About project

Voting application contains various frontend and backend microservices. These microservices are deployed and are available over ingress in Roost Cluster.

### Voter

An frontend application written in node to allow participants to vote.

Depends on: ballot and ecserver services

### Ballot

An backend app for voter written in Golang, to store the votes.

### Election Commission

An frontend to manage the election candidates and uses ecserver as backend to store candidates lists.

Depends on: ecserver

### ECserver

An backend app written in Golang for election-commission to store list of candidates.

## How to deploy

Right-click on [Makefile](./Makefile) and choose Run

## How to access application

Pattern: http://$namespace.$serviceName.$clusterPublicIP.nip.io

Voter: [default.voter.10.10.0.10.nip.io](http://default.voter.10.10.0.10.nip.io)

ElectionCommission: [default.ec.10.10.0.10.nip.io](default.ec.10.10.0.10.nip.io)

## How to test deployed app

Build and deploy service-test-suite in roost cluster.
Roost intelligently identifies service dependencies. So whenever dependent service is modified, specified test suite is triggered.
In event of building ballot image or restart of the ballot app, service test suite would be triggered and fitness events can be seen from event viewer ( Observability -> Service Fitness -> Fitness Event).
