package ext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/117503445/goutils"
	"github.com/117503445/vsc-init/pkg/assets"
	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

func queryExtsMeta() string {
	// var envExts string
	envExtsStr := os.Getenv("EXTS")
	if envExtsStr != "" {
		exts := make([]string, 0)
		envExts := strings.Split(envExtsStr, ",")
		for _, extName := range envExts {
			exts = append(exts, strings.TrimSpace(extName))
		}
		log.Info().Strs("exts", exts).Str("envExtsStr", envExtsStr).Msg("")
		assets.Exts = append(exts, assets.Exts...)
	}
	log.Info().Strs("assets.Exts", assets.Exts).Msg("")

	var filters []map[string]any
	for _, extName := range assets.Exts {
		filters = append(filters, map[string]any{
			"criteria": []map[string]any{
				{
					"filterType": 7,
					"value":      extName,
				},
			},
			"pageNumber": 1,
			"pageSize":   1,
		})
	}

	requestBody := map[string]any{
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
	// write config.yaml to prevent code-server log "Wrote default config file to /root/.config/code-server/config.yaml"
	err := goutils.WriteText("/root/.config/code-server/config.yaml", "")
	if err != nil {
		log.Fatal().Err(err).Msg("WriteText")
	}
	cmds := []string{"code-server", "--version"}
	outputByte, err := exec.Command(cmds[0], cmds[1:]...).Output()
	if err != nil {
		log.Fatal().Err(err).Msg("exec")
	}
	output := string(outputByte)

	var ver string

	if strings.Count(output, " ") == 4 {
		// 4.96.4 b7ef8f9bd70cb5b342fa8ec8a0086bad676d8124 with Code 1.96.4
		items := strings.Split(output, " ")
		ver = items[len(items)-1]
		ver = strings.TrimSpace(ver)
	} else {
		// 1.96.4
		// b7ef8f9bd70cb5b342fa8ec8a0086bad676d8124
		// x64
		items := strings.Split(output, "\n")
		ver = items[0]
		ver = strings.TrimSpace(ver)
	}

	log.Info().Str("vscodeEngine", ver).Msg("")

	return ver
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

	err = goutils.Download("https://github.com/117503445/vscode-key-runner/releases/latest/download/key-runner-0.0.1.vsix", "/tmp/exts/key-runner-0.0.1.vsix")
	if err != nil {
		log.Fatal().Err(err).Msg("DownloadFile")
	}
	cmds := []string{"code-server", "--install-extension", "/tmp/exts/key-runner-0.0.1.vsix"}
	cmd := exec.Command(cmds[0], cmds[1:]...)
	log.Info().Strs("cmds", cmds).Msg("")
	err = cmd.Run()
	if err != nil {
		log.Fatal().Err(err).Strs("cmds", cmds).Msg("exec")
	}

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

	getExtPath := func(ext string) string {
		return "/tmp/exts/" + getExtFileName(ext, latestExtVers[ext])
	}

	taskCh := make(chan string)
	var sg sync.WaitGroup
	for range 5 {
		sg.Add(1)
		go func() {
			defer sg.Done()
			for ext := range taskCh {
				url := fmt.Sprintf("https://ms-vscode.gallery.vsassets.io/_apis/public/gallery/publisher/%v/extension/%v/%v/assetbyname/Microsoft.VisualStudio.Services.VSIXPackage", latestExtPublishers[ext], latestExtNames[ext], latestExtVers[ext])

				extPath := getExtPath(ext)
				if !goutils.FileExists(extPath) {
					log.Info().Str("url", url).Str("extPath", extPath).Msg("Downloading")
					err = goutils.Download(url, extPath)
					if err != nil {
						log.Fatal().Err(err).Msg("DownloadFile")
					}
				}

			}
			log.Info().Msg("Done")
		}()
	}
	for _, ext := range assets.Exts {
		taskCh <- ext
	}
	close(taskCh)
	sg.Wait()

	for _, ext := range assets.Exts {
		extPath := getExtPath(ext)
		cmds := []string{"code-server", "--install-extension", extPath}
		cmd := exec.Command(cmds[0], cmds[1:]...)
		log.Info().Strs("cmds", cmds).Msg("")
		err := cmd.Run()
		if err != nil {
			log.Fatal().Err(err).Strs("cmds", cmds).Msg("exec")
		}
	}

}
func getExtFileName(extName string, ver string) string {
	return fmt.Sprintf("%s-%s.vsix", extName, ver)
}
