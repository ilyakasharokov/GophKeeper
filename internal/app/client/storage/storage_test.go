package storage

import (
	"gophkeeper/internal/common/models"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		fileStoragePath string
	}
	tests := []struct {
		name string
		args args
		want *Storage
	}{
		{name: "ok", args: struct{ fileStoragePath string }{fileStoragePath: ""}, want: &Storage{fileStoragePath: "", Check: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.fileStoragePath); !got.Check {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_AddItem(t *testing.T) {
	type fields struct {
		Data            []models.Note
		LastSyncDate    time.Time
		fileStoragePath string
		Check           bool
	}
	type args struct {
		title string
		body  string
		key   []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "small key size", fields: struct {
			Data            []models.Note
			LastSyncDate    time.Time
			fileStoragePath string
			Check           bool
		}{Data: make([]models.Note, 0), LastSyncDate: time.Now(), fileStoragePath: "", Check: true}, args: struct {
			title string
			body  string
			key   []byte
		}{title: "test", body: "body", key: []byte("key")}, wantErr: true},
		{name: "ok", fields: struct {
			Data            []models.Note
			LastSyncDate    time.Time
			fileStoragePath string
			Check           bool
		}{Data: make([]models.Note, 0), LastSyncDate: time.Now(), fileStoragePath: "", Check: true}, args: struct {
			title string
			body  string
			key   []byte
		}{title: "test", body: "body", key: []byte("keyqwekeyqwekeyqwekeyqwekeyqweqw")}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Data:            tt.fields.Data,
				LastSyncDate:    tt.fields.LastSyncDate,
				fileStoragePath: tt.fields.fileStoragePath,
				Check:           tt.fields.Check,
			}
			if err := s.AddItem(tt.args.title, tt.args.body, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_CheckFile(t *testing.T) {
	type fields struct {
		Data            []models.Note
		LastSyncDate    time.Time
		fileStoragePath string
		Check           bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "ok", fields: struct {
			Data            []models.Note
			LastSyncDate    time.Time
			fileStoragePath string
			Check           bool
		}{Data: nil, LastSyncDate: time.Now(), fileStoragePath: "", Check: true}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &Storage{
				Data:            tt.fields.Data,
				LastSyncDate:    tt.fields.LastSyncDate,
				fileStoragePath: tt.fields.fileStoragePath,
				Check:           tt.fields.Check,
			}
			if got := storage.CheckFile(); got != tt.want {
				t.Errorf("CheckFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Flush(t *testing.T) {
	type fields struct {
		Data            []models.Note
		LastSyncDate    time.Time
		fileStoragePath string
		Check           bool
	}
	type args struct {
		hash []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "ok", fields: struct {
			Data            []models.Note
			LastSyncDate    time.Time
			fileStoragePath string
			Check           bool
		}{Data: nil, LastSyncDate: time.Now(), fileStoragePath: "", Check: true}, args: struct{ hash []byte }{hash: []byte("qweqwew")}, wantErr: true},
		{name: "ok", fields: struct {
			Data            []models.Note
			LastSyncDate    time.Time
			fileStoragePath string
			Check           bool
		}{Data: nil, LastSyncDate: time.Now(), fileStoragePath: "", Check: true}, args: struct{ hash []byte }{hash: []byte("qweqasdsadsadsadsswew")}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &Storage{
				Data:            tt.fields.Data,
				LastSyncDate:    tt.fields.LastSyncDate,
				fileStoragePath: tt.fields.fileStoragePath,
				Check:           tt.fields.Check,
			}
			if err := storage.Flush(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Flush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_GetByIndex(t *testing.T) {
	type fields struct {
		Data            []models.Note
		LastSyncDate    time.Time
		fileStoragePath string
		Check           bool
	}
	type args struct {
		index int
	}
	notes := []models.Note{{
		ID: "123",
	}}
	flds := fields{Data: notes, LastSyncDate: time.Now()}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Note
		wantErr bool
	}{
		{name: "ok", fields: flds, args: struct{ index int }{index: 0}, want: notes[0], wantErr: false},
		{name: "not ok", fields: flds, args: struct{ index int }{index: 1}, want: models.Note{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Data:            tt.fields.Data,
				LastSyncDate:    tt.fields.LastSyncDate,
				fileStoragePath: tt.fields.fileStoragePath,
				Check:           tt.fields.Check,
			}
			got, err := s.GetByIndex(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetDataLen(t *testing.T) {
	type fields struct {
		Data            []models.Note
		LastSyncDate    time.Time
		fileStoragePath string
		Check           bool
	}
	flds := fields{
		Data: make([]models.Note, 5),
	}
	flds2 := fields{}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{name: "ok", fields: flds, want: 5},
		{name: "not ok", fields: flds2, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Data:            tt.fields.Data,
				LastSyncDate:    tt.fields.LastSyncDate,
				fileStoragePath: tt.fields.fileStoragePath,
				Check:           tt.fields.Check,
			}
			if got := s.GetDataLen(); got != tt.want {
				t.Errorf("GetDataLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetLastSyncDate(t *testing.T) {
	type fields struct {
		Data            []models.Note
		LastSyncDate    time.Time
		fileStoragePath string
		Check           bool
	}
	lsd := time.Now()
	flds := fields{
		LastSyncDate: lsd,
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{name: "ok", fields: flds, want: lsd},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Data:            tt.fields.Data,
				LastSyncDate:    tt.fields.LastSyncDate,
				fileStoragePath: tt.fields.fileStoragePath,
				Check:           tt.fields.Check,
			}
			if got := s.GetLastSyncDate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLastSyncDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetNonSyncedData(t *testing.T) {
	type fields struct {
		Data            []models.Note
		LastSyncDate    time.Time
		fileStoragePath string
		Check           bool
	}
	flds := fields{
		Data: []models.Note{
			{CreatedAt: time.Now().Add(1 * time.Minute)},
			{CreatedAt: time.Now().Add(5 * time.Minute)},
			{CreatedAt: time.Now().Add(-1 * time.Minute)},
		},
		LastSyncDate: time.Now(),
	}
	tests := []struct {
		name   string
		fields fields
		want   []models.Note
	}{
		{name: "ok", fields: flds, want: []models.Note{flds.Data[0], flds.Data[1]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Data:            tt.fields.Data,
				LastSyncDate:    tt.fields.LastSyncDate,
				fileStoragePath: tt.fields.fileStoragePath,
				Check:           tt.fields.Check,
			}
			if got := s.GetNonSyncedData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNonSyncedData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetNotes(t *testing.T) {
	type fields struct {
		Data            []models.Note
		LastSyncDate    time.Time
		fileStoragePath string
		Check           bool
	}
	type args struct {
		all bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []models.Note
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Data:            tt.fields.Data,
				LastSyncDate:    tt.fields.LastSyncDate,
				fileStoragePath: tt.fields.fileStoragePath,
				Check:           tt.fields.Check,
			}
			if got := s.GetNotes(tt.args.all); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Load(t *testing.T) {
	type fields struct {
		Data            []models.Note
		LastSyncDate    time.Time
		fileStoragePath string
		Check           bool
	}
	type args struct {
		hash []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &Storage{
				Data:            tt.fields.Data,
				LastSyncDate:    tt.fields.LastSyncDate,
				fileStoragePath: tt.fields.fileStoragePath,
				Check:           tt.fields.Check,
			}
			if err := storage.Load(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_SetDeleted(t *testing.T) {
	type fields struct {
		Data            []models.Note
		LastSyncDate    time.Time
		fileStoragePath string
		Check           bool
	}
	type args struct {
		index int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Data:            tt.fields.Data,
				LastSyncDate:    tt.fields.LastSyncDate,
				fileStoragePath: tt.fields.fileStoragePath,
				Check:           tt.fields.Check,
			}
			if err := s.SetDeleted(tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("SetDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_UpdateData(t *testing.T) {
	type fields struct {
		Data            []models.Note
		LastSyncDate    time.Time
		fileStoragePath string
		Check           bool
	}
	type args struct {
		newdata  []models.Note
		lastSync time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Data:            tt.fields.Data,
				LastSyncDate:    tt.fields.LastSyncDate,
				fileStoragePath: tt.fields.fileStoragePath,
				Check:           tt.fields.Check,
			}
			if err := s.UpdateData(tt.args.newdata, tt.args.lastSync); (err != nil) != tt.wantErr {
				t.Errorf("UpdateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newConsumer(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    *consumer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConsumer(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("newConsumer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newConsumer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newProducer(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    *producer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newProducer(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("newProducer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newProducer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
