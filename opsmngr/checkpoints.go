package opsmngr

import (
	"context"
	"fmt"
	"net/http"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	checkpoints = "groups/%s/clusters/%s/checkpoints"
)

// CheckpointsService is an interface for interfacing with the Checkpoint
// endpoints of the MongoDB Atlas API.
type CheckpointsService interface {
	List(context.Context, string, string, *atlas.ListOptions) (*atlas.Checkpoints, *atlas.Response, error)
	Get(context.Context, string, string, string) (*atlas.Checkpoint, *atlas.Response, error)
}

// CheckpointsServiceOp handles communication with the checkpoint related methods of the
// MongoDB Atlas API
type CheckpointsServiceOp struct {
	Client atlas.RequestDoer
}

var _ CheckpointsService = &CheckpointsServiceOp{}

func (s *CheckpointsServiceOp) List(ctx context.Context, groupID, clusterName string, listOptions *atlas.ListOptions) (*atlas.Checkpoints, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, atlas.NewArgError("clusterName", "must be set")
	}

	basePath := fmt.Sprintf(checkpoints, groupID, clusterName)
	path, err := setListOptions(basePath, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Checkpoints)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

func (s *CheckpointsServiceOp) Get(ctx context.Context, groupID, clusterName, checkpointID string) (*atlas.Checkpoint, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, atlas.NewArgError("clusterName", "must be set")
	}
	if checkpointID == "" {
		return nil, nil, atlas.NewArgError("checkpointID", "must be set")
	}

	basePath := fmt.Sprintf(checkpoints, groupID, clusterName)
	path := fmt.Sprintf("%s/%s", basePath, checkpointID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Checkpoint)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
