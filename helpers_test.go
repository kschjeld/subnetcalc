package subnetcalc

import (
	"reflect"
	"testing"
)

func Test_inetBToN(t *testing.T) {
	type args struct {
		ip4 []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty",
			args: args{
				ip4: []byte{0, 0, 0, 0},
			},
			want: 0,
		},
		{
			name: "max",
			args: args{
				ip4: []byte{255, 255, 255, 255},
			},
			want: 255<<24 + 255<<16 + 255<<8 + 255,
		},
		{
			name: "sample",
			args: args{
				ip4: []byte{10, 16, 24, 0},
			},
			want: 10<<24 + 16<<16 + 24<<8 + 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inetBToN(tt.args.ip4); got != tt.want {
				t.Errorf("inetBToN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inetNToA(t *testing.T) {
	type args struct {
		addrint int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				addrint: 0,
			},
			want: "0.0.0.0",
		},
		{
			name: "max",
			args: args{
				addrint: 255<<24 + 255<<16 + 255<<8 + 255,
			},
			want: "255.255.255.255",
		},
		{
			name: "sample",
			args: args{
				addrint: 10<<24 + 16<<16 + 24<<8 + 0,
			},
			want: "10.16.24.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inetNToA(tt.args.addrint); got != tt.want {
				t.Errorf("inetNToA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inetNToB(t *testing.T) {
	type args struct {
		addrint int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "empty",
			args: args{
				addrint: 0,
			},
			want: []byte{0, 0, 0, 0},
		},
		{
			name: "max",
			args: args{
				addrint: 255<<24 + 255<<16 + 255<<8 + 255,
			},
			want: []byte{255, 255, 255, 255},
		},
		{
			name: "sample",
			args: args{
				addrint: 10<<24 + 16<<16 + 24<<8 + 0,
			},
			want: []byte{10, 16, 24, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inetNToB(tt.args.addrint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inetNToB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inetNetworkAddress(t *testing.T) {
	type args struct {
		ip   int
		mask int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "equal",
			args: args{
				ip:   10<<24 + 0<<16 + 0<<8 + 0,
				mask: 16,
			},
			want: 10<<24 + 0<<16 + 255<<8 + 255,
		},
		{
			name: "masked",
			args: args{
				ip:   10<<24 + 16<<16 + 16<<8 + 255,
				mask: 24,
			},
			want: 10<<24 + 16<<16 + 16<<8 + 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inetNetworkAddress(tt.args.ip, tt.args.mask); got != tt.want {
				t.Errorf("inetNetworkAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inetSubnetAddresses(t *testing.T) {
	type args struct {
		mask int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty",
			args: args{
				mask: 0,
			},
			want: 255<<24 + 255<<16 + 255<<8 + 256,
		},
		{
			name: "full",
			args: args{
				mask: 31,
			},
			want: 0<<24 + 0<<16 + 0<<8 + 2,
		},
		{
			name: "half",
			args: args{
				mask: 16,
			},
			want: 0<<24 + 0<<16 + 256<<8 + 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inetSubnetAddresses(tt.args.mask); got != tt.want {
				t.Errorf("inetSubnetAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inetSubnetLastAddress(t *testing.T) {
	type args struct {
		subnet int
		mask   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "10.0.0.0/16",
			args: args{
				subnet: 10<<24 + 0<<16 + 0<<8 + 0,
				mask:   16,
			},
			want: 10<<24 + 0<<16 + 255<<8 + 255,
		},
		{
			name: "10.1.1.0/24",
			args: args{
				subnet: 10<<24 + 1<<16 + 1<<8 + 0,
				mask:   24,
			},
			want: 10<<24 + 1<<16 + 1<<8 + 255,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inetSubnetLastAddress(tt.args.subnet, tt.args.mask); got != tt.want {
				t.Errorf("inetSubnetLastAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inetSubnetNetmask(t *testing.T) {
	type args struct {
		mask int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "/0",
			args: args{
				mask: 0,
			},
			want: 0<<24 + 0<<16 + 0<<8 + 0,
		},
		{
			name: "/8",
			args: args{
				mask: 8,
			},
			want: 255<<24 + 0<<16 + 0<<8 + 0,
		},
		{
			name: "/16",
			args: args{
				mask: 16,
			},
			want: 255<<24 + 255<<16 + 0<<8 + 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inetSubnetNetmask(tt.args.mask); got != tt.want {
				t.Errorf("inetSubnetNetmask() = %v, want %v", got, tt.want)
			}
		})
	}
}
