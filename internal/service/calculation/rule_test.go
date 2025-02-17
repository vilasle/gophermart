package calculation

import (
	"regexp"
	"testing"

	"github.com/vilasle/gophermart/internal/service"
)

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
				exp:			 regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypePercent,
				value:		   5,
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
				exp:			 regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypeFixed,
				value:		   50,
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
				exp:			 regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationType(10),
				value:		   10,
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
				exp:			 regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypeFixed,
				value:		   10,
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
				exp:			 regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypeFixed,
				value:		   10,
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
				exp:			 regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypeFixed,
				value:		   10,
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
				exp:			 regexp.MustCompile("(?i)test"),
				calculationType: service.CalculationTypeFixed,
				value:		   10,
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
