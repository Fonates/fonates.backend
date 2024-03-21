package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	// Получаем путь к текущему исполняемому файлу
	// exePath, err := os.Executable()
	// if err != nil {
	// 	fmt.Println("Error getting executable path:", err)
	// 	return
	// }

	// Получаем путь к общей родительской директории
	parentDir := "/home/githubaction/actions-runner/_work"

	// Относительный путь к директории fonates.backend/fonates.backend относительно общей родительской директории
	fonatesBackendDir := filepath.Join(parentDir, "fonates.backend", "fonates.backend")

	// Преобразуем относительный путь в абсолютный путь
	absFonatesBackendDir, err := filepath.Abs(fonatesBackendDir)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	// Относительный путь к директории obs.alerts.plagin/obs.alerts.plagin относительно общей родительской директории
	obsAlertsPlaginDir := filepath.Join(parentDir, "obs.alerts.plagin", "obs.alerts.plagin")

	// Преобразуем относительный путь в абсолютный путь
	absObsAlertsPlaginDir, err := filepath.Abs(obsAlertsPlaginDir)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	fmt.Println("Directory fonates.backend/fonates.backend:", absFonatesBackendDir)
	fmt.Println("Directory obs.alerts.plagin/obs.alerts.plagin:", absObsAlertsPlaginDir)

	errFiles := filepath.Walk(absObsAlertsPlaginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error().Msgf("1: %s", err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Относительный путь файла внутри архива
		relativePath, err := filepath.Rel("obs.alerts.plagin", path)
		if err != nil {
			log.Error().Msgf("2: %s", err)
			return err
		}

		log.Info().Msgf("Relative path: %s", relativePath)

		// Если файл находится в поддиректории и его имя - main.min.js, то мы его изменяем
		if strings.Contains(relativePath, "scripts/main.min.js") {
			// Открываем файл для чтения
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Считываем содержимое файла
			content, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			// Ваш код для изменения содержимого файла main.min.js
			// Например, заменяем содержимое файла
			modifiedContent := []byte(strings.ReplaceAll(string(content), "<ton_wallet_address>", address))

			// Создаем файл в архиве и записываем в него измененное содержимое
			zipFile, err := zipWriter.Create(relativePath)
			if err != nil {
				return err
			}
			_, err = zipFile.Write(modifiedContent)
			if err != nil {
				return err
			}
		} else {
			excludeFiles := []string{"main.js", ".git", ".DS_Store"}

			// Проверяем, не содержится ли текущее имя файла в массиве исключений
			exclude := false
			for _, excluded := range excludeFiles {
				if info.Name() == excluded {
					exclude = true
					break
				}
			}

			if exclude {
				// Если имя файла содержится в массиве исключений, делаем что-то, например, пропускаем этот файл
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			zipFile, err := zipWriter.Create(relativePath)
			if err != nil {
				return err
			}

			_, err = io.Copy(zipFile, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if errFiles != nil {
		log.Error().Msgf("Error walking dir: %s", errFiles)
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error walking dir",
		})
		return
	}

	errClose := zipWriter.Close()
	if errClose != nil {
		log.Error().Msgf("Error closing zip writer: %s", errClose)
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error closing zip writer",
		})
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=obs.alerts.plagin.zip")
	w.Header().Set("Content-Length", string(buf.Len()))
	w.Write(buf.Bytes())
}
