package transaction

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	transactionid "github.com/giantswarm/microkit/transaction/context/id"
)

func Test_Executer_NoTransactionIDGiven(t *testing.T) {
	config := DefaultExecuterConfig()
	newExecuter, err := NewExecuter(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	var replayExecuted int
	var trialExecuted int

	replay := func(context context.Context, v interface{}) error {
		replayExecuted++
		return nil
	}
	trial := func(context context.Context) (interface{}, error) {
		trialExecuted++
		return nil, nil
	}

	var ctx context.Context
	var executeConfig ExecuteConfig
	{
		ctx = context.Background()

		executeConfig = newExecuter.ExecuteConfig()
		executeConfig.Replay = replay
		executeConfig.Trial = trial
		executeConfig.TrialID = "test-trial-ID"
	}

	// The first execution of the transaction causes the trial to be executed
	// once. The replay function must not be executed at all.
	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialExecuted != 1 {
			t.Fatal("expected", 1, "got", trialExecuted)
		}
		if replayExecuted != 0 {
			t.Fatal("expected", 0, "got", replayExecuted)
		}
	}

	// There is no transaction ID provided, so the trial is executed again and the
	// replay function is still untouched.
	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialExecuted != 2 {
			t.Fatal("expected", 2, "got", trialExecuted)
		}
		if replayExecuted != 0 {
			t.Fatal("expected", 0, "got", replayExecuted)
		}
	}

	// There is no transaction ID provided, so the trial is executed again and the
	// replay function is still untouched.
	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialExecuted != 3 {
			t.Fatal("expected", 3, "got", trialExecuted)
		}
		if replayExecuted != 0 {
			t.Fatal("expected", 0, "got", replayExecuted)
		}
	}
}

func Test_Executer_TransactionIDGiven(t *testing.T) {
	config := DefaultExecuterConfig()
	newExecuter, err := NewExecuter(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	var replayExecuted int
	var trialExecuted int

	replay := func(context context.Context, v interface{}) error {
		replayExecuted++
		return nil
	}
	trial := func(context context.Context) (interface{}, error) {
		trialExecuted++
		return nil, nil
	}

	var ctx context.Context
	var executeConfig ExecuteConfig
	{
		ctx = context.Background()
		ctx = transactionid.NewContext(ctx, "test-transaction-id")

		executeConfig = newExecuter.ExecuteConfig()
		executeConfig.Replay = replay
		executeConfig.Trial = trial
		executeConfig.TrialID = "test-trial-ID"
	}

	// The first execution of the transaction causes the trial to be executed
	// once. The replay function must not be executed at all.
	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialExecuted != 1 {
			t.Fatal("expected", 1, "got", trialExecuted)
		}
		if replayExecuted != 0 {
			t.Fatal("expected", 0, "got", replayExecuted)
		}
	}

	// There is a transaction ID provided, so the trial is not executed again and
	// the replay function is executed the first time.
	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialExecuted != 1 {
			t.Fatal("expected", 1, "got", trialExecuted)
		}
		if replayExecuted != 1 {
			t.Fatal("expected", 1, "got", replayExecuted)
		}
	}

	// There is a transaction ID provided, so the trial is still not executed
	// again and the replay function is executed the second time.
	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialExecuted != 1 {
			t.Fatal("expected", 1, "got", trialExecuted)
		}
		if replayExecuted != 2 {
			t.Fatal("expected", 2, "got", replayExecuted)
		}
	}
}

func Test_Executer_TransactionIDGiven_NoReplay(t *testing.T) {
	config := DefaultExecuterConfig()
	newExecuter, err := NewExecuter(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	var trialExecuted int

	trial := func(context context.Context) (interface{}, error) {
		trialExecuted++
		return nil, nil
	}

	var ctx context.Context
	var executeConfig ExecuteConfig
	{
		ctx = context.Background()
		ctx = transactionid.NewContext(ctx, "test-transaction-id")

		executeConfig = newExecuter.ExecuteConfig()
		executeConfig.Trial = trial
		executeConfig.TrialID = "test-trial-ID"
	}

	// The first execution of the transaction causes the trial to be executed
	// once. There is no replay function.
	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialExecuted != 1 {
			t.Fatal("expected", 1, "got", trialExecuted)
		}
	}

	// There is a transaction ID provided, so the trial is not executed again.
	// There is no replay function.
	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialExecuted != 1 {
			t.Fatal("expected", 1, "got", trialExecuted)
		}
	}

	// There is a transaction ID provided, so the trial is still not executed
	// again. There is no replay function.
	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialExecuted != 1 {
			t.Fatal("expected", 1, "got", trialExecuted)
		}
	}
}

func Test_Executer_TransactionResult_Byte(t *testing.T) {
	config := DefaultExecuterConfig()
	newExecuter, err := NewExecuter(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	var trialOutput []byte
	var replayInput []byte

	replay := func(context context.Context, v interface{}) error {
		replayInput = []byte(v.(string))
		return nil
	}
	trial := func(context context.Context) (interface{}, error) {
		trialOutput = []byte("hello world")
		return trialOutput, nil
	}

	var ctx context.Context
	var executeConfig ExecuteConfig
	{
		ctx = context.Background()
		ctx = transactionid.NewContext(ctx, "test-transaction-id")

		executeConfig = newExecuter.ExecuteConfig()
		executeConfig.Replay = replay
		executeConfig.Trial = trial
		executeConfig.TrialID = "test-trial-ID"
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if string(trialOutput) != "hello world" {
			t.Fatal("expected", "hello world", "got", string(trialOutput))
		}
		if string(replayInput) != "" {
			t.Fatal("expected", "", "got", string(replayInput))
		}
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if string(trialOutput) != "hello world" {
			t.Fatal("expected", "hello world", "got", string(trialOutput))
		}
		if string(replayInput) != "hello world" {
			t.Fatal("expected", "hello world", "got", string(replayInput))
		}
	}
}

func Test_Executer_TransactionResult_Float64(t *testing.T) {
	config := DefaultExecuterConfig()
	newExecuter, err := NewExecuter(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	var trialOutput float64
	var replayInput float64

	replay := func(context context.Context, v interface{}) error {
		var err error
		replayInput, err = strconv.ParseFloat(v.(string), 64)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		return nil
	}
	trial := func(context context.Context) (interface{}, error) {
		trialOutput = 4.3
		return trialOutput, nil
	}

	var ctx context.Context
	var executeConfig ExecuteConfig
	{
		ctx = context.Background()
		ctx = transactionid.NewContext(ctx, "test-transaction-id")

		executeConfig = newExecuter.ExecuteConfig()
		executeConfig.Replay = replay
		executeConfig.Trial = trial
		executeConfig.TrialID = "test-trial-ID"
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialOutput != 4.3 {
			t.Fatal("expected", 4.3, "got", trialOutput)
		}
		if replayInput != 0 {
			t.Fatal("expected", 0, "got", replayInput)
		}
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialOutput != 4.3 {
			t.Fatal("expected", 4.3, "got", trialOutput)
		}
		if replayInput != 4.3 {
			t.Fatal("expected", 4.3, "got", replayInput)
		}
	}
}

func Test_Executer_TransactionResult_Nil(t *testing.T) {
	config := DefaultExecuterConfig()
	newExecuter, err := NewExecuter(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	var replayInput interface{}

	replay := func(context context.Context, v interface{}) error {
		replayInput = v
		return nil
	}
	trial := func(context context.Context) (interface{}, error) {
		return nil, nil
	}

	var ctx context.Context
	var executeConfig ExecuteConfig
	{
		ctx = context.Background()
		ctx = transactionid.NewContext(ctx, "test-transaction-id")

		executeConfig = newExecuter.ExecuteConfig()
		executeConfig.Replay = replay
		executeConfig.Trial = trial
		executeConfig.TrialID = "test-trial-ID"
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if replayInput != nil {
			t.Fatal("expected", nil, "got", replayInput)
		}
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if replayInput != nil {
			t.Fatal("expected", nil, "got", replayInput)
		}
	}
}

func Test_Executer_TransactionResult_String(t *testing.T) {
	config := DefaultExecuterConfig()
	newExecuter, err := NewExecuter(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	var trialOutput string
	var replayInput string

	replay := func(context context.Context, v interface{}) error {
		replayInput = v.(string)
		return nil
	}
	trial := func(context context.Context) (interface{}, error) {
		trialOutput = "hello world"
		return trialOutput, nil
	}

	var ctx context.Context
	var executeConfig ExecuteConfig
	{
		ctx = context.Background()
		ctx = transactionid.NewContext(ctx, "test-transaction-id")

		executeConfig = newExecuter.ExecuteConfig()
		executeConfig.Replay = replay
		executeConfig.Trial = trial
		executeConfig.TrialID = "test-trial-ID"
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialOutput != "hello world" {
			t.Fatal("expected", "hello world", "got", trialOutput)
		}
		if replayInput != "" {
			t.Fatal("expected", "", "got", replayInput)
		}
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialOutput != "hello world" {
			t.Fatal("expected", "hello world", "got", trialOutput)
		}
		if replayInput != "hello world" {
			t.Fatal("expected", "hello world", "got", replayInput)
		}
	}
}

func Test_Executer_TransactionResult_EmptyString(t *testing.T) {
	config := DefaultExecuterConfig()
	newExecuter, err := NewExecuter(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	var trialOutput string
	var replayInput string

	replay := func(context context.Context, v interface{}) error {
		replayInput = v.(string)
		return nil
	}
	trial := func(context context.Context) (interface{}, error) {
		trialOutput = ""
		return trialOutput, nil
	}

	var ctx context.Context
	var executeConfig ExecuteConfig
	{
		ctx = context.Background()
		ctx = transactionid.NewContext(ctx, "test-transaction-id")

		executeConfig = newExecuter.ExecuteConfig()
		executeConfig.Replay = replay
		executeConfig.Trial = trial
		executeConfig.TrialID = "test-trial-ID"
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialOutput != "" {
			t.Fatal("expected", "", "got", trialOutput)
		}
		if replayInput != "" {
			t.Fatal("expected", "", "got", replayInput)
		}
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if trialOutput != "" {
			t.Fatal("expected", "", "got", trialOutput)
		}
		if replayInput != "" {
			t.Fatal("expected", "", "got", replayInput)
		}
	}
}

func Test_Executer_TransactionResult_Struct(t *testing.T) {
	config := DefaultExecuterConfig()
	newExecuter, err := NewExecuter(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	type testTrialOutput struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	var trialOutput testTrialOutput
	var replayInput string

	replay := func(context context.Context, v interface{}) error {
		replayInput = v.(string)
		return nil
	}
	trial := func(context context.Context) (interface{}, error) {
		trialOutput = testTrialOutput{Foo: "foo-val", Bar: 43}
		return trialOutput, nil
	}

	var ctx context.Context
	var executeConfig ExecuteConfig
	{
		ctx = context.Background()
		ctx = transactionid.NewContext(ctx, "test-transaction-id")

		executeConfig = newExecuter.ExecuteConfig()
		executeConfig.Replay = replay
		executeConfig.Trial = trial
		executeConfig.TrialID = "test-trial-ID"
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if !reflect.DeepEqual(trialOutput, testTrialOutput{Foo: "foo-val", Bar: 43}) {
			t.Fatal("expected", testTrialOutput{Foo: "foo-val", Bar: 43}, "got", trialOutput)
		}
		if replayInput != "" {
			t.Fatal("expected", "", "got", replayInput)
		}
	}

	{
		err := newExecuter.Execute(ctx, executeConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if !reflect.DeepEqual(trialOutput, testTrialOutput{Foo: "foo-val", Bar: 43}) {
			t.Fatal("expected", testTrialOutput{Foo: "foo-val", Bar: 43}, "got", trialOutput)
		}
		if replayInput != "{\"foo\":\"foo-val\",\"bar\":43}" {
			t.Fatal("expected", "{\"foo\":\"foo-val\",\"bar\":43}", "got", replayInput)
		}
	}
}
