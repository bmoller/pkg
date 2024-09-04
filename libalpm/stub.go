/*
 * stub.go
 *
 * Copyright (c) 2024 Brandon Moller
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package libalpm

/*
   #cgo pkg-config: libalpm
   #include <alpm.h>
*/
import "C"
import (
	"fmt"
)

/*
GetLocalPackages retrieves the list of locally-installed pkgs.
Keys are package names and values are their versions.
If an error is encountered it is returned in err.
*/
func GetLocalPackages(root, dbPath string) (pkgs map[string]string, err error) {
	e := C.alpm_errno_t(0)
	handle := C.alpm_initialize(C.CString(root), C.CString(dbPath), &e)
	if e != C.ALPM_ERR_OK {
		return nil, fmt.Errorf("failed to initialize alpm handle: %d", int(e))
	}
	defer C.alpm_release(handle)

	localDB := C.alpm_get_localdb(handle)
	pkgcache := C.alpm_db_get_pkgcache(localDB)
	pkgs = make(map[string]string)
	for pkg := pkgcache; pkg != nil; pkg = pkg.next {
		name := C.GoString(C.alpm_pkg_get_name((*C.alpm_pkg_t)(pkg.data)))
		version := C.GoString(C.alpm_pkg_get_version((*C.alpm_pkg_t)(pkg.data)))
		pkgs[name] = version
	}

	return
}

/*
GetSyncPackages retrieves the list of known packages for the specified repos.
Keys are package names and values are their versions.
If an error is encountered it is returned in err.
*/
func GetSyncPackages(root, dbPath string, repos []string) (pkgs map[string]string, err error) {
	e := C.alpm_errno_t(0)
	handle := C.alpm_initialize(C.CString(root), C.CString(dbPath), &e)
	if e != C.ALPM_ERR_OK {
		return nil, fmt.Errorf("failed to initialize alpm handle: %d", int(e))
	}
	defer C.alpm_release(handle)

	pkgs = make(map[string]string)
	for _, repo := range repos {
		syncDB := C.alpm_register_syncdb(handle, C.CString(repo), C.ALPM_SIG_PACKAGE&C.ALPM_SIG_DATABASE_OPTIONAL)
		pkgcache := C.alpm_db_get_pkgcache(syncDB)

		for pkg := pkgcache; pkg != nil; pkg = pkg.next {
			name := C.GoString(C.alpm_pkg_get_name((*C.alpm_pkg_t)(pkg.data)))
			version := C.GoString(C.alpm_pkg_get_name((*C.alpm_pkg_t)(pkg.data)))
			pkgs[name] = version
		}
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
