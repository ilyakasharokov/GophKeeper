package customerrors

import (
	"errors"
	"testing"
)

func TestCustomError_Error(t *testing.T) {
	type fields struct {
		Err       error
		ErrorText string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "ok", fields: fields{Err: errors.New("error"), ErrorText: "error"}, want: "error"},
		{name: "empty", fields: fields{Err: errors.New(""), ErrorText: "error"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &CustomError{
				Err:       tt.fields.Err,
				ErrorText: tt.fields.ErrorText,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomError_Unwrap(t *testing.T) {
	type fields struct {
		Err       error
		ErrorText string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "ok", fields: fields{Err: errors.New("error"), ErrorText: "error"}, wantErr: true},
		{name: "error", fields: fields{Err: errors.New("error"), ErrorText: "error"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err1 := &CustomError{
				Err:       tt.fields.Err,
				ErrorText: tt.fields.ErrorText,
			}
			if err := err1.Unwrap(); (err != nil) != tt.wantErr {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	type args struct {
		err       error
		errorText string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "ok", args: args{
			err:       errors.New("test"),
			errorText: "test",
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewError(tt.args.err, tt.args.errorText); (err != nil) != tt.wantErr {
				t.Errorf("NewError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "ok", args: args{err: errors.New("test")}, want: "internal server error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseError(tt.args.err); got != tt.want {
				t.Errorf("ParseError() = %v, want %v", got, tt.want)
			}
		})
	}
}
