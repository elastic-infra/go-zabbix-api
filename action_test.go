package zabbix_test

import (
	"testing"

	zapi "github.com/claranet/go-zabbix-api"
)

func testCreateAction(group *zapi.HostGroup, t *testing.T) *zapi.Action {
	actions := zapi.Actions{{
		Name:        "Register Linux servers",
		EventSource: zapi.AutoRegistrationEvent,
		Status:      zapi.Enabled,
		Filter: zapi.ActionFilter{
			EvaluationType: zapi.Custom,
			Formula:        "A or B",
			Conditions: zapi.ActionFilterConditions{
				{
					ConditionType: zapi.HostNameCondition,
					Operator:      zapi.Contains,
					Value:         "SRV",
					FormulaID:     "B",
				},
				{
					ConditionType: zapi.HostMetadataCondition,
					Operator:      zapi.Contains,
					Value:         "CentOS",
					FormulaID:     "A",
				},
			},
		},
		Operations: zapi.ActionOperations{
			{
				OperationType:  zapi.AddToHostGroup,
				EvaluationType: zapi.AndOr,
				HostGroups: zapi.ActionOperationHostGroups{
					{
						GroupID: group.GroupID,
					},
				},
			},
			{
				OperationType:  zapi.SetHostInventoryMode,
				EvaluationType: zapi.AndOr,
				Inventory: &zapi.ActionOperationInventory{
					InventoryMode: "1",
				},
			},
		},
	}}
	err := testGetAPI(t).ActionsCreate(actions)
	if err != nil {
		t.Fatal(err)
	}
	return &actions[0]
}

func testDeleteAction(action *zapi.Action, t *testing.T) {
	err := testGetAPI(t).ActionsDelete(zapi.Actions{*action})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAction(t *testing.T) {
	api := testGetAPI(t)

	hostGroup := testCreateHostGroup(t)
	defer testDeleteHostGroup(hostGroup, t)

	action := testCreateAction(hostGroup, t)

	getByIdAction, err := api.ActionGetByID(action.ActionID)
	if err != nil {
		t.Error(err)
	}

	getByIdAction.Name = "Register CentOS servers"
	// NOTE: EventSource can't be updated
	getByIdAction.EventSource = ""
	// NOTE: pause_suppressed set only TriggerEvent
	getByIdAction.PauseSuppressed = nil
	getByIdAction.Period = ""
	getByIdAction.Operations = nil
	err = api.ActionsUpdate(zapi.Actions{*getByIdAction})
	if err != nil {
		t.Error(err)
	}

	getAction, err := api.ActionsGet(zapi.Params{"actionids": action.ActionID})
	if err != nil {
		t.Error(err)
	}
	if len(getAction) != 1 {
		t.Errorf("Expecting one result and got : %d", len(getAction))
	}
	if getAction[0].Name != getByIdAction.Name {
		t.Errorf("Action name is %q and should be %q", getAction[0].Name, getByIdAction.Name)
	}

	testDeleteAction(action, t)
}

func TestActionWithRecoveryAndUpdateOperations(t *testing.T) {
	api := testGetAPI(t)

	hostGroup := testCreateHostGroup(t)
	defer testDeleteHostGroup(hostGroup, t)

	// Create action with recovery and update operations
	actions := zapi.Actions{{
		Name:        "Test Action with Recovery Operations",
		EventSource: zapi.TriggerEvent,
		Status:      zapi.Enabled,
		Filter: zapi.ActionFilter{
			EvaluationType: zapi.And,
			Conditions: zapi.ActionFilterConditions{
				{
					ConditionType: zapi.HostGroupCondition,
					Operator:      zapi.Equals,
					Value:         hostGroup.GroupID,
				},
			},
		},
		Operations: zapi.ActionOperations{
			{
				OperationType: zapi.SendMessage,
				Message: &zapi.ActionOperationMessage{
					DefaultMessage: "1",
					MediaTypeID:    "0",
					Subject:        "Problem",
					Message:        "Test message",
				},
			},
		},
		RecoveryOperations: zapi.ActionRecoveryOperations{
			{
				OperationType: zapi.NotifyRecoveryAllInvolved,
			},
		},
		UpdateOperations: zapi.ActionUpdateOperations{
			{
				OperationType: zapi.NotifyUpdateAllInvolved,
			},
		},
	}}

	err := api.ActionsCreate(actions)
	if err != nil {
		t.Fatal(err)
	}
	action := &actions[0]
	defer testDeleteAction(action, t)

	// Test reading back the action
	retrievedAction, err := api.ActionGetByID(action.ActionID)
	if err != nil {
		t.Fatal(err)
	}

	// Verify recovery operations
	if len(retrievedAction.RecoveryOperations) != 1 {
		t.Errorf("Expected 1 recovery operation, got %d", len(retrievedAction.RecoveryOperations))
	}
	if retrievedAction.RecoveryOperations[0].OperationType != zapi.NotifyRecoveryAllInvolved {
		t.Errorf("Expected NotifyRecoveryAllInvolved, got %d", retrievedAction.RecoveryOperations[0].OperationType)
	}

	// Verify update operations
	if len(retrievedAction.UpdateOperations) != 1 {
		t.Errorf("Expected 1 update operation, got %d", len(retrievedAction.UpdateOperations))
	}
	if retrievedAction.UpdateOperations[0].OperationType != zapi.NotifyUpdateAllInvolved {
		t.Errorf("Expected NotifyUpdateAllInvolved, got %d", retrievedAction.UpdateOperations[0].OperationType)
	}
}
