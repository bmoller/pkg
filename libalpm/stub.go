package libalpm

/*
   #cgo pkg-config: libalpm
   #include <string.h>
   #include <alpm.h>

   int compare_packages(const void *a, const void *b)
   {
      return strcmp(alpm_pkg_get_name(a), alpm_pkg_get_name(b));
   }
*/
import "C"
import (
	"fmt"
)

func GetForeignPackages(root, dbPath string, repos []string) (pkgs map[string]string, err error) {
	var e *C.alpm_errno_t
	handle := C.alpm_initialize(C.CString(root), C.CString(dbPath), e)
	if e != nil {
		return nil, fmt.Errorf("failed to initialize alpm handle: %d", int(*e))
	}
	defer C.alpm_release(handle)

	localDB := C.alpm_get_localdb(handle)
	foreign_packages := C.alpm_db_get_pkgcache(localDB)

	for _, repo := range repos {
		syncDB := C.alpm_register_syncdb(handle, C.CString(repo), C.ALPM_SIG_PACKAGE&C.ALPM_SIG_DATABASE_OPTIONAL)
		pkgCache := C.alpm_db_get_pkgcache(syncDB)
		foreign_packages = C.alpm_list_diff(foreign_packages, pkgCache, (*[0]byte)(C.compare_packages))
	}

	i := foreign_packages
	pkgs = make(map[string]string)
	for i != nil {
		name := C.GoString(C.alpm_pkg_get_name((*C.alpm_pkg_t)(i.data)))
		version := C.GoString(C.alpm_pkg_get_version((*C.alpm_pkg_t)(i.data)))
		pkgs[name] = version
		i = i.next
	}
	// only attempt to free list if it is not the original package cache reference
	if len(repos) > 0 {
		C.alpm_list_free(foreign_packages)
	}

	return
}

/*
CompareVersions uses libalpm's logic to compare two versions of arbitrary formats.
The return value is negative when a is less than b, zero when a and b are equal, and positive when b is greater than a.
*/
func CompareVersions(a, b string) int {
	return int(C.alpm_pkg_vercmp(C.CString(a), C.CString(b)))
}
