package utils

/*
   Copyright 2020 TheRedSpy15

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"
	"time"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/gocolly/colly"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/crypto/ssh/terminal"
)

// BytesToGigabytes - an internal utility function to convert bytes to gigabytes; used to clarify the output of PrintMemory()
func BytesToGigabytes(bytes uint64) float64 {
	const conversionFactor = 1000000000
	result := float64(bytes) / conversionFactor
	result = float64(math.Round(result*100) / 100)
	return result
}

// CheckTarget throws an error if the target is empty
func CheckTarget(target string) {
	if target == "" { // check if target is blank
		CheckErr(fmt.Errorf("target cannot be empty when performing this task"))
	}
}

// CheckErr prints any non-nil errors followed by exiting the application
func CheckErr(err error) {
	if err != nil { // check if there actually is any error
		ct.Foreground(ct.Red, true) // set texts color to bright red
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// CheckSudo throws an error if the current user is not sudo/root
func CheckSudo() {
	user, _ := user.Current()
	if !strings.Contains(user.Username, "root") {
		CheckErr(fmt.Errorf("cannot run this task without root/sudo"))
	}
}

// RunCmd - runs a command on the system and prints the result
// TODO: document
func RunCmd(command string, arg ...string) string {
	cmd := exec.Command(command)
	for _, arg := range arg {
		cmd.Args = append(cmd.Args, arg)
	}

	var o bytes.Buffer
	cmd.Stdout = &o // asign o to cmd's Stdout

	err := cmd.Run()
	CheckErr(err)

	return o.String()
}

// ReadFileIntoByte - is used for getting []byte of file
func ReadFileIntoByte(filename string) []byte {
	var data []byte                // specify type
	file, err := os.Open(filename) // make file object
	defer file.Close()             // close file on function end

	CheckErr(err)

	data, err = ioutil.ReadAll(file) // read all
	if err != nil {
		ct.Foreground(ct.Red, true) // set text color to bright red
		panic(err.Error())
	}

	return data // return file bytes
}

// GetPassword - securely gets password from a user
func GetPassword() string {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin)) // run password command, make var with result
	CheckErr(err)

	password := string(bytePassword) // cast to string var

	return password
}

func RandomColor() ct.Color {
	rand.Seed(time.Now().UTC().UnixNano())
	switch randInt(0, 5) {
	case 0:
		ct.Foreground(ct.Blue, true)
		return ct.Blue
	case 1:
		ct.Foreground(ct.Red, true)
		return ct.Red
	case 2:
		ct.Foreground(ct.Yellow, true)
		return ct.Yellow
	case 3:
		ct.Foreground(ct.Magenta, true)
		return ct.Magenta
	case 4:
		ct.Foreground(ct.Green, true)
		return ct.Green
	}
	return ct.White
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// DownloadFile downloads a file from a url
func DownloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// PrintBanner - displays the banner text
func PrintBanner() {
	RandomColor()

	fmt.Println(`
 __  __       _ _   _    ____
|  \/  |_   _| | |_(_)  / ___| ___
| |\/| | | | | | __| | | |  _ / _ \
| |  | | |_| | | |_| | | |_| | (_) |
|_|  |_|\__,_|_|\__|_|  \____|\___/`)
}

// CollyAddress - scrapes a website link
func CollyAddress(target string, savePage bool, ip bool) {
	if ip { // check if target is an IP address not URL
		target = "http://" + target + "/" // modify target to be valid address
	}

	c := colly.NewCollector() // make colly object
	c.IgnoreRobotsTxt = true  // ignore RobotsText

	// configuring colly/collector object
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) { // print error message on error
		ct.Foreground(ct.Red, true) // set text color to bright red
		log.Println("Something went wrong:", err)
		ct.ResetColor() // reset text color to default
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
		fmt.Println("Response:", r.StatusCode)
	})

	c.OnScraped(func(r *colly.Response) { // finished with site
		fmt.Println("Finished", r.Request.URL)

		if savePage { // check if save is enabled
			err := r.Save(r.FileName()) // saving data
			CheckErr(err)

			ct.Foreground(ct.Green, true) // set text color to bright red
			fmt.Println("Saved - ", r.FileName())
			ct.ResetColor() // reset text color to default color
		}
	})

	c.Visit(target) // actually using colly/collector object, and visiting target
}

// RandomString - returns a random string
// TODO: rewrite in my own code
// TODO: add more comments
// Util function - returns a random string
/* Original: https://stackoverflow.com/questions/22892120
/how-to-generate-a-random-string-of-a-fixed-length-in-go#31832326 */
func RandomString(length int) string {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // letters to use
		letterIdxBits = 6                                                      // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1                                   // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits                                     // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano()) // create random source

	b := make([]byte, length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// PrintCPU - prints CPU info
// TODO: add more info
func PrintCPU() {
	cpuCount, err1 := cpu.Counts(false)       // get cpu count total
	cpuCountLogical, err2 := cpu.Counts(true) // get cpu logical count
	cpuLoad, err3 := load.Avg()               // get current cpu load

	CheckErr(err1)
	CheckErr(err2)
	CheckErr(err3)

	ct.Foreground(ct.Red, true) // change text color to bright red
	fmt.Println("\n-- CPU --")
	ct.Foreground(ct.Yellow, false)                      // change text color to dark yellow
	fmt.Println("CPU Count: (logical)", cpuCountLogical) // cpu count logical
	fmt.Println("CPU Count:", cpuCount)                  // cpu count total
	fmt.Printf("CPU Usage:\n\tLast Minute: %d%%\n\tLast 5 Minutes: %d%%\n\tLast 15 Minutes: %d%%\n", int(cpuLoad.Load1*10), int(cpuLoad.Load5*10), int(cpuLoad.Load15*10))
}

// PrintMemory - prints info about system memory
// TODO: get physical memory instead of swap
// TODO: convert values to gigabytes
func PrintMemory() {
	memVirt, err1 := mem.VirtualMemory()
	CheckErr(err1)

	memSwap, err2 := mem.SwapMemory()
	CheckErr(err2)

	ct.Foreground(ct.Red, true) // change text color to bright red
	fmt.Println("\n-- Memory --")
	ct.Foreground(ct.Yellow, false)                                    // change text color to dark yellow
	fmt.Println("Memory Used (Gb):", BytesToGigabytes(memVirt.Used))   // used
	fmt.Println("Memory Free (Gb):", BytesToGigabytes(memVirt.Free))   // free
	fmt.Println("Memory Total (Gb):", BytesToGigabytes(memVirt.Total)) // total

	ct.Foreground(ct.Red, true) // change text color to bright red
	fmt.Println("\n-- Swap --")
	ct.Foreground(ct.Yellow, false)                                  // change text color to dark yellow
	fmt.Println("Swap Used (Gb):", BytesToGigabytes(memSwap.Used))   // used
	fmt.Println("Swap Free (Gb):", BytesToGigabytes(memSwap.Free))   // free
	fmt.Println("Swap Total (Gb):", BytesToGigabytes(memSwap.Total)) // total
}

// PrintHost - prints info about system host
func PrintHost() {
	hostInfo, err := host.Info() // get host info object
	CheckErr(err)

	ct.Foreground(ct.Red, true) // change text color to bright red
	fmt.Println("\n-- Host --")
	ct.Foreground(ct.Yellow, false)                            // change text color to dark yellow
	fmt.Println("Kernal Version:", hostInfo.KernelVersion)     // kernal version
	fmt.Println("Platform:", hostInfo.Platform)                // platform
	fmt.Println("Platform Family:", hostInfo.PlatformFamily)   // platform family
	fmt.Println("Platform Version:", hostInfo.PlatformVersion) // platform version
	fmt.Println("Uptime (seconds):", hostInfo.Uptime)          // uptime
	fmt.Println("Host Name:", hostInfo.Hostname)               // hostname
	fmt.Println("Host ID:", hostInfo.HostID)                   // host id
	fmt.Println("OS:", hostInfo.OS)                            // os
}
