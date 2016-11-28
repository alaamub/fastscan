/*
* Code written by Aladdin Mubaied
* 11/20/2016
 */
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// increase the worker as you wish
const workers = 10

/*
* This is the main function which fetches the resource from a specific url
 */
func requestResource(url string, ch chan int) {

	//defining a string
	dataParam := "hello"
	// Get the resource starting time.
	start := time.Now().UTC()
	// create the required parameters for building up the post
	message := strings.Join([]string{"data1=" + dataParam, ",", "data2=foobar", ",", "data3=" + dataParam}, "")
	// send the request and supply all the values in the POST and the headers
	client := http.Client{}
	// specify the query to be included in the POST request - here we use raw string so that we don't have to escape quotes, newlines, etc
	query := []byte(`id=` + dataParam + `&data4=` + message)
	// make a new POST request and include the query parameters
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		panic(err)
	}
	// send all the headers for the POST request
	req.Header.Add("custom-header", dataParam)
	// perform the POST request
	resp, _ := client.Do(req)
	// make sure to close the request
	defer resp.Body.Close()
	// read the response body
	body, _ := ioutil.ReadAll(resp.Body)
	// calculate end time and time spent on each run
	endTime := time.Now().UTC()
	elapsed := time.Since(start).Seconds()
	// return the data .. you can perform any checkings at this point
	fmt.Printf("Start Time: %s | End Time: %s | Time elapsed %.2fs | Content: %s | URL: %s\n", start, endTime, elapsed, string(body), url)
	// pass the number to the channel
	<-ch
}

func main() {
	// create a channel with number of workers
	ch := make(chan int, workers)
	// initiate the starting time
	startTime := time.Now().UTC()
	// read the list of the URLs from a file [ please note that you can also get the file name from stdin]
	file, err := os.Open("urls.txt")
	// make sure to close the file
	defer file.Close()
	if err != nil {
		panic(err)
	}
	// create a new scanner and read the file line by line (here i'm reading the list of urls from a file)
	line := bufio.NewScanner(file)
	// define i to zero so we can increase and control the number of concurrent requests
	i := 0
	for line.Scan() {
		// send each URL to our requestResource function
		go requestResource(line.Text(), ch)
		i++
		// pass the value of i to the channel indicating the finished work
		ch <- i
	}
	// drain remaining jobs in here
	for i := 0; i < workers; i++ {
		ch <- i
	}
	close(ch)
	// print out the total time for processing all the URLs
	elapsedTime := time.Since(startTime).Seconds()
	fmt.Printf("Script Total Execution Time %.2fs\n", elapsedTime)
}
