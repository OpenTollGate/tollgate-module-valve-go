# Tollgate Module - Valve (go)

This Tollgate module is responsible for opening and closing internet access for a mac address specified in Tollgate's Nostr session events.

# Compile for ATH79 (GL-AR300 NOR)

```bash
cd ./src
env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o valve -trimpath -ldflags="-s -w"

# Hint: copy to connected router 
scp valve root@119.201.26.1:/tmp/valve
```

## License
This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
