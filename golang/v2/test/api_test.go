// https://go.dev/doc/tutorial/add-a-test
// https://medium.com/nerd-for-tech/setup-and-teardown-unit-test-in-go-bd6fa1b785cd
package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"araali.proto/araali_api_service"

	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/araalinetworks/araali/third_party/api/golang/v2/api"
)

func setup(t *testing.T) {
	api.SetBackend("nightly.aws.araalinetworks.com")
	api.SetToken(os.Getenv("ARAALI_API_TOKEN"))
}

var verbose = 0

// TestAlerts calls ListAlerts since the beginning of time,
// checking for consistency with last baseline
func TestAlerts(t *testing.T) {
	setup(t)

	tenantID := "meta-tap"
	filter := araali_api_service.AlertFilter{
		Time: &araali_api_service.TimeSlice{
			StartTime: timestamppb.New(time.Date(1980, time.November, 0, 0, 0, 0, 0, time.UTC)),
			EndTime:   timestamppb.New(time.Now()),
		},
		ListAllAlerts:        false,
		OpenAlerts:           true,
		ClosedAlerts:         false,
		PerimeterIngress:     true,
		PerimeterEgress:      true,
		HomeNonAraaliIngress: true,
		HomeNonAraaliEgress:  true,
		AraaliToAraali:       true,
	}
	resp, err := api.ListAlerts(tenantID, &filter, 10, "")
	if verbose > 0 {
		fmt.Printf("\nR: %+v/%v\n", len(resp.Alerts), err)
	}
	if len(resp.Alerts) != 10 {
		t.Fatalf("ListAlerts() = %v, want 10", len(resp.Alerts))
	}
}

//Testing ListAlerts with specific Zone passed. Return only those zone specific alerts
func TestAlertsZoneFilter(t *testing.T) {
	setup(t)

	zone := "nightlycommon"
	tenantID := "meta-tap"
	filter := araali_api_service.AlertFilter{
		Time: &araali_api_service.TimeSlice{
			StartTime: timestamppb.New(time.Date(1980, time.November, 0, 0, 0, 0, 0, time.UTC)),
			EndTime:   timestamppb.New(time.Now()),
		},
		ListAllAlerts:        false,
		OpenAlerts:           true,
		ClosedAlerts:         false,
		PerimeterIngress:     true,
		PerimeterEgress:      true,
		HomeNonAraaliIngress: true,
		HomeNonAraaliEgress:  true,
		AraaliToAraali:       true,
		Zone:                 zone,
	}
	resp, err := api.ListAlerts(tenantID, &filter, 10, "")
	totalAlerts := len(resp.Alerts)

	if verbose > 0 {
		fmt.Printf("\nR: %+v/%v\n", totalAlerts, err)
	}

	counter := 0
	for _, s := range resp.Alerts {
		_, araaliEndpoint := s.Client.Info.(*araali_api_service.EndPoint_Araali)
		if araaliEndpoint {
			if s.Client.GetAraali().Zone == zone {
				counter++
			}
		}
	}

	//Validate whether all the alerts are corresponding to passed zone value
	if counter != totalAlerts {
		t.Fatalf("%v count in the response = %v, want %v", zone, counter, totalAlerts)
	}

}

//Testing ListAlerts with specific Zone passed which returns ZERO count
func TestAlertsZoneFilterNoResults(t *testing.T) {
	setup(t)

	zone := "intern-nightly"
	tenantID := "meta-tap"
	filter := araali_api_service.AlertFilter{
		Time: &araali_api_service.TimeSlice{
			StartTime: timestamppb.New(time.Date(2022, time.July, 25, 3, 8, 0, 0, time.UTC)),
			EndTime:   timestamppb.New(time.Date(2022, time.July, 25, 3, 9, 0, 0, time.UTC)),
		},
		ListAllAlerts:        false,
		OpenAlerts:           true,
		ClosedAlerts:         false,
		PerimeterIngress:     true,
		PerimeterEgress:      true,
		HomeNonAraaliIngress: true,
		HomeNonAraaliEgress:  true,
		AraaliToAraali:       true,
		Zone:                 zone,
	}
	resp, err := api.ListAlerts(tenantID, &filter, 10, "")
	if verbose > 0 {
		fmt.Printf("\nR: %+v/%v\n", len(resp.Alerts), err)
	}
	if len(resp.Alerts) != 0 {
		t.Fatalf("ListAlerts() = %v, want 0", len(resp.Alerts))
	}
}
