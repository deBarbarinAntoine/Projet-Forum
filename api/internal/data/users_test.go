package data

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDBtest() (*sql.DB, error) {
	db, err := sql.Open("mysql", "forum:F0rumAP1@/forum?parseTime=true")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(15 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func TestUserModel_GetByID(t *testing.T) {

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

	thorganCreatedAt, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", "2024-05-25 20:02:23 +0000 UTC")
	thorganBirthDate, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", "0001-01-01 00:00:00 +0000 UTC")

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "thorgan",
			fields: fields{
				DB: db,
			},
			args: args{
				id: 1,
			},
			want: &User{
				ID:        1,
				CreatedAt: thorganCreatedAt,
				Name:      "Thorgan",
				Email:     "thorgan@example.com",
				Password: password{
					plaintext: nil,
					hash:      []byte("sdlkfjeifjsoldgkhd"),
				},
				Role:            "normal",
				BirthDate:       thorganBirthDate,
				Bio:             "",
				Signature:       "",
				Avatar:          "/avatar.png",
				Status:          "to-confirm",
				Version:         1,
				FollowingTags:   nil,
				FavoriteThreads: nil,
				CategoriesOwned: nil,
				TagsOwned:       nil,
				ThreadsOwned:    nil,
				Posts:           nil,
				Friends:         nil,
				Invitations: struct {
					Received []Friend `json:"received"`
					Sent     []Friend `json:"sent"`
				}{
					Received: nil,
					Sent:     nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := UserModel{
				DB: tt.fields.DB,
			}
			var t1 = time.Now()
			got, _ := m.GetByID(tt.args.id)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() running time: %v, got = %+v, want %+v", time.Since(t1).String(), got, tt.want)
			}
		})
	}
}
