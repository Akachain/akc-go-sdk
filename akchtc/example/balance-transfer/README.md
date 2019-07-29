# Example Balance Transfer
## Problem?
  Robinson Credit Co. provides credit and financial services to large businesses. As such, their accounts are large, complex, and accessed by many people at once at any time of the day. They want to switch to blockchain, but are having trouble keeping up with the number of deposits and withdrawals happening at once on the same account. Additionally, they need to ensure users never withdraw more money than is available
  on an account, and transactions that do get rejected. The first problem is easy to solve, the second is more nuanced and requires a variety of strategies to accommodate high-throughput storage model design.

  To solve throughput, this new storage model is leveraged to allow every user performing transactions against the account to make that transaction in terms of a delta. For example, global e-commerce company America Inc. must be able to accept thousands of transactions an hour in order to keep up with their customer's demands. Rather than attempt to update a single row with the total amount of money in America Inc's account, Robinson Credit Co. accepts each transaction as an additive delta to America Inc's account. At the end of the day, America Inc's accounting department can quickly retrieve the total value in the account when the sums are aggregated.

  However, what happens when American Inc. now wants to pay its suppliers out of the same account, or a different account also on the blockchain?
  Robinson Credit Co. would like to be assured that America Inc.'s accounting department can't simply overdraw their account, which is difficult to do while at the same enabling transactions to happen quickly, as deltas are added to the ledger without any sort of bounds checking on the final aggregate value. There are a variety of solutions which can be used in combination to address this.

## Idea
  Demonstrates how to handle data in an application with a high transaction volume where the transactions all attempt to change the same key-value pair in the ledger. Such an application will have trouble as multiple transactions may read a value at a certain version, which will then be invalid when the first transaction updates the value to a new version, thus rejecting all other transactions until they're re-executed.
  Rather than relying on serialization of the transactions, which is slow, this application initializes a value and then accepts deltas of that value which are added as rows to the ledger. The actual value is then an aggregate of the initial value combined with all of the deltas. Additionally, a pruning function is provided which aggregates and deletes the deltas to update the initial value. This should be done during a maintenance window or when there is a lowered transaction volume, to avoid the proliferation of millions of rows of data. - `via Alexandre Pauwels from IBM.`

## How?
  This example is based on the problem and idea above, with the use of the Akachain SDK, idea isn't too difficult to implement.
  Do it.