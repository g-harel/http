package httpc

// Get makes an http GET request to the provided url.
func Get(url string, headers *Headers, log func(string)) (string, error) {
	log("get " + url)
	return "", nil
}

// Post makes an http POST request to the provided url.
func Post(url string, headers *Headers, data string, log func(string)) (string, error) {
	log("post " + url)
	log("data " + data)
	return "", nil
}
