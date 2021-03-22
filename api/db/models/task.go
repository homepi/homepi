package models

type Task int

const (
	TaskDoor Task = iota
	TaskLamp
	TaskToggle
	TaskFan
)