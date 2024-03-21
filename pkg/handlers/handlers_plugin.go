package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"fonates.backend/pkg/models"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (h *Handlers) GeneratePlugin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if address == "" {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "Address is required",
		})
		return
	}

	link, err := models.InitDonationLink().GetByAddress(h.Store, address)
	if err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error getting link",
		})
		return
	}

	keyActivation, err := models.InitKeysActivation(link.ID).GetByLinkID(h.Store, link.ID)
	if err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error getting key activation",
		})
		return
	}

	log.Info().Msgf("Key: %s", keyActivation.Key.String())

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Получаем абсолютный путь к директории obs.alerts.plagin/obs.alerts.plagin
	absObsAlertsPlaginDir := "/home/githubaction/actions-runner/_work/obs.alerts.plagin/obs.alerts.plagin"

	// Записываем все содержимое директории в zip-архив
	errGetDir := filepath.Walk(absObsAlertsPlaginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Относительный путь файла внутри архива
		relativePath, err := filepath.Rel(absObsAlertsPlaginDir, path)
		if err != nil {
			return err
		}

		// Создаем файл в архиве и копируем содержимое из оригинального файла
		zipFile, err := zipWriter.Create(relativePath)
		if err != nil {
			return err
		}

		// Если это директория, ничего не делаем
		if info.IsDir() {
			return nil
		}

		// Если это файл, копируем его содержимое в zip-архив
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipFile, file)
		if err != nil {
			return err
		}

		return nil
	})
	if errGetDir != nil {
		fmt.Println("Error walking directory:", err)
		return
	}

	// Найти и изменить файл scripts/main.min.js
	// mainJSPath := filepath.Join(absObsAlertsPlaginDir, "scripts", "main.min.js")
	content := []byte("console.log('Hello, world!');") // Новое содержимое файла

	// Создаем файл в архиве и записываем в него измененное содержимое
	zipFile, err := zipWriter.Create("scripts/main.min.js")
	if err != nil {
		fmt.Println("Error creating zip file:", err)
		return
	}
	_, err = zipFile.Write(content)
	if err != nil {
		fmt.Println("Error writing to zip file:", err)
		return
	}

	// Закрываем архив
	err = zipWriter.Close()
	if err != nil {
		fmt.Println("Error closing zip writer:", err)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=obs.alerts.plagin.zip")
	w.Header().Set("Content-Length", string(buf.Len()))
	w.Write(buf.Bytes())
}
