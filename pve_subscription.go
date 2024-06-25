package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const sharedKeyData = "kjfdlskfhiuewhfk947368"
const serverKeyFile = "/etc/ssh/ssh_host_rsa_key.pub"

func getTimestamp() int64 {
	return time.Now().Unix()
}

func md5Base64Perl(x string) string {
	hash := md5.Sum([]byte(x))
	return strings.TrimRight(base64.StdEncoding.EncodeToString(hash[:]), "=")
}

func md5Base64Rs(x string) string {
	hash := md5.Sum([]byte(x))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func generateServerID(key string) string {
	hash := md5.Sum([]byte(key))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

func dtString(format string, offsetSecs int) string {
	return time.Now().Add(time.Duration(offsetSecs) * time.Second).Format(format)
}

func generateSubscriptionPvePmg(key string, serverIDs []string) string {
	localInfo := map[string]interface{}{
		"checktime":      getTimestamp(),
		"status":         "Active",
		"key":            key,
		"validdirectory": strings.Join(serverIDs, ","),
		"productname":    "YajuuSenpai",
		"regdate":        dtString("2006-01-02 15:04:05", 0),
		"nextduedate":    dtString("2006-01-02", 1296000),
	}

	data, _ := json.Marshal(localInfo)
	encodedData := base64.StdEncoding.EncodeToString(data)
	cat := fmt.Sprintf("%d%s\n%s", localInfo["checktime"], encodedData, sharedKeyData)
	csum := md5Base64Perl(cat)

	return fmt.Sprintf("%s\n%s\n%s\n", key, csum, encodedData)
}

func activatePvePmg(key string, subscriptionFile string) {
	serverKey, err := os.ReadFile(serverKeyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading server key file: %v\n", err)
		os.Exit(1)
	}
	serverID := generateServerID(string(serverKey))
	subscription := generateSubscriptionPvePmg(key, []string{serverID})

	err = os.WriteFile(subscriptionFile, []byte(subscription), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing subscription file: %v\n", err)
		os.Exit(1)
	}
}

func generateSubscriptionPbs(key string, serverIDs []string) string {
	localInfo := map[string]interface{}{
		"status":      "active",
		"serverid":    strings.Join(serverIDs, ","),
		"checktime":   getTimestamp(),
		"key":         key,
		"message":     "Yajuu Senpai has got your back",
		"productname": "YajuuSenpai",
		"regdate":     dtString("2006-01-02 15:04:05", 0),
		"nextduedate": dtString("2006-01-02", 1296000),
		"url":         "https://github.com/Jamesits/pve-fake-subscription",
	}

	data, _ := json.Marshal(localInfo)
	encodedData := base64.StdEncoding.EncodeToString(data)
	cat := fmt.Sprintf("%d%s%s", localInfo["checktime"], encodedData, sharedKeyData)
	csum := md5Base64Rs(cat)

	return fmt.Sprintf("%s\n%s\n%s\n", key, csum, encodedData)
}

func activatePbs(key string, subscriptionFile string) {
	serverKey, err := ioutil.ReadFile(serverKeyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading server key file: %v\n", err)
		os.Exit(1)
	}
	serverID := generateServerID(string(serverKey))
	subscription := generateSubscriptionPbs(key, []string{serverID})

	err = os.WriteFile(subscriptionFile, []byte(subscription), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing subscription file: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if _, err := os.Stat("/etc/pve"); err == nil {
		fmt.Println("Activating Proxmox VE...")
		activatePvePmg("pve8p-1145141919", "/etc/subscription")
	}

	if _, err := os.Stat("/etc/pmg"); err == nil {
		fmt.Println("Activating Proxmox Mail Gateway...")
		activatePvePmg("pmgp-1145141919", "/etc/pmg/subscription")
	}

	if _, err := os.Stat("/etc/proxmox-backup"); err == nil {
		fmt.Println("Activating Proxmox Backup Server...")
		activatePbs("pbst-1145141919", "/etc/proxmox-backup/subscription")
	}
}
