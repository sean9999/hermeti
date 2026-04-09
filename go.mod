module github.com/sean9999/hermeti

go 1.26.0

require (
	github.com/sean9999/pear v0.0.5
	github.com/spf13/afero v1.15.0
)

require (
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	github.com/DataDog/gostackparse v0.7.0 // indirect
	github.com/sean9999/pkgalias v0.0.3 // indirect
	github.com/stretchr/testify v1.8.3 // indirect
	golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/telemetry v0.0.0-20260109210033-bd525da824e2 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
	golang.org/x/tools/go/expect v0.1.1-deprecated // indirect
	golang.org/x/tools/go/packages/packagestest v0.1.1-deprecated // indirect
	golang.org/x/vuln v1.1.4 // indirect
	honnef.co/go/tools v0.5.1 // indirect
)

tool (
	github.com/sean9999/pkgalias/cmd/pkgalias
	golang.org/x/vuln/cmd/govulncheck
	honnef.co/go/tools/cmd/staticcheck
)
