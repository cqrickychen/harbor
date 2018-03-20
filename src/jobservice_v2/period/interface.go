// Copyright 2018 The Harbor Authors. All rights reserved.

package period

import "github.com/vmware/harbor/src/jobservice_v2/models"

//Interface defines operations the periodic scheduler should have.
type Interface interface {
	//Schedule the specified cron job policy.
	//
	//jobName string           : The name of periodical job
	//params models.Parameters : The parameters required by the periodical job
	//cronSpec string          : The periodical settings with cron format
	//
	//Returns:
	//  The uuid of the cron job policy
	//  The latest next trigger time
	//  error if failed to schedule
	Schedule(jobName string, params models.Parameters, cronSpec string) (string, int64, error)

	//Unschedule the specified cron job policy.
	//
	//cronJobPolicyID string: The ID of cron job policy.
	//
	//Return:
	//  error if failed to unschedule
	UnSchedule(cronJobPolicyID string) error

	//Load data
	//
	//Return:
	//  error if failed to do
	Load() error

	//Clear all the cron job policies.
	//
	//Return:
	//  error if failed to do
	Clear() error

	//Start to serve
	Start() error
}