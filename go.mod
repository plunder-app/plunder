module github.com/plunder-app/plunder

go 1.12

require (
	github.com/AlecAivazis/survey/v2 v2.0.2 // indirect
	github.com/c4milo/gotoolkit v0.0.0-20190525173301-67483a18c17a // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/hooklift/assert v0.0.0-20170704181755-9d1defd6d214 // indirect
	github.com/hooklift/iso9660 v1.0.0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/krolaw/dhcp4 v0.0.0-20190531080455-7b64900047ae // indirect
	github.com/pkg/sftp v1.10.1 // indirect
	github.com/plunder-app/plunder/pkg/apiserver v0.0.0-20191105152536-b5c505aaf830
	github.com/plunder-app/plunder/pkg/certs v0.0.0-00010101000000-000000000000
	github.com/plunder-app/plunder/pkg/parlay v0.0.0-20191105152536-b5c505aaf830
	github.com/plunder-app/plunder/pkg/parlay/parlaytypes v0.0.0-20191105152536-b5c505aaf830
	github.com/plunder-app/plunder/pkg/plunderlogging v0.0.0-20191105152536-b5c505aaf830 // indirect
	github.com/plunder-app/plunder/pkg/services v0.0.0-20191105152536-b5c505aaf830
	github.com/plunder-app/plunder/pkg/ssh v0.0.0-00010101000000-000000000000
	github.com/plunder-app/plunder/pkg/utils v0.0.0-20191105152536-b5c505aaf830
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/thebsdbox/go-tftp v0.0.0-20190329154032-a7263f18c49c // indirect
	github.com/whyrusleeping/go-tftp v0.0.0-20180830013254-3695fa5761ee // indirect
	golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7 // indirect
	golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace (
	github.com/plunder-app/plunder/pkg/apiserver => ./pkg/apiserver
	github.com/plunder-app/plunder/pkg/certs => ./pkg/certs
	github.com/plunder-app/plunder/pkg/services => ./pkg/services
	github.com/plunder-app/plunder/pkg/ssh => ./pkg/ssh
)
