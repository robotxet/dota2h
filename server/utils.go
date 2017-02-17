package server

func getFormat(bytes []byte) (string) {
  if len(bytes) < 4 { return "" }
  if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 { return "png" }
  if bytes[0] == 0xFF && bytes[1] == 0xD8 { return "jpg" }
  if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 { return "gif" }
  if bytes[0] == 0x42 && bytes[1] == 0x4D { return "bmp" }
  return ""
}