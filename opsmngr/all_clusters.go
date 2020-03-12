package opsmngr

import (
	"context"
	"net/http"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	allClustersBasePath = "clusters"
)

// AllClustersService is an interface for interfacing with Clusters in MongoDB Ops Manager APIs
type AllClustersService interface {
	List(ctx context.Context) (*AllClustersProjects, *atlas.Response, error)
}

type AllClustersServiceOp struct {
	Client atlas.RequestDoer
}

type AllClustersProject struct {
	GroupName string               `json:"groupName"`
	OrgName   string               `json:"orgName"`
	PlanType  string               `json:"planType,omitempty"`
	GroupID   string               `json:"groupId"`
	OrgID     string               `json:"orgId"`
	Tags      []string             `json:"tags"`
	Clusters  []AllClustersCluster `json:"clusters"`
}

// AllClustersCluster represents MongoDB cluster.
type AllClustersCluster struct {
	ClusterID     string   `json:"clusterId"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Availability  string   `json:"availability"`
	Versions      []string `json:"versions"`
	BackupEnabled bool     `json:"backupEnabled"`
	AuthEnabled   bool     `json:"authEnabled"`
	SSLEnabled    bool     `json:"sslEnabled"`
	AlertCount    int64    `json:"alertCount"`
	DataSizeBytes int64    `json:"dataSizeBytes"`
	NodeCount     int64    `json:"nodeCount"`
}

type AllClustersProjects struct {
	Links      []*atlas.Link         `json:"links"`
	Results    []*AllClustersProject `json:"results"`
	TotalCount int                   `json:"totalCount"`
}

//List all clusters in the project.
func (s *AllClustersServiceOp) List(ctx context.Context) (*AllClustersProjects, *atlas.Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodGet, allClustersBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AllClustersProjects)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
