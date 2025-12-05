package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"videorepack/ffmpeg"
	"videorepack/mkv"
	"videorepack/naming"

	log "github.com/sirupsen/logrus"
)

const Name = "videorepack"

func main() {
	log.SetLevel(log.TraceLevel)

	if len(os.Args) < 2 {
		log.Fatalf("Uso: %s <input.mkv>", Name)
	}

	if strings.Index(os.Args[1], "*") == -1 {
		_ = convertVideo(os.Args[1])
	} else {
		// Walk files that match input pattern
		matches, err := filepath.Glob(os.Args[1])
		if err != nil {
			log.Fatalf("Error al buscar archivos: %v", err)
		}
		if len(matches) == 0 {
			log.Fatalf("No se encontraron archivos que coincidan con el patrón: %s", os.Args[1])
		}

		for _, file := range matches {
			log.Infof("Procesando archivo: %s", file)
			err := convertVideo(file)
			if err != nil {
				log.Errorf("Error al convertir %s: %v", file, err)
			}
		}
	}

}

func convertVideo(input string) error {
	outputPath := filepath.Dir(input)
	filesToDelete := []string{}
	defer func() {
		log.Info("Limpiando archivos intermedios...")
		for _, f := range filesToDelete {
			os.Remove(f)
		}
	}()

	log.Infof("Extrayendo pistas...")
	extracted, err := mkv.ExtractAll(input, "")
	if err != nil {
		log.Fatalf("Error extrayendo pistas: %v", err)
	}
	for _, t := range extracted.Tracks {
		filesToDelete = append(filesToDelete, t.FilePath)
		if len(t.TimeMapPath) > 0 {
			filesToDelete = append(filesToDelete, t.TimeMapPath)
		}
	}

	// Filtrar y modificar pistas
	originalLang, _ := mkv.FromIETFName("ja")
	onlyAudios := []string{"ja", "es", "es-ES", "gl", "gl-ES"}
	mainLang, _ := mkv.FromIETFName("gl")

	hasMainLangAudio := false
	hasMainLangSub := false
	var selected []mkv.ExtractedTrack
	for _, t := range extracted.Tracks {
		t.Info.Properties.TrackName = "" // Clear track name

		if t.Info.Type == "video" || t.Info.Type == "audio" {
			t.Info.Properties.TrackName = ""
			if t.Info.Properties.LanguageIETF == originalLang {
				t.Info.Properties.FlagOriginal = true
			}

			if t.Info.Type == "audio" {
				// Set default audio track to mainLang if exists
				if t.Info.Properties.LanguageIETF == mainLang && !hasMainLangAudio {
					t.Info.Properties.DefaultTrack = true
					hasMainLangAudio = true
				} else {
					t.Info.Properties.DefaultTrack = false
				}

				// Select only specified audio languages
				for _, lang := range onlyAudios {
					langIetf, _ := mkv.FromIETFName(lang)
					if t.Info.Properties.LanguageIETF == langIetf {
						selected = append(selected, t)
						break
					}
				}
			} else {
				selected = append(selected, t)
			}
		} else if t.Info.Type == "subtitles" {
			if hasMainLangSub {
				t.Info.Properties.DefaultTrack = false
			} else if t.Info.Properties.LanguageIETF == mainLang {
				if hasMainLangAudio {
					t.Info.Properties.DefaultTrack = t.Info.Properties.ForcedTrack
				} else {
					t.Info.Properties.DefaultTrack = true
				}
				if t.Info.Properties.DefaultTrack {
					hasMainLangSub = true
				}
			}

			selected = append(selected, t)
		}
	}

	// Convertir pistas con codecs no deseados
	for i := range selected {
		t := &selected[i]
		if t.Info.Type == "audio" {
			if strings.Index(t.Info.Properties.CodecID, "FLAC") != -1 {
				log.Infof("Convirtiendo pista de audio FLAC a EAC3 (%s) desde <%s>...", t.Info.Properties.LanguageIETF.String(), t.FilePath)
				targetFilePath := t.FilePath + ".eac3"
				err := ffmpeg.Convert(ffmpeg.ConvertOptions{
					Inputs: []ffmpeg.InputFile{{
						Path: t.FilePath,
					}},
					OutputPath: targetFilePath,
					Tracks: []ffmpeg.TrackConvertOptions{{
						Index:   "a",
						Encoder: ffmpeg.EncoderEAC3,
					}},
				})
				if err != nil {
					log.Warnf("Error al convertir pista de audio: %v. Se continua con la pista original.", err)
				} else {
					t.FilePath = targetFilePath
					filesToDelete = append(filesToDelete, targetFilePath)
					t.Info.Codec = "A_EAC3"
					t.Info.Properties.CodecID = "A_EAC3"
				}
			}
		}
	}

	// Construir nombre de archivo de salida
	parsedFileName := naming.Extract(filepath.Base(input))
	for i := range selected {
		t := &selected[i]
		if t.Info.Type == "video" {
			parsedFileName.VideoMetadata = t.Info.NamingMetadata()
		}
		if t.Info.Type == "audio" {
			parsedFileName.AudioMetadata = append(parsedFileName.AudioMetadata, t.Info.NamingMetadata())
		}
	}

	parsedFileName.Origin = "WEBDL"
	parsedFileName.Authors = []string{"Dussarax"}

	// Delay del español
	/*for i := range selected {
		t := &selected[i]
		if t.Info.Type == "audio" && t.Info.Properties.LanguageIETF.Tag == language.EuropeanSpanish {
			t.Operations.Delay = 6000
		}
	}*/

	// Escribir fichero de salida
	outputFile := path.Join(outputPath, parsedFileName.FileName())
	log.Infof("Empaquetando fichero de salida %s ...", outputFile)
	err = mkv.Merge(outputFile, mkv.ExtractedContainer{
		Tracks:   selected,
		Chapters: extracted.Chapters,
	})
	if err != nil {
		log.Errorf("Error al reempaquetar: %v", err)
	} else {
		log.Info("Proceso completado!")
	}

	return nil
}
