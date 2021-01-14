// Modified from source: https://gist.github.com/ankanch/8c8ec5aaf374039504946e7e2b2cdf7f

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"log"
)

func main() {

	// Requests to a local server running on VirtualBox sent from the (same) local machine result in X-FORWARDED-FOR header set to 127.0.0.1
	// The code below gets the IP address of the local machine and puts it into the X-FORWARDED-FOR header
	const url = "https://api.ipify.org?format=text"	// we are using a pulib IP API, we're using ipify here, below are some others
                                              // https://www.ipify.org
                                              // http://myexternalip.com
                                              // http://api.ident.me
											  // http://whatismyipaddress.com/api
											  
	fmt.Printf("Getting IP address from  ipify ...\n")

	respIp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer respIp.Body.Close()

	if respIp.StatusCode != http.StatusOK {
		log.Printf("%s returned a status code %d", url, respIp.StatusCode)
		return
	}

	ip, err := ioutil.ReadAll(respIp.Body)

	if err != nil {
		panic(err)
	}

	fmt.Printf("My IP is:%s\n", ip)

	// Client code to send a request to the weather server
	serverUrl := "http://localhost:8080"
	
	req, err := http.NewRequest("GET", serverUrl, nil)
    req.Header.Set("X-FORWARDED-FOR", string(ip))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		weather, _ := ioutil.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(weather, &result)
		fmt.Println(result["message"].(string))
	}

	weather, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(weather))
}
