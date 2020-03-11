package opsmngr

import (
	"context"
	"net/http"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	allClustersBasePath = "clusters"
)

// ClustersService is an interface for interfacing with Clusters in MongoDB Ops Manager APIs
type AllClustersService interface {
	List(ctx context.Context) ([]AllClustersProject, *atlas.Response, error)
}

type AllClustersServiceOp struct {
	Client atlas.RequestDoer
}

type AllClustersProject struct {
	GroupName string               `json:"groupName,omitempty"`
	OrgName   string               `json:"orgName,omitempty"`
	PlanType  string               `json:"planType,omitempty"`
	GroupID   string               `json:"groupId,omitempty"`
	OrgID     string               `json:"orgId,omitempty"`
	Tags      []string             `json:"tags,omitempty"`
	Clusters  []AllClustersCluster `json:"clusters,omitempty"`
}

// AllClustersCluster represents MongoDB cluster.
type AllClustersCluster struct {
	ClusterID     string   `json:"clusterId,omitempty"`
	Name          string   `json:"name,omitempty"`
	Type          string   `json:"type,omitempty"`
	Availability  string   `json:"availability,omitempty"`
	Versions      []string `json:"versions,omitempty"`
	BackupEnabled *bool    `json:"backupEnabled,omitempty"`
	AuthEnabled   *bool    `json:"authEnabled,omitempty"`
	SSLEnabled    *bool    `json:"sslEnabled,omitempty"`
	AlertCount    int64    `json:"alertCount,omitempty"`
	DataSizeBytes int64    `json:"dataSizeBytes,omitempty"`
	NodeCount     int64    `json:"nodeCount,omitempty"`
}

type allClustersResponse struct {
	Content []AllClustersProject `json:"content,omitempty"`
	Status  *int64               `json:"status,omitempty"`
}

//List all clusters in the project.
func (s *AllClustersServiceOp) List(ctx context.Context) ([]AllClustersProject, *atlas.Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodGet, allClustersBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(allClustersResponse)
	resp, err := s.Client.Do(ctx, req, root)

	return root.Content, resp, err
}
