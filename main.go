package main

import (
	"encoding/json"
	"errors"
	"log"
	"time"
	"net/http"
	"os"
	"strconv"

	ld "gopkg.in/launchdarkly/go-server-sdk.v4"
	"gopkg.in/launchdarkly/go-sdk-common.v1/ldvalue"
	"github.com/gorilla/mux"
)

var client *ld.LDClient

type UserRequest struct {
	User ld.User `json:"user"`
}

type CustomEventRequest struct {
	User        ld.User       `json:"user"`
	Key         string        `json:"key"`
	Data        ldvalue.Value `json:"data"`
	MetricValue *float64      `json:"metricValue"`
}

type RootResponse struct {
	Initialized bool `json:"initialized"`
}

type EvalFeatureRequest struct {
	User         ld.User       `json:"user"`
	DefaultValue ldvalue.Value `json:"defaultValue"`
	Detail       bool          `json:"detail"`
}

type EvalFeatureResponse struct {
	Key            string              `json:"key"`
	Result         ldvalue.Value       `json:"result"`
	VariationIndex *int                `json:"variationIndex,omitempty"`
	Reason         ld.EvaluationReason `json:"reason,omitempty"`
}

func getRootHandler(w http.ResponseWriter, r *http.Request) {
	responseBody := RootResponse{
		Initialized: client.Initialized(),
	}

	responseEncoded, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseEncoded)
}

func postFlushHandler(w http.ResponseWriter, req *http.Request) {
	client.Flush()
	w.WriteHeader(http.StatusNoContent)
}

func postEventHandler(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		http.Error(w, "expected a body", http.StatusBadRequest)
		return
	}

	var params CustomEventRequest
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if params.MetricValue == nil {
		err = client.Track(params.Key, params.User, params.Data)
	} else {
		err = client.TrackWithMetric(params.Key, params.User, params.Data, *params.MetricValue)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func postIdentifyHandler(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		http.Error(w, "expected a body", http.StatusBadRequest)
		return
	}

	var params UserRequest
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = client.Identify(params.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func postAllFlagsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		http.Error(w, "expected a body", http.StatusBadRequest)
		return
	}

	var params UserRequest
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allFlags := client.AllFlags(params.User)

	allFlagsEncoded, err := json.Marshal(allFlags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(allFlagsEncoded)
}

func PostVariationHandler(w http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]

	if req.Body == nil {
		http.Error(w, "expected a body", http.StatusBadRequest)
		return
	}

	var params EvalFeatureRequest
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := EvalFeatureResponse{Key: key}

	if params.Detail {
		value, detail, _ := client.JSONVariationDetail(key, params.User, params.DefaultValue)
		response.Result = value
		response.VariationIndex = detail.VariationIndex
		response.Reason = detail.Reason
	} else {
		value, _ := client.JSONVariation(key, params.User, params.DefaultValue)
		response.Result = value
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(encoded)
}

type Config struct {
	port uint16
	key  string
}

func makeDefaultConfig() Config {
	return Config{
		port: 8080,
		key:  "",
	}
}

func loadConfigFromEnvironment(config *Config) error {
	key := os.Getenv("SDK_KEY")
	if key == "" {
		return errors.New("SDK_KEY is required")
	} else {
		config.key = key
	}

	port := os.Getenv("PORT")
	if port != "" {
		if x, err := strconv.ParseUint(port, 10, 16); err == nil {
			config.port = uint16(x)
		} else {
			return err
		}
	}

	return nil
}

func main() {
	config := makeDefaultConfig()
	err := loadConfigFromEnvironment(&config)
	if err != nil {
		log.Fatal(err)
		return
	}

	client, _ = ld.MakeClient(config.key, 5 * time.Second)
	defer client.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", getRootHandler).Methods("GET")
	router.HandleFunc("/track", postEventHandler).Methods("POST")
	router.HandleFunc("/flush", postFlushHandler).Methods("POST")
	router.HandleFunc("/identify", postIdentifyHandler).Methods("POST")
	router.HandleFunc("/allFlags", postAllFlagsHandler).Methods("POST")
	router.HandleFunc("/feature/{key}/eval", PostVariationHandler).Methods("POST")

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(config.port)), nil))
}
