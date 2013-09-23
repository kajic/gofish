# gofish

Wrapper around [Two Fishes](http://twofishes.net/) for reversing points to countries. Supports single point reversal using a REST api or batch reversal given a file of points.

## Usage

Make sure you have the Two Fishes server running. Instructions for setting up Two Fishes can be found at http://twofishes.net.

Start the gofish server: 

```bash
$ go build
$ ./gofish -a :9090 
2013/09/23 11:12:08 server listening on :9090..
```

Now point your browser to http://localhost:9090/reverse?ll=59.329,18.068

## Batch reverse

Create a file of points called `points.txt`, one point per row ([lat]\t[lng]\n):

```
59.3137	18.0669
52.5206	13.4026
59.3277	18.0087
47.625	-122.515
```

Batch reverse the points, group by country:

```
$ ./gofish -batch.points points.txt
time 5.403775ms
2 Sweden
1 United States
1 Germany
```

Reverse without grouping:

```
$ ./gofish -batch.points points.txt -batch.group=false
59.3137 18.0669 Sweden
59.3277 18.0087 Sweden
52.5206 13.4026 Germany
47.625 -122.515 United States
```

As advertised at http://twofishes.net the performance is at least 1000 reversals/s. I've seen it come close to 2000 reversals/s.