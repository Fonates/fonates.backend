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

	errА := filepath.Walk(absObsAlertsPlaginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Получаем относительный путь файла внутри директории
		relativePath, err := filepath.Rel(absObsAlertsPlaginDir, path)
		if err != nil {
			return err
		}

		// Выводим размер файла или директории
		fmt.Printf("%s %d байт\n", relativePath, info.Size())

		// Создаем файл в архиве
		zipFile, err := zipWriter.Create(relativePath)
		if err != nil {
			return err
		}

		// Если это директория, то ничего не делаем
		if info.IsDir() {
			return nil
		}

		// Если это файл, то копируем его содержимое в архив
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
	if errА != nil {
		fmt.Println("Error walking directory:", err)
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
