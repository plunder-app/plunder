module github.com/plunder-app/plunder

go 1.12

require (
	github.com/AlecAivazis/survey/v2 v2.0.7 // indirect
	github.com/c4milo/gotoolkit v0.0.0-20190525173301-67483a18c17a // indirect
	github.com/coreos/go-etcd v2.0.0+incompatible // indirect
	github.com/cpuguy83/go-md2man v1.0.10 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/hooklift/assert v0.0.0-20170704181755-9d1defd6d214 // indirect
	github.com/hooklift/iso9660 v1.0.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/krolaw/dhcp4 v0.0.0-20190909130307-a50d88189771 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/sftp v1.11.0 // indirect
	github.com/plunder-app/BOOTy v0.0.0-20200513091117-6c4d474f5d95 // indirect
	github.com/plunder-app/plunder/pkg/apiserver v0.0.0-20200511161226-01273007741f
	github.com/plunder-app/plunder/pkg/certs v0.0.0-20200511161226-01273007741f
	github.com/plunder-app/plunder/pkg/parlay v0.0.0-20200511161226-01273007741f
	github.com/plunder-app/plunder/pkg/parlay/parlaytypes v0.0.0-20200511161226-01273007741f
	github.com/plunder-app/plunder/pkg/plunderlogging v0.0.0-20200511161226-01273007741f // indirect
	github.com/plunder-app/plunder/pkg/services v0.0.0-20200511161226-01273007741f
	github.com/plunder-app/plunder/pkg/ssh v0.0.0-20200511161226-01273007741f
	github.com/plunder-app/plunder/pkg/utils v0.0.0-20200511161226-01273007741f
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/thebsdbox/go-tftp v0.0.0-20190329154032-a7263f18c49c // indirect
	github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8 // indirect
	github.com/vishvananda/netlink v1.1.0 // indirect
	github.com/whyrusleeping/go-tftp v0.0.0-20180830013254-3695fa5761ee // indirect
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37 // indirect
	golang.org/x/net v0.0.0-20200506145744-7e3656a0809f // indirect
	golang.org/x/sys v0.0.0-20200513112337-417ce2331b5c // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace (
	github.com/plunder-app/plunder/pkg/apiserver => ./pkg/apiserver
	github.com/plunder-app/plunder/pkg/certs => ./pkg/certs
	github.com/plunder-app/plunder/pkg/services => ./pkg/services
	github.com/plunder-app/plunder/pkg/ssh => ./pkg/ssh
	github.com/plunder-app/plunder/pkg/utils => ./pkg/utils
)
