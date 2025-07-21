package systemchecks

import (
	"fmt"
	"os"
	"os/user"
	"syscall"

	"monitoring/types"
)

func CheckPermissions(cfg types.Config) interface{} {
	var allResults []map[string]interface{}

	for _, permBlock := range cfg.Permissions {
		if !permBlock.Enabled {
			continue

		}

		blockResult := map[string]interface{}{}

		users := make(map[string]*user.User)
		for _, username := range permBlock.CheckUserAccess {
			u, err := user.Lookup(username)
			if err != nil {
				blockResult[username] = map[string]interface{}{
					"error": fmt.Sprintf("could not find user %s: %v", username, err),
				}
				continue
			}
			users[username] = u
		}

		for _, path := range permBlock.Paths {
			info, err := os.Stat(path)
			if err != nil {
				blockResult[path] = map[string]interface{}{
					"error": err.Error(),
				}
				continue
			}

			mode := info.Mode()
			stat := info.Sys().(*syscall.Stat_t)

			entry := map[string]interface{}{
				"mode":        mode.String(),
			}

			if permBlock.ShowUserPermissions {
				ownerUID := stat.Uid
				groupGID := stat.Gid

				entry["owner_uid"] = ownerUID
				entry["group_gid"] = groupGID

				if ownerUser, err := user.LookupId(fmt.Sprint(ownerUID)); err == nil {
					entry["owner_user"] = ownerUser.Username
				}
				if group, err := user.LookupGroupId(fmt.Sprint(groupGID)); err == nil {
					entry["group_name"] = group.Name
				}
			}

			userAccess := map[string]map[string]bool{}
			for username, u := range users {
				userAccess[username] = checkAccess(stat, mode, u)
			}
			entry["user_access"] = userAccess

			blockResult[path] = entry
		}

		allResults = append(allResults, blockResult)
	}

	return allResults
}

func checkAccess(stat *syscall.Stat_t, mode os.FileMode, u *user.User) map[string]bool {
	uid := stat.Uid
	gid := stat.Gid

	userUID := u.Uid
	userGIDs, _ := u.GroupIds()

	userGroupMatch := false
	for _, g := range userGIDs {
		if g == fmt.Sprint(gid) {
			userGroupMatch = true
			break
		}
	}

	read, write, exec := false, false, false

	perm := mode.Perm()

	if fmt.Sprint(uid) == userUID {
		// Owner
		read = perm&0400 != 0
		write = perm&0200 != 0
		exec = perm&0100 != 0
	} else if userGroupMatch {
		// Group
		read = perm&0040 != 0
		write = perm&0020 != 0
		exec = perm&0010 != 0
	} else {
		// Other
		read = perm&0004 != 0
		write = perm&0002 != 0
		exec = perm&0001 != 0
	}

	return map[string]bool{
		"read":  read,
		"write": write,
		"exec":  exec,
	}
}
