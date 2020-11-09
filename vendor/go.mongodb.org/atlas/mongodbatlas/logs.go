package mongodbatlas

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const logsPath = "groups/%s/clusters/%s/logs/%s"

// LogsService is an interface for interfacing with the Logs
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/logs/
type LogsService interface {
	Get(context.Context, string, string, string, io.Writer, *DateRangetOptions) (*Response, error)
}

// LogsServiceOp handles communication with the Logs related methods of the
// MongoDB Atlas API
type LogsServiceOp struct {
	Client GZipRequestDoer
}

// DateRangetOptions specifies an optional date range query.
type DateRangetOptions struct {
	StartDate string `url:"startDate,omitempty"`
	EndDate   string `url:"endDate,omitempty"`
}

// Get gets a compressed (.gz) log file that contains a range of log messages for a particular host.
// Note: The input parameter out (io.Writer) is not closed by this function.
// See more: https://docs.atlas.mongodb.com/reference/api/logs/
func (s *LogsServiceOp) Get(ctx context.Context, groupID, hostName, logName string, out io.Writer, opts *DateRangetOptions) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}

	if hostName == "" {
		return nil, NewArgError("hostName", "must be set")
	}

	if logName == "" {
		return nil, NewArgError("logName", "must be set")
	}

	basePath := fmt.Sprintf(logsPath, groupID, hostName, logName)

	// Add query params
	path, err := setListOptions(basePath, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.Client.NewGZipRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, out)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
