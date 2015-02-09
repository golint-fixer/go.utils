package utils

import (
 "io"
 "io/ioutil"
 "os"
 "path/filepath"
)

// CopyFile copies `src` into `dst`. It preserves permissions of `src`.
// If `dst` is a directory, it copies `src` to this directory with source's
// filename. If `dst` points at `src`, no action is performed.
func CopyFile(dst, src string) error {
 if ok, err := IsTheSame(dst, src); err != nil {
   return err
 } else if ok {
   return nil
 }
 in, err := os.Open(src)
 if err != nil {
   return err
 }
 defer in.Close()
 fi, err := in.Stat()
 if err != nil {
   return err
 }
 tmp, err := ioutil.TempFile(filepath.Dir(dst), "")
 if err != nil {
   return err
 }
 if _, err = io.Copy(tmp, in); err != nil {
   tmp.Close()
   os.Remove(tmp.Name())
   return err
 }
 if err = tmp.Close(); err != nil {
   os.Remove(tmp.Name())
   return err
 }
 if dst, err = dstfix(dst, src); err != nil {
   os.Remove(tmp.Name())
   return err
 }
 if err = os.Rename(tmp.Name(), dst); err != nil {
   os.Remove(tmp.Name())
   return err
 }
 if err = os.Chmod(dst, fi.Mode()); err != nil {
   os.Remove(dst)
   return err
 }
 return nil
}

func dstfix(dst, src string) (string, error) {
 switch fi, err := os.Stat(dst); {
 case os.IsNotExist(err):
   return dst, nil
 case err != nil:
   return "", err
 case fi.IsDir():
   return filepath.Join(dst, filepath.Base(src)), nil
 default:
   return dst, nil
 }
}

// Exists returns a boolean indicating if `path` points to existing file/dir.
func Exists(path string) (bool, error) {
 _, err := os.Stat(path)
 if os.IsNotExist(err) {
   return false, nil
 }
 return true, err
}

// IsDir returns a boolean indicating if `path` points to a directory.
func IsDir(path string) (bool, error) {
 fi, err := os.Stat(path)
 if err != nil {
   return false, err
 }
 return fi.IsDir(), nil
}

// IsTheSame returns a boolean indicating if `lf` and `rf` points at the same
// file/dir without verifying if it really exists.
func IsTheSame(lf, rf string) (ok bool, err error) {
 if lf == rf {
   return true, nil
 }
 if lf, err = abs(lf); err != nil {
   return
 }
 if rf, err = abs(rf); err != nil {
   return
 }
 ok = lf == rf
 return
}

func abs(path string) (string, error) {
 if path = filepath.Clean(path); !filepath.IsAbs(path) {
   return filepath.Abs(path)
 }
 return path, nil
}
