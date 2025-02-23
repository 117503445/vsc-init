package assets

var Exts = []string{"alibaba-cloud.tongyi-lingma",
	"ms-azuretools.vscode-docker",
	"ms-ceintl.vscode-language-pack-zh-hans",
	"pkief.material-icon-theme",
	"jnoortheen.nix-ide",
	"golang.go",
	"tamasfe.even-better-toml",
	"bodil.prettier-toml",
	"humao.rest-client",
	"njzy.stats-bar",
	"zxh404.vscode-proto3",
	"redhat.vscode-xml",
	"redhat.vscode-yaml",
	"mhutchie.git-graph",
	"ms-python.python",
	"mechatroner.rainbow-csv",
	"eamodio.gitlens"}

var KeyBindings = `[
    {
        "key": "ctrl+alt+left",
        "command": "workbench.action.navigateBack",
        "when": "canNavigateBack"
    },
    {
        "key": "alt+left",
        "command": "-workbench.action.navigateBack",
        "when": "canNavigateBack"
    },
    {
        "key": "ctrl+alt+right",
        "command": "workbench.action.navigateForward",
        "when": "canNavigateForward"
    },
    {
        "key": "alt+right",
        "command": "-workbench.action.navigateForward",
        "when": "canNavigateForward"
    },
    {
        "key": "f5",
        "command": "key-runner.run"
    },
]`

var Settings = `{
  "editor.minimap.enabled": false,
  "editor.unicodeHighlight.allowedCharacters": {
    "！": true,
    "：": true,
    "，": true
  },
  "editor.wordWrap": "on",
  "explorer.confirmDelete": false,
  "explorer.confirmDragAndDrop": false,
  "editor.wordSeparators": "~!@#$%^&*()-=+[{]}\\|;:'\",.<>/！（）｛｝【】、；：’”，。《》？",
  "git.autofetch": true,
  "git.confirmSync": false,
  "git.ignoreMissingGitWarning": true,
  "git.enableSmartCommit": true,
  "go.testFlags": [
    "-v",
    "-count=1"
  ],
  "go.toolsManagement.autoUpdate": true,
  "terminal.integrated.copyOnSelection": true,
  "terminal.integrated.fontFamily": "UbuntuMono Nerd Font Mono,UbuntuMono NF,Microsoft YaHei Mono",
  "terminal.integrated.fontWeightBold": "normal",
  "terminal.integrated.defaultProfile.linux": "fish",
  "terminal.integrated.profiles.linux": {
    "fish": {
      "path": "fish"
    }
  },
  // "window.title": "${rootName}",
  //WINDOW
  "window.commandCenter": false,
  "workbench.layoutControl.enabled": false,
  "git.alwaysSignOff": true,
  "editor.inlineSuggest.enabled": true,
  "git.defaultBranchName": "master",
  "explorer.confirmPasteNative": false,
  "window.autoDetectColorScheme": true,
  "files.autoSave": "afterDelay",
  "workbench.iconTheme": "material-icon-theme",
  "redhat.telemetry.enabled": false,
  "terminal.integrated.commandsToSkipShell": [
    "key-runner.run"
  ],
  "terminal.integrated.scrollback": 10000
}`
