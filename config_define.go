package excelstreamer

type (
	ExcelDefine struct {
		Sheets []SheetDefine `yaml:"Sheets"`
	}

	SheetDefine struct {
		SheetName string        `yaml:"SheetName"`
		Anchor    string        `yaml:"Anchor"`
		Hide      bool          `yaml:"Hide"`
		Fields    []FieldDefine `yaml:"Fields"`
	}

	FieldDefine struct {
		Source      string `yaml:"Source"`
		Header      string `yaml:"Header"`
		Width       int    `yaml:"Width"`
		Format      string `yaml:"Format"`
		HeaderStyle string `yaml:"HeaderStyle"`
		DataStyle   string `yaml:"DataStyle"`
	}
)
