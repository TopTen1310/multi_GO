package tasks

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
	"fmt"
	"time"

	"github.com/TheRedSpy15/Multi-Go/utils"
	ct "github.com/daviddengcn/go-colortext"
)

// List will list all currently working tasks
func List() {
	optionColor := utils.RandomColor()

	ct.Foreground(ct.Red, true)
	fmt.Println("Available tasks:")
	time.Sleep(1 * time.Second)

	fmt.Println("\n-- Utility --")
	ct.Foreground(optionColor, false)
	fmt.Println("Scrape -r [URL]")
	fmt.Println("Email")
	fmt.Println("systemInfo")
	fmt.Println("Compress -r [file path]")
	fmt.Println("createUsb -r [usb location]")
	time.Sleep(1 * time.Second)

	ct.Foreground(ct.Red, true)
	fmt.Println("\n-- Security --")
	ct.Foreground(optionColor, false)
	fmt.Println("(sudo) Firewall -r [enable/disable/status]")
	fmt.Println("(sudo) Audit")
	fmt.Println("Hash -r [file path]")
	fmt.Println("encryptFile -r [file path]")
	fmt.Println("decryptFile -r [file path]")
	fmt.Println("pwnAccount -r [email]")
	fmt.Println("generatePassword -r [length]")
	fmt.Println("Bleach -r [file path]")
	time.Sleep(1 * time.Second)

	ct.Foreground(ct.Red, true)
	fmt.Println("\n-- Pentesting -- ")
	ct.Foreground(optionColor, false)
	fmt.Println("DOS -r [IP:PORT]")
	time.Sleep(1 * time.Second)

	ct.Foreground(ct.Red, true)
	fmt.Println("\n-- Other --")
	ct.Foreground(optionColor, false)
	fmt.Println("About")
	fmt.Println("cyberNews")
	fmt.Println(`Install -r [path or use "/bin/" for best result]`)
	fmt.Println("(sudo) Clean")
}
