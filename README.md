


# An Overview of Various Caching Methods

1. Cache-aside (Lazy Loading)
	•	Synonyms/Names: Lazy caching, cache-on-demand, pull-through cache, cache-aside pattern.
	•	Description: Loads data into the cache only when requested by the application. On a cache miss, data is fetched from the database and then added to the cache.
	•	Cache Write Behavior: Asynchronous. Data is only written to the cache if there’s a cache miss.
	•	On Write: The application directly updates the database. The cache is not updated immediately; instead, the entry may be invalidated, and the cache will fetch the updated value on the next read (lazy loading).
	•	Use Case: Best for read-heavy applications where cache consistency with the database can be slightly delayed.


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

2. Read-Through Cache

	•	Synonyms/Names: Read-through caching, read-through pattern, fetch-on-read cache.
	•	Description: The cache layer intercepts all read requests. If the data isn’t in the cache, the cache itself fetches it from the database and serves it to the application, automating cache population on a miss.
	•	Cache Write Behavior: Synchronous on cache miss. If data is not in the cache, the cache fetches it from the database.
	•	On Write: The application writes directly to the database, and the cache entry is either invalidated or updated. When the entry is next read, it will be reloaded from the database if it’s no longer in the cache.
	•	Use Case: Useful when the cache should automatically handle populating missing data, such as frequently accessed records in a read-heavy environment.


The cache is written to on writes too.

Feature of the database

Well synced with database
Infrequently requested data is written to the cache. A lot of unnecessary writes
Instagram example: 
99% of profiles don't get queried often.

3. Write-Through Cache

	•	Synonyms/Names: Write-through caching, synchronous cache write, write-through pattern.
	•	Description: All write operations go through the cache and are immediately propagated to the database. This keeps the cache and database synchronized at all times.
	•	Cache Write Behavior: Synchronous. Every write updates both the cache and database at the same time.
	•	On Write: Writes are made directly to the cache, and the cache then synchronously updates the database. This ensures both the cache and the database contain the latest data after every write.
	•	Use Case: Ideal for applications requiring high data consistency between the cache and database, especially in write-heavy scenarios.

4. Write-Back (Write-Behind) Cache
	•	Synonyms/Names: Write-behind caching, write-back caching, deferred write cache, asynchronous write cache.
	•	Description: Writes are performed on the cache first, with updates then propagated to the database in batches or based on a schedule. This reduces the number of writes to the database.
	•	Cache Write Behavior: Asynchronous. Writes to the database are delayed, occurring later based on a set policy (e.g., time interval or number of pending writes).
	•	On Write: The cache is updated directly on a write operation, but the database update is delayed. The cache writes back to the database asynchronously, usually in batches, making this approach more efficient for high-write scenarios. However, there’s a risk of stale data in the database during the delay.
	•	Use Case: Suited for applications with high write loads and where slight delays in data consistency are acceptable (e.g., analytics systems where up-to-the-second updates aren’t critical).

5. Database-Level Caching (Materialized Views, Query Caching)
	•	Synonyms/Names:
	•	Materialized Views: Precomputed views, persistent views, snapshot tables.
	•	Query Caching: Result caching, database query cache, statement caching.
	•	Description: Caching mechanisms are built into the database itself. Materialized views store precomputed query results, while query caching saves specific query results temporarily to avoid frequent re-evaluation.
	•	Cache Write Behavior:
	•	Materialized Views: Updated synchronously or asynchronously based on the configured refresh strategy.
	•	Query Caching: Typically synchronous; cached query results are returned until they expire or are invalidated.
	•	On Write:
	•	Materialized Views: Can be set to automatically refresh on database updates (synchronously) or at set intervals (asynchronously). With asynchronous refresh, the cached view may become temporarily stale.
	•	Query Caching: For cached queries, the results are automatically invalidated upon data change, ensuring consistent reads for users.
	•	Use Case: Useful for optimizing complex queries and aggregations, particularly on large datasets.

6. Hybrid Cache (Multi-Tiered)
	•	Synonyms/Names: Multi-tier cache, layered cache, hybrid caching, combined cache.
	•	Description: Combines multiple caching strategies, such as combining read-through for frequently accessed records and write-back for high-write scenarios, to meet diverse data performance and consistency needs.
	•	Cache Write Behavior: Varies by cache tier.
	•	Read-Through Cache: Asynchronous cache population.
	•	Write-Through Cache: Synchronous updates.
	•	Materialized Views: Configurable based on refresh strategy (synchronous or asynchronous).
	•	On Write: Each cache tier handles writes according to its specific strategy:
	•	Read-Through Cache: Writes directly to the database, then invalidates the cache entry, which will be repopulated on the next read.
	•	Write-Through Cache: Writes to both the cache and the database simultaneously.
	•	Materialized Views: Writes may trigger an immediate or scheduled refresh, depending on configuration.
	•	Use Case: Ideal for complex applications where different data types and operations require varied caching behaviors, providing flexibility and high performance.

# Cache eviction policies
### Least Recently Used (LRU)
### Least Frequently Used (LFU)
### Last In First Out (LIFO)
### Random Replacement (RR)

What if the cache goes down/has a network failure and contains stale data?
Redis SCAN
