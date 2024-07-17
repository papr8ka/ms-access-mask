package object

var (
	GlobalDescription = []string{
		"Generic Read",
		"Generic Write",
		"Generic Execute",
		"Generic All",
		"", //"Reserved",
		"", //"Reserved",
		"", //"Reserved",
		"Access ACL",
		"",
		"",
		"",
		"Synchronize",
		"Write Owner",
		"Write DAC",
		"Read Control",
		"Delete",
	}

	GlobalMask = uint32(0xF11F)
)
