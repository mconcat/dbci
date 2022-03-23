# Database-Blockchain Interface

Database layer built top of Tendermint ABCI

## Design choice

The atomic unit of state is Account. Each account has possibily infinite size of state, addressed by singleton key-value pairs, mappings, queues. Accounts could be sharded over multiple instances of DBCI. 

DBCI bridges between Tendermint and upper application layer. Application layer, which is also another process running on local/remote machine, can send a grpc call to do state get, set, update, iteration calls. The role of grpc call is similar to typical DB query language, but is more specialized on key-value store usage, adding merkle proof related operations and RxGO API support.

Assuming concurrent usage of DBCI, by default, the DBCI supports mutex-like lock operation on each accounts. For example, an ABCI app that batch processes transactions can have the list of account where the each transaction will access beforehand it initiates the block execution, in order to attempt Conservative 2-Phase Locking(especially selected to prevent deadlock). 

As the atomic access/concurrency unit is an Account, if the application wants to make concurrent access on single account(say, a contract account for AMM pool), the account need to make a child account first. A child account of an account is an exclusive subset of the account's state. Locking on an account will be effective on all of its children, but not in the opposite way. Allocating/Disallocating child account is only temporary for clustering data in a single block and does not change the state itself.

## Implementation

### Account

- Account has its own SMT, with the merkle root stored under the main AppHash tree.

### Data Structure

There are three common data structure in the current Cosmos SDK usage.

- Singleton(Simple key-value pair)
- Ordered indexed mapping(Stored some primary key that is associated the data)
- Ordered queue(Stored by timestamp or integer amount that is not associated with the data)

Each of them will be implemented in 

### gRPC Query

Query language for DBCI

- Adapts Rx/LINQ style functional query 