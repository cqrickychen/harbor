// Copyright 2018 The Harbor Authors. All rights reserved.

package opm

import "github.com/vmware/harbor/src/jobservice_v2/models"

//JobStatsManager defines the methods to handle stats of job.
type JobStatsManager interface {
	//Start to serve
	Start()

	//Shutdown the manager
	Shutdown()

	//Save the job stats
	//Async method to retry and improve performance
	//
	//jobStats models.JobStats : the job stats to be saved
	Save(jobStats models.JobStats)

	//Get the job stats from backend store
	//Sync method as we need the data
	//
	//Returns:
	//  models.JobStats : job stats data
	//  error           : error if meet any problems
	Retrieve(jobID string) (models.JobStats, error)

	//SetJobStatus will mark the status of job to the specified one
	//Async method to retry
	SetJobStatus(jobID string, status string)

	//Stop the job
	//
	//jobID string : ID of the being stopped job
	//
	//Returns:
	// error if meet any problems
	Stop(jobID string) error

	//Cancel the job
	//
	//jobID string : ID of the being cancelled job
	//
	//Returns:
	// error if meet any problems
	Cancel(jobID string) error

	//Retry the job
	//
	//jobID string : ID of the being retried job
	//
	//Returns:
	// error if meet any problems
	Retry(jobID string) error

	//CtlCommand checks if control command is fired for the specified job.
	//
	//jobID string : ID of the job
	//
	//Returns:
	//  the command if it was fired
	//  error if it was not fired yet to meet some other problems
	CtlCommand(jobID string) (string, error)

	//CheckIn message for the specified job like detailed progress info.
	//
	//jobID string   : ID of the job
	//message string : The message being checked in
	//
	CheckIn(jobID string, message string)

	//DieAt marks the failed jobs with the time they put into dead queue.
	//
	//jobID string   : ID of the job
	//message string : The message being checked in
	//
	DieAt(jobID string, dieAt int64)
}