package data

import (
	"database/sql"
	"reflect"
	"testing"
	"time"
)

func TestPostModel_GetByID(t *testing.T) {

	db, err := OpenDBtest()
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		id int
	}
	postCreatedAt, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", "2024-05-25 20:07:34 +0000 UTC")
	postUpdatedAt, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", "2024-05-25 20:07:34 +0000 UTC")
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Post
		wantErr bool
	}{
		{
			name: "Go tutorial",
			fields: fields{
				DB: db,
			},
			args: args{
				id: 1,
			},
			want: &Post{
				ID:        1,
				Content:   "Hello everyone, here is a beginner's course for the Go programming language!",
				CreatedAt: postCreatedAt,
				UpdatedAt: postUpdatedAt,
				Author: struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}{
					ID:   1,
					Name: "Thorgan",
				},
				IDParentPost: 0,
				Thread: struct {
					ID    int    `json:"id"`
					Title string `json:"title"`
				}{
					1,
					"Go programming language",
				},
				Version: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PostModel{
				DB: tt.fields.DB,
			}
			var t1 = time.Now()
			got, _ := p.GetByID(tt.args.id)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() running time: %v, got = %v, want %v", time.Since(t1).String(), got, tt.want)
			}
		})
	}
}
