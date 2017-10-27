package matchers

var (
	TypeMp4  = newType("mp4", "video/mp4")
	TypeM4v  = newType("m4v", "video/x-m4v")
	TypeMkv  = newType("mkv", "video/x-matroska")
	TypeWebm = newType("webm", "video/webm")
	TypeMov  = newType("mov", "video/quicktime")
	TypeAvi  = newType("avi", "video/x-msvideo")
	TypeWmv  = newType("wmv", "video/x-ms-wmv")
	TypeMpeg = newType("mpg", "video/mpeg")
	TypeFlv  = newType("flv", "video/x-flv")
)

var Video = Map{
	TypeMp4:  Mp4,
	TypeM4v:  M4v,
	TypeMkv:  Mkv,
	TypeMov:  Mov,
	TypeAvi:  Avi,
	TypeWmv:  Wmv,
	TypeMpeg: Mpeg,
	TypeFlv:  Flv,
}

func M4v(buf []byte) bool {
	return len(buf) > 10 &&
		buf[0] == 0x0 && buf[1] == 0x0 &&
		buf[2] == 0x0 && buf[3] == 0x1C &&
		buf[4] == 0x66 && buf[5] == 0x74 &&
		buf[6] == 0x79 && buf[7] == 0x70 &&
		buf[8] == 0x4D && buf[9] == 0x34 &&
		buf[10] == 0x56
}

func Mkv(buf []byte) bool {
	return (len(buf) > 15 &&
		buf[0] == 0x1A && buf[1] == 0x45 &&
		buf[2] == 0xDF && buf[3] == 0xA3 &&
		buf[4] == 0x93 && buf[5] == 0x42 &&
		buf[6] == 0x82 && buf[7] == 0x88 &&
		buf[8] == 0x6D && buf[9] == 0x61 &&
		buf[10] == 0x74 && buf[11] == 0x72 &&
		buf[12] == 0x6F && buf[13] == 0x73 &&
		buf[14] == 0x6B && buf[15] == 0x61) ||
		(len(buf) > 38 &&
			buf[31] == 0x6D && buf[32] == 0x61 &&
			buf[33] == 0x74 && buf[34] == 0x72 &&
			buf[35] == 0x6f && buf[36] == 0x73 &&
			buf[37] == 0x6B && buf[38] == 0x61)
}

func Webm(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x1A && buf[1] == 0x45 &&
		buf[2] == 0xDF && buf[3] == 0xA3
}

func Mov(buf []byte) bool {
	return len(buf) > 7 &&
		buf[0] == 0x0 && buf[1] == 0x0 &&
		buf[2] == 0x0 && buf[3] == 0x14 &&
		buf[4] == 0x66 && buf[5] == 0x74 &&
		buf[6] == 0x79 && buf[7] == 0x70
}

func Avi(buf []byte) bool {
	return len(buf) > 10 &&
		buf[0] == 0x52 && buf[1] == 0x49 &&
		buf[2] == 0x46 && buf[3] == 0x46 &&
		buf[8] == 0x41 && buf[9] == 0x56 &&
		buf[10] == 0x49
}

func Wmv(buf []byte) bool {
	return len(buf) > 9 &&
		buf[0] == 0x30 && buf[1] == 0x26 &&
		buf[2] == 0xB2 && buf[3] == 0x75 &&
		buf[4] == 0x8E && buf[5] == 0x66 &&
		buf[6] == 0xCF && buf[7] == 0x11 &&
		buf[8] == 0xA6 && buf[9] == 0xD9
}

func Mpeg(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x0 && buf[1] == 0x0 &&
		buf[2] == 0x1 && buf[3] >= 0xb0 &&
		buf[3] <= 0xbf
}

func Flv(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x46 && buf[1] == 0x4C &&
		buf[2] == 0x56 && buf[3] == 0x01
}

func Mp4(buf []byte) bool {
	return len(buf) > 27 &&
		(buf[0] == 0x0 && buf[1] == 0x0 && buf[2] == 0x0 &&
			((buf[3] == 0x18 || buf[3] == 0x20) && buf[4] == 0x66 &&
				buf[5] == 0x74 && buf[6] == 0x79 && buf[7] == 0x70) ||
			(buf[0] == 0x33 && buf[1] == 0x67 && buf[2] == 0x70 && buf[3] == 0x35) ||
			(buf[0] == 0x0 && buf[1] == 0x0 && buf[2] == 0x0 && buf[3] == 0x1C &&
				buf[4] == 0x66 && buf[5] == 0x74 && buf[6] == 0x79 && buf[7] == 0x70 &&
				buf[8] == 0x6D && buf[9] == 0x70 && buf[10] == 0x34 && buf[11] == 0x32 &&
				buf[16] == 0x6D && buf[17] == 0x70 && buf[18] == 0x34 && buf[19] == 0x31 &&
				buf[20] == 0x6D && buf[21] == 0x70 && buf[22] == 0x34 && buf[23] == 0x32 &&
				buf[24] == 0x69 && buf[25] == 0x73 && buf[26] == 0x6F && buf[27] == 0x6D))
}
