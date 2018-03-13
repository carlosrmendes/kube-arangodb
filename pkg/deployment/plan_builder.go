//
// DISCLAIMER
//
// Copyright 2018 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package deployment

import (
	api "github.com/arangodb/k8s-operator/pkg/apis/deployment/v1alpha"
	"github.com/rs/zerolog"
)

// createPlan considers the current specification & status of the deployment creates a plan to
// get the status in line with the specification.
// If a plan already exists, nothing is done.
func (d *Deployment) createPlan() error {
	// Create plan
	newPlan, changed := createPlan(d.deps.Log, d.status.Plan, d.apiObject.Spec, d.status)

	// If not change, we're done
	if !changed {
		return nil
	}

	// Save plan
	if len(newPlan) == 0 {
		// Nothing to do
		return nil
	}
	d.status.Plan = newPlan
	if err := d.updateCRStatus(); err != nil {
		return maskAny(err)
	}
	return nil
}

// createPlan considers the given specification & status and creates a plan to get the status in line with the specification.
// If a plan already exists, the given plan is returned with false.
// Otherwise the new plan is returned with a boolean true.
func createPlan(log zerolog.Logger, currentPlan api.Plan, spec api.DeploymentSpec, status api.DeploymentStatus) (api.Plan, bool) {
	if len(currentPlan) > 0 {
		// Plan already exists, complete that first
		return currentPlan, false
	}

	// Check for various scenario's
	var plan api.Plan

	// Check for scale up/down
	switch spec.Mode {
	case api.DeploymentModeSingle:
		// Never scale down
	case api.DeploymentModeResilientSingle:
		// Only scale singles
		plan = append(plan, createScalePlan(log, status.Members.Single, api.ServerGroupSingle, spec.Single.Count)...)
	case api.DeploymentModeCluster:
		// Scale dbservers, coordinators, syncmasters & syncworkers
		plan = append(plan, createScalePlan(log, status.Members.DBServers, api.ServerGroupDBServers, spec.DBServers.Count)...)
		plan = append(plan, createScalePlan(log, status.Members.Coordinators, api.ServerGroupCoordinators, spec.Coordinators.Count)...)
		plan = append(plan, createScalePlan(log, status.Members.SyncMasters, api.ServerGroupSyncMasters, spec.SyncMasters.Count)...)
		plan = append(plan, createScalePlan(log, status.Members.SyncWorkers, api.ServerGroupSyncWorkers, spec.SyncWorkers.Count)...)
	}

	// Return plan
	return plan, true
}

// createScalePlan creates a scaling plan for a single server group
func createScalePlan(log zerolog.Logger, members api.MemberStatusList, group api.ServerGroup, count int) api.Plan {
	var plan api.Plan
	if len(members) < count {
		// Scale up
		toAdd := count - len(members)
		for i := 0; i < toAdd; i++ {
			plan = append(plan, api.Action{Type: api.ActionTypeAddMember, Group: group})
		}
		log.Debug().
			Int("delta", toAdd).
			Str("role", group.AsRole()).
			Msg("Creating scale-up plan")
	} else if len(members) > count {
		// Note, we scale down 1 member as a time
		if m, err := members.SelectMemberToRemove(); err == nil {
			if group == api.ServerGroupDBServers {
				plan = append(plan,
					api.Action{Type: api.ActionTypeCleanOutMember, Group: group, MemberID: m.ID},
				)
			}
			plan = append(plan,
				api.Action{Type: api.ActionTypeShutdownMember, Group: group, MemberID: m.ID},
				api.Action{Type: api.ActionTypeRemoveMember, Group: group, MemberID: m.ID},
			)
			log.Debug().
				Str("role", group.AsRole()).
				Msg("Creating scale-down plan")
		}
	}
	return plan
}
