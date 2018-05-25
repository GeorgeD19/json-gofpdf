package jsongofpdf

import (
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/buger/jsonparser"
	"github.com/h2non/filetype"
	"github.com/spf13/cast"
)

func (p *JSONGOFPDF) GetBool(name string, logic string, fallback bool) (value bool) {
	result := fallback
	attribute, _, _, err := p.GetAttribute(name, logic, false)
	if err == nil {
		result, _ = jsonparser.ParseBoolean(attribute)
	}
	return result
}

func (p *JSONGOFPDF) GetFloat(name string, logic string, fallback float64) (value float64) {
	result := fallback
	attribute, _, _, err := p.GetAttribute(name, logic, false)
	if err == nil {
		result = cast.ToFloat64(cast.ToString(attribute))
	}
	return result
}

func (p *JSONGOFPDF) GetInt(name string, logic string, fallback int) (value int) {
	result := fallback
	attribute, _, _, err := p.GetAttribute(name, logic, false)

	if err == nil {
		result = cast.ToInt(cast.ToString(attribute))
	}
	return result
}

func (p *JSONGOFPDF) GetString(name string, logic string, fallback string) (value string) {
	return p.GetStringIndex(name, logic, fallback)
}

// GetImage returns a File type containing the hex version of the image with associated meta information
func GetImage(FileName string) (f ImageFile, Err error) {

	FoundFile := ImageFile{}

	// Get file contents
	FileData, _ := ioutil.ReadFile(FileName)

	FileMeta, err := filetype.Match(FileData)
	if err != nil {
		return FoundFile, err
	}

	// Compatible with MSSQL binary storage
	FileContent := hex.EncodeToString(FileData)
	FileContent = "0x" + FileContent

	FoundFile.Data = FileContent
	FoundFile.Type = FileMeta.Extension
	FoundFile.Mime = FileMeta.MIME.Value

	File, err := os.Open(FileName)
	defer File.Close()
	if err != nil {
		return FoundFile, err
	}

	head := make([]byte, 261)
	File.Read(head)

	// TODO Define err message
	if filetype.IsImage(head) == false {
		return FoundFile, err
	}

	// Only parse for supported functions
	switch FoundFile.Type {
	case "jpg":
		FoundFile.Width, FoundFile.Height = GetJpgDimensions(File)
	case "gif":
		FoundFile.Width, FoundFile.Height = GetGifDimensions(File)
	case "png":
		FoundFile.Width, FoundFile.Height = GetPngDimensions(File)
	case "bmp":
		FoundFile.Width, FoundFile.Height = GetBmpDimensions(File)
	default:
		return FoundFile, err
	}

	return FoundFile, nil
}

func GetJpgDimensions(file *os.File) (width int, height int) {
	fi, _ := file.Stat()
	fileSize := fi.Size()

	position := int64(4)
	bytes := make([]byte, 4)
	file.ReadAt(bytes[:2], position)
	length := int(bytes[0]<<8) + int(bytes[1])
	for position < fileSize {
		position += int64(length)
		file.ReadAt(bytes, position)
		length = int(bytes[2])<<8 + int(bytes[3])
		if (bytes[1] == 0xC0 || bytes[1] == 0xC2) && bytes[0] == 0xFF && length > 7 {
			file.ReadAt(bytes, position+5)
			width = int(bytes[2])<<8 + int(bytes[3])
			height = int(bytes[0])<<8 + int(bytes[1])
			return
		}
		position += 2
	}
	return 0, 0
}

func GetGifDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 4)
	file.ReadAt(bytes, 6)
	width = int(bytes[0]) + int(bytes[1])*256
	height = int(bytes[2]) + int(bytes[3])*256
	return
}

func GetBmpDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 8)
	file.ReadAt(bytes, 18)
	width = int(bytes[3])<<24 | int(bytes[2])<<16 | int(bytes[1])<<8 | int(bytes[0])
	height = int(bytes[7])<<24 | int(bytes[6])<<16 | int(bytes[5])<<8 | int(bytes[4])
	return
}

func GetPngDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 8)
	file.ReadAt(bytes, 16)
	width = int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])
	height = int(bytes[4])<<24 | int(bytes[5])<<16 | int(bytes[6])<<8 | int(bytes[7])
	return
}

func (p *JSONGOFPDF) GetStringIndex(name string, logic string, fallback string) (value string) {
	result := fallback
	attribute, _, _, err := p.GetAttributeIndex(name, logic, true)
	if err == nil {
		result = cast.ToString(attribute)
	}
	return result
}

func (p *JSONGOFPDF) GetAttribute(name string, logic string, debug bool) (value []byte, dataType jsonparser.ValueType, offset int, err error) {
	return p.GetAttributeIndex(name, logic, debug)
}

func (p *JSONGOFPDF) GetAttributeIndex(name string, logic string, debug bool) (value []byte, dataType jsonparser.ValueType, offset int, err error) {
	value, dataType, offset, err = jsonparser.Get([]byte(logic), name)
	// We should deliberately call this operation somehow... maybe if it is an object and you find a value called "func" or better we have a property called calculation which renders seperately
	// if dataType == jsonparser.Object {
	// 	value, dataType = p.RunObjectOperationsValue(string(value))
	// }
	return value, dataType, offset, err
}

// RunObjectOperationsValue entry point
func (p *JSONGOFPDF) RunObjectOperationsValue(logic string) (val []byte, dataType jsonparser.ValueType) {
	jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		val, dataType = p.RunValue(string(key), string(value))
		return nil
	})
	return val, dataType
}

func (p *JSONGOFPDF) RunValue(name string, logic string) (val []byte, dataType jsonparser.ValueType) {
	switch name {
	// case "logic":
	// 	// uses json-logic and can read from parser data
	// 	data := p.Parser.Parse()
	// 	result, _ := jsonlogic.Apply(logic, data)
	// 	switch v := result.(type) {
	// 	case bool:
	// 		return []byte(cast.ToString(v)), jsonparser.Boolean
	// 	case int:
	// 		return []byte(cast.ToString(v)), jsonparser.Number
	// 	case float64:
	// 		return []byte(cast.ToString(v)), jsonparser.Number
	// 	case string:
	// 		return []byte(cast.ToString(v)), jsonparser.String
	// 	default:
	// 		return []byte(cast.ToString(v)), jsonparser.String
	// 	}
	// 	break
	default:
		return nil, jsonparser.NotExist
	}

	return nil, jsonparser.NotExist
}
