package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

type constantHandler struct {
	Val string
}

type ApiResponse struct {
	ID     string   `json:"id"`
	URL    string   `json:"uri"`
	Status string   `json:"status"`
	Links  []string `json:"links"`
	stream int      `json:"streamable"`
}

type handleMagnet struct {
	Magnet string `json:"magnet"`
}

type env struct {
	apiKey string
	apiUrl string
}

var globalEnv *env

var err error

func NewEnv() (*env, error) {
	err := godotenv.Load()
	if err != nil {
		// Don't return error here, maybe ENV vars are set directly
		log.Println("Info/Warning: .env file not found or error loading:", err)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		// Decide if this is fatal. Let's assume it is for this example.
		return nil, fmt.Errorf("API_KEY environment variable not set")
	}

	apiUrl := os.Getenv("SERVER_HOST")
	if apiUrl == "" {
		// Decide if this is fatal. Let's assume it is for this example.
		return nil, fmt.Errorf("SERVER_HOST environment variable not set")
	}

	return &env{
		apiKey: apiKey,
		apiUrl: apiUrl,
	}, nil
}

func (ch constantHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	_, _ = rw.Write([]byte(ch.Val))
}

func (mg handleMagnet) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var magnet string

	fmt.Println(globalEnv.apiKey)
	fmt.Println(globalEnv.apiUrl)

	queryMagnet := req.URL.Query().Get("uri")
	if queryMagnet != "" {
		magnet = queryMagnet
	} else {
		if req.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		// Decode the magnet link from the request body
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&mg.Magnet); err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			log.Println("Error decoding body:", err)
			return
		}
		magnet = mg.Magnet
	}

	formData := url.Values{}

	formData.Set("magnet", magnet)

	encodedFormData := formData.Encode()

	payload := strings.NewReader(encodedFormData)

	// Prepare the HTTP request
	req, err = http.NewRequest("POST", globalEnv.apiUrl+"/torrents/addMagnet", payload)
	if err != nil {
		http.Error(w, "Error creating new request", http.StatusInternalServerError)
		log.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+globalEnv.apiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error making API request", http.StatusInternalServerError)
		log.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and handle the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		log.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the JSON response into the ApiResponse struct
	var apiResp1 ApiResponse
	if err := json.Unmarshal(respBody, &apiResp1); err != nil {
		http.Error(w, "Error unmarshaling response", http.StatusInternalServerError)
		return
	}

	// Extract the id from the response
	ID := apiResp1.ID

	// POST data
	formData2 := url.Values{}

	formData2.Set("files", "all")

	encodedFormData2 := formData2.Encode()

	payload2 := strings.NewReader(encodedFormData2)

	fmt.Println(payload2)

	selectFilesUrl := globalEnv.apiUrl + "/torrents/selectFiles/" + ID

	// Create request
	req2, err := http.NewRequest("POST", selectFilesUrl, payload2)
	if err != nil {
		http.Error(w, "Error creating new request", http.StatusInternalServerError)
		return
	}

	// Set headers
	req2.Header.Set("Authorization", "Bearer "+globalEnv.apiKey)

	// Send the request
	client2 := &http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil {
		http.Error(w, "Error making API request", http.StatusInternalServerError)
		return
	}
	defer resp2.Body.Close()

	// POST data

	selectFilesUrl2 := globalEnv.apiUrl + "/torrents/info/" + ID

	// Create request
	req3, err := http.NewRequest("GET", selectFilesUrl2, nil)
	if err != nil {
		http.Error(w, "Error creating new request", http.StatusInternalServerError)
		return
	}

	// Set headers
	req3.Header.Set("Authorization", "Bearer "+globalEnv.apiKey)

	// Send the request
	client3 := &http.Client{}
	resp3, err := client3.Do(req3)
	if err != nil {
		http.Error(w, "Error making API request", http.StatusInternalServerError)
		return
	}
	defer resp3.Body.Close()

	// Read and handle the response
	respBody3, err := io.ReadAll(resp3.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		log.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the JSON response into the ApiResponse struct
	var apiResp4 ApiResponse
	if err := json.Unmarshal(respBody3, &apiResp4); err != nil {
		http.Error(w, "Error unmarshaling response", http.StatusInternalServerError)
		return
	}

	// Extract the id from the response

	Link := apiResp4.Links[0]
	Status := apiResp4.Status
	Stream := apiResp4.stream

	// Extract the id from the response

	formData3 := url.Values{}

	formData3.Set("link", Link)

	encodedFormData3 := formData3.Encode()

	payload3 := strings.NewReader(encodedFormData3)

	selectFilesUrl3 := globalEnv.apiUrl + "/unrestrict/link"

	if Status == "downloaded" {

		// Prepare the HTTP request
		req, err = http.NewRequest("POST", selectFilesUrl3, payload3)
		if err != nil {
			http.Error(w, "Error creating new request", http.StatusInternalServerError)
			log.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Authorization", "Bearer "+globalEnv.apiKey)

		// Make the request
		client4 := &http.Client{}
		resp4, err := client4.Do(req)
		if err != nil {
			http.Error(w, "Error making API request", http.StatusInternalServerError)
			log.Println("Error making request:", err)
			return
		}
		defer resp4.Body.Close()

		// Read and handle the response
		respBody4, err := io.ReadAll(resp4.Body)
		if err != nil {
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			log.Println("Error reading response body:", err)
			return
		}

		// Unmarshal the JSON response into the ApiResponse struct
		var apiResp3 ApiResponse
		if err := json.Unmarshal(respBody4, &apiResp3); err != nil {
			http.Error(w, "Error unmarshaling response", http.StatusInternalServerError)
			return
		}

		// id4 := apiResp3.ID
		// streamURL := fmt.Sprintf("https://real-debrid.com/streaming-%s", id4)

		// Step 2: Unmarshal into a raw map to modify fields
		var raw map[string]json.RawMessage
		if err := json.Unmarshal(respBody4, &raw); err != nil {
			http.Error(w, "Error unmarshaling JSON", http.StatusInternalServerError)
			return
		}

		delete(raw, "streamable")
		delete(raw, "link")
		delete(raw, "id")
		delete(raw, "host")

		// raw["watch_now"] = json.RawMessage(fmt.Sprintf(`"%s"`, streamURL))

		// Step 5: Marshal the updated map
		updatedJSON, err := json.Marshal(raw)
		if err != nil {
			http.Error(w, "Error marshaling updated JSON", http.StatusInternalServerError)
			return
		}

		// Write response
		w.Write(updatedJSON)

	} else {
		w.Write([]byte("Torrent not downloaded yet"))
	}

}

func createMux(version string) *http.ServeMux {
	mux := http.NewServeMux()
	if version == "v1" {
		mux.Handle("/magnet", handleMagnet{})
		mux.Handle("/torrent", constantHandler{version + "torrent"})
	}

	return mux
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {

	globalEnv, err = NewEnv()
	if err != nil {
		os.Exit(1)
	}

	ip := GetOutboundIP()

	addressFormat := fmt.Sprintf("%s:%d", ip, 9000)

	w := os.Stderr

	slog.SetDefault(slog.New(tint.NewHandler(colorable.NewColorable(w), &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})))

	topMux := http.NewServeMux()
	topMux.Handle("/v1/", http.StripPrefix("/v1", createMux("v1")))

	slog.Info("Starting server", "addr", fmt.Sprintf("http://%s", addressFormat))

	http.ListenAndServe(addressFormat, topMux)

}
