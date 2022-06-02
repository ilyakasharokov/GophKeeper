package crypto

import (
	"reflect"
	"testing"
)

func TestDecrypt(t *testing.T) {
	type args struct {
		key  []byte
		text []byte
	}
	hash := Hash("ilya")
	text := []byte("mytext")
	emptyText := []byte("")
	enc, _ := Encrypt(hash, text)
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok", args: args{key: hash, text: enc}, want: text,
			wantErr: false,
		},
		{
			name: "err", args: args{key: hash, text: emptyText}, want: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.key, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		key  []byte
		text []byte
	}
	hash := Hash("ilya")
	text := []byte("mytext")
	emptyText := []byte("")
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "ok", args: args{key: hash, text: text}, wantErr: false},
		{name: "err", args: args{key: emptyText, text: emptyText}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Encrypt(tt.args.key, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestHash(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{name: "ok", args: args{s: "ilya"}, want: []byte{122, 100, 130, 151, 202, 198, 98, 8, 136, 159, 110, 213, 164, 250, 142, 67, 81, 100, 20, 186, 131, 229, 196, 217, 141, 5, 232, 89, 158, 237, 230, 181}},
		{name: "empty", args: args{s: ""}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
