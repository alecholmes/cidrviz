# cidrviz

`cidrviz` is a command line tool to help visualize overlapping CIDR blocks.

```
# A and B partially overlap
cidrviz A=1.2.3.4/8 B=1.1.1.1/10 C=2.2.2.2/20

1.0.0.0         AB-
1.63.255.255    AB-
1.255.255.255   A--
2.2.0.0         --C
2.2.15.255      --C
```

## Installation

To install `cidrviz`, run the following steps. These steps assume that Go is installed.

1. `git clone git@github.com:alecholmes/cidrviz.git`
1. `cd cidrviz`
1. `go install`
