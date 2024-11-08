


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

### Read Through
1. Try to fetch data from cache. If hit, return response.
2. If miss, the cache reads from the database and updates the cache entry.

Note: This is a feature part of the database, not a separate cache structure managed by the application.

**Pros**

Simplifies the application code

**Cons**
First request is always a cache miss.
Application can serve stale data


### Write Around Cache
\

Cache is part of the database (ex. Postgres).

1. Application writes to the database directly. The database populates the read cache on


Note: 

### Write Behind Cache
Application only interacts with the cache
The cache asynchronously persists the write to the database.

Cache consistency issues since the cache can drop out
Need cache replicas
Fastest approach, but can have inconsistency issues.

Eventual consistency vs immediate consistency
