package domain_test

import (
	"reflect"
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

func TestCreateBlock2MuteDomain(t *testing.T) {
	type args struct {
		NumberOfSuccess   uint
		SuccessTwitterIDs []string
	}
	tests := []struct {
		name string
		args args
		want entity.Block2Mute
	}{
		{
			name: "success",
			args: args{
				NumberOfSuccess:   1,
				SuccessTwitterIDs: []string{"1"},
			},
			want: *entity.NewBlock2Mute(
				1, []string{"1"},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entity.NewBlock2Mute(
				tt.args.NumberOfSuccess, tt.args.SuccessTwitterIDs,
			)
			if got.GetNumberOfSuccess() != tt.want.GetNumberOfSuccess() ||
				!reflect.DeepEqual(got.GetSuccessTwitterIDs(), tt.want.GetSuccessTwitterIDs()) {
				t.Errorf("createBlock2MuteDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
