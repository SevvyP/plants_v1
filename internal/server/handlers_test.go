package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/SevvyP/plants/internal/db"
	"github.com/SevvyP/plants/pkg"
	"github.com/gin-gonic/gin"
)

func TestServer_HandleCreatePlant(t *testing.T) {
	type fields struct {
		db *db.MockDB
	}
	type args struct {
		plant pkg.Plant
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pkg.Plant
		code   int
		checkReturn bool
	}{
		{
			name: "handle create plant fails if body is missing",
			fields: fields{
				db: nil,
			},
			code: 400,
		},
		{
			name: "handle create plant fails if db returns an error",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				plant: pkg.Plant{Name: "test", Description: "test"},
				err: errors.New("test"),
			},
			code: 500,
		},
		{
			name: "handle create plant is successful if db is successful",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				plant: pkg.Plant{Name: "test", Description: "test"},
				err: nil,
			},
			code: 200,
			want: pkg.Plant{Name: "test", Description: "test"},
			checkReturn: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = &http.Request{
				Header: make(http.Header),
			}
			c.Request.Method = "POST"
			c.Request.Header.Set("Content-Type", "application/json")
			jsonbytes, err := json.Marshal(tt.args.plant)
			if err != nil {
				t.Error(err)
			} 
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
			s := &Server{
				db: tt.fields.db,
			}
			if(tt.fields.db != nil) {
				tt.fields.db.On("CreatePlant", tt.args.plant, c).Return(tt.args.err)
			}
			s.HandleCreatePlant(c)
			if c.Writer.Status() != tt.code {
				t.Errorf("HandleCreatePlant response code: %d, expected %d", c.Writer.Status(), tt.code)
			}
			if tt.checkReturn {
				var got pkg.Plant
				err = json.Unmarshal(w.Body.Bytes(), &got)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Handle create plant returned %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestServer_HandleGetPlant(t *testing.T) {
	type fields struct {
		db *db.MockDB
	}
	type args struct {
		name string
		err error
		plant *pkg.Plant
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pkg.Plant
		code   int
		checkReturn bool
	}{
		{
			name: "handle get plant fails if body is missing",
			fields: fields{
				db: nil,
			},
			code: 400,
		},
		{
			name: "handle get plant fails if db returns an error",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				name: "test",
				err: errors.New("test"),
			},
			code: 500,
		},
		{
			name: "handle get plant fails if item is not found",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				name: "test",
				err: errors.New(db.ErrNotFound),
			},
			code: 404,
		},
		{
			name: "handle get plant is successful if db is successful",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				name: "test",
				err: nil,
				plant: &pkg.Plant{Name: "test", Description: "test"},
			},
			code: 200,
			want: pkg.Plant{Name: "test", Description: "test"},
			checkReturn: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = &http.Request{
				Header: make(http.Header),
			}
			c.Request.Method = "GET"
			c.Params = append(c.Params, gin.Param{Key: "name", Value: tt.args.name})
			s := &Server{
				db: tt.fields.db,
			}
			if(tt.fields.db != nil) {
				tt.fields.db.On("GetPlant", tt.args.name, c).Return(tt.args.plant, tt.args.err)
			}
			s.HandleGetPlant(c)
			if c.Writer.Status() != tt.code {
				t.Errorf("HandleCreatePlant response code: %d, expected %d", c.Writer.Status(), tt.code)
			}
			if tt.checkReturn {
				var got pkg.Plant
				err := json.Unmarshal(w.Body.Bytes(), &got)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Handle create plant returned %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestServer_HandleUpdatePlant(t *testing.T) {
	type fields struct {
		db *db.MockDB
	}
	type args struct {
		plant pkg.Plant
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pkg.Plant
		code   int
		checkReturn bool
	}{
		{
			name: "handle update plant fails if body is missing",
			fields: fields{
				db: nil,
			},
			code: 400,
		},
		{
			name: "handle update plant fails if db returns an error",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				plant: pkg.Plant{Name: "test", Description: "test"},
				err: errors.New("test"),
			},
			code: 500,
		},
		{
			name: "handle update plant is successful if db is successful",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				plant: pkg.Plant{Name: "test", Description: "test"},
				err: nil,
			},
			code: 200,
			want: pkg.Plant{Name: "test", Description: "test"},
			checkReturn: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = &http.Request{
				Header: make(http.Header),
			}
			c.Request.Method = "POST"
			c.Request.Header.Set("Content-Type", "application/json")
			jsonbytes, err := json.Marshal(tt.args.plant)
			if err != nil {
				t.Error(err)
			} 
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
			s := &Server{
				db: tt.fields.db,
			}
			if(tt.fields.db != nil) {
				tt.fields.db.On("UpdatePlant", tt.args.plant, c).Return(tt.args.err)
			}
			s.HandleUpdatePlant(c)
			if c.Writer.Status() != tt.code {
				t.Errorf("HandleUpdatePlant response code: %d, expected %d", c.Writer.Status(), tt.code)
			}
			if tt.checkReturn {
				var got pkg.Plant
				err = json.Unmarshal(w.Body.Bytes(), &got)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Handle update plant returned %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestServer_HandleDeletePlant(t *testing.T) {
	type fields struct {
		db *db.MockDB
	}
	type args struct {
		name string
		err error
		plant *pkg.Plant
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pkg.Plant
		code   int
		checkReturn bool
	}{
		{
			name: "handle delete plant fails if body is missing",
			fields: fields{
				db: nil,
			},
			code: 400,
		},
		{
			name: "handle delete plant fails if db returns an error",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				name: "test",
				err: errors.New("test"),
			},
			code: 500,
		},
		{
			name: "handle delete plant fails if item is not found",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				name: "test",
				err: errors.New(db.ErrNotFound),
			},
			code: 404,
		},
		{
			name: "handle delete plant is successful if db is successful",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				name: "test",
				err: nil,
				plant: &pkg.Plant{Name: "test", Description: "test"},
			},
			code: 200,
			want: pkg.Plant{Name: "test", Description: "test"},
			checkReturn: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = &http.Request{
				Header: make(http.Header),
			}
			c.Request.Method = "DELETE"
			c.Params = append(c.Params, gin.Param{Key: "name", Value: tt.args.name})
			s := &Server{
				db: tt.fields.db,
			}
			if(tt.fields.db != nil) {
				tt.fields.db.On("DeletePlant", tt.args.name, c).Return(tt.args.plant, tt.args.err)
			}
			s.HandleDeletePlant(c)
			if c.Writer.Status() != tt.code {
				t.Errorf("HandleDeletePlant response code: %d, expected %d", c.Writer.Status(), tt.code)
			}
			if tt.checkReturn {
				var got pkg.Plant
				err := json.Unmarshal(w.Body.Bytes(), &got)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Handle delete plant returned %v, want %v", got, tt.want)
				}
			}
		})
	}
}
