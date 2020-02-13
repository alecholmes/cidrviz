# cidrviz

`cidrviz` is a command line tool to help visualize overlapping subnet ranges.

```
# A and B partially overlap
cidrviz A=1.2.3.4/8 B=1.1.1.1/10 C=2.0.0.0/31 D=2.2.2.2/20

1.0.0.0         AB--
1.63.255.255    AB--
1.255.255.255   A---
2.0.0.0         --C-
2.0.0.1         --C-
                ---- (131070 IPs)
2.2.0.0         ---D
2.2.15.255      ---D
```

Interpreting the above:
* The first IP of subnet A is 1.0.0.0 and the last IP, inclusive, is 1.63.255.255.
* Subnet B overlaps part of subnet A, starting at the same IP of 1.0.0.0. The last IP of B is 1.63.255.255.
* Subnets A and C are adjacent. The first IP of C is the next IP after the end of subnet A.
* Subnets C and D are not adjacent. There are 131,070 IP addresses after C ends and before D begins.

## Installation

To install `cidrviz`, run the following steps. These steps assume that Go is installed.

1. `git clone git@github.com:alecholmes/cidrviz.git`
1. `cd cidrviz`
1. `go install`
