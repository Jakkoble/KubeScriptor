package internal

import "github.com/Jakkoble/HexaTask/src/Runner/pb"

type GrpcLogStreamer struct {
	stream pb.RunnerService_StreamLogsClient
	jobID  string
}

func NewGrpcLogStreamer(stream pb.RunnerService_StreamLogsClient, jobID string) *GrpcLogStreamer {
	return &GrpcLogStreamer{stream: stream, jobID: jobID}
}

func (s *GrpcLogStreamer) SendLog(msg string, isStderr bool) error {
	return s.stream.Send(&pb.LogMessage{
		JobId:   s.jobID,
		Log:     msg,
		IsError: isStderr,
	})
}

func (s *GrpcLogStreamer) Close(exitCode int) error {
	err := s.stream.Send(&pb.LogMessage{
		JobId:    s.jobID,
		Log:      "",
		IsError:  false,
		IsFinal:  true,
		ExitCode: int32(exitCode),
	})

	if err != nil {
		return err
	}

	_, err = s.stream.CloseAndRecv()
	return err
}
