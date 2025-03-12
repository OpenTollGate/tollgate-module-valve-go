package main

import (
	"flag"
	"context"
	"fmt"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"os/exec"
	"strconv"
	"time"
	"runtime/debug"
)

var (
	Version    string
	CommitHash string
	BuildTime  string
)

func getVersionInfo() string {
    if info, ok := debug.ReadBuildInfo(); ok {
        for _, setting := range info.Settings {
            switch setting.Key {
            case "vcs.revision":
                CommitHash = setting.Value[:7]
            case "vcs.time":
                BuildTime = setting.Value
            }
        }
    }
    return fmt.Sprintf("Version: %s\nCommit: %s\nBuild Time: %s", 
        Version, CommitHash, BuildTime)
}

func main() {
	fmt.Println("Starting Tollgate - Valve")

	listenForSessions()

	fmt.Println("Shutting down Tollgate - Valve")

	// Add a version flag
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		fmt.Println(getVersionInfo())
		return
	}
}

func listenForSessions() {

	var relayUrl = "ws://localhost:3334"
	var sessionEventKind = 66666
	//var merchantPubkey string = ""

	sk := nostr.GeneratePrivateKey()
	pk, _ := nostr.GetPublicKey(sk)
	nsec, _ := nip19.EncodePrivateKey(sk)
	npub, _ := nip19.EncodePublicKey(pk)

	fmt.Println("sk:", sk)
	fmt.Println("pk:", pk)
	fmt.Println(nsec)
	fmt.Println(npub)

	ctx := context.Background()
	relay, err := nostr.RelayConnect(ctx, relayUrl)

	if err != nil {
		panic(err)
	}

	var filters nostr.Filters
	if _, _, err := nip19.Decode(npub); err == nil {
		filters = []nostr.Filter{{
			Kinds: []int{sessionEventKind},
			//Authors: []string{merchantPubkey},
			//Limit: 10,
		}}
	} else {
		panic(err)
	}

	//ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	//defer cancel()

	sub, err := relay.Subscribe(ctx, filters)
	if err != nil {
		panic(err)
	}

	for ev := range sub.Events {
		// handle returned event.
		// channel will stay open until the ctx is cancelled (in this case, context timeout)
		handleSessionEvent(ev)
	}
}

func handleSessionEvent(event *nostr.Event) {
	fmt.Println(event)

	var macAddress = event.Tags.GetFirst([]string{"mac"}).Value()
	var sessionEndStr = event.Tags.GetFirst([]string{"session-end"}).Value()

	sessionEndUnix, err := strconv.ParseInt(sessionEndStr, 10, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println("mac: " + macAddress)
	fmt.Println("sessionEnd: " + string(sessionEndStr))

	var now = time.Now().Unix()

	var durationSeconds = sessionEndUnix - now

	if durationSeconds < 0 {
		fmt.Println(err, "Session for "+macAddress+" already ended, ignoring...")
		return
	}

	openGate(macAddress, durationSeconds)
}

func openGate(macAddress string, durationSeconds int64) {
	var durationMinutes int = int(durationSeconds / 60)

	// The minimum of this tollgate is 1 min, otherwise it would default to 24h
	if durationMinutes == 0 {
		durationMinutes = 1
	}

	fmt.Println("Opening gate for " + macAddress + " for the duration of " + strconv.Itoa(int(durationMinutes)) + " minute(s)")
	fmt.Println(durationMinutes)

	iwlistCmd := exec.Command("ndsctl", "auth", macAddress, strconv.FormatInt(int64(durationMinutes), 10))
	iwlistCmdOut, err := iwlistCmd.Output()
	if err != nil {
		fmt.Println(err, "Error when opening gate for macAddress "+macAddress)
	} else {
		fmt.Println(string(iwlistCmdOut))
	}
}
