package pkg

import "testing"

func TestHashShortening(t *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "#1",
			args: args{[]byte("https://ya.ru")},
			want: 2105327019,
		},
		{
			name: "#2",
			args: args{[]byte("yandex.ru")},
			want: 3785792127,
		},
		{
			name: "#3",
			args: args{[]byte("google.com")},
			want: 2006368837,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashShortening(tt.args.s); got != tt.want {
				t.Errorf("HashShortening() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLValidation(t *testing.T) {
	type args struct {
		inpURL string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "#1",
			args: args{inpURL: "https://ya.ru"},
			want: true,
		},
		{
			name: "#2",
			args: args{inpURL: "http://localhost:8080"},
			want: true,
		},
		{
			name: "#3",
			args: args{inpURL: "google"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URLValidation(tt.args.inpURL); got != tt.want {
				t.Errorf("URLValidation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHostOnly(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "#1",
			args: args{address: "https://ya.ru"},
			want: "https",
		},
		{
			name: "#2",
			args: args{address: "http://localhost:8080"},
			want: "http",
		},
		{
			name: "#3",
			args: args{address: "192.168.0.1:5001"},
			want: "192.168.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HostOnly(tt.args.address); got != tt.want {
				t.Errorf("HostOnly() = %v, want %v", got, tt.want)
			}
		})
	}
}
