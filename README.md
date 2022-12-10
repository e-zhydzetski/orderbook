# strips-tt

# Task 1

Check OS threads count in case of multiple long-lived outgoing network connections without explicit epoll.  

* [Measurements](https://github.com/e-zhydzetski/strips-tt/tree/master/threads)
* Useful links:
  * [Post about Golang and system threads](https://www.sobyte.net/post/2021-06/golang-number-of-threads-in-the-running-program/)

# Task 2

1. Implement simple "matching engine" that can execute limit and market orders.
2. The engine should have optimial store for orderbook (the structure of how you will store the orderbook, and which "indexes" you will implement on top of it the most interesting part)
3. The engine should have "rollback" feature, where you can start matching some order finish it, and then rollback book to the previous state. Rollback algo should be something more effective than just duplicate orderbook in memory and switch between them.