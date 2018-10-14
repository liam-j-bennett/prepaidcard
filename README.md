# prepaidcard

The datastore currently only uses postres.

Running ./run_service.sh should (famous last words) run a local instance of postgres and the app,
assuming you have Docker installed.

The app is exposed on port 8080, with endpoints on:

- /cards (POST) : Creates a new prepaid card and returns the object
- /cards/:cardId (GET) : Returns card object information about the card
- /cards/:cardId (POST) : Loads money onto the card, with JSON = {'amount': int64 in pence e.g. £100 == 10000}
- /cards/:cardId/spending : Returns a list of spending transactions on the card

- /transactions (POST) : Creates an auth transaction with JSON = {'merchantId': string (See main.go), 'card_id': string (card_number from card endpoints), 'amount': int64 auth amount}
- /transactions/:transactionId/capture : Captures funds already auth'ed, with JSON = {'amount': int64 MUST be less than auth_amount}
- /transactions/:transactionId/reverse : Reverses funds that have been auth'ed, with JSON with JSON = {'amount': int64 MUST be less than auth_amount}
- /transactions/:transactionId/refund : Refunds a capture, with JSON with JSON = {'amount': int64 MUST be less than captured amount}

The endpoints are located in server/handlers.go

Below is a snippet of python 3.6 using the requests library that: 

setup a card, add funds, make a transaction, capture some, reverse the rest and then refund the capture.

```
import requests

url = "http://localhost:8080"

card_id = requests.post(url + "/cards").json().get('card_number')

requests.post(url + f"/cards/{card_id}", json={"amount": 100})

transaction_id = requests.post(url + "/transactions", json={"amount":50, "card_number": card_id, "merchant_id": "amazon"}).json().get('id')

requests.patch(url + f"/transactions/{transaction_id}/capture", json={"amount": 25})

requests.patch(url + f"/transactions/{transaction_id}/reverse", json={"amount": 25})

requests.patch(url + f"/transactions/{transaction_id}/refund", json={"amount": 25})

list = requests.get(url + f"/cards/{card_id}/spending")
```
