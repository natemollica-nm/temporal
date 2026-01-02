package iplocate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.temporal.io/sdk/activity"
)

type HTTPGetter interface {
	Get(url string) (*http.Response, error)
}

type IPActivities struct {
	HTTPClient HTTPGetter
}

type IPInfo struct {
	Status      string  `json:"status"`
	City        string  `json:"city"`
	RegionName  string  `json:"regionName"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
	Query       string  `json:"query"`
	Timezone    string  `json:"timezone"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

// GetIP fetches the public IP address.
func (i *IPActivities) GetIP(ctx context.Context, scheduledTime int64) (string, error) {
	logger := activity.GetLogger(ctx)

	var err error
	metricsHandler := activity.GetMetricsHandler(ctx).WithTags(map[string]string{
		"stage": "GetIP",
	})
	metricsHandler = recordActivityStart(metricsHandler, "activity.get_ip", scheduledTime)
	startTime := time.Now()
	defer func() {
		recordActivityEnd(metricsHandler, startTime, err)
		logger.Info("GetIP activity completed")
	}()

	logger.Info("Getting IP address")
	resp, err := i.HTTPClient.Get("https://icanhazip.com")
	if err != nil {
		logger.Error("failed to obtain IP address err=", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read IP address response body err=", err)
		return "", err
	}

	ip := strings.TrimSpace(string(body))
	logger.Info("Got IP address", "ip=", ip)
	return ip, nil
}

func (i *IPActivities) retrieveIPAddressInfo(ip string) (IPInfo, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := i.HTTPClient.Get(url)
	if err != nil {
		return IPInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return IPInfo{}, err
	}

	var data IPInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		return IPInfo{}, err
	}
	return data, nil
}

// GetLocationInfo uses the IP address to fetch location information.
func (i *IPActivities) GetLocationInfo(ctx context.Context, ip string, scheduledTime int64) (string, error) {
	logger := activity.GetLogger(ctx)

	var err error
	metricsHandler := activity.GetMetricsHandler(ctx).WithTags(map[string]string{
		"stage": "GetLocationInfo",
	})
	metricsHandler = recordActivityStart(metricsHandler, "activity.get_location_info", scheduledTime)
	startTime := time.Now()
	defer func() {
		recordActivityEnd(metricsHandler, startTime, err)
		logger.Info("GetLocationInfo activity completed")
	}()

	info, err := i.retrieveIPAddressInfo(ip)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s, %s, %s", info.City, info.RegionName, info.Country), nil
}

func (i *IPActivities) GetInternetServiceProvider(ctx context.Context, ip string, scheduledTime int64) (string, error) {
	logger := activity.GetLogger(ctx)

	var err error
	metricsHandler := activity.GetMetricsHandler(ctx).WithTags(map[string]string{
		"stage": "GetInternetServiceProvider",
	})
	metricsHandler = recordActivityStart(metricsHandler, "activity.get_isp", scheduledTime)
	startTime := time.Now()
	defer func() {
		recordActivityEnd(metricsHandler, startTime, err)
		logger.Info("GetInternetServiceProvider activity completed")
	}()

	info, err := i.retrieveIPAddressInfo(ip)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", info.ISP), nil
}
