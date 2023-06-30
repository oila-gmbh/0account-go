### 0account-go

Golang library for 0account.

Please follow the instructions or official documentations to integrate 0account with your service.

#### Example 1 (Custom Engine functions):

We will be using redis for these examples: https://github.com/go-redis/redis.
However, feel free to use any other database.

```go
package main

import (
    ...
    "github.com/oilastudio/0account-go"
)

func main() {
	var redisClient = redis.NewClient(&redis.Options{})
	var zero = zeroaccount.New(
		zeroaccount.SetEngineSetter(func(ctx context.Context, k string, v []byte) error {
			// for best results the timeout should match the timeout 
			// set in frontend (updateInterval option, default: 3 minutes)
			return redisClient.Set(ctx, k, v, 3*time.Minute).Err()
		}),
		zeroaccount.SetEngineGetter(func(ctx context.Context, k string) ([]byte, error) {
			v, err := redisClient.Get(ctx, k).Result()
			if err != nil {
				return nil, err
			}
			return []byte(v), nil
		}),
	)
	prepareHeaders := func(header http.Header) map[string]string {
		headers := make(map[string]string)
		for k, v := range header {
			headers[k] = v[0]
		}
		return headers
	}

	// The route URL is the callback URL you have set when you created 0account app. 
	http.Handle("/zeroauth", func(w http.ResponseWriter, r *http.Request) {
		data, err := zero.Auth(context.Background(), prepareHeaders(r.Header), c.Body())
		if err != nil {
			return errs.New(fiber.StatusUnauthorized, "not authorized", err)
		}
		if data == nil {
			return c.SendStatus(fiber.StatusOK)
		}
		ar := dto.AuthRequest{}
		if err = json.Unmarshal(data, &ar); err != nil {
			return errs.New(fiber.StatusUnauthorized, "wrong credentials", err)
		}
	})
}
```
Now our authentication is production ready!

---

#### Example 1 (In Memory Engine):
`0account-go` by default uses in memory cache engine if a custom engine is not supplied.

For brevity, we will leave out comments for the following examples, 
if something is unclear please read the comments on the first example 
or refer to the official documentation. If things are still unclear please create an issue. 

```go
package main

import (
    "encoding/json"
    "net/http"
    "github.com/oilastudio/oneaccount-go"
)

func main() {
    var oa = oneaccount.New()
    // The route URL is the callback URL you have set when you created One account app.
    http.Handle("/oneaccountauth", oa.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !oneaccount.IsAuthenticated(r) {
            return
        }
        // user authenticated and you can implement any logic your application 
        // needs. As an example you can extract data sent by the user 
        // after successful authentication
        data := make(map[string]interface{})
        if err := json.Unmarshal(oneaccount.Data(r), &data); err != nil {
            // handle the error
        }
        // since One account doesn't differentiate between sign up and sign in, 
        // you can use userId to check if the user signed up on your website or not
        // userID, _ := data["userId"]
        // the same way you can access any other data you requested from the user:
        firstName, _ := data["firstName"]
        // or create a struct to extract the data to
        // any data returned here would be sent to onAuth function on front-end e.g.:
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        if err := json.NewEncoder(w).Encode(map[string]interface{}{"firstName": firstName}); err != nil {
            // handle the error
        }
    })))
}
```

#### Example 3 (Custom Engine):
```go
type OneaccountRedisEngine struct {
    client *redis.Client
}

func (ore OneaccountRedisEngine) Set(ctx context.Context, k string, v []byte) error {
    // for best results the timeout should match the timeout 
    // set in frontend (updateInterval option, default: 3 minutes) 
    return ore.client.Set(ctx, k, v, 3 * time.Minute).Err()
}

func (ore OneaccountRedisEngine) Get(ctx context.Context, k string) ([]byte, error) {
    v, err := ore.client.Get(ctx, k).Result()
    if err != nil {
        return nil, err
    }
    return []byte(v), ore.client.Del(ctx, k).Err()
}

func main() {
    var redisClient = redis.NewClient(&redis.Options{})
    var oa = oneaccount.New(
        oneaccount.SetEngine(&OneaccountRedisEngine{client: redisClient}),
    )
    // The route URL is the callback URL you have set when you created One account app.
    http.Handle("/oneaccountauth", oa.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !oneaccount.IsAuthenticated(r) {
            return
        }
    })))
}
```

This example is a little longer, but it allows a greater control 
and is easier to separate the logic into a separate file.
