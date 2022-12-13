# Order book

## Requirements

* Functional:
    * Execute [Market order](https://en.wikipedia.org/wiki/Order_(exchange)#Market_order)
    * Execute [Limit order](https://en.wikipedia.org/wiki/Order_(exchange)#Limit_order)
    * Limit order may be partially executed, no additional constraints
    * Rollback last executed order? TODO clarify use-case
* Non-functional:
    * Low-latency
    * Rollback without snapshotting
* Out of scope
    * Get [Market depth](https://en.wikipedia.org/wiki/Market_depth)
    * Cancel order