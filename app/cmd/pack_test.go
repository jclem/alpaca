package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"testing"

	"github.com/groob/plist"
	"github.com/jclem/alpaca/workflow"
	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPack(t *testing.T) {
	dir, err := filepath.Abs("./fixtures/pack_test")
	if err != nil {
		t.Fatal(err)
	}

	out := packWorkflow(dir)

	wfFile := filepath.Join(out, "pack_test.alfredworkflow")
	zipOut := unzip(wfFile)
	plistPath := filepath.Join(zipOut, "info.plist")
	var i workflow.Info
	if err := plist.Unmarshal(readFile(plistPath), &i); err != nil {
		t.Fatal(err)
	}

	// Test basic workflow metadata
	assert.Equal(t, "pack_test", i.Name)
	assert.Equal(t, "0.1.0", i.Version)
	assert.Equal(t, "Jonathan Clem <jonathan@jclem.net>", i.CreatedBy)
	assert.Equal(t, "com.jclem.alfred.alpaca-test.say-hello", i.BundleID)
	assert.Equal(t, "Says words", i.Description)
	assert.Equal(t, "https://github.com/jclem/alpaca/blob/master/app/tests/fixtures/pack_test", i.WebAddress)
	assert.Equal(t, "This is information about the workflow.\n", i.Readme)
	assert.Equal(t, readFile(filepath.Join(dir, "img/alpaca.png")), readFile(filepath.Join(zipOut, "icon.png")))
	assert.Equal(t, map[string]string{"FOO": "foo"}, i.Variables)
	assert.Equal(t, []string{"FOO"}, i.VariablesDontExport)

	// Test objects
	// Sort objects by type
	sortedObjs := make([]map[string]interface{}, len(i.Objects))
	copy(sortedObjs, i.Objects)
	sort.Slice(sortedObjs, func(i, j int) bool {
		iType := sortedObjs[i]["type"].(string)
		jType := sortedObjs[j]["type"].(string)
		return iType < jType
	})
	assert.Equal(t, 6, len(sortedObjs))

	applescript := sortedObjs[0]
	config := applescript["config"].(map[string]interface{})
	assert.True(t, config["cachescript"].(bool))
	assert.Equal(t, "the-applescript", config["applescript"].(string))
	assert.Equal(t, "alfred.workflow.action.applescript", applescript["type"])

	openurl := sortedObjs[1]
	config = openurl["config"].(map[string]interface{})
	assert.Equal(t, "alfred.workflow.action.openurl", openurl["type"])
	assert.Equal(t, "https://example.com", config["url"])

	script := sortedObjs[2]
	config = script["config"].(map[string]interface{})
	assert.Equal(t, "alfred.workflow.action.script", script["type"])
	assert.Equal(t, readFile(filepath.Join(dir, "img/alpaca.png")), readFile(filepath.Join(zipOut, script["uid"].(string)+".png")))
	assert.Equal(t, `echo "hi"`, config["script"])
	assert.Equal(t, uint64(0), config["type"])
	assert.Equal(t, uint64(1), config["scriptargtype"])

	keyword := sortedObjs[3]
	config = keyword["config"].(map[string]interface{})
	assert.Equal(t, "alfred.workflow.input.keyword", keyword["type"])
	assert.Equal(t, "the-keyword", config["keyword"])
	assert.True(t, config["withspace"].(bool))
	assert.Equal(t, "Keyword", config["text"])
	assert.Equal(t, "keyword", config["subtext"])
	assert.Equal(t, uint64(2), config["argumenttype"])

	scriptfilter := sortedObjs[4]
	config = scriptfilter["config"].(map[string]interface{})
	assert.Equal(t, "alfred.workflow.input.scriptfilter", scriptfilter["type"])
	assert.Equal(t, readFile(filepath.Join(dir, "img/alpaca.png")), readFile(filepath.Join(zipOut, scriptfilter["uid"].(string)+".png")))
	assert.Equal(t, uint64(1), config["argumenttype"])
	assert.Equal(t, uint64(1), config["argumenttrimmode"])
	assert.Equal(t, uint64(33), config["escaping"])
	assert.True(t, config["argumenttreatemptyqueryasnil"].(bool))
	assert.Equal(t, "filter", config["keyword"])
	assert.Equal(t, "Please wait...", config["runningsubtext"])
	assert.Equal(t, "Runs a script filter", config["subtext"])
	assert.Equal(t, "Run script-filter", config["title"])
	assert.True(t, config["withspace"].(bool))
	assert.Equal(t, uint64(1), config["scriptargtype"])
	assert.Equal(t, "scripts/script.js", config["scriptfile"])
	assert.True(t, config["alfredfiltersresults"].(bool))
	assert.Equal(t, uint64(2), config["alfredfiltersresultsmatchmode"])
	assert.True(t, config["queuedelayimmediatelyinitially"].(bool))
	assert.Equal(t, uint64(1), config["queuemode"])
	assert.Equal(t, uint64(1), config["queuedelaymode"])

	clipboard := sortedObjs[5]
	config = clipboard["config"].(map[string]interface{})
	assert.Equal(t, "alfred.workflow.output.clipboard", clipboard["type"])
	assert.Equal(t, "{query}", config["text"])


	// Test connections
	assert.Equal(t, applescript["uid"], i.Connections[clipboard["uid"].(string)][0].To)
	assert.Equal(t, applescript["uid"], i.Connections[keyword["uid"].(string)][0].To)
	assert.Equal(t, clipboard["uid"], i.Connections[scriptfilter["uid"].(string)][0].To)

	// Test UI data
	assert.Equal(t, int64(20), i.UIData[openurl["uid"].(string)].XPos)
	assert.Equal(t, int64(20), i.UIData[openurl["uid"].(string)].YPos)

	assert.Equal(t, int64(20), i.UIData[script["uid"].(string)].XPos)
	assert.Equal(t, int64(145), i.UIData[script["uid"].(string)].YPos)

	assert.Equal(t, int64(20), i.UIData[keyword["uid"].(string)].XPos)
	assert.Equal(t, int64(270), i.UIData[keyword["uid"].(string)].YPos)

	assert.Equal(t, int64(20), i.UIData[scriptfilter["uid"].(string)].XPos)
	assert.Equal(t, int64(395), i.UIData[scriptfilter["uid"].(string)].YPos)

	assert.Equal(t, int64(265), i.UIData[clipboard["uid"].(string)].XPos)
	assert.Equal(t, int64(20), i.UIData[clipboard["uid"].(string)].YPos)

	assert.Equal(t, int64(510), i.UIData[applescript["uid"].(string)].XPos)
	assert.Equal(t, int64(20), i.UIData[applescript["uid"].(string)].YPos)
}

func packWorkflow(dir string) string {
	out = mktemp()
	packCmd.Run(&cobra.Command{}, []string{dir})
	return out
}

func mktemp() string {
	temp, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	return temp
}

func unzip(path string) string {
	out := mktemp()
	zip := archiver.NewZip()
	if err := zip.Unarchive(path, out); err != nil {
		panic(err)
	}
	return out
}

func readFile(path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}

func printJson(x interface{}) {
	j, _ := json.Marshal(x)
	fmt.Println(string(j))
}
