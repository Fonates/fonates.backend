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

		excludeFiles := []string{"main.js", ".DS_Store", ".git", ".github"} // Добавьте другие файлы, которые нужно пропустить
		for _, excluded := range excludeFiles {
			if info.Name() == excluded {
				return nil
			}
		}

		// Получаем относительный путь файла внутри директории
		relativePath, err := filepath.Rel(absObsAlertsPlaginDir, path)
		if err != nil {
			return err
		}

		// Создаем запись для каждой директории в архиве
		if info.IsDir() {
			_, err := zipWriter.Create(relativePath + "/") // Добавляем "/" в конец, чтобы обозначить директорию
			if err != nil {
				return err
			}
			return nil
		}

		// Создаем файл в архиве
		zipFile, err := zipWriter.Create(relativePath)
		if err != nil {
			return err
		}

		if info.Name() == "main.min.js" {
			// Открываем файл для чтения
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Читаем содержимое файла
			content, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			// Вносим изменения в содержимое файла (заменяем <ton_wallet_address> на необходимое значение)
			modifiedContent := bytes.ReplaceAll(content, []byte("<ton_wallet_address>"), []byte(address))
			modifiedContent = bytes.ReplaceAll(modifiedContent, []byte("<key-activation>"), []byte(keyActivation.Key.String()))
			
			// Записываем измененное содержимое в архив
			_, err = zipFile.Write(modifiedContent)
			if err != nil {
				return err
			}
		} else {
			// Если это не файл main.min.js, копируем его содержимое в архив без изменений
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(zipFile, file)
			if err != nil {
				return err
			}
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
	w.Header().Set("Content-Disposition", "attachment; filename=fonates-plugin.zip")
	w.Header().Set("Content-Length", string(buf.Len()))
	w.Write(buf.Bytes())
}
