package types

type Communicator interface {
	SendStartProcess(jobID, fileName string) error
}
