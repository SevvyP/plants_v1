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

	"github.com/SevvyP/items/internal/db"
	"github.com/SevvyP/items/pkg"
	"github.com/gin-gonic/gin"
)

func TestServer_HandleCreateItem(t *testing.T) {
	type fields struct {
		db *db.MockDB
	}
	type args struct {
		item pkg.Item
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pkg.Item
		code   int
		checkReturn bool
	}{
		{
			name: "handle create item fails if body is missing",
			fields: fields{
				db: nil,
			},
			code: 400,
		},
		{
			name: "handle create item fails if db returns an error",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				item: pkg.Item{Name: "test", Description: "test"},
				err: errors.New("test"),
			},
			code: 500,
		},
		{
			name: "handle create item is successful if db is successful",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				item: pkg.Item{Name: "test", Description: "test"},
				err: nil,
			},
			code: 200,
			want: pkg.Item{Name: "test", Description: "test"},
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
			jsonbytes, err := json.Marshal(tt.args.item)
			if err != nil {
				t.Error(err)
			} 
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
			s := &Server{
				db: tt.fields.db,
			}
			if(tt.fields.db != nil) {
				tt.fields.db.On("CreateItem", tt.args.item, c).Return(tt.args.err)
			}
			s.HandleCreateItem(c)
			if c.Writer.Status() != tt.code {
				t.Errorf("HandleCreateItem response code: %d, expected %d", c.Writer.Status(), tt.code)
			}
			if tt.checkReturn {
				var got pkg.Item
				err = json.Unmarshal(w.Body.Bytes(), &got)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Handle create item returned %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestServer_HandleGetItem(t *testing.T) {
	type fields struct {
		db *db.MockDB
	}
	type args struct {
		name string
		err error
		item *pkg.Item
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pkg.Item
		code   int
		checkReturn bool
	}{
		{
			name: "handle get item fails if body is missing",
			fields: fields{
				db: nil,
			},
			code: 400,
		},
		{
			name: "handle get item fails if db returns an error",
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
			name: "handle get item fails if item is not found",
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
			name: "handle get item is successful if db is successful",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				name: "test",
				err: nil,
				item: &pkg.Item{Name: "test", Description: "test"},
			},
			code: 200,
			want: pkg.Item{Name: "test", Description: "test"},
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
				tt.fields.db.On("GetItem", tt.args.name, c).Return(tt.args.item, tt.args.err)
			}
			s.HandleGetItem(c)
			if c.Writer.Status() != tt.code {
				t.Errorf("HandleCreateItem response code: %d, expected %d", c.Writer.Status(), tt.code)
			}
			if tt.checkReturn {
				var got pkg.Item
				err := json.Unmarshal(w.Body.Bytes(), &got)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Handle create item returned %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestServer_HandleUpdateItem(t *testing.T) {
	type fields struct {
		db *db.MockDB
	}
	type args struct {
		item pkg.Item
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pkg.Item
		code   int
		checkReturn bool
	}{
		{
			name: "handle update item fails if body is missing",
			fields: fields{
				db: nil,
			},
			code: 400,
		},
		{
			name: "handle update item fails if db returns an error",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				item: pkg.Item{Name: "test", Description: "test"},
				err: errors.New("test"),
			},
			code: 500,
		},
		{
			name: "handle update item is successful if db is successful",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				item: pkg.Item{Name: "test", Description: "test"},
				err: nil,
			},
			code: 200,
			want: pkg.Item{Name: "test", Description: "test"},
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
			jsonbytes, err := json.Marshal(tt.args.item)
			if err != nil {
				t.Error(err)
			} 
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
			s := &Server{
				db: tt.fields.db,
			}
			if(tt.fields.db != nil) {
				tt.fields.db.On("UpdateItem", tt.args.item, c).Return(tt.args.err)
			}
			s.HandleUpdateItem(c)
			if c.Writer.Status() != tt.code {
				t.Errorf("HandleUpdateItem response code: %d, expected %d", c.Writer.Status(), tt.code)
			}
			if tt.checkReturn {
				var got pkg.Item
				err = json.Unmarshal(w.Body.Bytes(), &got)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Handle update item returned %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestServer_HandleDeleteItem(t *testing.T) {
	type fields struct {
		db *db.MockDB
	}
	type args struct {
		name string
		err error
		item *pkg.Item
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pkg.Item
		code   int
		checkReturn bool
	}{
		{
			name: "handle delete item fails if body is missing",
			fields: fields{
				db: nil,
			},
			code: 400,
		},
		{
			name: "handle delete item fails if db returns an error",
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
			name: "handle delete item fails if item is not found",
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
			name: "handle delete item is successful if db is successful",
			fields: fields{
				db: new(db.MockDB),
			},
			args: args{
				name: "test",
				err: nil,
				item: &pkg.Item{Name: "test", Description: "test"},
			},
			code: 200,
			want: pkg.Item{Name: "test", Description: "test"},
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
				tt.fields.db.On("DeleteItem", tt.args.name, c).Return(tt.args.item, tt.args.err)
			}
			s.HandleDeleteItem(c)
			if c.Writer.Status() != tt.code {
				t.Errorf("HandleDeleteItem response code: %d, expected %d", c.Writer.Status(), tt.code)
			}
			if tt.checkReturn {
				var got pkg.Item
				err := json.Unmarshal(w.Body.Bytes(), &got)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Handle delete item returned %v, want %v", got, tt.want)
				}
			}
		})
	}
}
