# gofish

gofish reverses points to countries. It supports single point reversal via a web api or batch reversal given a file of points.

## Usage

Make sure you have the [Two Fishes](http://twofishes.net) server running. 

Start the gofish server: 

```bash
$ go build
$ ./gofish -a :9090 
2013/09/23 11:12:08 server listening on :9090..
```

Point your browser to http://localhost:9090/reverse?ll=59.329,18.068

## Batch reverse

Create a file of points called `points.txt`, one point per row ([lat]\t[lng]\n):

```
59.3137 18.0669
52.5206 13.4026
59.3277 18.0087
47.625  -122.515
```

Batch reverse the points, group results by country:

```
$ ./gofish -batch.points points.txt
2 Sweden
1 United States
1 Germany
time 5.403775ms
```

Reverse without grouping:

```
$ ./gofish -batch.points points.txt -batch.group=false
59.3137 18.0669 Sweden
59.3277 18.0087 Sweden
52.5206 13.4026 Germany
47.625 -122.515 United States
time 5.577933ms
```

As advertised at http://twofishes.net the performance is at least 1000 reversals/s. I've even seen it come close to 2000 reversals/s.

## Geonames.org attribution

gofish uses country data from http://geonames.org to convert country codes to country names.