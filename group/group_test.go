package group

import (
	"context"
	"github.com/savageking-io/ogbuser/db"
	"github.com/savageking-io/ogbuser/perm"
	"github.com/savageking-io/ogbuser/schema"
	"reflect"
	"testing"
)

func TestGroup_GetId(t *testing.T) {
	type fields struct {
		raw        schema.GroupSchema
		db         *db.Database
		perms      *perm.Perm
		hasRawData bool
		hasId      bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int32
	}{
		{"No raw data", fields{hasRawData: false}, 0},
		{"With raw data", fields{hasRawData: true, raw: schema.GroupSchema{Id: 123}}, 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				raw:        tt.fields.raw,
				db:         tt.fields.db,
				perms:      tt.fields.perms,
				hasRawData: tt.fields.hasRawData,
				hasId:      tt.fields.hasId,
			}
			if got := g.GetId(); got != tt.want {
				t.Errorf("GetId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_GetName(t *testing.T) {
	type fields struct {
		raw        schema.GroupSchema
		db         *db.Database
		perms      *perm.Perm
		hasRawData bool
		hasId      bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"No raw data", fields{hasRawData: false}, ""},
		{"With raw data", fields{hasRawData: true, raw: schema.GroupSchema{Name: "test"}}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				raw:        tt.fields.raw,
				db:         tt.fields.db,
				perms:      tt.fields.perms,
				hasRawData: tt.fields.hasRawData,
				hasId:      tt.fields.hasId,
			}
			if got := g.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_Init(t *testing.T) {
	type fields struct {
		raw        schema.GroupSchema
		db         *db.Database
		perms      *perm.Perm
		hasRawData bool
		hasId      bool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Nil Db", fields{}, args{context.Background()}, true},
		{"Has Id, but no raw data", fields{db: &db.Database{}, hasId: true}, args{context.Background()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				raw:        tt.fields.raw,
				db:         tt.fields.db,
				perms:      tt.fields.perms,
				hasRawData: tt.fields.hasRawData,
				hasId:      tt.fields.hasId,
			}
			if err := g.Init(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewGroup(t *testing.T) {
	type args struct {
		db *db.Database
	}
	tests := []struct {
		name string
		args args
		want *Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGroup(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGroupFromId(t *testing.T) {
	type args struct {
		db *db.Database
		id int32
	}
	tests := []struct {
		name string
		args args
		want *Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGroupFromId(tt.args.db, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGroupFromId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGroupFromSchema(t *testing.T) {
	type args struct {
		db     *db.Database
		schema *schema.GroupSchema
	}
	tests := []struct {
		name string
		args args
		want *Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGroupFromSchema(tt.args.db, tt.args.schema); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGroupFromSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
