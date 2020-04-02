package opsmngr

import (
	"context"
	"fmt"
	"net/http"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	systemMeasurements = "groups/%s/hosts/%s/measurements"
)

// SystemMeasurementsService is an interface for interfacing with the System Measurements
// endpoints of the MongoDB Atlas API.
type SystemMeasurementsService interface {
	List(context.Context, string, string, *atlas.ProcessMeasurementListOptions) (*atlas.ProcessMeasurements, *atlas.Response, error)
}

// ServiceMeasurementsServiceOp handles communication with the system measurements related methods of the
// MongoDB Atlas API
type ServiceMeasurementsServiceOp struct {
	Client atlas.RequestDoer
}

var _ SystemMeasurementsService = &ServiceMeasurementsServiceOp{}

// List lists system and process measurements on the CPU usage of the hosts that run MongoDB.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/measures/get-host-process-system-measurements/
func (s *ServiceMeasurementsServiceOp) List(ctx context.Context, projectID, host string, listOptions *atlas.ProcessMeasurementListOptions) (*atlas.ProcessMeasurements, *atlas.Response, error) {
	if projectID == "" {
		return nil, nil, atlas.NewArgError("projectID", "must be set")
	}
	if host == "" {
		return nil, nil, atlas.NewArgError("host", "must be set")
	}

	basePath := fmt.Sprintf(systemMeasurements, projectID, host)
	path, err := setQueryParams(basePath, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.ProcessMeasurements)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
