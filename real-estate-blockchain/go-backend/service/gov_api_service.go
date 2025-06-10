package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type VerifyResponse struct {
	TotalCount int `json:"totalCount"`
}

func VerifyAgentLicense(agentName, registrationNumber string) (bool, error) {
	apiURL := "http://localhost:8081/verify"
	params := url.Values{}
	params.Add("agentName", agentName)
	params.Add("registrationNumber", registrationNumber)
	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return false, fmt.Errorf("목업 API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var apiResponse VerifyResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return false, fmt.Errorf("API 응답 JSON 파싱 실패: %w", err)
	}

	return apiResponse.TotalCount > 0, nil
}
