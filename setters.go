package zeroaccount

import "fmt"

func SetAppSecret(secret string) {
	appSecret = secret
	fmt.Printf("APP SECRET SET, appSecret: %s, secret %s\n", appSecret, secret)
}

func SetEngine(newEngine Engine) {
	setter = newEngine.Set
	getter = newEngine.Get
}

func SetEngineSetterAndGetter(newSetter Setter, newGetter Getter) {
	setter = newSetter
	getter = newGetter
}

func SetErrorListener(listener ErrorListener) {
	errorListener = listener
}
