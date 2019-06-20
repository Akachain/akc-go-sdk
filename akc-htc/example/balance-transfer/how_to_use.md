Balance transfer

This is the example chaincode for token transfer in the blockchain network.
This chaincode is written with golang.

In this chaincode, we will use the AKC SDK to solve the problem of multiple transfers at the same time.

I. Init
  1. Params: none;

II. Invoke
  1. Create new user
    - This is an extended function used to create new users, because there should be at least 2 users in database to able for transfer the token.
    - Invoke chaincode with:
      Function Name: createUser
      Arguments:
        [
          "<wallet_address>", 
          "<token_amount>"
        ]

  2. Token transfer
    - This function is used to perform the token transfer between user A and user B.
    - Invoke chaincode with:
      Function Name: transferToken
      Arguments:
      [
        "<wallet_address_send>",
        "<wallet_address_received>",
        "<token_amount_send>"
      ]

  3. Update user balance
    - This function used to update user token data, because when executing the `transferToken` function, token is not transfer immediately, it's pushed to high throughput. You need to run this fuction to update the token data into the database.
    - Invoke chaincode with:
      Function Name: updateUserBalance
      Arguments:
      [
        "<prune_type>" // 2 option available currently is PRUNE_FAST or PRUNE_SAFE.
      ]

  4. Get user balance
    - This function used to get user information via user wallet address.
    - Invoke chaincode with:
      Function Name: getUserToken
      Arguments:
      [
        "<wallet_address>"
      ]

  More specifically, read the source code and learn more.
