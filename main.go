package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

import bolt "go.etcd.io/bbolt"

type License struct {
	Key         string `json:"key"`
	Type        uint64 `json:"type"`
	UserLimit   uint64 `json:"user_limit"`
	DeviceLimit uint64 `json:"device_limit"`
	SiteLimit   uint64 `json:"site_limit"`
}

type Status struct {
	Hostname        string `json:"hostname"`
	Uptime          int64  `json:"uptime"`
	Version         string `json:"version"`
	UpdateAvailable bool   `json:"update_available"`
	UpdateChangelog string `json:"update_changelog"`
	UpdateVersion   string `json:"update_version"`
	UpdateDownload  string `json:"update_link"`
	Updating        bool   `json:"updating"`
	HostOS          string `json:"host_os"`
	Edition         string `json:"edition"`
	UserLimit       uint64 `json:"user_limit"`
	DeviceLimit     uint64 `json:"device_limit"`
	SiteLimit       uint64 `json:"site_limit"`
}

func main() {
	fmt.Println("[*] Hak5 C2 Licensing Toolkit by Pwn3rzs / CyberArsenal!")
	scanner := bufio.NewScanner(os.Stdin)
	db, err := bolt.Open("c2.db", 0600, nil)
	if err != nil {
		return
	}
	fmt.Println("[+] DB Opened!")
	defer db.Close()
	for {
		fmt.Println("[*] Enter a command (generate/decode/read/crack/exit):")
		scanner.Scan()
		command := scanner.Text()
		switch command {
		case "generate":
			generateHexCode()
		case "read":
			err = db.View(func(tx *bolt.Tx) error {
				fmt.Println("[*] Enter the bucket (setup/status):")
				scanner.Scan()
				var bucket string
				var key string
				var license License
				var status Status

				input := scanner.Text()
				switch input {
				case "setup":
					bucket = input
					key = "license"
				case "status":
					bucket = input
					key = "status"
				default:
					return fmt.Errorf("[!] Invalid bucket. Please enter 'setup', 'status'")
				}

				fmt.Printf("[+] Opening key '%s' in bucket '%s'\n", key, bucket)

				bucketz := tx.Bucket([]byte(bucket))
				if bucketz == nil {
					return fmt.Errorf("[!] Bucket '" + bucket + "' not found")
				}

				data := bucketz.Get([]byte(key))
				if data == nil {
					return fmt.Errorf("[!] Data for '" + key + "' not found")
				}

				decoder := gob.NewDecoder(bytes.NewReader(data))
				if key == "license" {
					if err := decoder.Decode(&license); err != nil {
						return err
					}
					prettyLicense, _ := json.MarshalIndent(license, "", "    ")

					fmt.Printf("[+] Decoded License struct: %+v\n", string(prettyLicense))
				} else {
					if err := decoder.Decode(&status); err != nil {
						return err
					}
					prettyStatus, _ := json.MarshalIndent(status, "", "    ")
					fmt.Printf("[+] Decoded Status struct: %+v\n", string(prettyStatus))
				}

				return nil
			})

			if err != nil {
				fmt.Println("[!] Error decoding data:", err)
			}
		case "decode":
			fmt.Println("[*] Enter the struct (license/status):")
			scanner.Scan()
			mode := scanner.Text()
			fmt.Println("[*] Enter the hex string:")
			scanner.Scan()
			hexData := scanner.Text()
			decodeHex(hexData, mode)
		case "crack":
			_ = db.Update(func(tx *bolt.Tx) error {
				license := License{
					Key:         "Pwn3rzs",
					Type:        2,
					UserLimit:   uint64(10000),
					DeviceLimit: uint64(10000),
					SiteLimit:   uint64(10000),
				}

				var buf bytes.Buffer
				encoder := gob.NewEncoder(&buf)
				err := encoder.Encode(license)
				if err != nil {
					fmt.Println("[!] Encoding error:", err)
					return err
				}
				//var status Status

				bucketz := tx.Bucket([]byte("setup"))
				if bucketz == nil {
					return fmt.Errorf("[!] Bucket 'setup' not found")
				}

				data := bucketz.Put([]byte("license"), buf.Bytes())
				if data != nil {
					return fmt.Errorf("[!] Data for 'license' not found")
				}
				var statusView Status
				bucketz = tx.Bucket([]byte("status"))
				v := bucketz.Get([]byte("status"))
				decoder := gob.NewDecoder(bytes.NewReader(v))
				if err := decoder.Decode(&statusView); err != nil {
					return err
				}
				prettyLicense, _ := json.MarshalIndent(statusView, "", "    ")
				fmt.Printf("[+] Decoded Status struct: %+v\n", string(prettyLicense))
				statusView.DeviceLimit = uint64(10000)
				statusView.UserLimit = uint64(10000)
				statusView.Edition = "teams"
				statusView.SiteLimit = uint64(10000)
				encoder = gob.NewEncoder(&buf)
				_ = encoder.Encode(statusView)
				_ = bucketz.Put([]byte("status"), buf.Bytes())

				return nil
			})

			_ = db.View(func(tx *bolt.Tx) error {
				var licenseView License
				b := tx.Bucket([]byte("setup"))
				v := b.Get([]byte("license"))
				decoder := gob.NewDecoder(bytes.NewReader(v))
				if err := decoder.Decode(&licenseView); err != nil {
					return err
				}
				prettyLicense, _ := json.MarshalIndent(licenseView, "", "    ")
				fmt.Printf("[+] Decoded License struct: %+v\n", string(prettyLicense))
				return nil
			})

			_ = db.View(func(tx *bolt.Tx) error {
				var statusView Status
				b := tx.Bucket([]byte("status"))
				v := b.Get([]byte("status"))
				decoder := gob.NewDecoder(bytes.NewReader(v))
				if err := decoder.Decode(&statusView); err != nil {
					return err
				}
				prettyStatus, _ := json.MarshalIndent(statusView, "", "    ")
				fmt.Printf("[+] Decoded Status struct: %+v\n", string(prettyStatus))
				return nil
			})

			fmt.Println("[+] DB Values edited")
			fmt.Println("[*] Patching application")

		case "exit":
			fmt.Println("[!] Exiting...")
			return
		default:
			fmt.Println("[!] Invalid command. Please enter 'generate', 'decode', 'read', or 'exit'.")
		}
	}
}

func generateHexCode() {
	license := License{
		Key:         "Pwn3rzs",
		Type:        1,
		UserLimit:   500,
		DeviceLimit: 500,
		SiteLimit:   500,
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(license)
	if err != nil {
		fmt.Println("[!] Encoding error:", err)
		return
	}
	fmt.Println("[+] Generated Hex Code:", hex.EncodeToString(buf.Bytes()))
}

func decodeHex(hexData string, mode string) {
	var license License
	var status Status

	switch mode {
	case "license":
		hexDecoded, err := hex.DecodeString(hexData)
		if err != nil {
			fmt.Println("Error decoding hex string:", err)
			return
		}
		dec := gob.NewDecoder(bytes.NewReader(hexDecoded))

		if err := dec.Decode(&license); err != nil {
			fmt.Println("Error decoding gob data:", err)
			return
		}
		fmt.Printf("Decoded License struct: %+v\n", license)
	case "status":
		hexDecoded, err := hex.DecodeString(hexData)
		if err != nil {
			fmt.Println("Error decoding hex string:", err)
			return
		}
		dec := gob.NewDecoder(bytes.NewReader(hexDecoded))
		if err := dec.Decode(&status); err != nil {
			fmt.Println("Error decoding gob data:", err)
			return
		}
		fmt.Printf("Decoded Status struct: %+v\n", status)
	default:
		fmt.Println("Invalid command. Please enter 'create', 'decode', or 'exit'.")
	}
}
