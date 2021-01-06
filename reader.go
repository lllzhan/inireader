package inireader

import (
	"io/ioutil"
	"regexp"
	"strings"
)

// IReader 类
type IReader struct {
	FilePath  string
	FileData  string
	ConfigMap map[string]map[string]string
}

// InitIReader IReader构造函数
func InitIReader(fp string) (*IReader, error) {
	iReader := IReader{
		FilePath:  fp,
		ConfigMap: make(map[string]map[string]string, 0),
	}

	err := iReader.read()
	if err != nil {
		return nil, err
	}

	err = iReader.parse()
	if err != nil {
		return nil, err
	}
	return &iReader, nil
}

func (i *IReader) read() error {
	file, err := ioutil.ReadFile(i.FilePath)
	if err != nil {
		return err
	}
	i.FileData = string(file)
	return nil
}

func (i *IReader) parse() error {
	lines := strings.Split(i.FileData, "\n")
	roleBlock := regexp.MustCompile(`\[(.*?)\]`)
	roleLine := regexp.MustCompile(`(.*?)=(.*)`)
	var nowBlock string = ""
	for _, eline := range lines {
		if roleBlock.MatchString(eline) {
			temp := roleBlock.FindStringSubmatch(eline)
			nowBlock = temp[1]
			i.ConfigMap[nowBlock] = make(map[string]string, 0)
		} else if roleLine.MatchString(eline) {
			kvs := roleLine.FindStringSubmatch(eline)
			i.ConfigMap[nowBlock][string.TrimSpace(kvs[1])] = string.TrimSpace(kvs[2])
		}

	}
	return nil
}

// Key 获取配置参数（例：redis/host）
func (i *IReader) Key(name string) string {
	keys := strings.Split(name, "/")
	if len(keys) < 1 {
		return ""
	}

	if kv, ok := i.ConfigMap[keys[0]]; ok {
		if value, ok := kv[keys[1]]; ok {
			return value
		}
	}
	return ""
}
