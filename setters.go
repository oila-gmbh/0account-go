package zeroaccount

import "fmt"

func SetAppSecret(secret string) {
	appSecret = secret
}

func SetEngine(newEngine Engine) {
	setter = newEngine.Set
	getter = newEngine.Get
}

func SetEngineSetterAndGetter(newSetter Setter, newGetter Getter) {
	fmt.Println("---------- SETTING SETTER AND GETTER")
	setter = newSetter
	getter = newGetter
}

func SetErrorListener(listener ErrorListener) {
	errorListener = listener
}
