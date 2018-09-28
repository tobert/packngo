package packngo

import (
	"encoding/json"
	"fmt"
)

const batchBasePath = "/batches"

// BatchService interface defines available batch methods
type BatchService interface {
	Get(batchID string, listOpt *ListOptions) (*Batch, *Response, error)
	List(ProjectID string, listOpt *ListOptions) ([]Batch, *Response, error)
	Create(projectID string, batches *BatchDeviceCreateRequest) ([]Batch, *Response, error)
	Delete(string, bool) (*Response, error)
}

// Batch type
type Batch struct {
	ID         string     `json:"id"`
	State      string     `json:"state,omitempty"`
	Quantity   int32      `json:"quantity,omitempty"`
	CreatedAt  *Timestamp `json:"created_at,omitempty"`
	Href       string     `json:"href,omitempty"`
	Project    Href       `json:"project,omitempty"`
	Facilities []Facility `json:"facilities,omitempty"`
	Devices    []Device   `json:"devices,omitempty"`
}

//BatchesList represents collection of batches
type batchesList struct {
	Batches []Batch `json:"batches,omitempty"`
}

// BatchDeviceCreateRequest type used to create batch of device instances
type BatchDeviceCreateRequest struct {
	Batches []BatchCreateDevice `json:"batches"`
}

// BatchCreateDevice type used to describe batch instances
type BatchCreateDevice struct {
	DeviceCreateRequest
	Quantity               int32 `json:"quantity"`
	FacilityDiversityLevel int32 `json:"facility_diversity_level,omitempty"`
}

// MarshalJSON custom marshalling to handle DeviceCreateRequest and additional fields
func (bcd *BatchCreateDevice) MarshalJSON() ([]byte, error) {
	dcr, err := bcd.DeviceCreateRequest.MarshalJSON()
	if err != nil {
		return nil, err
	}
	temp := make(map[string]interface{})
	err = json.Unmarshal(dcr, &temp)
	if err != nil {
		return nil, err
	}
	temp["quantity"] = bcd.Quantity
	if bcd.FacilityDiversityLevel != 0 {
		temp["facility_diversity_level"] = bcd.FacilityDiversityLevel
	}
	return json.Marshal(temp)
}

// BatchServiceOp implements BatchService
type BatchServiceOp struct {
	client *Client
}

// Get returns batch details
func (s *BatchServiceOp) Get(batchID string, listOpt *ListOptions) (*Batch, *Response, error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}
	path := fmt.Sprintf("%s/%s?%s", batchBasePath, batchID, params)
	batch := new(Batch)

	resp, err := s.client.DoRequest("GET", path, nil, batch)
	if err != nil {
		return nil, resp, err
	}

	return batch, resp, err
}

// List returns batches on a project
func (s *BatchServiceOp) List(projectID string, listOpt *ListOptions) (batches []Batch, resp *Response, err error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}
	path := fmt.Sprintf("%s/%s%s?%s", projectBasePath, projectID, batchBasePath, params)
	subset := new(batchesList)
	resp, err = s.client.DoRequest("GET", path, nil, subset)
	if err != nil {
		return nil, resp, err
	}

	batches = append(batches, subset.Batches...)
	return batches, resp, err
}

// Create function to create batch of device instances
func (s *BatchServiceOp) Create(projectID string, request *BatchDeviceCreateRequest) ([]Batch, *Response, error) {
	path := fmt.Sprintf("%s/%s/devices/batch", projectBasePath, projectID)

	batches := new(batchesList)
	resp, err := s.client.DoRequest("POST", path, request, batches)

	if err != nil {
		return nil, resp, err
	}

	return batches.Batches, resp, err
}

// Delete function to remove an instance batch
func (s *BatchServiceOp) Delete(id string, removeDevices bool) (*Response, error) {
	path := fmt.Sprintf("%s/%s?remove_associated_instances=%t", batchBasePath, id, removeDevices)

	return s.client.DoRequest("DELETE", path, nil, nil)
}
