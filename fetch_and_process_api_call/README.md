# Fetch and process API call

This example shows how to use the Blocky Attestation Service (Blocky AS) to
fetch data from an API and process it. The example uses the Coingecko API to
fetch the current price of Bitcoin in USD on the Binance market.

## Setup

- Install the Blocky AS CLI by following the
[setup instructions](https://blocky-docs.redocly.app/attestation-service/setup)
in the Blocky AS documentation. 
- Make sure you also have
[Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
your system.
- [Get a key for the CoinGecko API](https://docs.coingecko.com/reference/setting-up-your-api-key) 
and set it in `fn-call.json` in the `api_key` field. 

## Run 

To run this example, all you need to do is call:

```bash
make
```

Your output should show the attested price of Bitcoin in USD on the Binance 
market:

```json
{
  "Success": true,
  "Value": {
    "market": "Binance",
    "coin_id": "BTC",
    "currency": "USD",
    "price": 97151,
    "timestamp": "2025-02-14T23:09:00Z"
  },
  "Error": ""
}
```
