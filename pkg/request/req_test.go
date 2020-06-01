package request

import (
	"gin-web/models"
	"reflect"
	"testing"
)

func TestUpdateIncrementalIdsRequestStruct_GetIncremental(t *testing.T) {
	menus := make([]models.SysMenu, 0)
	for i := 0; i < 10; i++ {
		menus = append(menus, models.SysMenu{
			Model: models.Model{
				Id: uint(i + 1),
			},
		})
	}
	type fields struct {
		Create []uint
		Delete []uint
	}
	type args struct {
		oldList []uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []uint
	}{
		{
			name: "case1",
			fields: fields{
				Create: []uint{1, 2, 3, 4, 5},
				Delete: []uint{6, 7, 8, 9, 10},
			},
			args: args{[]uint{6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
			want: []uint{11, 12, 13, 14, 15, 1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UpdateIncrementalIdsRequestStruct{
				Create: tt.fields.Create,
				Delete: tt.fields.Delete,
			}
			if got := s.GetIncremental(tt.args.oldList, menus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIncremental() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 性能测试
func BenchmarkUpdateIncrementalIdsRequestStruct_GetIncremental(b *testing.B) {
	menus := make([]models.SysMenu, 0)
	for i := 0; i < 10; i++ {
		menus = append(menus, models.SysMenu{
			Model: models.Model{
				Id: uint(i + 1),
			},
		})
	}
	s := &UpdateIncrementalIdsRequestStruct{
		Create: []uint{1, 2, 3, 4, 5},
		Delete: []uint{6, 7, 8, 9, 10},
	}
	oldList := []uint{6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s.GetIncremental(oldList, menus)
	}
}
