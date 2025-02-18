package calculation

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vilasle/gophermart/internal/logger"
	"github.com/vilasle/gophermart/internal/repository/calculation"
	"github.com/vilasle/gophermart/internal/service"
)

type MockLoggerWriter struct {
	msgs []string
}

func NewMockLoggerWriter() *MockLoggerWriter {
	return &MockLoggerWriter{
		msgs: make([]string, 0),
	}
}

func (l *MockLoggerWriter) Write(p []byte) (n int, err error) {
	l.msgs = append(l.msgs, string(p))
	return len(p), nil
}

func TestCalculationService_fillRules(t *testing.T) {
	type fields struct {
		repCalc  repository.CalculationRepository
		repRules repository.CalculationRules
		mxRules  *sync.Mutex
		rules    map[int16]rule
		manager  *EventManager
	}
	type args struct {
		rs []repository.RuleInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fill rules",
			fields: fields{
				repCalc:  &MockCalculationRepository{},
				repRules: &MockCalculationRules{},
				mxRules:  &sync.Mutex{},
				rules:    make(map[int16]rule),
				manager:  NewEventManager(context.Background()),
			},
			args: args{
				rs: []repository.RuleInfo{
					{
						ID:              1,
						Match:           "test",
						Point:           1234,
						CalculationType: service.CalculationTypeFixed,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "wrong type of calculation type",
			fields: fields{
				repCalc:  &MockCalculationRepository{},
				repRules: &MockCalculationRules{},
				mxRules:  &sync.Mutex{},
				rules:    make(map[int16]rule),
				manager:  NewEventManager(context.Background()),
			},
			args: args{
				rs: []repository.RuleInfo{
					{
						ID:              1,
						Match:           "test",
						Point:           1234,
						CalculationType: -21321,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "wrong pattern",
			fields: fields{
				repCalc:  &MockCalculationRepository{},
				repRules: &MockCalculationRules{},
				mxRules:  &sync.Mutex{},
				rules:    make(map[int16]rule),
				manager:  NewEventManager(context.Background()),
			},
			args: args{
				rs: []repository.RuleInfo{
					{
						ID:              1,
						Match:           `test\`,
						Point:           1234,
						CalculationType: service.CalculationTypeFixed,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CalculationService{
				repCalc:  tt.fields.repCalc,
				repRules: tt.fields.repRules,
				mxRules:  tt.fields.mxRules,
				rules:    tt.fields.rules,
				manager:  tt.fields.manager,
			}
			err := c.fillRules(tt.args.rs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, c.rules, len(tt.args.rs))
			}

		})
	}
}

func TestCalculationService_readRule(t *testing.T) {
	type behavior func(*MockCalculationRules, context.Context, repository.RuleFilter)

	type fields struct {
		repCalc repository.CalculationRepository
		mxRules *sync.Mutex
		rules   map[int16]rule
		manager *EventManager
	}
	type args struct {
		ctx context.Context
		dto repository.RuleFilter
		behavior
		event Event
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		resultMsg string
	}{
		{
			name: "read specific rule",
			fields: fields{
				repCalc: &MockCalculationRepository{},
				mxRules: &sync.Mutex{},
				rules:   make(map[int16]rule),
				manager: NewEventManager(context.Background()),
			},
			args: args{
				ctx: context.Background(),
				event: Event{
					Type: NewRule,
					Data: int16(1),
				},
				behavior: func(mcr *MockCalculationRules, ctx context.Context, dto repository.RuleFilter) {
					var err error
					result := []repository.RuleInfo{
						{
							ID:              dto.ID,
							Match:           "test",
							Point:           1234,
							CalculationType: service.CalculationTypeFixed,
						},
					}
					mcr.EXPECT().Rules(ctx, dto).Return(result, err)
				},
				dto: repository.RuleFilter{ID: 1},
			},
			resultMsg: "",
		},
		{
			name: "event has wrong data",
			fields: fields{
				repCalc: &MockCalculationRepository{},
				mxRules: &sync.Mutex{},
				rules:   make(map[int16]rule),
				manager: NewEventManager(context.Background()),
			},
			args: args{
				ctx: context.Background(),
				event: Event{
					Type: NewRule,
					Data: "1",
				},
				dto:      repository.RuleFilter{},
				behavior: func(mcr *MockCalculationRules, ctx context.Context, dto repository.RuleFilter) {},
			},
			resultMsg: "was raise event with wrong data",
		},
		{
			name: "repository does not have rules by id",
			fields: fields{
				repCalc: &MockCalculationRepository{},
				mxRules: &sync.Mutex{},
				rules:   make(map[int16]rule),
				manager: NewEventManager(context.Background()),
			},
			args: args{
				ctx: context.Background(),
				event: Event{
					Type: NewRule,
					Data: int16(1),
				},
				dto: repository.RuleFilter{ID: 1},
				behavior: func(mcr *MockCalculationRules, ctx context.Context, dto repository.RuleFilter) {
					var err = errors.New("repository does not have rule")
					result := []repository.RuleInfo{}
					mcr.EXPECT().Rules(ctx, dto).Return(result, err)
				},
			},
			resultMsg: "getting specific rule",
		},
		{
			name: "wrong match expression",
			fields: fields{
				repCalc: &MockCalculationRepository{},
				mxRules: &sync.Mutex{},
				rules:   make(map[int16]rule),
				manager: NewEventManager(context.Background()),
			},
			args: args{
				ctx: context.Background(),
				event: Event{
					Type: NewRule,
					Data: int16(1),
				},
				dto: repository.RuleFilter{ID: 1},
				behavior: func(mcr *MockCalculationRules, ctx context.Context, dto repository.RuleFilter) {
					var err error
					result := []repository.RuleInfo{
						{
							ID:              1,
							Match:           `test\`,
							Point:           1234,
							CalculationType: service.CalculationTypeFixed,
						},
					}
					mcr.EXPECT().Rules(ctx, dto).Return(result, err)
				},
			},
			resultMsg: "preparing rules for using",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//we have only log message on error because will check messages which writer stored
			wrt := NewMockLoggerWriter()

			logger.Init(wrt, logger.DebugLevel)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repRules := NewMockCalculationRules(ctrl)

			tt.args.behavior(repRules, tt.args.ctx, tt.args.dto)

			c := CalculationService{
				repCalc:  tt.fields.repCalc,
				repRules: repRules,
				mxRules:  tt.fields.mxRules,
				rules:    tt.fields.rules,
				manager:  tt.fields.manager,
			}

			c.readRule(tt.args.ctx, tt.args.event)

			var lenMsg int
			if tt.resultMsg != "" {
				lenMsg += 1
			}

			assert.Len(t, wrt.msgs, lenMsg, "writer had to get message")

			if len(wrt.msgs) > 0 {
				assert.True(t, strings.Contains(wrt.msgs[0], tt.resultMsg), "logger got wrong message")
			}
		})
	}
}

func TestCalculationService_readAllRules(t *testing.T) {
	type behavior func(*MockCalculationRules, context.Context)

	type fields struct {
		repCalc repository.CalculationRepository
		mxRules *sync.Mutex
		rules   map[int16]rule
		manager *EventManager
	}
	type args struct {
		ctx context.Context
		behavior
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "read all rules",
			fields: fields{
				repCalc: &MockCalculationRepository{},
				mxRules: &sync.Mutex{},
				rules:   make(map[int16]rule),
				manager: NewEventManager(context.Background()),
			},
			args: args{
				ctx: context.Background(),
				behavior: func(mcr *MockCalculationRules, ctx context.Context) {
					var err error
					result := []repository.RuleInfo{
						{
							ID:              1,
							Match:           "test",
							Point:           1234,
							CalculationType: service.CalculationTypeFixed,
						},
					}
					mcr.EXPECT().Rules(ctx, repository.RuleFilter{}).Return(result, err)
				},
			},
			wantErr: false,
		},
		{
			name: "repository error",
			fields: fields{
				repCalc: &MockCalculationRepository{},
				mxRules: &sync.Mutex{},
				rules:   make(map[int16]rule),
				manager: NewEventManager(context.Background()),
			},
			args: args{
				ctx: context.Background(),
				behavior: func(mcr *MockCalculationRules, ctx context.Context) {
					var err = errors.New("repository error")
					result := []repository.RuleInfo{}
					mcr.EXPECT().Rules(ctx, repository.RuleFilter{}).Return(result, err)
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//we have only log message on error because will check messages which writer stored
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repRules := NewMockCalculationRules(ctrl)

			tt.args.behavior(repRules, tt.args.ctx)

			c := CalculationService{
				repCalc:  tt.fields.repCalc,
				repRules: repRules,
				mxRules:  tt.fields.mxRules,
				rules:    tt.fields.rules,
				manager:  tt.fields.manager,
			}

			err := c.readAllRules()

			if tt.wantErr {
				assert.Error(t, err, "readAllRules() should return an error")
			} else {
				assert.NoError(t, err, "readAllRules() should not return an error")
			}
		})
	}
}

func TestCalculationService_fillCalculatedDto(t *testing.T) {
	type fields struct {
		repCalc  repository.CalculationRepository
		repRules repository.CalculationRules
		mxRules  *sync.Mutex
		rules    map[int16]rule
		manager  *EventManager
	}
	type args struct {
		orderNumber string
		value       float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   repository.AddCalculationResult
	}{
		{
			name: "fill successful calculated dto",
			fields: fields{
				repCalc:  &MockCalculationRepository{},
				repRules: &MockCalculationRules{},
				mxRules:  &sync.Mutex{},
				rules:    make(map[int16]rule),
				manager:  NewEventManager(context.Background()),
			},
			args: args{
				orderNumber: "1234567890",
				value:       100,
			},
			want: repository.AddCalculationResult{
				OrderNumber: "1234567890",
				Value:       100,
				Status:      repository.Processed,
			},
		},
		{
			name: "fill invalid calculated dto",
			fields: fields{
				repCalc:  &MockCalculationRepository{},
				repRules: &MockCalculationRules{},
				mxRules:  &sync.Mutex{},
				rules:    make(map[int16]rule),
				manager:  NewEventManager(context.Background()),
			},
			args: args{
				orderNumber: "1234567890",
				value:       0,
			},
			want: repository.AddCalculationResult{
				OrderNumber: "1234567890",
				Value:       0,
				Status:      repository.Invalid,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CalculationService{
				repCalc:  tt.fields.repCalc,
				repRules: tt.fields.repRules,
				mxRules:  tt.fields.mxRules,
				rules:    tt.fields.rules,
				manager:  tt.fields.manager,
			}
			if got := c.fillCalculatedDto(tt.args.orderNumber, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculationService.fillCalculatedDto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculationService_calculateProduct(t *testing.T) {
	type fields struct {
		repCalc  repository.CalculationRepository
		repRules repository.CalculationRules
		mxRules  *sync.Mutex
		rules    map[int16]rule
		manager  *EventManager
	}
	type args struct {
		product service.ProductRow
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "calculate product",
			fields: fields{
				repCalc:  &MockCalculationRepository{},
				repRules: &MockCalculationRules{},
				mxRules:  &sync.Mutex{},
				rules: map[int16]rule{
					1: {
						calculationType: service.CalculationTypeFixed,
						value:           10,
						exp:             regexp.MustCompile("(?i)test"),
					},
					2: {
						calculationType: service.CalculationTypePercent,
						value:           10,
						exp:             regexp.MustCompile("(?i)test"),
					},
				},
				manager: NewEventManager(context.Background()),
			},
			args: args{
				product: service.ProductRow{
					Name:  "test",
					Price: 100,
				},
			},
			want: 10,
		},
		{
			name: "not found matched rules",
			fields: fields{
				repCalc:  &MockCalculationRepository{},
				repRules: &MockCalculationRules{},
				mxRules:  &sync.Mutex{},
				rules: map[int16]rule{
					1: {
						calculationType: service.CalculationTypeFixed,
						value:           10,
						exp:             regexp.MustCompile("(?i)test"),
					},
					2: {
						calculationType: service.CalculationTypePercent,
						value:           10,
						exp:             regexp.MustCompile("(?i)test"),
					},
				},
				manager: NewEventManager(context.Background()),
			},
			args: args{
				product: service.ProductRow{
					Name:  "ping-pong",
					Price: 100,
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CalculationService{
				repCalc:  tt.fields.repCalc,
				repRules: tt.fields.repRules,
				mxRules:  tt.fields.mxRules,
				rules:    tt.fields.rules,
				manager:  tt.fields.manager,
			}
			if got := c.calculateProduct(tt.args.product); got != tt.want {
				t.Errorf("CalculationService.calculateProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculationService_calculateOrder(t *testing.T) {
	type behavior func(*MockCalculationRepository, context.Context, repository.AddCalculationResult, error)

	type fields struct {
		repRules repository.CalculationRules
		mxRules  *sync.Mutex
		rules    map[int16]rule
		manager  *EventManager
	}
	type args struct {
		ctx   context.Context
		event Event
		behavior
		dto repository.AddCalculationResult
		err error
	}

	baseFields := fields{
		repRules: &MockCalculationRules{},
		mxRules:  &sync.Mutex{},
		rules: map[int16]rule{
			1: {
				calculationType: service.CalculationTypeFixed,
				value:           10,
				exp:             regexp.MustCompile("(?i)test"),
			},
			2: {
				calculationType: service.CalculationTypePercent,
				value:           10,
				exp:             regexp.MustCompile("(?i)test"),
			},
		},
		manager: NewEventManager(context.Background()),
	}

	type result struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		result
	}{
		{
			name:   "calculate order: normal situation",
			fields: baseFields,
			args: args{
				ctx: context.Background(),
				event: Event{
					Type: NewOrder,
					Data: service.RegisterCalculationRequest{
						OrderNumber: "123445",
						Products: []service.ProductRow{
							{
								Name:  "test-product",
								Price: 100,
							},
							{
								Name:  "product1",
								Price: 100,
							},
							{
								Name:  "product2",
								Price: 100,
							},
						},
					},
				},
				dto: repository.AddCalculationResult{
					OrderNumber: "123445",
					Value:       10,
					Status:      repository.Processed,
				},
				behavior: func(mcr *MockCalculationRepository, ctx context.Context, dto repository.AddCalculationResult, err error) {
					mcr.EXPECT().SaveCalculationResult(ctx, dto).Return(err)
				},
			},
		},
		{
			name:   "calculate order: wrong input",
			fields: baseFields,
			args: args{
				ctx: context.Background(),
				event: Event{
					Type: NewOrder,
					Data: 13456,
				},
				dto:      repository.AddCalculationResult{},
				behavior: func(*MockCalculationRepository, context.Context, repository.AddCalculationResult, error) {},
			},
			result: result{
				msg: "was raise event with wrong data",
			},
		},
		{
			name:   "calculate order: saving error",
			fields: baseFields,
			args: args{
				ctx: context.Background(),
				event: Event{
					Type: NewOrder,
					Data: service.RegisterCalculationRequest{
						OrderNumber: "123445",
						Products: []service.ProductRow{
							{
								Name:  "test-product",
								Price: 100,
							},
							{
								Name:  "product1",
								Price: 100,
							},
							{
								Name:  "product2",
								Price: 100,
							},
						},
					},
				},
				dto: repository.AddCalculationResult{
					OrderNumber: "123445",
					Value:       10,
					Status:      repository.Processed,
				},
				behavior: func(mcr *MockCalculationRepository, ctx context.Context, dto repository.AddCalculationResult, err error) {
					mcr.EXPECT().SaveCalculationResult(ctx, dto).Return(err)
				},
				err: errors.New("saving error"),
			},
			result: result{
				msg: "saving calculation result",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rep := NewMockCalculationRepository(ctrl)

			tt.args.behavior(rep, tt.args.ctx, tt.args.dto, tt.args.err)

			c := CalculationService{
				repCalc:  rep,
				repRules: tt.fields.repRules,
				mxRules:  tt.fields.mxRules,
				rules:    tt.fields.rules,
				manager:  tt.fields.manager,
			}
			//we have only log message on error because will check messages which writer stored
			wrt := NewMockLoggerWriter()
			logger.Init(wrt, logger.DebugLevel)

			c.calculateOrder(tt.args.ctx, tt.args.event)

			var lenMsg int
			if tt.result.msg != "" {
				lenMsg += 1
			}

			assert.Len(t, wrt.msgs, lenMsg, "writer had to get message")

			if len(wrt.msgs) > 0 {
				assert.True(t, strings.Contains(wrt.msgs[0], tt.result.msg), "logger got wrong message")
			}

		})
	}
}

func TestCalculationService_Calculation(t *testing.T) {
	type behavior func(*MockCalculationRepository, context.Context, repository.CalculationFilter, []repository.CalculationInfo, error)

	type fields struct {
		repRules repository.CalculationRules
		mxRules  *sync.Mutex
		rules    map[int16]rule
		manager  *EventManager
	}

	baseFields := fields{
		repRules: NewMockCalculationRules(gomock.NewController(t)),
		mxRules:  &sync.Mutex{},
		rules:    make(map[int16]rule),
		manager:  NewEventManager(context.Background()),
	}
	type args struct {
		behavior
		ctx       context.Context
		dto       service.CalculationFilterRequest
		dtoResult []repository.CalculationInfo
		err       error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []service.CalculationInfo
		wantErr bool
	}{
		{
			name: "getting calculation, there are row",
			args: args{
				behavior: func(mcr *MockCalculationRepository, ctx context.Context, dto repository.CalculationFilter, dtoResult []repository.CalculationInfo, err error) {
					mcr.EXPECT().Calculations(ctx, dto).Return(dtoResult, err)

				},
				ctx: context.Background(),
				dto: service.CalculationFilterRequest{
					OrderNumber: "123456",
				},
				dtoResult: []repository.CalculationInfo{
					{
						OrderNumber: "123456",
						Value:       10,
						Status:      repository.Processed,
					},
				},
			},
			fields: baseFields,
			want: []service.CalculationInfo{
				{
					OrderNumber: "123456",
					Accrual:     10,
					Status:      "PROCESSED",
				},
			},
		},
		{
			name: "getting calculation, with empty filter",
			args: args{
				behavior: func(mcr *MockCalculationRepository, ctx context.Context, dto repository.CalculationFilter, dtoResult []repository.CalculationInfo, err error) {
					mcr.EXPECT().Calculations(ctx, dto).Return(dtoResult, err)

				},
				ctx: context.Background(),
				dto: service.CalculationFilterRequest{},
				dtoResult: []repository.CalculationInfo{
					{
						OrderNumber: "123456",
						Value:       10,
						Status:      repository.Processed,
					},
					{
						OrderNumber: "54321",
						Value:       10,
						Status:      repository.Processed,
					},
					{
						OrderNumber: "097653",
						Value:       10,
						Status:      repository.Processing,
					},
					{
						OrderNumber: "34365",
						Value:       10,
						Status:      repository.Invalid,
					},
					{
						OrderNumber: "67453903",
						Value:       10,
						Status:      7,
					},
				},
			},
			fields: baseFields,
			want: []service.CalculationInfo{
				{
					OrderNumber: "123456",
					Accrual:     10,
					Status:      "PROCESSED",
				},
				{
					OrderNumber: "54321",
					Accrual:     10,
					Status:      "PROCESSED",
				},
				{
					OrderNumber: "097653",
					Accrual:     10,
					Status:      "PROCESSING",
				},
				{
					OrderNumber: "34365",
					Accrual:     10,
					Status:      "INVALID",
				},
				{
					OrderNumber: "67453903",
					Accrual:     10,
					Status:      "",
				},
			},
		},
		{
			name: "getting calculation, there not rows",
			args: args{
				behavior: func(mcr *MockCalculationRepository, ctx context.Context, dto repository.CalculationFilter, dtoResult []repository.CalculationInfo, err error) {
					mcr.EXPECT().Calculations(ctx, dto).Return(dtoResult, err)

				},
				ctx: context.Background(),
				dto: service.CalculationFilterRequest{
					OrderNumber: "123456",
				},
				dtoResult: []repository.CalculationInfo{},
			},
			fields:  baseFields,
			want:    []service.CalculationInfo{},
			wantErr: true,
		},
		{
			name: "getting calculation, error on repository",
			args: args{
				behavior: func(mcr *MockCalculationRepository, ctx context.Context, dto repository.CalculationFilter, dtoResult []repository.CalculationInfo, err error) {
					mcr.EXPECT().Calculations(ctx, dto).Return(dtoResult, err)

				},
				ctx: context.Background(),
				dto: service.CalculationFilterRequest{
					OrderNumber: "123456",
				},
				dtoResult: []repository.CalculationInfo{},
				err:       errors.New("error on repository"),
			},
			fields:  baseFields,
			want:    []service.CalculationInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rep := NewMockCalculationRepository(ctrl)

			tt.args.behavior(rep, tt.args.ctx, repository.CalculationFilter{OrderNumber: tt.args.dto.OrderNumber}, tt.args.dtoResult, tt.args.err)

			c := CalculationService{
				repCalc:  rep,
				repRules: tt.fields.repRules,
				mxRules:  tt.fields.mxRules,
				rules:    tt.fields.rules,
				manager:  tt.fields.manager,
			}
			got, err := c.Calculation(tt.args.ctx, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculationService.Calculation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculationService.Calculation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculationService_Register(t *testing.T) {
	type behavior func(*MockCalculationRepository, context.Context, []repository.AddingCalculation, error)
	type behaviorRules func(*MockCalculationRules, context.Context, repository.RuleFilter)

	type fields struct {
		repRules repository.CalculationRules
		manager  *EventManager
	}

	baseFields := fields{
		repRules: &MockCalculationRules{},
		manager:  NewEventManager(context.Background()),
	}

	type args struct {
		behavior
		ctx         context.Context
		dto         service.RegisterCalculationRequest
		dtoExpected []repository.AddingCalculation
		err         error
		behaviorRules
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "success registration",
			fields: baseFields,
			args: args{
				behavior: func(mcr *MockCalculationRepository, ctx context.Context, ac []repository.AddingCalculation, err error) {
					mcr.EXPECT().AddCalculationToQueue(ctx, ac).Return(err)
				},
				behaviorRules: func(mcr *MockCalculationRules, ctx context.Context, rf repository.RuleFilter) {
					mcr.EXPECT().Rules(ctx, rf).Return([]repository.RuleInfo{}, nil)
				},
				ctx: context.Background(),
				dto: service.RegisterCalculationRequest{
					OrderNumber: "123456",
					Products: []service.ProductRow{
						{
							Name:  "test",
							Price: 100,
						},
					},
				},
				dtoExpected: []repository.AddingCalculation{
					{
						OrderNumber: "123456",
						ProductName: "test",
						Price:       100,
					},
				},
				err: nil,
			},
			wantErr: false,
		},
		{
			name:   "failed registration",
			fields: baseFields,
			args: args{
				behavior: func(mcr *MockCalculationRepository, ctx context.Context, ac []repository.AddingCalculation, err error) {
					mcr.EXPECT().AddCalculationToQueue(ctx, ac).Return(err)
				},
				behaviorRules: func(mcr *MockCalculationRules, ctx context.Context, rf repository.RuleFilter) {
					mcr.EXPECT().Rules(ctx, rf).Return([]repository.RuleInfo{}, nil)
				},
				ctx: context.Background(),
				dto: service.RegisterCalculationRequest{
					OrderNumber: "123456",
					Products: []service.ProductRow{
						{
							Name:  "test",
							Price: 100,
						},
					},
				},
				dtoExpected: []repository.AddingCalculation{
					{
						OrderNumber: "123456",
						ProductName: "test",
						Price:       100,
					},
				},
				err: errors.New("error on repository"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rep := NewMockCalculationRepository(ctrl)

			tt.args.behavior(rep, tt.args.ctx, tt.args.dtoExpected, tt.args.err)

			c := CalculationService{
				repCalc:  rep,
				repRules: &MockCalculationRules{},
				mxRules:  &sync.Mutex{},
				rules:    make(map[int16]rule),
				manager:  tt.fields.manager,
			}
			if err := c.Register(tt.args.ctx, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("CalculationService.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
