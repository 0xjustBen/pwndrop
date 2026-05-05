package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/kgretzky/pwndrop/config"
	"github.com/kgretzky/pwndrop/log"
	"github.com/kgretzky/pwndrop/storage"
)

type ApiResponse struct {
	ErrorCode int         `json:"error_code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

var Cfg *config.Config = nil

func SaveUploadedFile(file multipart.File, fhead *multipart.FileHeader, save_path string) error {
	f, err := os.OpenFile(save_path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}
	return nil
}

func DumpResponse(w http.ResponseWriter, message string, http_status int, error_code int, o interface{}) {
	resp := &ApiResponse{
		ErrorCode: error_code,
		Message:   message,
		Data:      o,
	}

	d, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "corrupted response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http_status)
	w.Header().Set("content-type", "application/json")
	w.Write(d)
}

func getMachineName(r *http.Request) string {
	machineName := r.FormValue("machine_name")
	if machineName != "" {
		return machineName
	}
	return r.Header.Get("X-Machine-Name")
}

func realIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return fwd
	}
	return r.RemoteAddr
}

func AuditEvent(uid int, action, status, fileName, machineName string, r *http.Request) {
	username := fmt.Sprintf("uid:%d", uid)
	if u, err := storage.UserGet(uid); err == nil {
		username = u.Name
	}
	if machineName == "" {
		machineName = "-"
	}
	if fileName == "" {
		fileName = "-"
	}
	log.Info("AUDIT action=%s status=%s user=%q file=%q machine=%q url=%s ip=%s",
		action, status, username, fileName, machineName, r.URL.Path, realIP(r))
}

func SetConfig(cfg *config.Config) {
	Cfg = cfg
}
