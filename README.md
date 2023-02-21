# Papapi

A go client to communicate with the low-level Post Affiliate Pro API. There's nothing fancy about this.

_Heavily inspired by https://github.com/JSBizon/papapi_

Read more at [Post Affiliate Pro Low level API Documentation](https://support.qualityunit.com/919534-Overview)

## ⚠️ Limitation
This API is limited to only work with affiliate accounts and grid requests. It's not intended to be a fully fledged wrapper. 

## Example usage

```golang
import "github.com/voke/papapi"

sess := papapi.NewSession("https://login.network.com/scripts/server.php", papapi.Affiliate)
err := sess.Login("john.doe@example.com", "secret")

if err != nil {
    panic(err)
}

req := papapi.NewGridRequest("Pap_Affiliates_Reports_TransactionsGrid", "getRows", sess)

req.AddColumn("commission")
req.AddColumn("orderid")
req.AddColumn("campaignname")

req.AddFilter("dateinserted", papapi.DateGreater, "2022-02-12")
req.AddFilter("dateinserted", papapi.DateLower, "2022-02-20")

req.SetLimit(30)
req.SetOffset(0)

res, err := req.Do()

if err != nil {
    panic(err)
}

for _, rec := range res.Records() {
    // Do something with rec
}

```