


# An Overview of Various Caching Methods

### Read Strategies

**Write Around (Cache Aside, Lazy Loading)**
1. Try to fetch data from cache. If cache hit, return data to caller.
2. If miss, application will get data from the database.
3. Update the cache with the data from the database.
Note: The application calls the database directly.



**Pros**
Resilient to cache failures and falls back to DB
Cache and DB schemas can differ (?)
**Cons**
Cache can become out of sync with the database.
Application can serve stale data.

Reactive: only update the cache on reads. Writes don't do anything (except invaldiate the cache).

Cache gets filled only after a cache miss. 3 trips
**Best For**
Works best for read-heavy workloads
Cache will be only filled with frequently read data.

### Write Through Cache
Application writes to the cache and returns response to the application.
Data is persisted synchronously after cache is written.

The cache is written to on writes too.

Feature of the database

Well synced with database
Infrequently requested data is written to the cache. A lot of unnecessary writes
Instagram example: 
99% of profiles don't get queried often.

### Write-Back Cache
Reads and writes go through the cache
Asynchronously write to persistent database after writing to cache


# Cache eviction policies
### Least Recently Used (LRU)
### Least Frequently Used (LFU)
### Last In First Out (LIFO)
### Random Replacement (RR)
