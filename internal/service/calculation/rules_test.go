package calculation

import (
	"context"
	"regexp"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/vilasle/gophermart/internal/repository/calculation"
	"github.com/vilasle/gophermart/internal/service"
)

func TestRuleService_Register(t *testing.T) {
	type behavior func(rep *MockCalculationRules, ctx context.Context, dto service.RegisterCalculationRuleRequest, id int16) error

	type args struct {
		ctx context.Context
		dto service.RegisterCalculationRuleRequest
		behavior
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: service.RegisterCalculationRuleRequest{
					Match: "test",
					Point: 5,
					Type:  service.CalculationTypePercent,
				},
				behavior: func(rep *MockCalculationRules, ctx context.Context, dto service.RegisterCalculationRuleRequest, id int16) error {
					repDto := repository.AddingRule{Match: dto.Match, Point: dto.Point, CalculationType: dto.Type}

					rep.EXPECT().AddRules(ctx, repDto).Return(id, nil)

					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rep := NewMockCalculationRules(ctrl)

			tt.args.behavior(rep, tt.args.ctx, tt.args.dto, 1)

			em := NewEventManager(tt.args.ctx)

			cfg := RuleServiceConfig{
				Repository: rep,
				EventManager:     em,
			}

			s := NewRuleService(cfg)

			err := s.Register(tt.args.ctx, tt.args.dto)
			if (err != nil) && !tt.wantErr {
				t.Errorf("RuleService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if (err == nil) && tt.wantErr {
				t.Errorf("RuleService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_rule_calculate(t *testing.T) {
	type fields struct {
		exp             *regexp.Regexp
		calculationType service.CalculationType
		value           float64
	}
	type args struct {
		name  string
		price float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "percent type",
			fields: fields{
				exp:             regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypePercent,
				value:           5,
			},
			args: args{
				name:  "test",
				price: 100,
			},
			want: 5,
		},
		{
			name: "fixed type",
			fields: fields{
				exp:             regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypeFixed,
				value:           50,
			},
			args: args{
				name:  "test",
				price: 100,
			},
			want: 50,
		},
		{
			name: "unknown type",
			fields: fields{
				exp:             regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationType(10),
				value:           10,
			},
			args: args{
				name:  "test",
				price: 100,
			},
			want: 0,
		},
		{
			name: "upper case",
			fields: fields{
				exp:             regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypeFixed,
				value:           10,
			},
			args: args{
				name:  "TEST",
				price: 100,
			},
			want: 10,
		},
		{
			name: "match in middle",
			fields: fields{
				exp:             regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypeFixed,
				value:           10,
			},
			args: args{
				name:  "fvdfg23432TEST32423fvfdv23",
				price: 100,
			},
			want: 10,
		},
		{
			name: "match in start",
			fields: fields{
				exp:             regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypeFixed,
				value:           10,
			},
			args: args{
				name:  "testfvdfg23432TEST32423fvfdv23",
				price: 100,
			},
			want: 10,
		},
		{
			name: "match in finish",
			fields: fields{
				exp:             regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypeFixed,
				value:           10,
			},
			args: args{
				name:  "fvdfg2343232423fvfdv23test",
				price: 100,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &rule{
				exp:             tt.fields.exp,
				calculationType: tt.fields.calculationType,
				value:           tt.fields.value,
			}
			if got := r.calculate(tt.args.name, tt.args.price); got != tt.want {
				t.Errorf("rule.calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
