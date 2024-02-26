package zeroaccount

// TODO: add setters for custom marshaller and unmarshaller

func SetAppSecret(secret string) {
	appSecret = secret
}

func SetEngine(newEngine Engine) {
	SetEngineSetterAndGetter(newEngine.Set, newEngine.Get)
}

func SetEngineSetterAndGetter(newSetter Setter, newGetter Getter) {
	setter = newSetter
	getter = newGetter
}

func SetErrorListener(listener ErrorListener) {
	errorListener = listener
}
