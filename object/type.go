package object

type Type int

const (
	File Type = iota
	Directory
	Pipe

	Types
)

var (
	Name = [Types]string{
		"File",
		"Directory",
		"Pipe",
	}
	Description = [Types][]string{
		append(GlobalDescription, FileDescription...),
		append(GlobalDescription, DirectoryDescription...),
		append(GlobalDescription, PipeDescription...),
	}
	Mask = [Types]uint32{
		GlobalMask<<16 | FileMask,
		GlobalMask<<16 | DirectoryMask,
		GlobalMask<<16 | PipeMask,
	}
)
