package ext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/117503445/goutils"
	"github.com/117503445/vsc-init/pkg/assets"
	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

func queryExtsMeta() string {
	var filters []map[string]interface{}
	for _, extName := range assets.Exts {
		filters = append(filters, map[string]interface{}{
			"criteria": []map[string]interface{}{
				{
					"filterType": 7,
					"value":      extName,
				},
			},
			"pageNumber": 1,
			"pageSize":   1,
		})
	}

	requestBody := map[string]interface{}{
		"filters": filters,
		"flags":   17, // 1(IncludeVersions) + 16(IncludeVersionProperties)
		// https://github.com/microsoft/vscode/blob/12ae331012923024bedaf873ba4259a8c64db020/src/vs/platform/extensionManagement/common/extensionGalleryService.ts
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal().Err(err).Msg("Marshal")
	}

	req, err := http.NewRequest("POST", "https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal().Err(err).Msg("NewRequest")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json;api-version=3.0-preview.1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msg("Do")
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err).Msg("ReadAll")
	}

	response := string(respBytes)
	return response
}
func getVscodeEngine() string {
	// TODO: get vscode engine version by exec
	const vscodeEngine = "1.91.1"
	return vscodeEngine
}

func isEngineValid(engine string, constraint string) (bool, error) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false, err
	}
	v, err := semver.NewVersion(engine)
	if err != nil {
		return false, err
	}
	return c.Check(v), nil
}

func InstallLatestExts() {
	var err error

	response := queryExtsMeta()
	vscodeEngine := getVscodeEngine()

	var versions []string
	var publishers []string
	var names []string
	for _, result := range gjson.Get(response, "results").Array() {
		publisher := result.Get("extensions.0.publisher.publisherName").String()
		publishers = append(publishers, publisher)

		name := result.Get("extensions.0.extensionName").String()
		names = append(names, name)

		var version string
		constraints := []string{}
		for _, versionResult := range result.Get("extensions.0.versions").Array() {
			constraint := ""
			// 遍历 versionResult.properties, 如果 property.key == "Microsoft.VisualStudio.Code.Engine", 那么 constraint = property.value
			for _, property := range versionResult.Get("properties").Array() {
				if property.Get("key").String() == "Microsoft.VisualStudio.Code.Engine" {
					constraint = property.Get("value").String()
					break
				}
			}
			if constraint == "" {
				log.Fatal().Msg("constraint is empty")
			}
			valid, err := isEngineValid(vscodeEngine, constraint)
			if err != nil {
				log.Fatal().Err(err).Str("constraint", constraint).Str("vscodeEngine", vscodeEngine).Msg("isEngineValid")
			}

			if valid {
				ver := versionResult.Get("version").String()
				version = ver
				break
			}
		}
		if version == "" {
			log.Fatal().Str("extName", name).Strs("constraints", constraints).Msg("version is empty")
		}
		versions = append(versions, version)

	}
	// log.Info().Strs("versions", versions).Msg("")

	latestExtVers := map[string]string{}
	latestExtPublishers := map[string]string{}
	latestExtNames := map[string]string{}
	for i, extName := range assets.Exts {
		latestExtVers[extName] = versions[i]
		latestExtPublishers[extName] = publishers[i]
		latestExtNames[extName] = names[i]
	}
	log.Info().Interface("latestExtVers", latestExtVers).Interface("latestExtPublishers", latestExtPublishers).Interface("latestExtNames", latestExtNames).Msg("")

	// var localExtVers map[string]string
	// if goutils.FileExists(fileExtVers) {
	// 	err = goutils.ReadJSON(fileExtVers, &localExtVers)
	// 	if err != nil {
	// 		log.Fatal().Err(err).Msg("ReadJSON")
	// 	}
	// } else {
	// 	localExtVers = map[string]string{}
	// }
	// log.Info().Interface("localExtVers", localExtVers).Msg("")

	for _, ext := range assets.Exts {
		url := fmt.Sprintf("https://ms-vscode.gallery.vsassets.io/_apis/public/gallery/publisher/%v/extension/%v/%v/assetbyname/Microsoft.VisualStudio.Services.VSIXPackage", latestExtPublishers[ext], latestExtNames[ext], latestExtVers[ext])

		extPath := "/tmp/exts/" + getExtFileName(ext, latestExtVers[ext])
		if !goutils.FileExists(extPath) {
			log.Info().Str("url", url).Str("extPath", extPath).Msg("Downloading")
			err = goutils.Download(url, extPath)
			if err != nil {
				log.Fatal().Err(err).Msg("DownloadFile")
			}
		}
	}
}
func getExtFileName(extName string, ver string) string {
	return fmt.Sprintf("%s-%s.vsix", extName, ver)
}
