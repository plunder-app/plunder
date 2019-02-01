package parlay

import (
	"testing"
)

const expectedHaProxycfg = `frontend k8s-api
    bind 192.168.0.1:6443
    bind 127.0.0.1:6443
    mode tcp
    option tcplog
    default_backend k8s-api

backend k8s-api
    mode tcp
    option tcplog
    option tcp-check
    balance roundrobin
    default-server inter 10s downinter 5s rise 2 fall 2 slowstart 60s maxconn 250 maxqueue 256 weight 100

    server etcd01 192.168.0.2:6443 check
    server etcd02 192.168.0.3:6443 check
    server etcd03 192.168.0.4:6443 check
`

func Test_loadBalancer_buildHAProxycfg(t *testing.T) {
	type fields struct {
		LBHostname        string
		LBPort            int
		EndPointHostnames []string
		EndPointAddresses []string
		EndPointPorts     []int
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Sample HA Proxy example",
			fields: fields{
				LBHostname:        "192.168.0.1",
				LBPort:            6443,
				EndPointAddresses: []string{"192.168.0.2", "192.168.0.3", "192.168.0.4"},
				EndPointHostnames: []string{"etcd01", "etcd02", "etcd03"},
				EndPointPorts:     []int{6443, 6443, 6443},
			},
			wantErr: false,
			want:    expectedHaProxycfg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &loadBalancer{
				LBHostname:        tt.fields.LBHostname,
				LBPort:            tt.fields.LBPort,
				EndPointHostnames: tt.fields.EndPointHostnames,
				EndPointAddresses: tt.fields.EndPointAddresses,
				EndPointPorts:     tt.fields.EndPointPorts,
			}
			got, err := l.buildHAProxycfg()
			if (err != nil) != tt.wantErr {
				t.Errorf("loadBalancer.buildHAProxycfg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("loadBalancer.buildHAProxycfg() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
