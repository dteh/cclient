module github.com/dteh/cclient

go 1.22.0

toolchain go1.22.2

require (
	github.com/dteh/fhttp v0.0.0
	github.com/refraction-networking/utls v1.6.7
	golang.org/x/net v0.30.0
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/cloudflare/circl v1.5.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)

replace github.com/dteh/fhttp => ../fhttp
