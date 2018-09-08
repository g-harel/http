package httpc

// Get makes an http GET request to the provided url.
func Get(url string, headers *Headers, log func(string)) (string, error) {
	log("get " + url)
	return "", nil
}
