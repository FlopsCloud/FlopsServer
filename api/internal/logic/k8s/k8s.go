package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Code      int      `json:"code"`
	Message   string   `json:"message"`
	Namespace string   `json:"namespace"`
	FcbPods   []FcbPod `json:"fcbPods"`
}

type FcbPod struct {
	Name    string `json:"name"`
	Port    int    `json:"port"`
	Cpu     int    `json:"cpu"`
	Memory  int    `json:"memory"`
	Storage int    `json:"storage"`
	Gpu     int    `json:"gpu"`
	Image   string `json:"image"`
	Mount   string `json:"mount"`
}
type CreateFcbPodRequest struct {
	Fcbpod FcbPod `json:"fcbpod"`
}

const Host = "http://127.0.0.1:8766"

func listVhost(token string) (*Response, error) {
	url := Host + "/vhost/list"
	return sendRequest("GET", url, token, nil)
}

func createVhost(token string, fcbPod FcbPod) (*Response, error) {
	url := Host + "/vhost/create"
	request := CreateFcbPodRequest{
		Fcbpod: fcbPod,
	}
	return sendRequest("POST", url, token, request)
}

type UpdateFcbPodRequest struct {
	Fcbpod FcbPod `json:"fcbpod"`
}

func updateVhost(token string, fcbPod FcbPod) (*Response, error) {
	url := Host + "/vhost/update"
	request := UpdateFcbPodRequest{
		Fcbpod: fcbPod,
	}
	return sendRequest("POST", url, token, request)
}

type DeleteRequest struct {
	Name string `json:"name"`
}

func deleteVhost(token string, name string) (*Response, error) {
	url := Host + "/vhost/del"
	request := DeleteRequest{
		Name: name,
	}
	return sendRequest("POST", url, token, request)
}

func DEMO() {
}
func Testmain() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBldGFAMTYzLmNvbSIsImV4cCI6MTc2NTY4NDUxNywiaWF0IjoxNzM0MTQ4NTE3LCJyb2xlIjoic3VwZXJhZG1pbiIsInVpZCI6MX0.nulwgyj-LZz7_V2SllwfMaBgJqk8Mbvg7HXqvTtESlk"
	// 调用获取列表接口
	listResp, listErr := listVhost(token)
	if listErr != nil {
		fmt.Println("获取列表失败:", listErr)
		return
	}
	fmt.Println("获取列表成功:", listResp)

	// 调用创建接口
	createFcbPod := FcbPod{
		Name:    "jupyter42",
		Port:    8888,
		Cpu:     1,
		Memory:  2,
		Storage: 20,
		Gpu:     1,
		Image:   "quay.io/jupyter/pytorch-notebook:cuda11-python-3.12.7",
		Mount:   "/home/jovyan/work",
	}
	createResp, createErr := createVhost(token, createFcbPod)
	if createErr != nil {
		fmt.Println("创建失败:", createErr)
		return
	}
	fmt.Println("创建成功:", createResp)

	// 调用更新接口
	updateFcbPod := FcbPod{
		Name:    "jupyter4",
		Port:    8888,
		Cpu:     2,
		Memory:  2,
		Storage: 20,
		Gpu:     1,
		Image:   "quay.io/jupyter/pytorch-notebook:cuda11-python-3.12.7",
		Mount:   "/home/jovyan/workx",
	}
	updateResp, updateErr := updateVhost(token, updateFcbPod)
	if updateErr != nil {
		fmt.Println("更新失败:", updateErr)
		return
	}
	fmt.Println("更新成功:", updateResp)

	// 调用删除接口
	delName := "jupyter4"
	delResp, delErr := deleteVhost(token, delName)
	if delErr != nil {
		fmt.Println("删除失败:", delErr)
		return
	}
	fmt.Println("删除成功:", delResp)
}
func sendRequest(method, url string, token string, body interface{}) (*Response, error) {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	client := &http.Client{}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
