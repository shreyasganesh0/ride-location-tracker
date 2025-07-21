# Redis GeoSpatial Commands

## GEOADD And GEORADIUS
```
GEOADD key [NX | XX] [CH] longitude latitude member [longitude latitude member ...]
```
- adds (longtidue latitude name) to the key
- O(log n) complexity for insertion
- queried with GEOSEARCH
- (NO GEODEL command since it uses sorted set and use ZREM to remove elements)

- valid longitudes -180 to 180
- valid latitudes -85.05112878 to 85.05112878 degrees

- OPTIONS
    - XX - only update exisiting keys no adding new elements
    - NX - dont update elements add new elements
    - CH - modify return value of the command
            - default is number of elements added
            - makes it  number of elements changed 
                - new elements added + existing elements modified

- Uses a Geohash to populate the sorted set
    - latitude and longitude interleaved to form a unique 52 bits
    - allows for bounding box and radius querying by checking 1 + 8 areas
    - areas checked by calculating range of box covered
        - remove enough bits from the less significant part of the ss score
        - copmute score range tot query sorted set for each area
- uses haversine distance
    - assumes sphereical earth
    - worst case error upto 0.5%
```
GEORADIUS key longitude latitude radius <M | KM | FT | MI>
    [WITHCOORD] [WITHDIST] [WITHHASH] [COUNT count [ANY]] [ASC | DESC]
    [STORE key | STOREDIST key]
```
- returns geospatial information added to sorted set using GEOADD
- O(N + logM) where N is the number of elements and M is number of items in the index
- retrieve items nead a specified point not further than a given amount of distance units
    - <M | KM | ML | FT>
    - M - meters
    - KM - kilometers
    - ML - miles
    - FT - feet
- additional return info flags
    - WITHDIST - distance of returned items from specified center
    - WITHCOORD - return longitude latitude coords of matching items
    - WITHHASH - return the sorted set score (52 bit int) // usually debugging purposes
- default return is unsorted items
    - can also invoke with sorting
    - ASC and DESC
- we can limit using [COUNT count <ANY>]
    - when any is specified command will return when it gets as many points
      as specified (may not be the closest)
    - with ANY queries on large areas may still be slow even with a small count

```
GEOSEARCH key <FROMMEMBER member | FROMLONLAT longitude latitude>
    <BYRADIUS radius <M | KM | FT | MI> | BYBOX width height <M | KM |
    FT | MI>> [ASC | DESC] [COUNT count [ANY]] [WITHCOORD] [WITHDIST]
    [WITHHASH]
```
- return members of the sorted set added to by GEOADD within borders of given shape
- extends GEORADIUS command to support rectangle areas

- center to be used is given by
    - FROMMEMBER - position of existing member
    - FROMLONLAT - use long lat position
- query shape
    - BYRADIUS similar to GEORADIUS search inside circular area
    - BYBOX search inside rectangle height and width
- command returns additional information
    - WITHDIST - also return distance of items from point
    - WITHCOORD - return longitude and latitude
    - WITHHASH - return raw geohash-encoded sorted set score
- matching items unsorted by default
    - ASC and DESC in terms of distance
- COUNT count [ANY]
    - limit number of returned items same as GEORADIUS count details

