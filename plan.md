## Implementation

We can follow three main approaches to support high availability and failover.

- Dumb Proxy Approach - This is to support an existing beanstalk server scale to load and failover easily. Beanstalkg will act as a proxy to
multiple beanstalkd servers allowing their state to be replicated. Master will be elected based on configuration. This approach requires 
the proxy to keep long running connections with all the beanstalks for each client that connects to it.
- Simpler/Less Performance/Expensive Approach - Save all job state to a backing store such as dynamodb/mongodb/mysql that support high availability. Then we can use
this store to coordinate the multi beanstalkg server setup.
- Complex/High Peformance/Cheaper Approach - Make all beanstalkg servers act as proxies to each other while using some algorithm such as Raft to coordinate leadership and failover. 
Once a leader has been elected all the operations will be verified with it before execution. This might require the priority queue to be implemented with a merge-able heap.

