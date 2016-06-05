# goson

Simple CLI tool for getting any data from __anywhere__ (ambitious).

This project is only a work in progress as of now.

## Current Status (Support for JSON)

Currently `goson` works only with JSON encoded data. But the intention is to
make it work with almost all kinds of standard serialization formats.

## Design

* The reason for existence of `goson` is to provide easy tooling to work with
  standard data formats like JSON, YAML, TOML, etc. and make it really easy to
  use them with exiting CLI tooling available to us today.
* Faster iteration without writing actual code.
* As of now, `goson` reads data only from `STDIN`. The primary reason it was
  made was that it would be fed data by some other command like `curl` making
  API calls to other services. However, this may change in the future.

## Installation

	go get github.com/grofers/goson

## Examples

Imagine you have to use AWS CLI to find the list of instances belonding to a
particular load balancer.

    $ aws elb describe-load-balancers --load-balancer-names api-elb \
        | goson --foreach LoadBalancerDescriptions --asitem lb \
            --foreach lb.Instances --asitem instance \
            get instance.InstanceId
    i-6c6bcca0
    i-19088b97
    i-1a088b94
    i-e6216b2a
    i-b8a36777

Another more practical example, which gets the list of public IP address of all
the nodes behind an ELB:

    $ aws ec2 describe-instances \
        --instance-ids $(aws elb describe-load-balancers \
                            --load-balancer-names api-elb \
                                | goson --foreach LoadBalancerDescriptions \
                                    --asitem lb --foreach lb.Instances \
                                    --asitem instance get instance.InstanceId \
                                | tr '\n' ' ') \
        | goson --foreach Reservations --asitem reservation \
            --foreach reservation.Instances --asitem instance \
            get instance.PrivateIpAddress
    172.31.253.74
    172.31.252.123
